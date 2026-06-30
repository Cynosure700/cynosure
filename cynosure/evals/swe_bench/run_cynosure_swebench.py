#!/usr/bin/env python3
from __future__ import annotations

import argparse
import json
import os
import signal
import shutil
import subprocess
import sys
import time
from pathlib import Path

from datasets import load_dataset

ROOT = Path(__file__).resolve().parents[2]
LIB = ROOT / "evals" / "lib"
sys.path.insert(0, str(LIB))

from eval_utils import (  # noqa: E402
    build_swebench_prompt,
    load_cynosure_settings,
    redact_settings,
    write_bypass_permissions,
    write_predictions_jsonl,
)


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Run Cynosure inference on SWE-bench Lite cases.")
    parser.add_argument("--binary", required=True, type=Path, help="Path to cynosure-eval-agent binary")
    parser.add_argument("--settings", default=Path.home() / ".cynosure" / "settings.json", type=Path)
    parser.add_argument("--output-dir", required=True, type=Path)
    parser.add_argument("--n", default=30, type=int)
    parser.add_argument("--timeout-sec", default=3600, type=int)
    parser.add_argument("--keep-workspaces", action="store_true")
    parser.add_argument("--resume", action="store_true")
    parser.add_argument("--max-patch-bytes", default=250_000, type=int)
    parser.add_argument("--dataset-name", default="princeton-nlp/SWE-bench_Lite")
    return parser.parse_args()


def run(cmd: list[str], *, cwd: Path | None = None, env: dict[str, str] | None = None, timeout: int | None = None) -> subprocess.CompletedProcess:
    return subprocess.run(cmd, cwd=cwd, env=env, timeout=timeout, text=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)


def run_agent(cmd: list[str], *, env: dict[str, str], timeout: int, log_path: Path) -> tuple[int | None, bool]:
    log_path.parent.mkdir(parents=True, exist_ok=True)
    log_file = log_path.open("w")
    proc = subprocess.Popen(
        cmd,
        env=env,
        text=True,
        stdout=log_file,
        stderr=subprocess.STDOUT,
        start_new_session=True,
    )
    deadline = time.monotonic() + timeout
    while proc.poll() is None and time.monotonic() < deadline:
        time.sleep(1)
    timed_out = proc.poll() is None
    if timed_out:
        os.killpg(proc.pid, signal.SIGTERM)
        grace_deadline = time.monotonic() + 10
        while proc.poll() is None and time.monotonic() < grace_deadline:
            time.sleep(0.2)
        if proc.poll() is None:
            os.killpg(proc.pid, signal.SIGKILL)
    proc.wait()
    log_file.close()
    return proc.returncode, timed_out


def prepare_instance(instance: dict, workspace: Path) -> Path:
    repo_dir = workspace / "repo"
    repo = instance["repo"]
    base_commit = instance["base_commit"]
    url = f"https://github.com/{repo}.git"
    clone = run(["git", "clone", "--filter=blob:none", url, str(repo_dir)], timeout=900)
    if clone.returncode != 0:
        raise RuntimeError(f"git clone failed for {repo}: {clone.stdout}")
    checkout = run(["git", "checkout", base_commit], cwd=repo_dir, timeout=300)
    if checkout.returncode != 0:
        raise RuntimeError(f"git checkout failed for {instance['instance_id']}: {checkout.stdout}")
    write_bypass_permissions(repo_dir)
    return repo_dir


