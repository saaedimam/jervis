#!/bin/bash
# Populate API Registry with exported types from Event Bus and Permissions

API_DB="5e2dad61-5186-46f7-be6b-e7e5c3715f04"

api_counter=1

add_api() {
  local id=$(printf "API-%03d" $api_counter)
  local name="$1"
  local pkg="$2"
  local pkg_id="$3"
  local type="$4"
  local sig="$5"
  local spec="$6"
  local spec_id="$7"
  local file="$8"
  local file_id="$9"
  local version="1.0"
  local breaking="false"
  local status="Stable"
  
  if [ "$type" = "Interface" ]; then
    breaking="true"
  fi
  
  json=$(cat <<EOF
{
  "parent": {"database_id": "$API_DB"},
  "properties": {
    "API ID": {"title": [{"text": {"content": "$id"}}]},
    "Name": {"rich_text": [{"text": {"content": "$name"}}]},
    "Package": {"rich_text": [{"text": {"content": "$pkg"}}]},
    "Package ID": {"rich_text": [{"text": {"content": "$pkg_id"}}]},
    "Type": {"select": {"name": "$type"}},
    "Signature": {"rich_text": [{"text": {"content": "$sig"}}]},
    "Version": {"rich_text": [{"text": {"content": "$version"}}]},
    "Breaking": {"checkbox": $breaking},
    "Coverage": {"rich_text": [{"text": {"content": "100%"}}]},
    "Specification": {"rich_text": [{"text": {"content": "$spec"}}]},
    "Spec ID": {"rich_text": [{"text": {"content": "$spec_id"}}]},
    "Status": {"select": {"name": "$status"}},
    "Source File": {"rich_text": [{"text": {"content": "$file"}}]},
    "File ID": {"rich_text": [{"text": {"content": "$file_id"}}]}
  }
}
EOF
)
  
  curl -s -X POST "https://api.notion.com/v1/pages" \
    -H "Authorization: Bearer $NOTION_API_KEY" \
    -H "Notion-Version: 2025-09-03" \
    -H "Content-Type: application/json" \
    -d "$json" | jq -r '.id'
  
  api_counter=$((api_counter + 1))
}

# Event Bus APIs
add_api "Publisher" "eventbus/contracts" "PKG-007" "Interface" "type Publisher interface { Publish(*Envelope) }" "Event Bus Spec" "SPEC-001" "interfaces.go" "FILE-0001"
add_api "Subscriber" "eventbus/contracts" "PKG-007" "Interface" "type Subscriber interface { Subscribe(eventType string, handler Handler) SubscriptionID }" "Event Bus Spec" "SPEC-001" "interfaces.go" "FILE-0001"
add_api "Handler" "eventbus/contracts" "PKG-007" "Interface" "type Handler interface { Handle(*Envelope) error }" "Event Bus Spec" "SPEC-001" "interfaces.go" "FILE-0001"
add_api "Dispatcher" "eventbus/contracts" "PKG-007" "Interface" "type Dispatcher interface { Dispatch(*Envelope) error }" "Dispatcher Spec" "SPEC-004" "interfaces.go" "FILE-0001"
add_api "Validator" "eventbus/contracts" "PKG-007" "Interface" "type Validator interface { Validate(*Envelope) error }" "Event Bus Spec" "SPEC-001" "interfaces.go" "FILE-0001"
add_api "Middleware" "eventbus/contracts" "PKG-007" "Interface" "type Middleware interface { Execute(*Envelope, Handler) error }" "Middleware Spec" "SPEC-005" "interfaces.go" "FILE-0001"
add_api "Envelope" "eventbus/events" "PKG-008" "Type" "type Envelope struct" "Event Model" "SPEC-002" "envelope.go" "FILE-0002"
add_api "NewEnvelope" "eventbus/events" "PKG-008" "Function" "func NewEnvelope(eventType string, payload interface{}) *Envelope" "Event Model" "SPEC-002" "envelope.go" "FILE-0002"
add_api "Priority" "eventbus/events" "PKG-008" "Constant" "const (PriorityLow, PriorityNormal, PriorityHigh, PriorityCritical)" "Event Model" "SPEC-002" "envelope.go" "FILE-0002"
add_api "Registry" "eventbus/registry" "PKG-011" "Type" "type Registry struct" "" "" "registry.go" "FILE-0005"
add_api "NewRegistry" "eventbus/registry" "PKG-011" "Function" "func NewRegistry() *Registry" "" "" "registry.go" "FILE-0005"

