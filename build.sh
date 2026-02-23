#!/usr/bin/env bash
set -e

APP="passgen"
VERSION="${1:-1.0.0}"
OUT="dist"

rm -rf "$OUT" && mkdir -p "$OUT"

targets=(
  "darwin   arm64  ${APP}-${VERSION}-macos-arm64"
  "darwin   amd64  ${APP}-${VERSION}-macos-intel"
  "linux    amd64  ${APP}-${VERSION}-linux-amd64"
  "linux    arm64  ${APP}-${VERSION}-linux-arm64"
  "windows  amd64  ${APP}-${VERSION}-windows-amd64"
)

for target in "${targets[@]}"; do
  read -r goos goarch name <<< "$target"

  binary="$name"
  [[ "$goos" == "windows" ]] && binary="${name}.exe"

  printf "  Building %-40s" "$name ..."
  GOOS=$goos GOARCH=$goarch go build -ldflags="-s -w -X main.version=${VERSION}" -o "$OUT/$binary" .

  # Package
  if [[ "$goos" == "windows" ]]; then
    (cd "$OUT" && zip -q "${name}.zip" "${binary}" && rm "${binary}")
  else
    (cd "$OUT" && tar czf "${name}.tar.gz" "${binary}" && rm "${binary}")
  fi
  echo "done"
done

echo ""
echo "Builds written to ./${OUT}/"
ls -lh "$OUT"
