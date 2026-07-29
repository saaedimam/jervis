#!/bin/bash
# Add ADRs to Notion

ADRS_DB="6b1d8415-c863-4943-bbe4-4381672f48f0"

curl -s -X POST "https://api.notion.com/v1/pages" \
  -H "Authorization: Bearer $NOTION_API_KEY" \
  -H "Notion-Version: 2025-09-03" \
  -H "Content-Type: application/json" \
  -d "{
    \"parent\": {\"database_id\": \"$ADRS_DB\"},
    \"properties\": {
      \"ADR ID\": {\"title\": [{\"text\": {\"content\": \"ADR-0001\"}}]},
      \"Title\": {\"rich_text\": [{\"text\": {\"content\": \"Initial Architecture Baseline\"}}]},
      \"Status\": {\"select\": {\"name\": \"Superseded\"}},
      \"Date\": {\"date\": {\"start\": \"2026-07-28\"}},
      \"Summary\": {\"rich_text\": [{\"text\": {\"content\": \"Initial architecture baseline established.\"}}]}
    }
  }" | jq -r '.id'

curl -s -X POST "https://api.notion.com/v1/pages" \
  -H "Authorization: Bearer $NOTION_API_KEY" \
  -H "Notion-Version: 2025-09-03" \
  -H "Content-Type: application/json" \
  -d "{
    \"parent\": {\"database_id\": \"$ADRS_DB\"},
    \"properties\": {
      \"ADR ID\": {\"title\": [{\"text\": {\"content\": \"ADR-0002\"}}]},
      \"Title\": {\"rich_text\": [{\"text\": {\"content\": \"Architecture Reconciliation - Runtime Ownership & AI Decoupling\"}}]},
      \"Status\": {\"select\": {\"name\": \"Accepted\"}},
      \"Date\": {\"date\": {\"start\": \"2026-07-28\"}},
      \"Summary\": {\"rich_text\": [{\"text\": {\"content\": \"Replaced AI-centric runtime with canonical 5-tier architecture. Selected Go (Golang 1.22+). Enforced 15 mandatory design rules.\"}}]}
    }
  }" | jq -r '.id'
