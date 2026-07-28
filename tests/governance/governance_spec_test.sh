#!/bin/bash
# tests/governance/governance_spec_test.sh
# Verifies that invalid governance specifications are rejected.

set -e

print_test() { echo -e "\033[1;36m[TEST]\033[0m $1"; }
pass() { echo -e "\033[1;32m[PASS]\033[0m $1"; }
fail() { echo -e "\033[1;31m[FAIL]\033[0m $1"; exit 1; }

# Mock setup script behavior for validation
validate_spec() {
    local file=$1
    if [[ ! -f "$file" ]]; then return 1; fi
    # Simple check for required fields as per setup_github.sh logic
    yq eval 'has("apiVersion") and has("kind") and has("metadata")' "$file" | grep -q "true"
}

print_test "Running Governance Spec Tests..."

# Test 1: Valid Spec
if validate_spec ".github/governance.yaml"; then
    pass "Valid spec correctly identified"
else
    fail "Valid spec rejected"
fi

# Test 2: Invalid Spec (Missing fields)
cat > tmp_invalid.yaml <<EOF
repository:
  name: test
EOF

if validate_spec "tmp_invalid.yaml"; then
    fail "Invalid spec incorrectly accepted"
else
    pass "Invalid spec correctly rejected"
fi

rm tmp_invalid.yaml
echo "All Spec Tests Passed."
