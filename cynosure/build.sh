#!/usr/bin/env bash
set -euo pipefail

ROOTDIR=$(cd "$(dirname "$0")" || exit 1; pwd)
OUTPUT_DIR=${OUTPUT_DIR:-"${ROOTDIR}/output"}
GO_BINARY=${GO_BINARY:-go}
APP_BINARY_NAME=${APP_BINARY_NAME:-cynosure}

echo "打包 cynosure..."

rm -rf "${OUTPUT_DIR}"
mkdir -p "${OUTPUT_DIR}/bin"

(
  cd "${ROOTDIR}"
  "${GO_BINARY}" build -o "${OUTPUT_DIR}/bin/${APP_BINARY_NAME}" .
)

if [ -d "${ROOTDIR}/cmd" ]; then
  for command_dir in "${ROOTDIR}"/cmd/*; do
    [ -d "${command_dir}" ] || continue
    command_name=$(basename "${command_dir}")
    [ -f "${command_dir}/main.go" ] || continue
    (
      cd "${ROOTDIR}"
      "${GO_BINARY}" build -o "${OUTPUT_DIR}/bin/${command_name}" "./cmd/${command_name}"
    )
  done
fi

echo "打包完成：${OUTPUT_DIR}"
echo "system_prompt.md 与内置 skills 已嵌入二进制（go:embed），无需随包分发。"
