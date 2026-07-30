#!/bin/zsh
# scripts/populate_files.sh
export PATH="/opt/homebrew/bin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin"
source .env

echo "Generating Page ID maps..."
mkdir -p scratch/maps

get_map() {
  local db_id=$1
  local name=$2
  local output=$3
  local prop_type=$4

  echo "Fetching $name map..."
  curl -s -X POST "https://api.notion.com/v1/databases/$db_id/query" \
    -H "Authorization: Bearer $NOTION_TOKEN" \
    -H "Notion-Version: 2022-06-28" \
    -H "Content-Type: application/json" \
    | jq -r ".results[] | {id: .id, name: (.properties.\"$prop_type\".title[0].text.content // .properties.\"$prop_type\".rich_text[0].text.content // \"\") } | \"\(.name):\(.id)\"" > "scratch/maps/$output.map"
}

get_map "$ARCHITECTURE_DB" "Architecture" "arch" "Subsystem"
get_map "$PACKAGES_DB" "Packages" "pkg" "Package"
get_map "$SPECIFICATIONS_DB" "Specifications" "spec" "Specification"

echo "Maps generated in scratch/maps/"
