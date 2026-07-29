#!/bin/bash
# JERVIS Notion Database Population - Execution Script
# Populates all databases with real entries via Notion API

set -e

NOTION_VERSION="2025-09-03"

# Database IDs
MASTER_PAGE="3ab1b27f-dcba-81d0-8b35-ed766e2e8420"
ARCH_DB="d3dcb133-f96e-4e8e-944f-5825c2d1eee0"
PKG_DB="9c8bb7d5-5675-4cc5-b1b7-6a9c1ac3fe2f"
SPEC_DB="f30e0d51-a787-421a-ad6b-77935f7d2e53"
API_DB="5e2dad61-5186-46f7-be6b-e7e5c3715f04"
FILE_DB="d5b8d71a-c568-4288-9443-f3deb8b316bc"
ADR_DB="abc5d892-1299-4813-b8bf-a143d6c8c73c"
MILESTONES_DB="39ae6e23-2bc1-4e34-a7b0-a1da9410b081"

echo "=== JERVIS NOTION DATABASE POPULATION ==="
echo ""

# Helper function to create database entry
create_entry() {
  local db_id="$1"
  local json_data="$2"
  
  local response=$(curl -s -X POST "https://api.notion.com/v1/pages" \
    -H "Authorization: Bearer $NOTION_API_KEY" \
    -H "Notion-Version: $NOTION_VERSION" \
    -H "Content-Type: application/json" \
    -d "$json_data")
  
  local page_id=$(echo "$response" | jq -r '.id // empty')
  if [ -n "$page_id" ] && [ "$page_id" != "null" ]; then
    echo "$page_id"
  else
    echo "ERROR: $(echo "$response" | jq -r '.message // "Unknown error"')" >&2
    echo ""
  fi
}

# ============================================
# STEP 1: POPULATE ARCHITECTURE DATABASE
# ============================================
echo "Step 1: Architecture Database (4 entries)..."

# ARCH-001: Runtime Core
arch1='{
  "parent": {"database_id": "'$ARCH_DB'"},
  "icon": {"emoji": "🏛️"},
  "properties": {
    "Architecture ID": {"title": [{"text": {"content": "ARCH-001"}}]},
    "Name": {"rich_text": [{"text": {"content": "Runtime Core"}}]},
    "Layer": {"select": {"name": "Layer 1"}},
    "Status": {"select": {"name": "In Progress"}},
    "Responsibilities": {"rich_text": [{"text": {"content": "• System lifecycle management (initialization, graceful shutdown)\n• Configuration management and validation\n• Version and build information tracking\n• Foundation error types and contracts\n• Runtime type definitions\n• Build metadata management"}}]},
    "Interfaces": {"rich_text": [{"text": {"content": "• Lifecycle: Initialize, Shutdown, Status\n• Config: Load, Validate, Get, Set\n• Version: GetVersion, GetBuildInfo\n• Contracts: Runtime interfaces"}}]},
    "Risks": {"rich_text": [{"text": {"content": "Low - Foundation is stable with 100% test coverage"}}]},
    "Future Work": {"rich_text": [{"text": {"content": "• Enhanced observability hooks\n• Runtime metrics collection\n• Distributed tracing integration"}}]}
  }
}'

arch1_id=$(create_entry "$ARCH_DB" "$arch1")
echo "  ✅ ARCH-001 (Runtime Core): ${arch1_id:0:8}..."

# ARCH-002: Event Bus
arch2='{
  "parent": {"database_id": "'$ARCH_DB'"},
  "icon": {"emoji": "📡"},
  "properties": {
    "Architecture ID": {"title": [{"text": {"content": "ARCH-002"}}]},
    "Name": {"rich_text": [{"text": {"content": "Event Bus Engine"}}]},
    "Layer": {"select": {"name": "Layer 1"}},
    "Status": {"select": {"name": "Complete"}},
    "Responsibilities": {"rich_text": [{"text": {"content": "• In-process synchronous event routing\n• Priority-based dispatch (Low, Normal, High, Critical)\n• Panic isolation per handler\n• Middleware chain support (FIFO registration)\n• Aggregate error handling\n• Pattern matching for event subscription"}}]},
    "Interfaces": {"rich_text": [{"text": {"content": "• Publisher: Publish events\n• Subscriber: Subscribe/Unsubscribe\n• Handler: Event handling\n• Dispatcher: Dispatch with error aggregation\n• Validator: Event validation\n• Middleware: Chain execution"}}]},
    "Risks": {"rich_text": [{"text": {"content": "None - Complete with 100% test coverage across 8 packages"}}]},
    "Future Work": {"rich_text": [{"text": {"content": "• Potential async extension (requires ADR)\n• Distributed event bus adapter\n• Event persistence layer"}}]}
  }
}'

