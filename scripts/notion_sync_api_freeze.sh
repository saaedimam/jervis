#!/bin/bash
# Add API Freeze entries to Notion

API_FREEZE_DB="61b41663-2c66-45a1-87c9-5f04b3fd82b6"

add_api_freeze() {
  local phase="$1"
  local packages="$2"
  
  curl -s -X POST "https://api.notion.com/v1/pages" \
    -H "Authorization: Bearer $NOTION_API_KEY" \
    -H "Notion-Version: 2025-09-03" \
    -H "Content-Type: application/json" \
    -d "{
      \"parent\": {\"database_id\": \"$API_FREEZE_DB\"},
      \"properties\": {
        \"Phase\": {\"title\": [{\"text\": {\"content\": \"$phase\"}}]},
        \"Packages\": {\"rich_text\": [{\"text\": {\"content\": \"$packages\"}}]},
        \"Status\": {\"select\": {\"name\": \"Frozen\"}},
        \"ADR Required\": {\"checkbox\": true}
      }
    }" | jq -r '.id'
}

add_api_freeze "Phase 1.1" "runtime/contracts, runtime/types, runtime/errors, runtime/version, runtime/buildinfo, runtime/config, runtime/lifecycle"
add_api_freeze "Phase 1.2.1" "eventbus/contracts, eventbus/events, eventbus/errors"
add_api_freeze "Phase 1.2.2" "eventbus/subscription, eventbus/registry"
add_api_freeze "Phase 1.2.4" "eventbus/dispatcher"
add_api_freeze "Phase 1.2.6" "eventbus/middleware"
add_api_freeze "Phase 1.2.8" "eventbus (Facade)"
add_api_freeze "Phase 1.3.1" "permissions/contracts, capability, decision, validator, errors"
add_api_freeze "Phase 1.3.2" "permissions/rule, permissions/policy"
add_api_freeze "Phase 1.3.3" "permissions/registry"
add_api_freeze "Phase 1.3.4" "permissions/engine"
add_api_freeze "Phase 1.4.1" "observer/contracts, observer/notification, observer/errors"
