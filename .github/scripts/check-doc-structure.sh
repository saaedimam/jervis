#!/usr/bin/env bash
set -e

ALLOWED="^./README.md$|^./CONTRIBUTING.md$|^./LICENSE$|^./CHANGELOG.md$|^./SECURITY.md$|^./SUPPORT.md$|^./CODE_OF_CONDUCT.md$"

VIOLATIONS=$(find . -maxdepth 1 -name "*.md" | grep -vE "$ALLOWED" || true)

if [ -n "$VIOLATIONS" ]; then
    echo "❌ Error: Root markdown files violate clean documentation structure:"
    echo "$VIOLATIONS"
    echo "Please place specification and architecture docs in appropriate subdirectories under docs/"
    exit 1
fi

echo "✅ Root documentation structure verified."
