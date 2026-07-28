#!/bin/bash
# scripts/setup_github.sh
# Idempotent setup for the jervis GitHub repository and local remote.

set -e

REPO="ioriimasu/jervis"

print_step() { echo -e "\033[1;34m[STEP]\033[0m $1"; }
print_ok() { echo -e "\033[1;32m[OK]\033[0m $1"; }
print_warn() { echo -e "\033[1;33m[WARN]\033[0m $1"; }

# 1. Idempotent Repository Creation
print_step "Checking GitHub repository status..."
if ! gh repo view $REPO &>/dev/null; then
    print_step "Creating repository $REPO..."
    gh repo create $REPO --public --description "Local-first runtime and context operating system for AI agents."
else
    print_ok "Repository $REPO already exists."
fi

# 2. Idempotent Remote Configuration
print_step "Checking local git remote..."
if ! git remote get-url origin &>/dev/null; then
    print_step "Adding origin remote..."
    git remote add origin https://github.com/$REPO.git
else
    print_ok "Remote 'origin' already configured."
fi

# 3. Synchronize Core Settings
print_step "Configuring Repository Features..."
gh repo edit $REPO \
    --enable-discussions \
    --enable-wiki \
    --enable-projects \
    --delete-branch-on-merge \
    --allow-squash-merge \
    --allow-rebase-merge=false

# 4. Labels & Milestones
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

# 5. Branch Protections
print_step "Applying Branch Protection: main..."
gh api -X PUT repos/$REPO/branches/main/protection \
  -H "Accept: application/vnd.github+json" \
  --input - <<EOF
{
  "required_status_checks": { "strict": true, "contexts": ["test", "lint", "security"] },
  "enforce_admins": true,
  "required_pull_request_reviews": { "dismiss_stale_reviews": true, "require_code_owner_reviews": true, "required_approving_review_count": 1 },
  "required_linear_history": true,
  "allow_force_pushes": false,
  "allow_deletions": false,
  "required_conversation_resolution": true
}
EOF

print_ok "GitHub Ecosystem Setup Complete!"
