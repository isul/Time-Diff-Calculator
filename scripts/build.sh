#!/usr/bin/env bash
# timediff 프로덕션 빌드 (테스트 후 wails build)
#
# 기본:
#   - macOS 호스트: windows/amd64, linux/amd64, darwin/amd64, darwin/arm64
#   - Linux/Windows 호스트: windows/amd64, linux/amd64
#     (Wails 제약으로 비-macOS에서 macOS 크로스 빌드는 기본 제외)
#
# 사용:
#   ./scripts/build.sh
#   ./scripts/build.sh --native              # 현재 OS만 (빠른 로컬 빌드)
#   ./scripts/build.sh --all-platforms        # macOS 타깃 포함 강제 시도
#   ./scripts/build.sh --no-test
#   ./scripts/build.sh --clean
#   ./scripts/build.sh -platform windows/amd64   # 지정 플랫폼만
#   ./scripts/build.sh -debug                    # 모든 플랫폼 + 디버그
#
# 참고: Wails v2는 Linux/Windows -> macOS 크로스 빌드를 지원하지 않습니다.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$ROOT"

NO_TEST=0
CLEAN=0
NATIVE=0
ALL_PLATFORMS=0
WAILS_ARGS=()

case "$(uname -s)" in
  Darwin)
    DEFAULT_PLATFORMS="windows/amd64,linux/amd64,darwin/amd64,darwin/arm64"
    ;;
  *)
    DEFAULT_PLATFORMS="windows/amd64,linux/amd64"
    ;;
esac

while [[ $# -gt 0 ]]; do
  case "$1" in
    --no-test) NO_TEST=1; shift ;;
    --clean)   CLEAN=1; shift ;;
    --native)  NATIVE=1; shift ;;
    --all-platforms) ALL_PLATFORMS=1; shift ;;
    *)         WAILS_ARGS+=("$1"); shift ;;
  esac
done

has_explicit_platform() {
  local a
  for a in "$@"; do
    if [[ "$a" == "-platform" || "$a" == -platform=* ]]; then
      return 0
    fi
  done
  return 1
}

if [[ "$CLEAN" -eq 1 ]]; then
  rm -rf "$ROOT/build/bin"
  echo "[build] cleaned build/bin"
fi

if [[ "$NO_TEST" -eq 0 ]]; then
  echo "[build] go test ./..."
  go test ./...
fi

if [[ "$ALL_PLATFORMS" -eq 1 ]]; then
  DEFAULT_PLATFORMS="windows/amd64,linux/amd64,darwin/amd64,darwin/arm64"
  echo "[build] --all-platforms: macOS 타깃 포함 시도"
  if [[ "$(uname -s)" != "Darwin" ]]; then
    echo "[build] 경고: 현재 OS에서는 macOS 타깃이 Wails 제약으로 실패/건너뛰기 될 수 있습니다."
  fi
fi

if [[ "$NATIVE" -eq 1 ]] || has_explicit_platform "${WAILS_ARGS[@]}"; then
  echo "[build] wails build${WAILS_ARGS[*]:+ ${WAILS_ARGS[*]}}"
  wails build "${WAILS_ARGS[@]}"
else
  echo "[build] wails build -platform $DEFAULT_PLATFORMS${WAILS_ARGS[*]:+ ${WAILS_ARGS[*]}}"
  wails build -platform "$DEFAULT_PLATFORMS" "${WAILS_ARGS[@]}"
fi

echo "[build] 완료: $ROOT/build/bin/"
