#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TARGET_DIR="$ROOT_DIR/src/app"

PATTERN="snackBar\\.open\\(\\s*'[^']+'"

violations="$(rg -n "$PATTERN" "$TARGET_DIR" --glob '*.ts' || true)"

if [[ -n "$violations" ]]; then
  echo "Found hardcoded snackbar messages. Use i18n keys instead:"
  echo "$violations"

  if [[ "${CI_STRICT_HARDCODED_MESSAGES:-false}" == "true" ]]; then
    exit 1
  fi

  echo "Non-blocking mode: set CI_STRICT_HARDCODED_MESSAGES=true to fail the check."
  exit 0
fi

echo "No hardcoded snackbar messages found."
