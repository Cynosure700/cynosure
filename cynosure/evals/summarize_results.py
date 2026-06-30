#!/usr/bin/env python3
from __future__ import annotations

import json
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]


def load_terminal_rows() -> list[dict]:
    rows = []
    for run in ["tb30-20260629", "tb30b-20260629"]:
        path = ROOT / "evals" / "runs" / run / "terminal_bench" / run / "results.json"
        data = json.loads(path.read_text())
        for item in data["results"]:
            rows.append(
                {
                    "run_id": run,
                    "task_id": item["task_id"],
                    "resolved": item["is_resolved"],
                    "failure_mode": item["failure_mode"],
                    "agent_started_at": item.get("agent_started_at"),
                    "agent_ended_at": item.get("agent_ended_at"),
                    "test_started_at": item.get("test_started_at"),
                    "test_ended_at": item.get("test_ended_at"),
                }
            )
    return rows


def load_swe_rows() -> tuple[list[dict], list[dict], dict]:
    swe_root = ROOT / "evals" / "runs" / "swe30-20260629" / "swe_bench"
    inference = json.loads((swe_root / "inference_summary.json").read_text())
    inference_rows = []
    for item in inference:
        patch_path = swe_root / "cases" / item["instance_id"] / "patch.diff"
        inference_rows.append(
            {
                "instance_id": item["instance_id"],
                "agent_completed": item["completed"],
                "duration_sec": item["duration_sec"],
                "patch_bytes": patch_path.stat().st_size if patch_path.exists() else 0,
                "error": item.get("error", ""),
            }
        )

    official_rows = []
    report_root = ROOT / "logs" / "run_evaluation" / "swe30-20260629" / "cynosure__glm-5.2"
    for report_path in sorted(report_root.glob("*/report.json")):
        data = json.loads(report_path.read_text())
        instance_id = next(iter(data))
        report = data[instance_id]
        official_rows.append(
            {
                "instance_id": instance_id,
                "patch_applied": report.get("patch_successfully_applied"),
                "resolved": report.get("resolved"),
            }
        )

    prebuilt = json.loads((ROOT / "cynosure__glm-5.2.swe30-20260629-prebuilt.json").read_text())
    return inference_rows, official_rows, prebuilt


def main() -> int:
    terminal_rows = load_terminal_rows()
    swe_inference, swe_official, swe_prebuilt = load_swe_rows()
    summary = {
        "model": "glm-5.2",
        "terminal_bench": {
            "dataset": "terminal-bench-core==0.1.1",
            "cases": len(terminal_rows),
            "resolved": sum(1 for row in terminal_rows if row["resolved"] is True),
            "unresolved": sum(1 for row in terminal_rows if row["resolved"] is not True),
            "rows": terminal_rows,
        },
        "swe_bench_lite": {
            "dataset": "princeton-nlp/SWE-bench_Lite",
            "selected_cases": 30,
            "predictions": len(swe_inference),
            "agent_completed": sum(1 for row in swe_inference if row["agent_completed"]),
            "non_empty_patches": sum(1 for row in swe_inference if row["patch_bytes"] > 0),
            "official_reports_completed": len(swe_official),
            "official_reports_resolved": sum(1 for row in swe_official if row["resolved"] is True),
            "prebuilt_harness_summary": swe_prebuilt,
            "inference_rows": swe_inference,
            "official_rows": swe_official,
        },
    }
    out = ROOT / "evals" / "runs" / "benchmark_summary_20260629.json"
    out.write_text(json.dumps(summary, ensure_ascii=False, indent=2))
    print(out)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
