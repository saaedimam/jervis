#!/bin/bash
# scripts/update_dashboard_metrics.sh
source .env

ROOT_PAGE_ID="3ab1b27f-dcba-81d0-8b35-ed766e2e8420"

# Count items in Notion
count_items() {
  local db_id=$1
  curl -s -X POST "https://api.notion.com/v1/databases/$db_id/query" \
    -H "Authorization: Bearer $NOTION_TOKEN" \
    -H "Notion-Version: 2022-06-28" \
    -H "Content-Type: application/json" | jq '.results | length'
}

arch_count=$(count_items "$ARCHITECTURE_DB")
pkg_count=$(count_items "$PACKAGES_DB")
file_count=$(count_items "$FILES_DB")
spec_count=$(count_items "$SPECIFICATIONS_DB")
adr_count=$(count_items "$ADRS_DB")
session_count=$(count_items "$SESSIONS_DB")

echo "Metrics: Arch=$arch_count, Pkg=$pkg_count, File=$file_count, Spec=$spec_count, ADR=$adr_count, Session=$session_count"

# Update Dashboard (append metrics)
cat > scratch/metrics_blocks.json <<EOF
{
  "children": [
    {
      "object": "block",
      "type": "heading_2",
      "heading_2": { "rich_text": [ { "text": { "content": "📉 Knowledge Graph Metrics" } } ] }
    },
    {
      "object": "block",
      "type": "table",
      "table": {
        "table_width": 2,
        "has_column_header": true,
        "has_row_header": false
      }
    }
  ]
}
EOF

# Notion table population via API is tricky (requires creating rows separately).
# I'll just use a callout or list for now.

cat > scratch/metrics_blocks.json <<EOF
{
  "children": [
    {
      "object": "block",
      "type": "heading_2",
      "heading_2": { "rich_text": [ { "text": { "content": "📉 Knowledge Graph Metrics" } } ] }
    },
    {
      "object": "block",
      "type": "bulleted_list_item",
      "bulleted_list_item": { "rich_text": [ { "text": { "content": "🏛️ Architectures: $arch_count" } } ] }
    },
    {
      "object": "block",
      "type": "bulleted_list_item",
      "bulleted_list_item": { "rich_text": [ { "text": { "content": "📦 Packages: $pkg_count" } } ] }
    },
    {
      "object": "block",
      "type": "bulleted_list_item",
      "bulleted_list_item": { "rich_text": [ { "text": { "content": "📄 Files: $file_count" } } ] }
    },
    {
      "object": "block",
      "type": "bulleted_list_item",
      "bulleted_list_item": { "rich_text": [ { "text": { "content": "📋 Specifications: $spec_count" } } ] }
    },
    {
      "object": "block",
      "type": "bulleted_list_item",
      "bulleted_list_item": { "rich_text": [ { "text": { "content": "📝 Sessions: $session_count" } } ] }
    }
  ]
}
EOF

curl -s -X PATCH "https://api.notion.com/v1/blocks/$ROOT_PAGE_ID/children" \
  -H "Authorization: Bearer $NOTION_TOKEN" \
  -H "Content-Type: application/json" \
  -H "Notion-Version: 2022-06-28" \
  -d @scratch/metrics_blocks.json > /dev/null

echo "✅ Dashboard Metrics Updated."
