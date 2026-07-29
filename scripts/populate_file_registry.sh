#!/bin/bash
# Batch File Registry Population Script
# Populates key Go files into Notion File Registry

FILE_DB="d5b8d71a-c568-4288-9443-f3deb8b316bc"

populate_file() {
  local file_id="$1"
  local path="$2"
  local pkg="$3"
  local pkg_id="$4"
  local arch="$5"
  local arch_id="$6"
  local lang="$7"
  local exports="$8"
  local coverage="$9"
  local frozen="${10}"
  local lines="${11}"
  
  curl -s -X POST "https://api.notion.com/v1/pages" \
    -H "Authorization: Bearer $NOTION_API_KEY" \
    -H "Notion-Version: 2025-09-03" \
    -H "Content-Type: application/json" \
    -d "{
      \"parent\": {\"database_id\": \"$FILE_DB\"},
      \"properties\": {
        \"File ID\": {\"rich_text\": [{\"text\": {\"content\": \"$file_id\"}}]},
        \"Path\": {\"rich_text\": [{\"text\": {\"content\": \"$path\"}}]},
        \"Package\": {\"rich_text\": [{\"text\": {\"content\": \"$pkg\"}}]},
        \"Package ID\": {\"rich_text\": [{\"text\": {\"content\": \"$pkg_id\"}}]},
        \"Architecture\": {\"rich_text\": [{\"text\": {\"content\": \"$arch\"}}]},
        \"Arch ID\": {\"rich_text\": [{\"text\": {\"content\": \"$arch_id\"}}]},
        \"Language\": {\"select\": {\"name\": \"$lang\"}},
        \"Exports\": {\"rich_text\": [{\"text\": {\"content\": \"$exports\"}}]},
        \"Coverage\": {\"rich_text\": [{\"text\": {\"content\": \"$coverage\"}}]},
        \"Frozen\": {\"checkbox\": $frozen},
        \"Lines\": {\"number\": $lines},
        \"Status\": {\"select\": {\"name\": \"Active\"}}
      }
    }" | jq -r '.id'
  
  sleep 0.3
}

echo "Populating File Registry..."

# Event Bus Files (ARCH-002)
echo "  Adding Event Bus files..."
populate_file "FILE-0001" "internal/runtime/eventbus/contracts/interfaces.go" "eventbus/contracts" "PKG-007" "Event Bus" "ARCH-002" "Go" "Publisher, Subscriber, Handler, Dispatcher, Validator, Middleware" "100%" "true" "45"
populate_file "FILE-0002" "internal/runtime/eventbus/events/envelope.go" "eventbus/events" "PKG-008" "Event Bus" "ARCH-002" "Go" "Envelope, Builder, Priority" "100%" "true" "120"
populate_file "FILE-0003" "internal/runtime/eventbus/events/errors.go" "eventbus/errors" "PKG-009" "Event Bus" "ARCH-002" "Go" "ErrInvalidEvent, ErrHandlerFailure" "100%" "true" "30"
populate_file "FILE-0004" "internal/runtime/eventbus/subscription/subscription.go" "eventbus/subscription" "PKG-010" "Event Bus" "ARCH-002" "Go" "Subscription, SubscriptionID" "100%" "true" "80"
populate_file "FILE-0005" "internal/runtime/eventbus/registry/registry.go" "eventbus/registry" "PKG-011" "Event Bus" "ARCH-002" "Go" "Registry, Register, Unregister, Lookup" "100%" "true" "150"
populate_file "FILE-0006" "internal/runtime/eventbus/dispatcher/dispatcher.go" "eventbus/dispatcher" "PKG-012" "Event Bus" "ARCH-002" "Go" "Dispatcher, AggregateError, Dispatch" "100%" "true" "200"
populate_file "FILE-0007" "internal/runtime/eventbus/middleware/middleware.go" "eventbus/middleware" "PKG-013" "Event Bus" "ARCH-002" "Go" "Chain, Middleware, Use, Execute" "100%" "true" "90"
populate_file "FILE-0008" "internal/runtime/eventbus/eventbus.go" "eventbus" "PKG-014" "Event Bus" "ARCH-002" "Go" "EventBus, New, Publish, Subscribe" "100%" "true" "110"

