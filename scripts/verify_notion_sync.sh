#!/bin/bash
# Verification gate for Notion synchronization
# Performs pre-flight checks to ensure safe execution

set -e

REPO_PATH="/Users/ioriimasu/dev/jervis"
CONFIG_FILE="$REPO_PATH/config/notion_databases.json"
LOG_FILE="$REPO_PATH/.verification_log"

# Normalize credential variables - NOTION_TOKEN is canonical
: "${NOTION_TOKEN:?NOTION_TOKEN is required}"
export NOTION_API_KEY="$NOTION_TOKEN"
export NOTION_API_TOKEN="$NOTION_TOKEN"

log() {
    echo "[$(date -u +%Y-%m-%dT%H:%M:%SZ)] $1" >> "$LOG_FILE"
    echo "[VERIFY] $1"
}

error() {
    echo "[$(date -u +%Y-%m-%dT%H:%M:%SZ)] ERROR: $1" >> "$LOG_FILE"
    echo "[ERROR] $1" >&2
    exit 1
}

# 1. Check authentication
log "Checking Notion authentication..."
if ! ntn auth token >/dev/null 2>&1; then
    error "Notion authentication failed. Check NOTION_TOKEN."
fi
log "Authentication successful"

# 2. Load and validate configuration
log "Loading configuration from $CONFIG_FILE"
if [ ! -f "$CONFIG_FILE" ]; then
    error "Configuration file not found: $CONFIG_FILE"
fi

# Extract required database IDs from config
DB_PARENT_PAGE=$(jq -r '.parent_page // empty' "$CONFIG_FILE")
DB_ARCHITECTURE=$(jq -r '.architecture // empty' "$CONFIG_FILE")
DB_PACKAGES=$(jq -r '.packages // empty' "$CONFIG_FILE")
DB_FILES=$(jq -r '.files // empty' "$CONFIG_FILE")
DB_SPECIFICATIONS=$(jq -r '.specifications // empty' "$CONFIG_FILE")
DB_ADRS=$(jq -r '.adrs // empty' "$CONFIG_FILE")
DB_APIS=$(jq -r '.apis // empty' "$CONFIG_FILE")
DB_TIMELINE=$(jq -r '.timeline // empty' "$CONFIG_FILE")
DB_GATES=$(jq -r '.gates // empty' "$CONFIG_FILE")
DB_HANDOFFS=$(jq -r '.handoffs // empty' "$CONFIG_FILE")
DB_COMMITS=$(jq -r '.commits // empty' "$CONFIG_FILE")
DB_MEMORY=$(jq -r '.memory // empty' "$CONFIG_FILE")
DB_MILESTONES=$(jq -r '.milestones // empty' "$CONFIG_FILE")
DB_DASHBOARD=$(jq -r '.dashboard // empty' "$CONFIG_FILE")
DB_DEPENDENCIES=$(jq -r '.dependencies // empty' "$CONFIG_FILE")
DB_TASKS=$(jq -r '.tasks // empty' "$CONFIG_FILE")
DB_RELEASES=$(jq -r '.releases // empty' "$CONFIG_FILE")
DB_BUGS=$(jq -r '.bugs // empty' "$CONFIG_FILE")
DB_RISKS=$(jq -r '.risks // empty' "$CONFIG_FILE")
DB_TECH_DEBT=$(jq -r '.tech_debt // empty' "$CONFIG_FILE")
DB_RESEARCH=$(jq -r '.research // empty' "$CONFIG_FILE")
DB_PROMPT_LIBRARY=$(jq -r '.prompt_library // empty' "$CONFIG_FILE")
DB_LESSONS_LEARNED=$(jq -r '.lessons_learned // empty' "$CONFIG_FILE")
DB_PATTERNS=$(jq -r '.patterns // empty' "$CONFIG_FILE")
DB_PERFORMANCE_BENCHMARKS=$(jq -r '.performance_benchmarks // empty' "$CONFIG_FILE")
DB_CODING_STANDARDS=$(jq -r '.coding_standards // empty' "$CONFIG_FILE")
DB_GLOSSARY=$(jq -r '.glossary // empty' "$CONFIG_FILE")
DB_ASSETS=$(jq -r '.assets // empty' "$CONFIG_FILE")
DB_EXTERNAL_DEPENDENCIES=$(jq -r '.external_dependencies // empty' "$CONFIG_FILE")
DB_ROADMAP=$(jq -r '.roadmap // empty' "$CONFIG_FILE")
DB_ENGINEERING_DECISIONS=$(jq -r '.engineering_decisions // empty' "$CONFIG_FILE")
DB_SESSIONS=$(jq -r '.sessions // empty' "$CONFIG_FILE")

# Check that all required IDs are present
missing=0
for var in DB_PARENT_PAGE DB_ARCHITECTURE DB_PACKAGES DB_FILES DB_SPECIFICATIONS DB_ADRS DB_APIS DB_TIMELINE DB_GATES DB_HANDOFFS DB_COMMITS DB_MEMORY DB_MILESTONES DB_DASHBOARD DB_DEPENDENCIES DB_TASKS DB_RELEASES DB_BUGS DB_RISKS DB_TECH_DEBT DB_RESEARCH DB_PROMPT_LIBRARY DB_LESSONS_LEARNED DB_PATTERNS DB_PERFORMANCE_BENCHMARKS DB_CODING_STANDARDS DB_GLOSSARY DB_ASSETS DB_EXTERNAL_DEPENDENCIES DB_ROADMAP DB_ENGINEERING_DECISIONS DB_SESSIONS; do
    if [ -z "${!var}" ] || [ "${!var}" = "null" ]; then
        echo "Missing or empty required configuration key: $var" >&2
        ((missing++))
    fi
