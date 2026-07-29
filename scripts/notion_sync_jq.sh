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
  content_json=$(jq -Rs . "$file")
  
  # Create the JSON payload
  cat > payload.json <<EOF
{
  "parent": { "database_id": "$db_id" },
  "properties": {
    "$prop_name": { "title": [ { "text": { "content": "$title" } } ] },
    "Repository Path": { "rich_text": [ { "text": { "content": "$file" } } ] }
  },
  "children": [
    {
      "object": "block",
      "type": "code",
      "code": {
        "language": "markdown",
        "rich_text": [ { "text": { "content": $content_json } } ]
      }
    }
  ]
}
EOF

  # Send via ntn api
  ntn api /v1/pages -d @payload.json
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
find context/sessions -maxdepth 1 -name "*.md" | while read -r f; do
  sync_file "$f" "$SESSION_DB" "Session Logs" "Session ID"
done
