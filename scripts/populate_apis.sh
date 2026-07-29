#!/bin/zsh
# scripts/populate_apis.sh
export PATH="/opt/homebrew/bin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin"
source .env

echo "🚀 Starting API Registry Population..."

# Load maps
typeset -A pkg_map
while IFS=: read -r name id; do pkg_map["$name"]="$id"; done < scratch/maps/pkg.map

# Fetch File Map (Path to ID)
echo "Fetching File map..."
typeset -A file_map
while IFS=: read -r f_path id; do file_map["$f_path"]="$id"; done < <(curl -s -X POST "https://api.notion.com/v1/databases/$FILES_DB/query" \
  -H "Authorization: Bearer $NOTION_TOKEN" \
  -H "Notion-Version: 2022-06-28" \
  -H "Content-Type: application/json" \
  | jq -r ".results[] | {id: .id, path: .properties.Path.rich_text[0].text.content} | \"\(.path):\(.id)\"")

create_api_entry() {
  local api_name=$1
  local api_type=$2
  local file_path=$3
  local line_num=$4
  
  local file_id=${file_map["$file_path"]}
  local pkg_name=$(basename "$(dirname "$file_path")")
  local pkg_id=${pkg_map["$pkg_name"]}

  echo "Syncing API: $api_name ($api_type) in $file_path"

  props="{
    \"API ID\": { \"title\": [ { \"text\": { \"content\": \"$api_name\" } } ] },
    \"API Name\": { \"rich_text\": [ { \"text\": { \"content\": \"$api_name\" } } ] },
    \"Type\": { \"select\": { \"name\": \"$api_type\" } },
    \"Line Number\": { \"number\": $line_num },
    \"Status\": { \"select\": { \"name\": \"Active\" } }
  }"

  if [[ -n "$file_id" ]]; then
    props=$(echo "$props" | jq ". + { \"File\": { \"relation\": [ { \"id\": \"$file_id\" } ] } }")
  fi
  if [[ -n "$pkg_id" ]]; then
    props=$(echo "$props" | jq ". + { \"Package\": { \"relation\": [ { \"id\": \"$pkg_id\" } ] } }")
  fi

  curl -s -X POST "https://api.notion.com/v1/pages" \
    -H "Authorization: Bearer $NOTION_TOKEN" \
    -H "Content-Type: application/json" \
    -H "Notion-Version: 2022-06-28" \
    -d "{\"parent\": { \"database_id\": \"$API_REGISTRY_DB\" }, \"properties\": $props}" > /dev/null
}

# Scan Go files for exports
find internal/runtime internal/memory -name "*.go" ! -name "*_test.go" | while read -r f; do
  # Exported functions
  grep -n "^func [A-Z]" "$f" | while IFS=: read -r line content; do
    name=$(echo "$content" | sed 's/func //g' | cut -d'(' -f1 | cut -d' ' -f1)
    create_api_entry "$name" "Function" "$f" "$line"
  done
  # Exported structs
  grep -n "^type [A-Z].*struct" "$f" | while IFS=: read -r line content; do
    name=$(echo "$content" | sed 's/type //g' | cut -d' ' -f1)
    create_api_entry "$name" "Struct" "$f" "$line"
  done
  # Exported interfaces
  grep -n "^type [A-Z].*interface" "$f" | while IFS=: read -r line content; do
    name=$(echo "$content" | sed 's/type //g' | cut -d' ' -f1)
    create_api_entry "$name" "Interface" "$f" "$line"
  done
done

echo "✅ API Registry Population Complete."