done
if [ $missing -gt 0 ]; then
    error "Missing or empty required configuration keys"
fi

log "All required configuration keys present"

# 3. Verify database accessibility and schemas
log "Verifying database accessibility and schemas..."

# Check critical databases that we know are used
check_database() {
    local name="$1"
    local db_id="$2"
    if [ -z "$db_id" ]; then
        error "Database ID not provided for $name"
    fi
    log "Checking database '$name' ($db_id)..."

    # Check if database exists and is accessible
    if ! db_info=$(ntn api "/v1/databases/$db_id" 2>/dev/null); then
        error "Cannot access database '$name' (ID: $db_id). Check if it's shared with the integration."
    fi

    # Extract title for verification
    db_title=$(echo "$db_info" | jq -r '.title[0].plain_text // empty' 2>/dev/null)
    if [ -z "$db_title" ]; then
        error "Database '$name' (ID: $db_id) has no title or is inaccessible"
    fi

    # For Context database, verify required properties exist
    if [ "$name" = "Context" ]; then
        # Get the data source ID
        ds_id=$(echo "$db_info" | jq -r '.data_sources[0].id // empty')
        if [ -z "$ds_id" ] || [ "$ds_id" = "null" ]; then
            error "Context database has no data source"
        fi
        export DS_ID="$ds_id"  # Store for later use

        # Query the data source to get schema (get actual page to see properties)
        log "Checking schema for Context database data source..."
        # Get a sample page to inspect properties
        if ! sample_page=$(ntn datasources query "$ds_id" 2>/dev/null | head -1); then
            error "Cannot sample data source for Context database"
        fi

        # Extract the page ID from the sample data
        page_id=$(echo "$sample_page" | cut -f1)
        if [ -z "$page_id" ] || [ "$page_id" = "" ]; then
            error "Could not get sample page ID from Context database data source"
        fi

        # Get the actual page to see properties
        if ! page_data=$(ntn api "/v1/pages/$page_id" 2>/dev/null); then
            error "Cannot get sample page from Context database data source"
        fi

        # Check for required properties
        required_props=("Context Name" "Repository Path")
        missing_props=0

        for prop in "${required_props[@]}"; do
            # Check if property exists in the results
            if ! echo "$page_data" | grep -q "\"$prop\""; then
                error "Required property '$prop' not found in Context database schema"
            fi
        done

        log "Context database schema validation passed"
    fi
}

# Check all critical databases
check_database "Context" "$DB_MEMORY"
check_database "Specifications" "$DB_SPECIFICATIONS"
check_database "Handoffs" "$DB_HANDOFFS"
check_database "Commit Registry" "$DB_COMMITS"
check_database "Session Logs" "$DB_SESSIONS"
check_database "Tasks" "$DB_TASKS"
check_database "Projects/Packages" "$DB_PACKAGES"

log "All critical databases accessible and schemas validated"

# 4. Check for duplicate Repository Path records (sample check)
log "Checking for potential duplicate Repository Path entries..."
# For Context database, check a few records to see if we see duplicate paths
if [ -n "${ds_id:-}" ]; then
    # Get a sample of records
    sample_data=$(ntn datasources query "$ds_id" 2>/dev/null || echo "")
    if [ -n "$sample_data" ]; then
        # Extract Repository Path values (third column)
        paths=$(echo "$sample_data" | cut -f3 | sort | uniq -d)
        if [ -n "$paths" ]; then
            log "WARNING: Found duplicate Repository Path values in sample data:"
            echo "$paths" | while read -r dup; do
                log "  Duplicate path: $dup"
            done
        else
            log "No duplicate Repository Path values found in sample"
        fi
    fi
fi

# 5. Validate sync state consistency (if state files exist)
log "Checking sync state consistency..."
JERVIS_SYNC_STATE="$REPO_PATH/.jervis_sync_state"
if [ -f "$JERVIS_SYNC_STATE" ]; then
    line_count=$(wc -l < "$JERVIS_SYNC_STATE")
    if [ "$line_count" -gt 0 ]; then
        log "Sync state file contains $line_count entries"

        # Verify each file in state exists and has reasonable hash
        while IFS=: read -r file_path file_hash; do
            if [ -z "$file_path" ] || [ -z "$file_hash" ]; then
                continue
            fi

            # Handle both absolute and relative paths in state file
            if [[ "$file_path" = /* ]]; then
                # Absolute path
                full_path="$file_path"
            else
                # Relative path - prepend repo path
                full_path="$REPO_PATH/$file_path"
            fi

            if [ ! -f "$full_path" ]; then
                error "Tracked file does not exist: $file_path"
            fi

            # Verify hash format (32 hex chars for MD5)
            if ! [[ "$file_hash" =~ ^[a-f0-9]{32}$ ]]; then
                error "Invalid hash format for file $file_path: $file_hash"
            fi

        done < "$JERVIS_SYNC_STATE"

        log "Sync state validation passed"
    else
        log "Sync state file is empty (no tracked files)"
    fi
else
    log "No sync state file found - will create on first run"
fi

log "All verification checks passed"
echo "VERIFICATION_PASSED=true" > "$REPO_PATH/.verification_status"
exit 0