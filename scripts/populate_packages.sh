#!/bin/bash
# scripts/populate_packages.sh
source .env

create_pkg() {
  local name=$1
  local status=$2
  local purpose=$3
  local path=$4
  local coverage=$5

  echo "Populating Package: $name"
  
  # Check if page exists
  page_id=$(curl -s -X POST "https://api.notion.com/v1/databases/$PACKAGES_DB/query" \
    -H "Authorization: Bearer $NOTION_TOKEN" \
    -H "Notion-Version: 2022-06-28" \
    -H "Content-Type: application/json" \
    -d "{
      \"filter\": {
        \"property\": \"Package\",
        \"title\": { \"equals\": \"$name\" }
      }
    }" | jq -r '.results[0].id // empty')

  if [ -n "$page_id" ]; then
    method="PATCH"
    url="https://api.notion.com/v1/pages/$page_id"
    parent=""
  else
    method="POST"
    url="https://api.notion.com/v1/pages"
    parent="\"parent\": { \"database_id\": \"$PACKAGES_DB\" },"
  fi

  curl -s -X $method "$url" \
    -H "Authorization: Bearer $NOTION_TOKEN" \
    -H "Content-Type: application/json" \
    -H "Notion-Version: 2022-06-28" \
    -d "{
      $parent
      \"properties\": {
        \"Package\": { \"title\": [ { \"text\": { \"content\": \"$name\" } } ] },
        \"Status\": { \"select\": { \"name\": \"$status\" } },
        \"Purpose\": { \"rich_text\": [ { \"text\": { \"content\": \"$purpose\" } } ] },
        \"Repository Path\": { \"rich_text\": [ { \"text\": { \"content\": \"$path\" } } ] },
        \"Coverage\": { \"number\": $coverage }
      }
    }" > /dev/null
}

# Scan internal/ and pkg/
find internal pkg -type d -maxdepth 5 | while read -r d; do
  # Skip if no .go files in the directory directly
  if ls "$d"/*.go >/dev/null 2>&1; then
    pkg_name=$(basename "$d")
    # Try to find a doc.go or a summary in the code
    purpose=$(grep -r "Package $pkg_name" "$d" | head -n 1 | sed 's/\/\/ //g' | sed 's/Package '"$pkg_name"' //g' | cut -c1-100)
    if [ -z "$purpose" ]; then purpose="Package $pkg_name implementation."; fi
    
    # Check coverage (mocking for now, or use go tool cover if possible)
    # For now, if tests exist, set coverage to 1.0, else 0.0
    if ls "$d"/*_test.go >/dev/null 2>&1; then coverage=1.0; else coverage=0.0; fi
    
    create_pkg "$pkg_name" "Frozen" "$purpose" "$d" $coverage
  fi
done
