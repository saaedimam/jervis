#!/bin/bash
# Create Specification Registry entries with canonical IDs
# Jervis - Specification Registry

SPEC_DB="f30e0d51-a787-421a-ad6b-77935f7d2e53"

add_spec() {
  local id="$1"
  local name="$2"
  local version="$3"
  local status="$4"
  local pkg="$5"
  local arch="$6"
  local frozen="$7"
  local file="$8"
  
  local json=$(cat <<EOF
{
  "parent": {"database_id": "$SPEC_DB"},
  "properties": {
    "Spec ID": {"title": [{"text": {"content": "$id"}}]},
    "Name": {"rich_text": [{"text": {"content": "$name"}}]},
    "Version": {"rich_text": [{"text": {"content": "$version"}}]},
    "Status": {"select": {"name": "$status"}},
    "Package": {"rich_text": [{"text": {"content": "$pkg"}}]},
    "Architecture": {"rich_text": [{"text": {"content": "$arch"}}]},
    "Frozen": {"checkbox": $frozen},
    "Markdown File": {"url": "https://github.com/saaedimam/jervis/blob/main/$file"}
  }
}
EOF
)
  
  curl -s -X POST "https://api.notion.com/v1/pages" \
    -H "Authorization: Bearer $NOTION_API_KEY" \
    -H "Notion-Version: 2025-09-03" \
    -H "Content-Type: application/json" \
    -d "$json" | jq -r '.id'
}

# Architecture Specifications
add_spec "SPEC-001" "Event Bus Specification" "1.0" "Frozen" "eventbus/*" "ARCH-002" "true" "EVENT_BUS_SPECIFICATION.md"
add_spec "SPEC-002" "Event Model Specification" "1.0" "Frozen" "eventbus/events" "ARCH-002" "true" "EVENT_MODEL.md"
add_spec "SPEC-003" "Event Contracts" "1.0" "Frozen" "eventbus/contracts" "ARCH-002" "true" "EVENT_CONTRACTS.md"
add_spec "SPEC-004" "Dispatcher Specification" "1.0" "Frozen" "eventbus/dispatcher" "ARCH-002" "true" "DISPATCHER_SPECIFICATION.md"
add_spec "SPEC-005" "Middleware Specification" "1.0" "Frozen" "eventbus/middleware" "ARCH-002" "true" "MIDDLEWARE_SPECIFICATION.md"
add_spec "SPEC-006" "Bus Specification" "1.0" "Frozen" "eventbus" "ARCH-002" "true" "BUS_SPECIFICATION.md"

# Permission Specifications
add_spec "SPEC-010" "Permission Engine Specification" "1.0" "Frozen" "permissions/*" "ARCH-003" "true" "PERMISSION_ENGINE_SPECIFICATION.md"
add_spec "SPEC-011" "Permission Model" "1.0" "Frozen" "permissions/rule,policy" "ARCH-003" "true" "PERMISSION_MODEL.md"
add_spec "SPEC-012" "Permission Contracts" "1.0" "Frozen" "permissions/contracts" "ARCH-003" "true" "PERMISSION_CONTRACTS.md"

# Observer Specifications
add_spec "SPEC-020" "Observer Specification" "1.0" "Frozen" "observer/contracts" "ARCH-004" "true" "internal/runtime/observer/OBSERVER_SPECIFICATION.md"
add_spec "SPEC-021" "Observer Model" "1.0" "Frozen" "observer/notification" "ARCH-004" "true" "internal/runtime/observer/OBSERVER_MODEL.md"
add_spec "SPEC-022" "Observer Contracts" "1.0" "Frozen" "observer/contracts" "ARCH-004" "true" "internal/runtime/observer/OBSERVER_CONTRACTS.md"

echo "Specification registry populated"
