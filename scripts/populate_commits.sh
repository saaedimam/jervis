#!/bin/zsh
# scripts/populate_commits.sh
export PATH="/opt/homebrew/bin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin"
source .env

echo "🚀 Starting Commit Registry Population..."

# Fetch File Map (Path to ID)
echo "Fetching File map..."
typeset -A file_map
while IFS=: read -r f_path id; do file_map["$f_path"]="$id"; done < <(curl -s -X POST "https://api.notion.com/v1/databases/$FILES_DB/query" \
  -H "Authorization: Bearer $NOTION_TOKEN" \
  -H "Notion-Version: 2022-06-28" \
  -H "Content-Type: application/json" \
  | jq -r ".results[] | {id: .id, path: .properties.Path.rich_text[0].text.content} | \"\(.path):\(.id)\"")

# Fetch Session Map (Session ID to Page ID)
echo "Fetching Session map..."
typeset -A session_map
while IFS=: read -r name id; do session_map["$name"]="$id"; done < <(curl -s -X POST "https://api.notion.com/v1/databases/$SESSIONS_DB/query" \
  -H "Authorization: Bearer $NOTION_TOKEN" \
  -H "Notion-Version: 2022-06-28" \
  -H "Content-Type: application/json" \
  | jq -r ".results[] | {id: .id, name: .properties.\"Session ID\".title[0].text.content} | \"\(.name):\(.id)\"")

sync_commit() {
  local hash=$1
  local author=$2
  local date=$3
  local message=$4
  local branch=$5
  
  echo "Syncing commit: $hash"

  props="{
    \"Commit ID\": { \"title\": [ { \"text\": { \"content\": \"$hash\" } } ] },
    \"Hash\": { \"rich_text\": [ { \"text\": { \"content\": \"$hash\" } } ] },
    \"Author\": { \"rich_text\": [ { \"text\": { \"content\": \"$author\" } } ] },
    \"Date\": { \"date\": { \"start\": \"$date\" } },
    \"Message\": { \"rich_text\": [ { \"text\": { \"content\": \"$message\" } } ] },
    \"Branch\": { \"rich_text\": [ { \"text\": { \"content\": \"$branch\" } } ] }
  }"

  # Link touched files
  local files=$(git show --name-only --format="" "$hash")
  local file_relations="[]"
  while read -r f; do
    local f_id=${file_map["$f"]}
    if [[ -n "$f_id" ]]; then
      file_relations=$(echo "$file_relations" | jq ". + [ { \"id\": \"$f_id\" } ]")
    fi
  done <<< "$files"
  
  if [[ "$file_relations" != "[]" ]]; then
    props=$(echo "$props" | jq ". + { \"Files\": { \"relation\": $file_relations } }")
  fi

  # Link Session (heuristic: find session with same date or latest)
  local day=$(echo "$date" | cut -d'T' -f1)
  local session_name="${day}-session-"
  # This is hard. I'll skip session linking for now or just link manually.

  curl -s -X POST "https://api.notion.com/v1/pages" \
    -H "Authorization: Bearer $NOTION_TOKEN" \
    -H "Content-Type: application/json" \
    -H "Notion-Version: 2022-06-28" \
    -d "{\"parent\": { \"database_id\": \"$COMMITS_DB\" }, \"properties\": $props}" > /dev/null
}

# Get last 10 commits
git log -n 10 --pretty=format:"%H|%an|%cI|%s" | while read -r line; do
  IFS='|' read -r hash author date message <<< "$line"
  branch=$(git branch --show-current)
  sync_commit "$hash" "$author" "$date" "$message" "$branch"
done

echo "✅ Commit Registry Population Complete."
