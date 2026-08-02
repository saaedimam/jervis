#!/bin/bash
# Batch File Registry Population Script
# Populates key files into Notion File Registry

set -e

REPO_PATH="/Users/ioriimasu/dev/jervis"
source "$REPO_PATH/scripts/load_config.sh"
load_config

# 1. Get all files from the Notion File Registry
get_notion_files() {
  local next_cursor=""
  while : ; do
    local response=$(curl -s -X POST "https://api.notion.com/v1/databases/$FILE_DB/query" \
      -H "Authorization: Bearer $NOTION_TOKEN" \
      -H "Notion-Version: 2022-06-28" \
      -H "Content-Type: application/json" \
      -d $([ -z "$next_cursor" ] && echo '{"page_size": 100}' || echo "{\"page_size\": 100, \"start_cursor\": \"$next_cursor\"}"))

    echo "$response" | jq -r '.results[].properties.Path.rich_text[0].text.content'

    next_cursor=$(echo "$response" | jq -r '.next_cursor')
    [ "$next_cursor" == "null" ] && break
  done
}

NOTION_FILES=$(get_notion_files)

# 2. Get all files from the local git repository
GIT_FILES=$(git ls-files)

# 3. Find the files that are in git but not in Notion
MISSING_FILES=$(comm -13 <(echo "$NOTION_FILES" | sort) <(echo "$GIT_FILES" | sort))

populate_file() {
  local file_id="$1"
  local path="$2"

  # Get the package name from the file path
  local pkg=$(dirname "$path")

  # Get the language from the file extension
  local lang=""
  case "$path" in
    *.go) lang="Go" ;;
    *.md) lang="Markdown" ;;
    *.yml|*.yaml) lang="YAML" ;;
    *.json) lang="JSON" ;;
    *.sh) lang="Shell" ;;
    *) lang="Unknown" ;;
  esac

  curl -s -X POST "https://api.notion.com/v1/pages" \
    -H "Authorization: Bearer $NOTION_TOKEN" \
    -H "Notion-Version: 2022-06-28" \
    -H "Content-Type: application/json" \
    -d "{
      \"parent\": {\"database_id\": \"$FILE_DB\"},
      \"properties\": {
        \"File ID\": {\"title\": [{\"text\": {\"content\": \"$file_id\"}}]},
        \"Path\": {\"rich_text\": [{\"text\": {\"content\": \"$path\"}}]},
        \"Package\": {\"rich_text\": [{\"text\": {\"content\": \"$pkg\"}}]},
        \"Language\": {\"select\": {\"name\": \"$lang\"}},
        \"Status\": {\"select\": {\"name\": \"Active\"}}
      }
    }" | jq -r '.id'
}

echo "Populating File Registry..."

i=1
for file in $MISSING_FILES; do
  file_id=$(printf "FILE-%04d" $i)
  echo "  Adding $file_id: $file..."
  populate_file "$file_id" "$file"
  i=$((i+1))
done

echo ""
echo "File Registry population complete!"