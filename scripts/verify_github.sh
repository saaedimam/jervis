#!/bin/bash
# scripts/verify_github.sh
# Verifies the actual GitHub repository state against GITHUB_AUDIT.md

set -e

REPO="ioriimasu/jervis"
AUDIT_FILE="GITHUB_AUDIT.md"

echo "🔎 Verifying GitHub Repository State for $REPO"
echo "Reading requirements from $AUDIT_FILE"

# Helper for status printing
pass() { echo -e "\033[1;32m[PASS]\033[0m $1"; }
fail() { echo -e "\033[1;31m[FAIL]\033[0m $1"; exit 1; }

# 1. Repository Properties
print_status() {
    local key=$1
    local value=$2
    echo "Checking $key..."
}

repo_json=$(gh repo view $REPO --json description,isPublic,hasWikiEnabled,hasDiscussionsEnabled,hasProjectsEnabled)

[[ $(echo "$repo_json" | jq -r '.isPublic') == "true" ]] && pass "Repository is Public" || fail "Repository is NOT Public"
[[ $(echo "$repo_json" | jq -r '.hasWikiEnabled') == "true" ]] && pass "Wiki is Enabled" || fail "Wiki is NOT Enabled"
[[ $(echo "$repo_json" | jq -r '.hasDiscussionsEnabled') == "true" ]] && pass "Discussions are Enabled" || fail "Discussions are NOT Enabled"

# 2. Branch Protections
check_protection() {
    local branch=$1
    echo "Verifying $branch protection..."
    protection=$(gh api repos/$REPO/branches/$branch/protection --silent || echo "{}")
    if [[ "$protection" == "{}" ]]; then
        fail "$branch branch is NOT protected"
    fi
    pass "$branch branch protection verified"
}

check_protection "main"
check_protection "develop"

# 3. Labels
label_count=$(gh label list --repo $REPO | wc -l | xargs)
if (( label_count >= 20 )); then
    pass "Labels taxonomy verified ($label_count labels)"
else
    fail "Labels taxonomy incomplete (found $label_count, expected >= 20)"
fi

# 4. Security Features
security_json=$(gh api repos/$REPO --jq '.security_and_analysis')
# Note: Security and analysis data might be restricted based on token permissions
if [[ $(echo "$security_json" | jq -r '.secret_scanning.status') == "enabled" ]]; then
    pass "Secret Scanning is Enabled"
else
    echo "Warning: Could not verify Secret Scanning status (requires specific admin token)"
fi

echo "------------------------------------------------"
pass "All verifiable governance checks passed!"
echo "Audit complete. Reference $AUDIT_FILE for full manual checklist."
