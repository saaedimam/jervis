#!/bin/bash
# Live Schema Reconciliation Script
# Performs read-only inspection of all databases referenced in config/notion_databases.json
# Retrieves and records live schemas, creates expected-vs-live schema diff

set -e

REPO_PATH="/Users/ioriimasu/dev/jervis"
CONFIG_FILE="$REPO_PATH/config/notion_databases.json"
SCHEMA_OUTPUT_DIR="$REPO_PATH/context"
LIVE_SCHEMA_FILE="$SCHEMA_OUTPUT_DIR/live_schemas.json"
EXPECTED_SCHEMA_FILE="$SCHEMA_OUTPUT_DIR/expected_schemas.json"
DIFF_REPORT="$SCHEMA_OUTPUT_DIR/schema_diff_report.md"

# Normalize credential variables - NOTION_TOKEN is canonical
: "${NOTION_TOKEN:?NOTION_TOKEN is required}"
export NOTION_API_KEY="$NOTION_TOKEN"
export NOTION_API_TOKEN="$NOTION_TOKEN"

log() {
    echo "[$(date -u +%Y-%m-%dT%H:%M:%SZ)] [SCHEMA_RECON] $1"
}

error() {
    echo "[$(date -u +%Y-%m-%dT%H:%M:%SZ)] [ERROR] $1" >&2
    exit 1
}

# Check if ntn CLI is available
if ! command -v ntn &> /dev/null; then
    error "ntn CLI not found. Please install notion CLI."
fi

# Initialize output files
echo "{}" > "$LIVE_SCHEMA_FILE"
echo "{}" > "$EXPECTED_SCHEMA_FILE"
echo "# Schema Reconciliation Report" > "$DIFF_REPORT"
echo "" >> "$DIFF_REPORT"
echo "Generated: $(date -u +%Y-%m-%dT%H:%M:%SZ)" >> "$DIFF_REPORT"
echo "" >> "$DIFF_REPORT"

# Load configuration
log "Loading configuration from $CONFIG_FILE"
if [ ! -f "$CONFIG_FILE" ]; then
    error "Configuration file not found: $CONFIG_FILE"
fi

# Get all database keys from config (excluding parent_page and non-database entries)
DB_KEYS=$(jq -r 'to_entries[] | select(.key | test("^(parent_page|architecture|packages|files|specifications|adrs|apis|timeline|gates|handoffs|commits|memory|milestones|dashboard|dependencies|tasks|releases|bugs|risks|tech_debt|research|prompt_library|lessons_learned|patterns|performance_benchmarks|coding_standards|glossary|assets|external_dependencies|roadmap|engineering_decisions|sessions)$")) | .key' "$CONFIG_FILE")

log "Found $(echo "$DB_KEYS" | wc -l) database keys to reconcile"

