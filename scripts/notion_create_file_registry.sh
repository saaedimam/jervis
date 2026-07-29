#!/bin/bash
# Create File Registry database in Notion if it doesn't exist

set -e

REPO_PATH="/Users/ioriimasu/dev/jervis"
source "$REPO_PATH/scripts/load_config.sh"
load_config

if [ -n "$FILE_DB" ] && [ "$FILE_DB" != "null" ]; then
  echo "File Registry already configured: $FILE_DB"
  exit 0
fi

# Create File Registry database
response=$(curl -s -X POST "https://api.api.notion.com/v1/databases" \
  -H "Authorization: Bearer $NOTION_API_KEY" \
  -H "Notion-Version: 2025-09-03" \
  -H "Content-Type: application/json" \
  -d "{
    \"parent\": {\"type\": \"page_id\", \"page_id\": \"$JERVIS_PAGE\"},
    \"title\": [{\"type\": \"text\", \"text\": {\"content\": \"📄 File Registry\"}}],
    \"properties\": {
      \"File ID\": {\"title\": {}},
      \"Path\": {\"rich_text\": {}},
      \"Package\": {\"relation\": {\"database_id\": \"$PKG_DB\", \"single_property\": {}}},
      \"Architecture\": {\"relation\": {\"database_id\": \"$ARCH_DB\", \"single_property\": {}}},
      \"Specification\": {\"relation\": {\"database_id\": \"$SPEC_DB\", \"single_property\": {}}},
      \"Language\": {\"select\": {\"options\": [
        {\"name\": \"Go\", \"color\": \"blue\"},
        {\"name\": \"Markdown\", \"color\": \"gray\"},
        {\"name\": \"YAML\", \"color\": \"yellow\"},
        {\"name\": \"JSON\", \"color\": \"green\"},
        {\"name\": \"Shell\", \"color\": \"orange\"},
        {\"name\": \"Unknown\", \"color\": \"default\"}
      ]}},
      \"Exports\": {\"rich_text\": {}},
      \"Imports\": {\"rich_text\": {}},
      \"Coverage\": {\"number\": {\"format\": \"percent\"}},
      \"Frozen\": {\"checkbox\": {}},
      \"Status\": {\"select\": {\"options\": [
        {\"name\": \"Active\", \"color\": \"green\"},
        {\"name\": \"Missing\", \"color\": \"red\"}
      ]}},
      \"Last Commit\": {\"relation\": {\"database_id\": \"$COMMIT_DB\", \"single_property\": {}}},
      \"Last Session\": {\"relation\": {\"database_id\": \"$HANDOFF_DB\", \"single_property\": {}}},
      \"Owner\": {\"rich_text\": {}},
      \"API Count\": {\"number\": {\"format\": \"number\"}},
      \"Test File\": {\"checkbox\": {}},
      \"Generated\": {\"checkbox\": {}},
      \"Notes\": {\"rich_text\": {}}
    }
  }")

new_db_id=$(echo "$response" | jq -r '.id')
if [ -z "$new_db_id" ] || [ "$new_db_id" == "null" ]; then
  echo "Failed to create File Registry"
  echo "$response"
  exit 1
fi

echo "Created File Registry: $new_db_id"

# Update config
tmp=$(mktemp)
jq --arg id "$new_db_id" '.files = $id' "$REPO_PATH/config/notion_databases.json" > "$tmp"
mv "$tmp" "$REPO_PATH/config/notion_databases.json"
echo "Updated config with new File Registry ID."
