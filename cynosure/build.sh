#!/usr/bin/env bash
set -euo pipefail

ROOTDIR=$(cd "$(dirname "$0")" || exit 1; pwd)
OUTPUT_DIR=${OUTPUT_DIR:-"${ROOTDIR}/output"}
GO_BINARY=${GO_BINARY:-go}
APP_BINARY_NAME=${APP_BINARY_NAME:-cynosure}

echo "打包 cynosure..."

rm -rf "${OUTPUT_DIR}"
mkdir -p "${OUTPUT_DIR}/bin"
mkdir -p "${OUTPUT_DIR}/skills"
mkdir -p "${OUTPUT_DIR}/cmd"
mkdir -p "${OUTPUT_DIR}/logs"
mkdir -p "${OUTPUT_DIR}/workspace"

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
  cp -R "${ROOTDIR}/cmd/." "${OUTPUT_DIR}/cmd/"
fi

if [ -f "${ROOTDIR}/config.json" ]; then
  cp "${ROOTDIR}/config.json" "${OUTPUT_DIR}/config.json"
fi

if [ -f "${ROOTDIR}/system_prompt.md" ]; then
  cp "${ROOTDIR}/system_prompt.md" "${OUTPUT_DIR}/system_prompt.md"
fi

if [ -d "${ROOTDIR}/skills" ]; then
  cp -R "${ROOTDIR}/skills/." "${OUTPUT_DIR}/skills/"
fi

echo "打包完成：${OUTPUT_DIR}"
