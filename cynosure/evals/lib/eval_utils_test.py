import json
import os
import sys
import tempfile
import unittest
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent))

from eval_utils import (
    build_swebench_prompt,
    load_cynosure_settings,
    redact_settings,
    write_bypass_permissions,
    write_predictions_jsonl,
)


class EvalUtilsTest(unittest.TestCase):
    def test_load_cynosure_settings_requires_env_fields(self):
        with tempfile.TemporaryDirectory() as tmp:
            path = Path(tmp) / "settings.json"
            path.write_text(
                json.dumps(
                    {
                        "env": {
                            "open_auth_token": "token",
                            "open_model": "model",
                            "open_base_url": "https://example.test",
                        }
                    }
                )
            )

            settings = load_cynosure_settings(path)

        self.assertEqual(settings["env"]["open_auth_token"], "token")
        self.assertEqual(settings["env"]["open_model"], "model")
        self.assertEqual(settings["env"]["open_base_url"], "https://example.test")

    def test_redact_settings_removes_secret_but_keeps_model_context(self):
        settings = {
            "env": {
                "open_auth_token": "secret-token",
                "open_model": "glm-5.2",
                "open_base_url": "https://example.test",
            }
        }

        redacted = redact_settings(settings)

        self.assertEqual(redacted["env"]["open_auth_token"], "***")
        self.assertEqual(redacted["env"]["open_model"], "glm-5.2")
        self.assertEqual(redacted["env"]["open_base_url"], "https://example.test")

    def test_write_predictions_jsonl_uses_swebench_contract(self):
        rows = [
            {
                "instance_id": "repo__name-1",
                "model_name_or_path": "cynosure/glm-5.2",
                "model_patch": "diff --git a/a.py b/a.py\n",
            }
        ]
        with tempfile.TemporaryDirectory() as tmp:
            path = Path(tmp) / "predictions.jsonl"
            write_predictions_jsonl(path, rows)
            payload = json.loads(path.read_text().strip())

        self.assertEqual(payload["instance_id"], "repo__name-1")
        self.assertEqual(payload["model_name_or_path"], "cynosure/glm-5.2")
        self.assertEqual(payload["model_patch"], "diff --git a/a.py b/a.py\n")

    def test_build_swebench_prompt_contains_issue_and_required_exit_contract(self):
        prompt = build_swebench_prompt(
            instance_id="astropy__astropy-12907",
            problem_statement="Fix separability_matrix for nested compound models.",
            hints_text="Look at separable.py",
        )

        self.assertIn("astropy__astropy-12907", prompt)
        self.assertIn("Fix separability_matrix", prompt)
        self.assertIn("Look at separable.py", prompt)
        self.assertIn("Do not commit", prompt)
        self.assertIn("git diff", prompt)

    def test_write_bypass_permissions_preserves_existing_workspace_settings(self):
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp)
            settings_path = root / ".cynosure" / "settings.json"
            settings_path.parent.mkdir()
            settings_path.write_text(json.dumps({"other": {"keep": True}}))

            write_bypass_permissions(root)

            payload = json.loads(settings_path.read_text())

        self.assertEqual(payload["other"], {"keep": True})
        self.assertEqual(payload["permissions"]["defaultMode"], "bypassPermissions")


if __name__ == "__main__":
    unittest.main()
