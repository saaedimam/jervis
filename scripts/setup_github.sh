#!/bin/bash
set -e

REPO="ioriimasu/jervis"

echo "🚀 Starting Idempotent GitHub Ecosystem Setup for $REPO"

# Helper for status printing
print_step() { echo -e "\033[1;34m[STEP]\033[0m $1"; }
print_ok() { echo -e "\033[1;32m[OK]\033[0m $1"; }
print_warn() { echo -e "\033[1;33m[WARN]\033[0m $1"; }

# 1. Enable Core Features
print_step "Configuring Repository Features..."
gh repo edit $REPO --enable-discussions --enable-wiki --enable-projects --description "Local-first runtime and context operating system for AI agents."

# 2. Security Configuration
print_step "Configuring Security Features..."
gh api -X PATCH repos/$REPO/import/settings -f private_vulnerability_reporting=enabled --silent || print_warn "Private Vulnerability Reporting already enabled or unsupported."
gh repo edit $REPO --enable-github-actions=true

# 3. Labels Taxonomy
print_step "Synchronizing Labels..."
labels=(
  "kind/bug:#d73a4a"
  "kind/feature:#0075ca"
  "kind/refactor:#e99695"
  "kind/docs:#008672"
  "kind/security:#cfd3d7"
  "kind/performance:#a2eeef"
  "kind/test:#fef2c0"
  "kind/architecture:#7057ff"
  "priority/high:#d93f0b"
  "priority/medium:#fbca04"
  "priority/low:#0e8a16"
  "status/blocked:#000000"
  "status/in-review:#cc317c"
  "status/needs-info:#5319e7"
  "phase/runtime:#bfd4f2"
  "phase/memory:#bfdadc"
  "phase/services:#c5def5"
  "phase/providers:#d4c5f9"
  "phase/interfaces:#f9d0c4"
)

existing_labels=$(gh label list --repo $REPO --limit 100 | awk '{print $1}')

for label in "${labels[@]}"; do
  IFS=":" read -r name color <<< "$label"
  if echo "$existing_labels" | grep -q "^$name$"; then
    gh label edit "$name" --color "$color" --repo $REPO --silent
  else
    gh label create "$name" --color "$color" --repo $REPO --silent
  fi
done

# 4. Version-based Milestones
print_step "Synchronizing Milestones..."
milestones=(
  "v0.1.0 Runtime"
  "v0.2.0 Memory"
  "v0.3.0 Services"
  "v0.4.0 Providers"
  "v0.5.0 Interfaces"
  "v1.0.0 Stable"
)

existing_milestones=$(gh api repos/$REPO/milestones --jq '.[].title')

for m in "${milestones[@]}"; do
  if echo "$existing_milestones" | grep -q "^$m$"; then
    print_ok "Milestone '$m' exists."
  else
    gh api -X POST repos/$REPO/milestones -f title="$m" --silent
  fi
done

# 5. Branch Protection (main)
print_step "Applying Branch Protection: main..."
gh api -X PUT repos/$REPO/branches/main/protection \
  -H "Accept: application/vnd.github+json" \
  --input - <<EOF
{
  "required_status_checks": {
    "strict": true,
    "contexts": ["test", "lint", "security"]
  },
  "enforce_admins": true,
  "required_pull_request_reviews": {
    "dismiss_stale_reviews": true,
    "require_code_owner_reviews": true,
    "required_approving_review_count": 1
  },
  "restrictions": null,
  "required_linear_history": true,
  "allow_force_pushes": false,
  "allow_deletions": false,
  "required_conversation_resolution": true
}
EOF

# 6. Branch Protection (develop)
print_step "Applying Branch Protection: develop..."
gh api -X PUT repos/$REPO/branches/develop/protection \
  -H "Accept: application/vnd.github+json" \
  --input - <<EOF
{
  "required_status_checks": {
    "strict": true,
    "contexts": ["test"]
  },
  "enforce_admins": false,
  "required_pull_request_reviews": {
    "required_approving_review_count": 1
  },
  "restrictions": null,
  "required_linear_history": true,
  "allow_force_pushes": true,
  "allow_deletions": false
}
EOF

print_step "Final Audit Summary"
echo "------------------------------------------------"
echo "Repository: $REPO"
echo "Wiki:       Enabled"
echo "Discussions: Enabled"
echo "Labels:     $(gh label list --repo $REPO | wc -l | xargs) created"
echo "Milestones: $(gh api repos/$REPO/milestones --jq '. | length') created"
echo "------------------------------------------------"
print_ok "GitHub Ecosystem Setup Complete!"
