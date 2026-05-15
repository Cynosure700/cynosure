#!/usr/bin/env bash
set -euo pipefail

ROOTDIR=$(cd "$(dirname "$0")" || exit 1; pwd)
OUTPUT_DIR=${OUTPUT_DIR:-"${ROOTDIR}/output"}
GO_BINARY=${GO_BINARY:-go}
APP_BINARY_NAME=${APP_BINARY_NAME:-go-agent}

echo "打包 go-agent..."

rm -rf "${OUTPUT_DIR}"
mkdir -p "${OUTPUT_DIR}/bin"
mkdir -p "${OUTPUT_DIR}/workspace/bin"
mkdir -p "${OUTPUT_DIR}/workspace/cmd"
mkdir -p "${OUTPUT_DIR}/workspace/skills"

build_bin() {
  local output_name="$1"
  local package_path="$2"
  "${GO_BINARY}" build -o "${output_name}" "${package_path}"
}

build_bin "${OUTPUT_DIR}/bin/${APP_BINARY_NAME}" "${ROOTDIR}"

for main_file in "${ROOTDIR}"/cmd/*/main.go; do
  if [ ! -f "${main_file}" ]; then
    continue
  fi

  cmd_dir=$(basename "$(dirname "${main_file}")")
  if [ "${cmd_dir}" = "build-artifacts" ]; then
    continue
  fi

  build_bin "${OUTPUT_DIR}/workspace/bin/${cmd_dir}" "${ROOTDIR}/cmd/${cmd_dir}"
done

while IFS= read -r -d '' script_file; do
  rel_path=${script_file#"${ROOTDIR}/cmd/"}
  target_path="${OUTPUT_DIR}/workspace/cmd/${rel_path}"
  mkdir -p "$(dirname "${target_path}")"
  cp "${script_file}" "${target_path}"
done < <(find "${ROOTDIR}/cmd" -type f \( -name "*.py" -o -name "*.sh" -o -name "*.rb" -o -name "*.pl" \) -print0)

if [ -f "${ROOTDIR}/config.json" ]; then
  cp "${ROOTDIR}/config.json" "${OUTPUT_DIR}/config.json"
fi

if [ -d "${ROOTDIR}/skills" ]; then
  cp -R "${ROOTDIR}/skills/." "${OUTPUT_DIR}/workspace/skills/"
fi

echo "打包完成：${OUTPUT_DIR}"
