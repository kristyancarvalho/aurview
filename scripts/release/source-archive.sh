#!/usr/bin/env bash

set -euo pipefail

if [[ "$#" -ne 2 ]]; then
  echo "usage: $0 <version> <output-path>" >&2
  exit 1
fi

version="$1"
output="$2"
repo_root="$(git rev-parse --show-toplevel)"
caller_dir="$PWD"

if [[ "$output" != /* ]]; then
  output="${caller_dir}/${output}"
fi

stage_dir="$(mktemp -d)"
trap 'rm -rf "$stage_dir"' EXIT

prefix="aurview-${version}"
mkdir -p "${stage_dir}/${prefix}"

cd "$repo_root"
cp -R README.md LICENSE CHANGELOG.md CONTRIBUTING.md go.mod go.sum cmd internal docs "${stage_dir}/${prefix}/"
mkdir -p "$(dirname "$output")"
tar \
  --sort=name \
  --mtime='UTC 1970-01-01' \
  --owner=0 \
  --group=0 \
  --numeric-owner \
  -C "$stage_dir" \
  -czf "$output" \
  "$prefix"
