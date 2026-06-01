#!/usr/bin/env bash
set -euo pipefail

ROOTDIR=$(cd "$(dirname "$0")" || exit 1; pwd)
OUTPUT_DIR=${OUTPUT_DIR:-"${ROOTDIR}/output"}
GO_BINARY=${GO_BINARY:-go}
APP_BINARY_NAME=${APP_BINARY_NAME:-go-agent}

echo "打包 go-agent..."

rm -rf "${OUTPUT_DIR}"
mkdir -p "${OUTPUT_DIR}/bin"
mkdir -p "${OUTPUT_DIR}/workspace/skills"

(
  cd "${ROOTDIR}"
  "${GO_BINARY}" build -o "${OUTPUT_DIR}/bin/${APP_BINARY_NAME}" .
)

(
  cd "${ROOTDIR}"
  "${GO_BINARY}" run ./cmd/build-artifacts \
    --app-home "${ROOTDIR}" \
    --command-source "${ROOTDIR}/cmd" \
    --command-bin-dir "${OUTPUT_DIR}/workspace/bin" \
    --command-script-dir "${OUTPUT_DIR}/workspace/cmd" \
    --go-binary "${GO_BINARY}"
)

if [ -f "${ROOTDIR}/config.json" ]; then
  cp "${ROOTDIR}/config.json" "${OUTPUT_DIR}/config.json"
fi

if [ -f "${ROOTDIR}/system_prompt.md" ]; then
  cp "${ROOTDIR}/system_prompt.md" "${OUTPUT_DIR}/system_prompt.md"
fi

if [ -d "${ROOTDIR}/skills" ]; then
  cp -R "${ROOTDIR}/skills/." "${OUTPUT_DIR}/workspace/skills/"
fi

echo "打包完成：${OUTPUT_DIR}"
