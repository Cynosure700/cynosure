from __future__ import annotations

import copy
import json
from pathlib import Path
from typing import Iterable, Mapping


REQUIRED_SETTINGS_ENV = ("open_auth_token", "open_model", "open_base_url")


def load_cynosure_settings(path: Path) -> dict:
    settings = json.loads(path.read_text())
    env = settings.get("env")
    if not isinstance(env, dict):
        raise ValueError(f"{path} must contain an env object")
    missing = [key for key in REQUIRED_SETTINGS_ENV if not str(env.get(key, "")).strip()]
    if missing:
        raise ValueError(f"{path} is missing env fields: {', '.join(missing)}")
    return settings


def redact_settings(settings: Mapping) -> dict:
    redacted = copy.deepcopy(dict(settings))
    env = redacted.get("env")
    if isinstance(env, dict) and "open_auth_token" in env:
        env["open_auth_token"] = "***"
    return redacted


def write_predictions_jsonl(path: Path, rows: Iterable[Mapping[str, str]]) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    with path.open("w") as f:
        for row in rows:
            payload = {
                "instance_id": row["instance_id"],
                "model_name_or_path": row["model_name_or_path"],
                "model_patch": row.get("model_patch", ""),
            }
            f.write(json.dumps(payload, ensure_ascii=False) + "\n")


def write_bypass_permissions(workspace_root: Path) -> None:
    settings_path = workspace_root / ".cynosure" / "settings.json"
    if settings_path.exists():
        payload = json.loads(settings_path.read_text())
    else:
        payload = {}
    if not isinstance(payload, dict):
        payload = {}
    permissions = payload.get("permissions")
    if not isinstance(permissions, dict):
        permissions = {}
    permissions["defaultMode"] = "bypassPermissions"
    payload["permissions"] = permissions
    settings_path.parent.mkdir(parents=True, exist_ok=True)
    settings_path.write_text(json.dumps(payload, ensure_ascii=False, indent=2))


def build_swebench_prompt(
    *,
    instance_id: str,
    problem_statement: str,
    hints_text: str = "",
) -> str:
    hints = hints_text.strip() or "(none)"
    return f"""You are evaluating Cynosure on SWE-bench Lite instance {instance_id}.

Fix the bug described below in the current repository. Make the smallest correct code change.
Run the relevant tests when practical. Do not commit changes.
At the end, briefly summarize what changed and mention the git diff is ready.

Problem statement:
{problem_statement.strip()}

Hints:
{hints}
"""