def run_instance(binary: Path, settings: Path, output_dir: Path, instance: dict, timeout_sec: int, keep_workspace: bool, max_patch_bytes: int) -> dict:
    instance_id = instance["instance_id"]
    case_dir = output_dir / "cases" / instance_id
    workspace = case_dir / "workspace"
    case_dir.mkdir(parents=True, exist_ok=True)
    if workspace.exists():
        shutil.rmtree(workspace)
    workspace.mkdir(parents=True)

    started = time.time()
    result = {
        "instance_id": instance_id,
        "repo": instance["repo"],
        "base_commit": instance["base_commit"],
        "completed": False,
        "patch_path": str(case_dir / "patch.diff"),
        "transcript_path": str(case_dir / "agent.log"),
        "duration_sec": None,
        "error": "",
    }
    try:
        repo_dir = prepare_instance(instance, workspace)
        prompt = build_swebench_prompt(
            instance_id=instance_id,
            problem_statement=instance["problem_statement"],
            hints_text=instance.get("hints_text", ""),
        )
        prompt_path = case_dir / "prompt.txt"
        prompt_path.write_text(prompt)
        home = case_dir / "home"
        (home / ".cynosure").mkdir(parents=True, exist_ok=True)
        shutil.copy2(settings, home / ".cynosure" / "settings.json")
        env = os.environ.copy()
        env["HOME"] = str(home)
        returncode, timed_out = run_agent(
            [str(binary), "--cwd", str(repo_dir), "--prompt-file", str(prompt_path)],
            env=env,
            timeout=timeout_sec,
            log_path=case_dir / "agent.log",
        )
        if timed_out:
            result["error"] = f"agent timed out after {timeout_sec}s"
        elif returncode != 0:
            result["error"] = f"agent exited {returncode}"
        diff = run(["git", "-c", "core.fileMode=false", "diff"], cwd=repo_dir, timeout=120)
        patch = diff.stdout if diff.returncode == 0 else ""
        if len(patch.encode()) > max_patch_bytes:
            (case_dir / "patch.too-large.diff").write_text(patch)
            patch = ""
            result["error"] = (result["error"] + "; " if result["error"] else "") + "patch exceeded max size"
        (case_dir / "patch.diff").write_text(patch)
        result["completed"] = returncode == 0
        result["model_patch"] = patch
    except subprocess.TimeoutExpired as exc:
        result["error"] = f"timeout after {exc.timeout}s"
        if "repo_dir" in locals():
            diff = run(["git", "-c", "core.fileMode=false", "diff"], cwd=repo_dir, timeout=120)
            patch = diff.stdout if diff.returncode == 0 else ""
            if len(patch.encode()) > max_patch_bytes:
                (case_dir / "patch.too-large.diff").write_text(patch)
                patch = ""
            result["model_patch"] = patch
        else:
            result["model_patch"] = ""
    except Exception as exc:
        result["error"] = str(exc)
        result["model_patch"] = ""
    finally:
        result["duration_sec"] = round(time.time() - started, 2)
        if not keep_workspace and workspace.exists():
            shutil.rmtree(workspace)
    return result


def main() -> int:
    args = parse_args()
    if not args.binary.is_file():
        raise SystemExit(f"binary not found: {args.binary}")
    settings = load_cynosure_settings(args.settings)
    args.output_dir.mkdir(parents=True, exist_ok=True)
    (args.output_dir / "settings.redacted.json").write_text(json.dumps(redact_settings(settings), ensure_ascii=False, indent=2))

    dataset = load_dataset(args.dataset_name, split="test")
    selected = [dataset[i] for i in range(args.n)]
    selected_path = args.output_dir / "selected_instances.json"
    selected_path.write_text(json.dumps([row["instance_id"] for row in selected], indent=2))

    rows = []
    summaries = []
    completed_ids = set()
    if args.resume and (args.output_dir / "predictions.jsonl").exists():
        with (args.output_dir / "predictions.jsonl").open() as f:
            for line in f:
                if not line.strip():
                    continue
                row = json.loads(line)
                rows.append(row)
                completed_ids.add(row["instance_id"])
    if args.resume and (args.output_dir / "inference_summary.json").exists():
        summaries = json.loads((args.output_dir / "inference_summary.json").read_text())
    model_name = f"cynosure/{settings['env']['open_model']}"
    for idx, instance in enumerate(selected, start=1):
        print(f"[{idx}/{len(selected)}] {instance['instance_id']}", flush=True)
        if instance["instance_id"] in completed_ids:
            print(f"  skipping completed {instance['instance_id']}", flush=True)
            continue
        summary = run_instance(args.binary, args.settings, args.output_dir, instance, args.timeout_sec, args.keep_workspaces, args.max_patch_bytes)
        summaries.append({k: v for k, v in summary.items() if k != "model_patch"})
        rows.append(
            {
                "instance_id": instance["instance_id"],
                "model_name_or_path": model_name,
                "model_patch": summary.get("model_patch", ""),
            }
        )
        write_predictions_jsonl(args.output_dir / "predictions.jsonl", rows)
        (args.output_dir / "inference_summary.json").write_text(json.dumps(summaries, ensure_ascii=False, indent=2))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
