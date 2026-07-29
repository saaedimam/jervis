#!/bin/zsh
# scripts/populate_files.sh
export PATH="/opt/homebrew/bin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin"
source .env

echo "🚀 Starting File Registry Population..."

# Load maps
typeset -A pkg_map
while IFS=: read -r name id; do pkg_map["$name"]="$id"; done < scratch/maps/pkg.map

typeset -A arch_map
while IFS=: read -r name id; do arch_map["$name"]="$id"; done < scratch/maps/arch.map

typeset -A dir_arch
while IFS=: read -r dir arch; do dir_arch["$dir"]="$arch"; done < scratch/dir_to_arch.map

create_file_entry() {
  local f_path=$1
  local name=$(basename "$f_path")
  local dir=$(dirname "$f_path")
  local ext="${f_path##*.}"
  local lang="Other"
  case "$ext" in
    go) lang="Go" ;;
    md) lang="Markdown" ;;
    yaml|yml) lang="YAML" ;;
    json) lang="JSON" ;;
    sh) lang="Shell" ;;
  esac

  local line_count=$(wc -l < "$f_path" | xargs)
  
  # Determine Package and Architecture
  local pkg_name=$(basename "$dir")
  local pkg_id=${pkg_map["$pkg_name"]}
  
  local arch_name=""
  while IFS=: read -r d a; do
    if [[ "$f_path" == "$d"* ]]; then
      arch_name="$a"
      break
    fi
  done < scratch/dir_to_arch.map
  local arch_id=${arch_map["$arch_name"]}

  echo "Syncing file: $f_path (Pkg: $pkg_name, Arch: $arch_name)"

  # Build Properties
  props="{
    \"File ID\": { \"title\": [ { \"text\": { \"content\": \"$name\" } } ] },
    \"Path\": { \"rich_text\": [ { \"text\": { \"content\": \"$f_path\" } } ] },
    \"Language\": { \"select\": { \"name\": \"$lang\" } },
    \"Line Count\": { \"number\": $line_count },
    \"Status\": { \"select\": { \"name\": \"Active\" } }
  }"

  if [[ -n "$pkg_id" ]]; then
    props=$(echo "$props" | jq ". + { \"Package\": { \"relation\": [ { \"id\": \"$pkg_id\" } ] } }")
  fi
  if [[ -n "$arch_id" ]]; then
    props=$(echo "$props" | jq ". + { \"Architecture\": { \"relation\": [ { \"id\": \"$arch_id\" } ] } }")
  fi

  # Check if exists
  page_id=$(curl -s -X POST "https://api.notion.com/v1/databases/$FILES_DB/query" \
    -H "Authorization: Bearer $NOTION_TOKEN" \
    -H "Notion-Version: 2022-06-28" \
    -H "Content-Type: application/json" \
    -d "{
      \"filter\": {
        \"property\": \"Path\",
        \"rich_text\": { \"equals\": \"$f_path\" }
      }
    }" | jq -r '.results[0].id // empty')

  if [[ -n "$page_id" ]]; then
    method="PATCH"
    url="https://api.notion.com/v1/pages/$page_id"
    payload="{\"properties\": $props}"
  else
    method="POST"
    url="https://api.notion.com/v1/pages"
    payload="{\"parent\": { \"database_id\": \"$FILES_DB\" }, \"properties\": $props}"
  fi

  curl -s -X $method "$url" \
    -H "Authorization: Bearer $NOTION_TOKEN" \
    -H "Content-Type: application/json" \
    -H "Notion-Version: 2022-06-28" \
    -d "$payload" > /dev/null
}

# Scan for files
find internal/runtime internal/memory -maxdepth 3 -type f \( -name "*.go" -o -name "*.md" \) | while read -r f; do
  create_file_entry "$f"
done

echo "✅ File Registry Population Complete."
