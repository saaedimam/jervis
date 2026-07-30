#!/bin/bash
# Populate Package Registry with canonical IDs
# Jervis - Package Registry

PACKAGE_DB="9c8bb7d5-5675-4cc5-b1b7-6a9c1ac3fe2f"

counter=1

add_pkg() {
  local name="$1"
  local module="$2"
  local purpose="$3"
  local status="$4"
  local coverage="$5"
  local frozen="$6"
  local arch="$7"
  local pkg_id=$(printf "PKG-%03d" $counter)
  
  local json=$(cat <<EOF
{
  "parent": {"database_id": "$PACKAGE_DB"},
  "properties": {
    "Package ID": {"title": [{"text": {"content": "$pkg_id"}}]},
    "Name": {"rich_text": [{"text": {"content": "$name"}}]},
    "Module": {"rich_text": [{"text": {"content": "$module"}}]},
    "Purpose": {"rich_text": [{"text": {"content": "$purpose"}}]},
    "Status": {"select": {"name": "$status"}},
    "Coverage": {"rich_text": [{"text": {"content": "$coverage"}}]},
    "Frozen": {"checkbox": $frozen},
    "Source Path": {"url": "https://github.com/ioriimasu/jervis/tree/main/$name"},
    "Architecture": {"rich_text": [{"text": {"content": "$arch"}}]}
  }
}
EOF
)
  
  curl -s -X POST "https://api.notion.com/v1/pages" \
    -H "Authorization: Bearer $NOTION_API_KEY" \
    -H "Notion-Version: 2025-09-03" \
    -H "Content-Type: application/json" \
    -d "$json" | jq -r '.id'
  
  counter=$((counter + 1))
}

# Layer 1: Runtime
add_pkg "internal/runtime/contracts" "runtime" "Core runtime interface contracts" "Complete" "100%" "true" "ARCH-001"
add_pkg "internal/runtime/types" "runtime" "Core runtime type definitions" "Complete" "100%" "true" "ARCH-001"
add_pkg "internal/runtime/errors" "runtime" "Canonical runtime error definitions" "Complete" "100%" "true" "ARCH-001"
add_pkg "internal/runtime/version" "runtime" "Version information" "Complete" "100%" "true" "ARCH-001"
add_pkg "internal/runtime/buildinfo" "runtime" "Build information" "Complete" "100%" "true" "ARCH-001"
add_pkg "internal/runtime/config" "runtime" "Configuration management" "Complete" "100%" "true" "ARCH-001"
add_pkg "internal/runtime/lifecycle" "runtime" "Application lifecycle management" "Complete" "100%" "true" "ARCH-001"

# Layer 1.2: Event Bus
add_pkg "internal/runtime/eventbus/contracts" "eventbus" "Event bus interface contracts" "Complete" "100%" "true" "ARCH-002"
add_pkg "internal/runtime/eventbus/events" "eventbus" "Event envelope and builder" "Complete" "100%" "true" "ARCH-002"
add_pkg "internal/runtime/eventbus/errors" "eventbus" "Event bus error definitions" "Complete" "100%" "true" "ARCH-002"
add_pkg "internal/runtime/eventbus/subscription" "eventbus" "Subscription management" "Complete" "100%" "true" "ARCH-002"
add_pkg "internal/runtime/eventbus/registry" "eventbus" "Subscriber registry with pattern matching" "Complete" "100%" "true" "ARCH-002"
add_pkg "internal/runtime/eventbus/dispatcher" "eventbus" "Synchronous event dispatcher" "Complete" "100%" "true" "ARCH-002"
add_pkg "internal/runtime/eventbus/middleware" "eventbus" "Middleware chain execution" "Complete" "100%" "true" "ARCH-002"

# Layer 1.3: Permissions
add_pkg "internal/runtime/permissions/contracts" "permissions" "Permission engine contracts" "Complete" "100%" "true" "ARCH-003"
add_pkg "internal/runtime/permissions/capability" "permissions" "Capability value objects" "Complete" "100%" "true" "ARCH-003"
add_pkg "internal/runtime/permissions/decision" "permissions" "Decision value objects" "Complete" "100%" "true" "ARCH-003"
add_pkg "internal/runtime/permissions/validator" "permissions" "Capability validation" "Complete" "100%" "true" "ARCH-003"
add_pkg "internal/runtime/permissions/errors" "permissions" "Permission errors" "Complete" "100%" "true" "ARCH-003"
add_pkg "internal/runtime/permissions/rule" "permissions" "Permission rule domain model" "Complete" "100%" "true" "ARCH-003"
add_pkg "internal/runtime/permissions/policy" "permissions" "Policy domain model" "Complete" "100%" "true" "ARCH-003"
add_pkg "internal/runtime/permissions/registry" "permissions" "Policy registry" "Complete" "100%" "true" "ARCH-003"
add_pkg "internal/runtime/permissions/engine" "permissions" "Permission evaluation engine" "Complete" "100%" "true" "ARCH-003"

# Layer 1.4: Observer (In Progress)
add_pkg "internal/runtime/observer/contracts" "observer" "Observer interface contracts" "Complete" "100%" "true" "ARCH-004"
add_pkg "internal/runtime/observer/notification" "observer" "Notification wrapper for events" "Complete" "100%" "true" "ARCH-004"
add_pkg "internal/runtime/observer/errors" "observer" "Observer error definitions" "Complete" "100%" "true" "ARCH-004"
add_pkg "internal/runtime/observer/registry" "observer" "Observer registry (in progress)" "In Progress" "TBD" "false" "ARCH-004"
add_pkg "internal/runtime/observer/dispatcher" "observer" "Observer dispatcher (planned)" "Planned" "TBD" "false" "ARCH-004"
add_pkg "internal/runtime/observer/observer" "observer" "Observer facade (planned)" "Planned" "TBD" "false" "ARCH-004"

echo "Package registry populated with canonical IDs"
