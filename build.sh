#!/usr/bin/env bash
# build.sh — GoコードをWebAssemblyにコンパイルし、必要なファイルを配置する

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
GO_SRC_DIR="${SCRIPT_DIR}/go"
OUTPUT_WASM="${SCRIPT_DIR}/main.wasm"

echo "==> GoをWasmにコンパイル中..."
(
  cd "${GO_SRC_DIR}"
  GOOS=js GOARCH=wasm go build -o "${OUTPUT_WASM}" .
)
echo "    出力: ${OUTPUT_WASM}"

echo "==> wasm_exec.js をコピー中..."
# GoのバージョンによってパスがGoRoot/misc/wasm/ または GoRoot/lib/wasm/ に変わる
GOROOT="$(go env GOROOT)"
WASM_EXEC_CANDIDATES=(
  "${GOROOT}/misc/wasm/wasm_exec.js"
  "${GOROOT}/lib/wasm/wasm_exec.js"
)

WASM_EXEC_SRC=""
for candidate in "${WASM_EXEC_CANDIDATES[@]}"; do
  if [[ -f "${candidate}" ]]; then
    WASM_EXEC_SRC="${candidate}"
    break
  fi
done

if [[ -z "${WASM_EXEC_SRC}" ]]; then
  echo "エラー: wasm_exec.js が見つかりませんでした。"
  echo "探索した場所:"
  for candidate in "${WASM_EXEC_CANDIDATES[@]}"; do
    echo "  ${candidate}"
  done
  exit 1
fi

cp "${WASM_EXEC_SRC}" "${SCRIPT_DIR}/wasm_exec.js"
echo "    ソース: ${WASM_EXEC_SRC}"
echo "    出力:   ${SCRIPT_DIR}/wasm_exec.js"

echo ""
echo "ビルド完了！"
echo "Chrome の拡張機能ページで '${SCRIPT_DIR}' を読み込んでください。"