# Process each database
first=true
> "$LIVE_SCHEMA_FILE"  # Clear the file
echo '{' > "$LIVE_SCHEMA_FILE"
first=true
for DB_KEY in $DB_KEYS; do
    DB_ID=$(jq -r --arg key "$DB_KEY" '.[$key]' "$CONFIG_FILE")

    if [ -z "$DB_ID" ] || [ "$DB_ID" = "null" ] || [ "$DB_ID" = "" ]; then
        log "Skipping $DB_KEY: empty or null ID"
        continue
    fi

    log "Processing $DB_KEY ($DB_ID)"

    # Verify database accessibility
    if ! db_info=$(ntn api "/v1/databases/$DB_ID" 2>/dev/null); then
        log "WARNING: Cannot access database $DB_KEY ($DB_ID)"
        if [ "$first" = true ]; then
            first=false
            echo -n "\"$DB_KEY\":{\"title\":\"\",\"data_source_id\":\"\",\"error\":\"Cannot access database\"}" >> "$LIVE_SCHEMA_FILE"
        else
            echo -n ",\"$DB_KEY\":{\"title\":\"\",\"data_source_id\":\"\",\"error\":\"Cannot access database\"}" >> "$LIVE_SCHEMA_FILE"
        fi
        continue
    fi

    # Extract database title
    db_title=$(echo "$db_info" | jq -r '.title[0].plain_text // empty')
    if [ -z "$db_title" ]; then
        db_title="Untitled"
    fi

    # Get data source ID for schema inspection
    ds_id=$(echo "$db_info" | jq -r '.data_sources[0].id // empty')
    if [ -z "$ds_id" ] || [ "$ds_id" = "null" ]; then
        log "WARNING: Database $DB_KEY ($DB_ID) has no data source"
        if [ "$first" = true ]; then
            first=false
            echo -n "\"$DB_KEY\":{\"title\":\"$db_title\",\"data_source_id\":\"\",\"error\":\"No data source found\"}" >> "$LIVE_SCHEMA_FILE"
        else
            echo -n ",\"$DB_KEY\":{\"title\":\"$db_title\",\"data_source_id\":\"\",\"error\":\"No data source found\"}" >> "$LIVE_SCHEMA_FILE"
        fi
        continue
    fi

    # Get a sample page to inspect properties
    if ! sample_page=$(ntn datasources query "$ds_id" 2>/dev/null | head -1); then
        log "WARNING: Cannot sample data source for $DB_KEY ($DB_ID)"
        if [ "$first" = true ]; then
            first=false
            echo -n "\"$DB_KEY\":{\"title\":\"$db_title\",\"data_source_id\":\"$ds_id\",\"error\":\"Cannot sample data source\"}" >> "$LIVE_SCHEMA_FILE"
        else
            echo -n ",\"$DB_KEY\":{\"title\":\"$db_title\",\"data_source_id\":\"$ds_id\",\"error\":\"Cannot sample data source\"}" >> "$LIVE_SCHEMA_FILE"
        fi
        continue
    fi

    # Extract page ID from sample data
    page_id=$(echo "$sample_page" | cut -f1)
    if [ -z "$page_id" ] || [ "$page_id" = "" ]; then
        log "WARNING: Could not get sample page ID from data source for $DB_KEY"
        if [ "$first" = true ]; then
            first=false
            echo -n "\"$DB_KEY\":{\"title\":\"$db_title\",\"data_source_id\":\"$ds_id\",\"error\":\"Could not get sample page ID\"}" >> "$LIVE_SCHEMA_FILE"
        else
            echo -n ",\"$DB_KEY\":{\"title\":\"$db_title\",\"data_source_id\":\"$ds_id\",\"error\":\"Could not get sample page ID\"}" >> "$LIVE_SCHEMA_FILE"
        fi
        continue
    fi

    # Get the actual page to see properties
    if ! page_data=$(ntn api "/v1/pages/$page_id" 2>/dev/null); then
        log "WARNING: Cannot get sample page for $DB_KEY ($DB_ID)"
        if [ "$first" = true ]; then
            first=false
            echo -n "\"$DB_KEY\":{\"title\":\"$db_title\",\"data_source_id\":\"$ds_id\",\"error\":\"Cannot get sample page\"}" >> "$LIVE_SCHEMA_FILE"
        else
            echo -n ",\"$DB_KEY\":{\"title\":\"$db_title\",\"data_source_id\":\"$ds_id\",\"error\":\"Cannot get sample page\"}" >> "$LIVE_SCHEMA_FILE"
        fi
        continue
    fi

    # Extract all property names and types from the page using jq
    # Build properties JSON using jq directly
    props_json=$(echo "$page_data" | jq -r '.properties | to_entries | map( { (.key): {type: .value.type} } ) | add' 2>/dev/null || echo '{}')

    # Store live schema
    if [ "$first" = true ]; then
        first=false
        echo -n "\"$DB_KEY\":{\"title\":\"$db_title\",\"data_source_id\":\"$ds_id\",\"properties\":$props_json}" >> "$LIVE_SCHEMA_FILE"
    else
        echo -n ",\"$DB_KEY\":{\"title\":\"$db_title\",\"data_source_id\":\"$ds_id\",\"properties\":$props_json}" >> "$LIVE_SCHEMA_FILE"
    fi

    log "Successfully retrieved schema for $DB_KEY: $db_title"
