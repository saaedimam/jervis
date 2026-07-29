#!/bin/bash
# Create Architecture Registry database in Notion

JERVIS_PAGE="3ab1b27f-dcba-81d0-8b35-ed766e2e8420"

# Create Architecture Registry database
curl -s -X POST "https://api.notion.com/v1/databases" \
  -H "Authorization: Bearer $NOTION_API_KEY" \
  -H "Notion-Version: 2025-09-03" \
  -H "Content-Type: application/json" \
  -d "{
    \"parent\": {\"type\": \"page_id\", \"page_id\": \"$JERVIS_PAGE\"},
    \"title\": [{\"type\": \"text\", \"text\": {\"content\": \"🏛️ Architecture Registry\"}}],
    \"properties\": {
      \"Architecture ID\": {\"title\": {}},
      \"Name\": {\"rich_text\": {}},
      \"Layer\": {\"select\": {\"options\": [
        {\"name\": \"Layer 1: Runtime\", \"color\": \"red\"},
        {\"name\": \"Layer 2: Memory\", \"color\": \"orange\"},
        {\"name\": \"Layer 3: Services\", \"color\": \"yellow\"},
        {\"name\": \"Layer 4: AI Provider\", \"color\": \"green\"},
        {\"name\": \"Layer 5: Interfaces\", \"color\": \"blue\"}
      ]}},
      \"Status\": {\"select\": {\"options\": [
        {\"name\": \"Complete\", \"color\": \"green\"},
        {\"name\": \"In Progress\", \"color\": \"yellow\"},
        {\"name\": \"Planned\", \"color\": \"gray\"}
      ]}},
      \"Owner\": {\"rich_text\": {}},
      \"Responsibilities\": {\"rich_text\": {}},
      \"Packages\": {\"relation\": {\"database_id\": \"PLACEHOLDER\", \"single_property\": {}}},
      \"Specifications\": {\"relation\": {\"database_id\": \"PLACEHOLDER\", \"single_property\": {}}},
      \"ADRs\": {\"relation\": {\"database_id\": \"PLACEHOLDER\", \"single_property\": {}}}
    }
  }" | jq -r '.id'
