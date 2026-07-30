#!/bin/bash
# Graph validation

set -e

REPO_PATH="/Users/ioriimasu/dev/jervis"
GRAPH_FILE="$REPO_PATH/graph.json"

if [ ! -f "$GRAPH_FILE" ]; then
  echo "Error: graph.json not found."
  exit 1
fi

echo "Validating Engineering Knowledge Graph..."

# 1. Duplicate File IDs
duplicates=$(jq -r '.[].id' "$GRAPH_FILE" | sort | uniq -d)
if [ -n "$duplicates" ]; then
  echo "ERROR: Duplicate FILE IDs found:"
  echo "$duplicates"
  exit 1
fi

# 3. Missing relations (Architecture, Package)
orphans=$(jq -r '.[] | select(.architecture == "" or .architecture == null or .package == "" or .package == null) | .id' "$GRAPH_FILE")
if [ -n "$orphans" ]; then
  echo "ERROR: Missing Architecture or Package relation for:"
  echo "$orphans"
  exit 1
fi

# Check YAML stores are not empty
for yaml in "$REPO_PATH/data/"*.yaml; do
  if [ -s "$yaml" ]; then
    # file is not empty
    count=$(grep -c "id:" "$yaml" || true)
    if [ "$count" -eq 0 ]; then
      echo "ERROR: No IDs found in $yaml"
      exit 1
    fi
  else
    echo "ERROR: YAML store $yaml is empty"
    exit 1
  fi
done

echo "Validation Passed: No duplicates, orphans, or empty databases."
exit 0
