#!/bin/bash
# scripts/notion_sync.sh

source .env
export NOTION_API_KEY="$NOTION_TOKEN"

sync_file() {
  local file=$1
  local db_id=$2
  local prop_name=$3

  echo "Syncing $file to database $db_id..."
  
  title=$(basename "$file" .md)
  
  response=$(curl -s -X POST "https://api.notion.com/v1/pages" \
    -H "Authorization: Bearer $NOTION_TOKEN" \
    -H "Content-Type: application/json" \
    -H "Notion-Version: 2022-06-28" \
    -d "{
      \"parent\": { \"database_id\": \"$db_id\" },
      \"properties\": {
        \"$prop_name\": { \"title\": [ { \"text\": { \"content\": \"$title\" } } ] },
        \"Repository Path\": { \"rich_text\": [ { \"text\": { \"content\": \"$file\" } } ] }
      }
    }")
  
  page_id=$(echo "$response" | jq -r '.id // empty')

  if [[ -n "$page_id" && "$page_id" != "null" ]]; then
    # 2. Populate the page content using the python script (which handles chunking)
    python3 scripts/notion_populate_pages.py "$page_id" "$file"
  else
    echo "Failed to create page for $file. Response: $response"
  fi
}

# 1. Sync Specifications
find . -maxdepth 1 -name "*SPECIFICATION.md" | while read -r f; do
  sync_file "$f" "$SPECIFICATIONS_DB" "Specification"
done

# 2. Sync Context & Memory
if [ -d "context" ]; then
  find context -maxdepth 1 -name "*.md" | while read -r f; do
    sync_file "$f" "$MEMORY_DB" "Context Name"
  done
fi

# 3. Sync Sessions
if [ -d "context/sessions" ]; then
  find context/sessions -maxdepth 1 -name "*.md" | while read -r f; do
    sync_file "$f" "$SESSIONS_DB" "Session ID"
  done
fi

# 4. Sync Standards
if [ -f "04_CODING_STANDARD.md" ]; then
  sync_file "04_CODING_STANDARD.md" "$STANDARDS_DB" "Standard ID"
fi

# 5. Sync Everything else in root to Engineering Memory
for f in *.md; do
  case "$f" in
    *SPECIFICATION.md|README.md|CONTRIBUTING.md|04_CODING_STANDARD.md)
      continue
      ;;
    *)
      sync_file "$f" "$MEMORY_DB" "Context Name"
      ;;
  esac
done