done

# Close JSON object
echo "" >> "$LIVE_SCHEMA_FILE"
echo "}" >> "$LIVE_SCHEMA_FILE"

# Ensure expected_schemas.json is valid JSON
echo "{}" > "$EXPECTED_SCHEMA_FILE"

# Generate expected schemas from notion_schemas.json (if exists)
if [ -f "$SCHEMA_OUTPUT_DIR/notion_schemas.json" ]; then
    log "Generating expected schemas from notion_schemas.json"
    cp "$SCHEMA_OUTPUT_DIR/notion_schemas.json" "$EXPECTED_SCHEMA_FILE"
else
    log "WARNING: notion_schemas.json not found, creating empty expected schemas"
    echo "{}" > "$EXPECTED_SCHEMA_FILE"
fi

# Generate diff report
log "Generating schema diff report"
{
    echo "## Schema Reconciliation Results"
    echo ""
    echo "### Summary"
    echo ""

    # Count databases processed
    live_count=$(jq 'length' "$LIVE_SCHEMA_FILE" 2>/dev/null || echo "0")
    expected_count=$(jq 'length' "$EXPECTED_SCHEMA_FILE" 2>/dev/null || echo "0")

    echo "- Live schemas retrieved: $live_count"
    echo "- Expected schemas available: $expected_count"
    echo ""

    echo "### Detailed Comparison"
    echo ""

    # Process each database in live schemas
    if [ "$live_count" -gt 0 ] && [ "$live_count" -gt 0 ] 2>/dev/null; then
        jq -r 'keys[]' "$LIVE_SCHEMA_FILE" 2>/dev/null | while read db_key; do
            live_info=$(jq -r ".$db_key // empty" "$LIVE_SCHEMA_FILE" 2>/dev/null)
            if [ "$live_info" = "null" ] || [ -z "$live_info" ]; then
                continue
            fi

            live_title=$(echo "$live_info" | jq -r '.title // "Unknown"')
            live_error=$(echo "$live_info" | jq -r '.error // empty')

            if [ -n "$live_error" ]; then
                echo "#### ❌ $db_key ($live_title)"
                echo ""
                echo "**Status:** ERROR - $live_error"
                echo ""
                continue
            fi

            live_props=$(echo "$live_info" | jq -r '.properties | to_entries | map("\(.key): \(.value.type)") | join(", ") // "None"')

            # Get expected schema if available
            expected_info=$(jq -r ".$db_key // empty" "$EXPECTED_SCHEMA_FILE" 2>/dev/null || echo "")
            if [ -n "$expected_info" ] && [ "$expected_info" != "null" ]; then
                expected_title=$(echo "$expected_info" | jq -r '.title // "Unknown"')
                expected_props=$(echo "$expected_info" | jq -r '.properties | to_entries | map("\(.key): \(.value.type)") | join(", ") // "None"')

                echo "#### 🔄 $db_key"
                echo ""
                echo "**Live Title:** $live_title"
                echo "**Expected Title:** $expected_title"
                echo ""
                echo "**Live Properties:**"
                echo "$live_props"
                echo ""
                echo "**Expected Properties:**"
                echo "$expected_props"
                echo ""

                # Simple diff
                if [ "$live_title" != "$expected_title" ]; then
                    echo "⚠️ **Title Mismatch**"
                    echo ""
                fi

                if [ "$live_props" != "$expected_props" ]; then
                    echo "⚠️ **Properties Mismatch**"
                    echo ""
                fi
            else
                echo "#### 🆕 $db_key ($live_title)"
                echo ""
                echo "**Status:** NEW - Not in expected schemas"
                echo ""
                echo "**Live Properties:**"
                echo "$live_props"
                echo ""
            fi
            echo "-->"
            echo ""
        done
    else
        echo "No live schemas retrieved for comparison." >> "$DIFF_REPORT"
        echo "" >> "$DIFF_REPORT"
    fi
} >> "$DIFF_REPORT"

log "Schema reconciliation complete. Report saved to $DIFF_REPORT"
echo "Schema reconciliation complete. Report saved to $DIFF_REPORT"