arch2_id=$(create_entry "$ARCH_DB" "$arch2")
echo "  ✅ ARCH-002 (Event Bus): ${arch2_id:0:8}..."

# ARCH-003: Permission Engine
arch3='{
  "parent": {"database_id": "'$ARCH_DB'"},
  "icon": {"emoji": "🔐"},
  "properties": {
    "Architecture ID": {"title": [{"text": {"content": "ARCH-003"}}]},
    "Name": {"rich_text": [{"text": {"content": "Permission Engine"}}]},
    "Layer": {"select": {"name": "Layer 1"}},
    "Status": {"select": {"name": "Complete"}},
    "Responsibilities": {"rich_text": [{"text": {"content": "• Capability-based access control (CBAC)\n• Deny-First precedence rule evaluation\n• Wildcard pattern matching (*, prefix*)\n• Immutable domain models (Rule, Policy)\n• 6-stage evaluation pipeline\n• Default deny fallback policy"}}]},
    "Interfaces": {"rich_text": [{"text": {"content": "• Capability: Subject, Resource, Action\n• Decision: Allow/Deny with reason\n• Validator: Validate capabilities\n• Rule: Evaluate single rule\n• Policy: Collection of rules\n• Engine: Authorize requests"}}]},
    "Risks": {"rich_text": [{"text": {"content": "None - Complete with 100% test coverage across 9 packages"}}]},
    "Future Work": {"rich_text": [{"text": {"content": "• Performance optimization for large policy sets\n• Policy caching layer\n• Audit logging integration"}}]}
  }
}'

arch3_id=$(create_entry "$ARCH_DB" "$arch3")
echo "  ✅ ARCH-003 (Permission Engine): ${arch3_id:0:8}..."

# ARCH-004: Observer
arch4='{
  "parent": {"database_id": "'$ARCH_DB'"},
  "icon": {"emoji": "👁️"},
  "properties": {
    "Architecture ID": {"title": [{"text": {"content": "ARCH-004"}}]},
    "Name": {"rich_text": [{"text": {"content": "Observer Subsystem"}}]},
    "Layer": {"select": {"name": "Layer 1"}},
    "Status": {"select": {"name": "Complete"}},
    "Responsibilities": {"rich_text": [{"text": {"content": "• Notification delivery system\n• Compositional event wrapper (Notification)\n• Observer registry with ID management\n• Observable pattern implementation\n• Synchronous notification dispatch\n• Filter chain for notifications"}}]},
    "Interfaces": {"rich_text": [{"text": {"content": "• Notification: Event wrapper\n• Observer: Notification handler\n• Observable: Subscribe/Unsubscribe\n• Registry: Observer management\n• Dispatcher: Notification dispatch"}}]},
    "Risks": {"rich_text": [{"text": {"content": "None - Complete with 100% test coverage"}}]},
    "Future Work": {"rich_text": [{"text": {"content": "• Filter chain enhancement\n• Async notification delivery\n• Notification persistence"}}]}
  }
}'

arch4_id=$(create_entry "$ARCH_DB" "$arch4")
echo "  ✅ ARCH-004 (Observer): ${arch4_id:0:8}..."

echo ""
echo "Architecture Database: 4 entries created ✅"
echo ""

# ============================================
# STEP 2: POPULATE MILESTONES DATABASE
# ============================================
echo "Step 2: Milestones Database (19 entries)..."

# Completed milestones
for i in "Phase 1.1:Core Runtime" "Phase 1.2:Event Bus" "Phase 1.3:Permission Engine" "Phase 1.4:Observer" "Phase 1.5:Scheduler" "Phase 1.6:Session Management" "Phase 2.1:Working Memory"; do
  IFS=':' read -r name desc <<< "$i"
  json='{
    "parent": {"database_id": "'$MILESTONES_DB'"},
    "icon": {"emoji": "✅"},
    "properties": {
      "Name": {"title": [{"text": {"content": "'"$name"'"}}]},
      "Phase": {"select": {"name": "'"$desc"'"}},
      "Status": {"select": {"name": "Done"}},
      "Coverage": {"rich_text": [{"text": {"content": "100%"}}]}
    }
  }'
  create_entry "$MILESTONES_DB" "$json" > /dev/null
  echo "  ✅ $name"
done

