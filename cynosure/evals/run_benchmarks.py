#!/usr/bin/env python3
from __future__ import annotations

import argparse
import json
import os
import subprocess
import sys
import time
from datetime import datetime
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
LIB = ROOT / "evals" / "lib"
sys.path.insert(0, str(LIB))

from eval_utils import load_cynosure_settings, redact_settings  # noqa: E402


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Run Cynosure Terminal-Bench and SWE-bench Lite evaluations.")
    parser.add_argument("--run-id", default=datetime.now().strftime("%Y%m%d-%H%M%S"))
    parser.add_argument("--terminal-n", default=30, type=int)
    parser.add_argument("--terminal-task-id", action="append", default=[])
    parser.add_argument("--swe-n", default=30, type=int)
    parser.add_argument("--settings", default=Path.home() / ".cynosure" / "settings.json", type=Path)
    parser.add_argument("--skip-terminal", action="store_true")
    parser.add_argument("--skip-swe", action="store_true")
    parser.add_argument("--dry-run", action="store_true")
    parser.add_argument("--terminal-timeout-sec", default=1800, type=int)
    parser.add_argument("--terminal-test-timeout-sec", default=900, type=int)
    parser.add_argument("--swe-timeout-sec", default=3600, type=int)
    return parser.parse_args()


def run(cmd: list[str], *, cwd: Path, env: dict[str, str] | None = None) -> int:
    print("$ " + " ".join(str(part) for part in cmd), flush=True)
    proc = subprocess.run(cmd, cwd=cwd, env=env)
    return proc.returncode


def docker_goarch() -> str:
    try:
        proc = subprocess.run(
            ["docker", "info", "--format", "{{.Architecture}}"],
            text=True,
            stdout=subprocess.PIPE,
            stderr=subprocess.DEVNULL,
            check=True,
        )
    except Exception:
        return "amd64"
    arch = proc.stdout.strip().lower()
    if arch in {"aarch64", "arm64"}:
        return "arm64"
    return "amd64"


def main() -> int:
    args = parse_args()
    settings = load_cynosure_settings(args.settings)
    run_root = ROOT / "evals" / "runs" / args.run_id
    run_root.mkdir(parents=True, exist_ok=True)
    (run_root / "settings.redacted.json").write_text(json.dumps(redact_settings(settings), ensure_ascii=False, indent=2))
    metadata = {
        "run_id": args.run_id,
        "started_at": datetime.now().isoformat(),
        "terminal_n": args.terminal_n,
        "swe_n": args.swe_n,
        "model": settings["env"]["open_model"],
        "base_url": settings["env"]["open_base_url"],
        "commands": [],
    }

    goarch = docker_goarch()
    host_binary = run_root / "cynosure-eval-agent-host"
    linux_binary = run_root / f"cynosure-eval-agent-linux-{goarch}"
    host_build_cmd = [
        "go",
        "build",
        "-o",
        str(host_binary),
        "./evals/cmd/cynosure-eval-agent",
    ]
    linux_build_cmd = [
        "go",
        "build",
        "-o",
        str(linux_binary),
        "./evals/cmd/cynosure-eval-agent",
    ]
    env = os.environ.copy()
    linux_env = os.environ.copy()
    linux_env.update({"GOOS": "linux", "GOARCH": goarch, "CGO_ENABLED": "0"})
    metadata["commands"].append({"name": "build-host-agent", "cmd": host_build_cmd})
    metadata["commands"].append({"name": "build-linux-agent", "cmd": linux_build_cmd})
    if not args.dry_run and run(host_build_cmd, cwd=ROOT, env=env) != 0:
        return 1
    if not args.dry_run and run(linux_build_cmd, cwd=ROOT, env=linux_env) != 0:
        return 1

    if not args.skip_terminal:
        terminal_output = run_root / "terminal_bench"
        cmd = [
            "tb",
            "run",
            "--dataset",
            "terminal-bench-core==0.1.1",
            "--n-tasks",
            str(args.terminal_n),
            "--n-concurrent",
            "1",
            "--output-path",
            str(terminal_output),
            "--run-id",
            args.run_id,
            "--agent-import-path",
            "cynosure_agent:CynosureTerminalBenchAgent",
            "--agent-kwarg",
            f"binary_path={linux_binary}",
            "--agent-kwarg",
            f"settings_path={args.settings}",
            "--agent-kwarg",
            f"timeout_sec={args.terminal_timeout_sec}",
            "--global-test-timeout-sec",
            str(args.terminal_test_timeout_sec),
            "--cleanup",
            "--no-upload-results",
        ]
        if args.terminal_task_id:
            task_args = []
            for task_id in args.terminal_task_id:
                task_args.extend(["--task-id", task_id])
            idx = cmd.index("--n-tasks")
            del cmd[idx : idx + 2]
            cmd[idx:idx] = task_args
        terminal_env = os.environ.copy()
        terminal_env["PYTHONPATH"] = str(ROOT / "evals" / "terminal_bench")
        metadata["commands"].append({"name": "terminal-bench", "cmd": cmd})
        if not args.dry_run and run(cmd, cwd=ROOT, env=terminal_env) != 0:
            return 1

    if not args.skip_swe:
        swe_output = run_root / "swe_bench"
        inference_cmd = [
            sys.executable,
            str(ROOT / "evals" / "swe_bench" / "run_cynosure_swebench.py"),
            "--binary",
            str(host_binary),
            "--settings",
            str(args.settings),
            "--output-dir",
            str(swe_output),
            "--n",
            str(args.swe_n),
            "--timeout-sec",
            str(args.swe_timeout_sec),
        ]
        metadata["commands"].append({"name": "swe-inference", "cmd": inference_cmd})
        if args.dry_run:
            selected = [f"dry-run-instance-{i}" for i in range(1, args.swe_n + 1)]
        elif run(inference_cmd, cwd=ROOT) != 0:
            return 1
        else:
            selected = json.loads((swe_output / "selected_instances.json").read_text())
        eval_cmd = [
            sys.executable,
            "-m",
            "swebench.harness.run_evaluation",
            "--dataset_name",
            "princeton-nlp/SWE-bench_Lite",
            "--predictions_path",
            str(swe_output / "predictions.jsonl"),
            "--instance_ids",
            *selected,
            "--max_workers",
            "1",
            "--timeout",
            "1800",
            "--run_id",
            args.run_id,
            "--report_dir",
            str(swe_output),
            "--clean",
            "true",
            "--namespace",
            "none",
        ]
        metadata["commands"].append({"name": "swe-evaluation", "cmd": eval_cmd})
        if not args.dry_run and run(eval_cmd, cwd=ROOT) != 0:
            return 1

    metadata["finished_at"] = datetime.now().isoformat()
    (run_root / "run_metadata.json").write_text(json.dumps(metadata, ensure_ascii=False, indent=2))
    return 0


if __name__ == "__main__":
    started = time.time()
    code = main()
    print(f"elapsed_sec={time.time() - started:.2f}")
    raise SystemExit(code)
