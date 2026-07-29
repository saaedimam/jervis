#!/bin/bash
# scripts/notion_sync.sh

export NOTION_API_TOKEN="$NOTION_TOKEN"

SPEC_DB="3ab1b27f-dcba-8178-8316-c73d93b940e7"
CONTEXT_DB="3ab1b27f-dcba-813c-848e-d922a96996d9"
SESSION_DB="3ab1b27f-dcba-817d-9488-c0b11b753f7c"

sync_file() {
  local file=$1
  local db_id=$2
  local prop_name=$4

  echo "Syncing $file to database $db_id..."
  
  title=$(basename "$file" .md)
  content=$(cat "$file" | sed 's/"/\\"/g' | sed ':a;N;$!ba;s/\n/\\n/g')
  
  # Create the page with content as a code block (to preserve formatting without complex parsing)
  ntn api /v1/pages \
    "parent[database_id]=$db_id" \
    "properties[$prop_name][title][0][text][content]=$title" \
    "properties[Repository Path][rich_text][0][text][content]=$file" \
    "children[0][object]=block" \
    "children[0][type]=code" \
    "children[0][code][language]=markdown" \
    "children[0][code][rich_text][0][text][content]=$content"
}

# 1. Sync Specifications
find . -maxdepth 1 -name "*SPECIFICATION.md" | while read -r f; do
  sync_file "$f" "$SPEC_DB" "Specifications" "Specification"
done

# 2. Sync Context
find context -maxdepth 1 -name "*.md" | while read -r f; do
  sync_file "$f" "$CONTEXT_DB" "Context" "Context Name"
done

# 3. Sync Sessions
find context/sessions -name "*.md" | while read -r f; do
  sync_file "$f" "$SESSION_DB" "Session Logs" "Session ID"
done