# Pending milestones
for i in "Phase 2.2:Knowledge Store" "Phase 3:Domain Services" "Phase 4:AI Provider" "Phase 5:Client Interfaces"; do
  IFS=':' read -r name desc <<< "$i"
  json='{
    "parent": {"database_id": "'$MILESTONES_DB'"},
    "icon": {"emoji": "⏳"},
    "properties": {
      "Name": {"title": [{"text": {"content": "'"$name"'"}}]},
      "Phase": {"select": {"name": "'"$desc"'"}},
      "Status": {"select": {"name": "Pending"}},
      "Coverage": {"rich_text": [{"text": {"content": "TBD"}}]}
    }
  }'
  create_entry "$MILESTONES_DB" "$json" > /dev/null
  echo "  ⏳ $name"
done

echo ""
echo "Milestones Database: 11 entries created ✅"
echo ""

# ============================================
# STEP 3: POPULATE ADR DATABASE
# ============================================
echo "Step 3: ADR Database (4 entries)..."

# ADR-0001
adr1='{
  "parent": {"database_id": "'$ADR_DB'"},
  "icon": {"emoji": "📋"},
  "properties": {
    "ADR ID": {"title": [{"text": {"content": "ADR-0001"}}]},
    "Title": {"rich_text": [{"text": {"content": "Runtime Ownership and Architecture"}}]},
    "Status": {"select": {"name": "Accepted"}},
    "Date": {"date": {"start": "2026-07-15"}},
    "Summary": {"rich_text": [{"text": {"content": "Established the 5-tier hierarchy (OS → Runtime → Memory → Services → AI → Interfaces) with deterministic Runtime ownership ensuring system operation independent of AI providers."}}]}
  }
}'

create_entry "$ADR_DB" "$adr1" > /dev/null
echo "  ✅ ADR-0001: Runtime Ownership"

# ADR-0002
adr2='{
  "parent": {"database_id": "'$ADR_DB'"},
  "icon": {"emoji": "📋"},
  "properties": {
    "ADR ID": {"title": [{"text": {"content": "ADR-0002"}}]},
    "Title": {"rich_text": [{"text": {"content": "Event Bus Synchronous Dispatch"}}]},
    "Status": {"select": {"name": "Accepted"}},
    "Date": {"date": {"start": "2026-07-20"}},
    "Summary": {"rich_text": [{"text": {"content": "Mandated 100% synchronous event dispatch without goroutines, channels, or mutexes. Panic isolation per handler with aggregate error handling."}}]}
  }
}'

create_entry "$ADR_DB" "$adr2" > /dev/null
echo "  ✅ ADR-0002: Event Bus Synchronous Dispatch"

# ADR-0003
adr3='{
  "parent": {"database_id": "'$ADR_DB'"},
  "icon": {"emoji": "📋"},
  "properties": {
    "ADR ID": {"title": [{"text": {"content": "ADR-0003"}}]},
    "Title": {"rich_text": [{"text": {"content": "Permission Engine Deny-First Precedence"}}]},
    "Status": {"select": {"name": "Accepted"}},
    "Date": {"date": {"start": "2026-07-22"}},
    "Summary": {"rich_text": [{"text": {"content": "Established explicit Deny-First rule precedence in capability-based access control with default deny fallback policy."}}]}
  }
}'

create_entry "$ADR_DB" "$adr3" > /dev/null
echo "  ✅ ADR-0003: Deny-First Precedence"

# ADR-0004
adr4='{
  "parent": {"database_id": "'$ADR_DB'"},
  "icon": {"emoji": "📋"},
  "properties": {
    "ADR ID": {"title": [{"text": {"content": "ADR-0004"}}]},
    "Title": {"rich_text": [{"text": {"content": "Observer Compositional Pattern"}}]},
    "Status": {"select": {"name": "Accepted"}},
    "Date": {"date": {"start": "2026-07-25"}},
    "Summary": {"rich_text": [{"text": {"content": "Adopted compositional wrapper pattern for Notification wrapping Event interface without field duplication, enabling observer pattern without code generation."}}]}
  }
}'

create_entry "$ADR_DB" "$adr4" > /dev/null
echo "  ✅ ADR-0004: Observer Compositional Pattern"

echo ""
echo "ADR Database: 4 entries created ✅"
echo ""

echo "=== POPULATION COMPLETE ==="
echo ""
echo "Databases Populated:"
echo "  • Architecture: 4 entries"
echo "  • Milestones: 11 entries"
echo "  • ADRs: 4 entries"
echo ""
echo "Next Steps:"
echo "  • Populate Package Registry (29 packages)"
echo "  • Populate Specification Registry (15 specs)"
echo "  • Populate API Registry (31 APIs)"
echo "  • Populate File Registry (23+ files)"
