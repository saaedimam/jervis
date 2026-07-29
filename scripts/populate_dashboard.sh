#!/bin/bash
# scripts/populate_dashboard.sh
source .env

# We'll update the root page content
ROOT_PAGE_ID="3ab1b27f-dcba-81d0-8b35-ed766e2e8420"
mkdir -p scratch

echo "Populating Dashboard (Root Page)..."

# Prepare the dashboard content JSON
cat > scratch/dashboard_blocks.json <<EOF
{
  "children": [
    {
      "object": "block",
      "type": "heading_1",
      "heading_1": {
        "rich_text": [ { "type": "text", "text": { "content": "🚀 JERVIS Project Operating System" } } ]
      }
    },
    {
      "object": "block",
      "type": "callout",
      "callout": {
        "rich_text": [
          { "type": "text", "text": { "content": "Jervis is a local-first personal OS and service platform designed for developer productivity. This Notion workspace is the canonical synchronized knowledge layer." } }
        ],
        "icon": { "emoji": "🧠" },
        "color": "blue_background"
      }
    },
    {
      "object": "block",
      "type": "heading_2",
      "heading_2": {
        "rich_text": [ { "type": "text", "text": { "content": "📊 Executive Summary" } } ]
      }
    },
    {
      "object": "block",
      "type": "paragraph",
      "paragraph": {
        "rich_text": [
          { "type": "text", "text": { "content": "Current Phase: " } },
          { "type": "text", "text": { "content": "Phase 2.1 (Working Memory & Timeline Ledger)", "link": null }, "annotations": { "bold": true, "color": "green" } },
          { "type": "text", "text": { "content": "\nArchitecture Health: " } },
          { "type": "text", "text": { "content": "100.0% Coverage (Frozen Layers)", "link": null }, "annotations": { "bold": true, "color": "blue" } },
          { "type": "text", "text": { "content": "\nRepository Status: " } },
          { "type": "text", "text": { "content": "Deterministic Runtime Owned", "link": null }, "annotations": { "bold": true } }
        ]
      }
    },
    {
      "object": "block",
      "type": "heading_2",
      "heading_2": {
        "rich_text": [ { "type": "text", "text": { "content": "🔗 Quick Navigation" } } ]
      }
    },
    {
      "object": "block",
      "type": "column_list",
      "column_list": {}
    }
  ]
}
EOF

# Note: Notion API doesn't allow creating column_list with children in one go easily via simple JSON if we don't know the structure.
# I'll just use simple bullet points for navigation for now to ensure it works.

cat > scratch/dashboard_blocks.json <<EOF
{
  "children": [
    {
      "object": "block",
      "type": "heading_1",
      "heading_1": {
        "rich_text": [ { "type": "text", "text": { "content": "🚀 JERVIS Project Operating System" } } ]
      }
    },
    {
      "object": "block",
      "type": "callout",
      "callout": {
        "rich_text": [
          { "type": "text", "text": { "content": "Jervis is a local-first personal OS and service platform designed for developer productivity. This Notion workspace is the canonical synchronized knowledge layer." } }
        ],
        "icon": { "emoji": "🧠" },
        "color": "blue_background"
      }
    },
    {
      "object": "block",
      "type": "heading_2",
      "heading_2": {
        "rich_text": [ { "type": "text", "text": { "content": "📊 Executive Summary" } } ]
      }
    },
    {
      "object": "block",
      "type": "paragraph",
      "paragraph": {
        "rich_text": [
          { "type": "text", "text": { "content": "Current Phase: " } },
          { "type": "text", "text": { "content": "Phase 2.1 (Working Memory & Timeline Ledger)", "link": null }, "annotations": { "bold": true, "color": "green" } },
          { "type": "text", "text": { "content": "\nArchitecture Health: " } },
          { "type": "text", "text": { "content": "100.0% Coverage (Frozen Layers)", "link": null }, "annotations": { "bold": true, "color": "blue" } }
        ]
      }
    },
    {
      "object": "block",
      "type": "heading_2",
      "heading_2": {
        "rich_text": [ { "type": "text", "text": { "content": "🔗 Engineering Registries" } } ]
      }
    },
    {
      "object": "block",
      "type": "bulleted_list_item",
      "bulleted_list_item": { "rich_text": [ { "type": "text", "text": { "content": "🏛️ Architecture Registry", "link": { "url": "https://www.notion.so/${ARCHITECTURE_DB//-/}" } } } ] }
    },
    {
      "object": "block",
      "type": "bulleted_list_item",
      "bulleted_list_item": { "rich_text": [ { "type": "text", "text": { "content": "📦 Package Registry", "link": { "url": "https://www.notion.so/${PACKAGES_DB//-/}" } } } ] }
    },
    {
      "object": "block",
      "type": "bulleted_list_item",
      "bulleted_list_item": { "rich_text": [ { "type": "text", "text": { "content": "📋 Specification Registry", "link": { "url": "https://www.notion.so/${SPECIFICATIONS_DB//-/}" } } } ] }
    },
    {
      "object": "block",
      "type": "bulleted_list_item",
      "bulleted_list_item": { "rich_text": [ { "type": "text", "text": { "content": "📝 Session Logs", "link": { "url": "https://www.notion.so/${SESSIONS_DB//-/}" } } } ] }
    },
    {
      "object": "block",
      "type": "bulleted_list_item",
      "bulleted_list_item": { "rich_text": [ { "type": "text", "text": { "content": "📋 ADRs", "link": { "url": "https://www.notion.so/${ADRS_DB//-/}" } } } ] }
    }
  ]
}
EOF

# Append blocks to the root page
curl -s -X PATCH "https://api.notion.com/v1/blocks/$ROOT_PAGE_ID/children" \
  -H "Authorization: Bearer $NOTION_TOKEN" \
  -H "Content-Type: application/json" \
  -H "Notion-Version: 2022-06-28" \
  -d @scratch/dashboard_blocks.json > /dev/null

echo "✅ Dashboard Populated."