# Permission APIs
add_api "Capability" "permissions/capability" "PKG-016" "Type" "type Capability struct" "Permission Model" "SPEC-011" "capability.go" "FILE-0012"
add_api "NewCapability" "permissions/capability" "PKG-016" "Function" "func NewCapability(subject, resource, action string) Capability" "Permission Model" "SPEC-011" "capability.go" "FILE-0012"
add_api "Decision" "permissions/decision" "PKG-017" "Type" "type Decision struct" "" "" "decision.go" "FILE-0013"
add_api "NewAllow" "permissions/decision" "PKG-017" "Function" "func NewAllow(reason string) Decision" "" "" "decision.go" "FILE-0013"
add_api "NewDeny" "permissions/decision" "PKG-017" "Function" "func NewDeny(reason string) Decision" "" "" "decision.go" "FILE-0013"
add_api "Rule" "permissions/rule" "PKG-019" "Type" "type Rule struct" "Permission Model" "SPEC-011" "rule.go" "FILE-0015"
add_api "NewRule" "permissions/rule" "PKG-019" "Function" "func NewRule(id string, capability Capability, effect Effect) Rule" "Permission Model" "SPEC-011" "rule.go" "FILE-0015"
add_api "Policy" "permissions/policy" "PKG-020" "Type" "type Policy struct" "Permission Model" "SPEC-011" "policy.go" "FILE-0016"
add_api "NewPolicy" "permissions/policy" "PKG-020" "Function" "func NewPolicy(id, name string) Policy" "Permission Model" "SPEC-011" "policy.go" "FILE-0016"
add_api "Engine" "permissions/engine" "PKG-022" "Type" "type Engine struct" "Permission Engine Spec" "SPEC-010" "engine.go" "FILE-0018"
add_api "NewEngine" "permissions/engine" "PKG-022" "Function" "func NewEngine(registry *registry.Registry) *Engine" "Permission Engine Spec" "SPEC-010" "engine.go" "FILE-0018"
add_api "Authorize" "permissions/engine" "PKG-022" "Method" "func (e *Engine) Authorize(cap Capability) Decision" "Permission Engine Spec" "SPEC-010" "engine.go" "FILE-0018"

# Observer APIs
add_api "Notification" "observer/notification" "PKG-024" "Type" "type Notification struct" "Observer Model" "SPEC-021" "notification.go" "FILE-0020"
add_api "NewNotification" "observer/notification" "PKG-024" "Function" "func NewNotification(event Event) Notification" "Observer Model" "SPEC-021" "notification.go" "FILE-0020"
add_api "Observer" "observer/contracts" "PKG-023" "Interface" "type Observer interface { Observe(Notification) }" "Observer Spec" "SPEC-020" "interfaces.go" "FILE-0019"
add_api "Observable" "observer/contracts" "PKG-023" "Interface" "type Observable interface { Attach(Observer) error; Detach(Observer) error }" "Observer Spec" "SPEC-020" "interfaces.go" "FILE-0019"

# Event Bus Facade
add_api "EventBus" "eventbus" "PKG-014" "Type" "type EventBus struct" "Bus Spec" "SPEC-006" "eventbus.go" "FILE-0010"
add_api "New" "eventbus" "PKG-014" "Function" "func New() *EventBus" "Bus Spec" "SPEC-006" "eventbus.go" "FILE-0010"
add_api "Publish" "eventbus" "PKG-014" "Method" "func (eb *EventBus) Publish(eventType string, payload interface{}) error" "Bus Spec" "SPEC-006" "eventbus.go" "FILE-0010"
add_api "Subscribe" "eventbus" "PKG-014" "Method" "func (eb *EventBus) Subscribe(eventType string, handler Handler) SubscriptionID" "Bus Spec" "SPEC-006" "eventbus.go" "FILE-0010"

echo "API Registry populated with $((api_counter - 1)) exported APIs"
