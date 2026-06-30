from __future__ import annotations

import shlex
from pathlib import Path

from terminal_bench.agents.base_agent import AgentResult, BaseAgent
from terminal_bench.terminal.models import TerminalCommand
from terminal_bench.terminal.tmux_session import TmuxSession


class CynosureTerminalBenchAgent(BaseAgent):
    def __init__(
        self,
        binary_path: str,
        settings_path: str,
        timeout_sec: float = 1800.0,
        **kwargs,
    ):
        super().__init__(**kwargs)
        self.binary_path = Path(binary_path).expanduser().resolve()
        self.settings_path = Path(settings_path).expanduser().resolve()
        self.timeout_sec = float(timeout_sec)
        if not self.binary_path.is_file():
            raise FileNotFoundError(f"cynosure eval binary not found: {self.binary_path}")
        if not self.settings_path.is_file():
            raise FileNotFoundError(f"cynosure settings not found: {self.settings_path}")

    @staticmethod
    def name() -> str:
        return "cynosure"

    def perform_task(
        self,
        instruction: str,
        session: TmuxSession,
        logging_dir: Path | None = None,
    ) -> AgentResult:
        session.copy_to_container(
            self.binary_path,
            container_dir="/installed-agent",
            container_filename="cynosure-eval-agent",
        )
        session.copy_to_container(
            self.settings_path,
            container_dir="/installed-agent",
            container_filename="settings.json",
        )
        session.container.exec_run(["chmod", "+x", "/installed-agent/cynosure-eval-agent"])
        quoted_instruction = shlex.quote(self._render_instruction(instruction))
        session.container.exec_run(
            [
                "sh",
                "-c",
                f"printf %s {quoted_instruction} > /installed-agent/instruction.txt",
            ]
        )

        for command in self._run_agent_commands(instruction):
            session.send_command(command)

        return AgentResult()

    def _run_agent_commands(self, instruction: str) -> list[TerminalCommand]:
        command = (
            "mkdir -p /tmp/cynosure-home/.cynosure .cynosure && "
            "cp /installed-agent/settings.json /tmp/cynosure-home/.cynosure/settings.json && "
            "printf '%s\n' '{\"permissions\":{\"defaultMode\":\"bypassPermissions\"}}' "
            "> .cynosure/settings.json && "
            "HOME=/tmp/cynosure-home /installed-agent/cynosure-eval-agent "
            '--cwd "$PWD" --prompt-file /installed-agent/instruction.txt'
        )
        return [
            TerminalCommand(
                command=command,
                min_timeout_sec=0.0,
                max_timeout_sec=self.timeout_sec,
                block=True,
                append_enter=True,
            )
        ]