# Permission Engine Files (ARCH-003)
echo "  Adding Permission Engine files..."
populate_file "FILE-0009" "internal/runtime/permissions/contracts/contracts.go" "permissions/contracts" "PKG-015" "Permission Engine" "ARCH-003" "Go" "Capability, Decision, Validator" "100%" "true" "60"
populate_file "FILE-0010" "internal/runtime/permissions/capability/capability.go" "permissions/capability" "PKG-016" "Permission Engine" "ARCH-003" "Go" "Capability, New, Subject, Resource, Action" "100%" "true" "85"
populate_file "FILE-0011" "internal/runtime/permissions/decision/decision.go" "permissions/decision" "PKG-017" "Permission Engine" "ARCH-003" "Go" "Decision, NewAllow, NewDeny, IsAllowed" "100%" "true" "75"
populate_file "FILE-0012" "internal/runtime/permissions/validator/validator.go" "permissions/validator" "PKG-018" "Permission Engine" "ARCH-003" "Go" "Validator, Validate" "100%" "true" "50"
populate_file "FILE-0013" "internal/runtime/permissions/rule/rule.go" "permissions/rule" "PKG-019" "Permission Engine" "ARCH-003" "Go" "Rule, New, Evaluate, ID, Subject, Resource, Action, Effect" "100%" "true" "140"
populate_file "FILE-0014" "internal/runtime/permissions/policy/policy.go" "permissions/policy" "PKG-020" "Permission Engine" "ARCH-003" "Go" "Policy, New, Rules, Count, Validate" "100%" "true" "100"
populate_file "FILE-0015" "internal/runtime/permissions/registry/registry.go" "permissions/registry" "PKG-021" "Permission Engine" "ARCH-003" "Go" "Registry, Register, Unregister, Get, Policies" "100%" "true" "95"
populate_file "FILE-0016" "internal/runtime/permissions/engine/engine.go" "permissions/engine" "PKG-022" "Permission Engine" "ARCH-003" "Go" "Engine, New, Authorize, Registry" "100%" "true" "180"

# Observer Files (ARCH-004 - In Progress)
echo "  Adding Observer files..."
populate_file "FILE-0017" "internal/runtime/observer/contracts/interfaces.go" "observer/contracts" "PKG-023" "Observer" "ARCH-004" "Go" "Notification, Observer, Observable, Registry" "100%" "true" "55"
populate_file "FILE-0018" "internal/runtime/observer/notification/notification.go" "observer/notification" "PKG-024" "Observer" "ARCH-004" "Go" "Notification, New, Event, ObservedAt, IsZero" "100%" "true" "70"
populate_file "FILE-0019" "internal/runtime/observer/errors/errors.go" "observer/errors" "PKG-025" "Observer" "ARCH-004" "Go" "ErrInvalidNotification, ErrObserverNotFound" "100%" "true" "35"

# Runtime Core Files (ARCH-001)
echo "  Adding Runtime Core files..."
populate_file "FILE-0020" "internal/runtime/contracts/contracts.go" "runtime/contracts" "PKG-001" "Runtime" "ARCH-001" "Go" "Runtime interface contracts" "100%" "true" "40"
populate_file "FILE-0021" "internal/runtime/types/types.go" "runtime/types" "PKG-002" "Runtime" "ARCH-001" "Go" "Timestamp, Duration" "100%" "true" "50"
populate_file "FILE-0022" "internal/runtime/errors/errors.go" "runtime/errors" "PKG-004" "Runtime" "ARCH-001" "Go" "ErrRuntimeNotStarted, ErrInvalidState" "100%" "true" "30"
populate_file "FILE-0023" "internal/runtime/config/config.go" "runtime/config" "PKG-006" "Runtime" "ARCH-001" "Go" "Config, Load, Validate" "100%" "true" "80"

echo ""
echo "File Registry population complete!"
echo "Total files added: 23"
