import tempfile
import unittest
from pathlib import Path

from cynosure_agent import CynosureTerminalBenchAgent


class CynosureTerminalBenchAgentTest(unittest.TestCase):
    def test_run_command_uses_copied_binary_and_instruction_file(self):
        with tempfile.TemporaryDirectory() as tmp:
            binary = Path(tmp) / "cynosure-linux-amd64"
            settings = Path(tmp) / "settings.json"
            binary.write_text("binary")
            settings.write_text(
                '{"env":{"open_auth_token":"secret","open_model":"model","open_base_url":"https://example.test"}}'
            )

            agent = CynosureTerminalBenchAgent(
                binary_path=str(binary),
                settings_path=str(settings),
                timeout_sec=123,
            )
            commands = agent._run_agent_commands("fix it")

        self.assertEqual(len(commands), 1)
        command = commands[0].command
        self.assertIn("/installed-agent/cynosure-eval-agent", command)
        self.assertIn('--cwd "$PWD"', command)
        self.assertIn("--prompt-file /installed-agent/instruction.txt", command)
        self.assertEqual(commands[0].max_timeout_sec, 123)
        self.assertNotIn("secret", command)

    def test_missing_binary_is_rejected(self):
        with tempfile.TemporaryDirectory() as tmp:
            settings = Path(tmp) / "settings.json"
            settings.write_text(
                '{"env":{"open_auth_token":"secret","open_model":"model","open_base_url":"https://example.test"}}'
            )

            with self.assertRaises(FileNotFoundError):
                CynosureTerminalBenchAgent(
                    binary_path=str(Path(tmp) / "missing"),
                    settings_path=str(settings),
                )


if __name__ == "__main__":
    unittest.main()
