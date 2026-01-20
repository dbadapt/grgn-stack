#!/bin/bash
# GRGN Stack Template Initialization Script
# Run this script to set up your new project from the template

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
NC='\033[0m' # No Color

info() { echo -e "${CYAN}[INFO]${NC} $1"; }
success() { echo -e "${GREEN}[OK]${NC} $1"; }
warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
error() { echo -e "${RED}[ERROR]${NC} $1"; }

show_help() {
    cat << EOF
GRGN Stack Template Initialization Script

USAGE:
    ./init-template.sh [OPTIONS]

OPTIONS:
    -p, --project-name    Your project name (e.g., "MyApp")
    -u, --github-user     Your GitHub username
    -r, --repo-name       Repository name (defaults to project name lowercase)
    --skip-git            Skip git initialization
    -h, --help            Show this help message

EXAMPLES:
    ./init-template.sh
    ./init-template.sh -p "MyApp" -u "johndoe"
    ./init-template.sh -p "MyApp" -u "johndoe" -r "my-app"

EOF
}

# Parse arguments
PROJECT_NAME=""
GITHUB_USER=""
REPO_NAME=""
SKIP_GIT=false

while [[ $# -gt 0 ]]; do
    case $1 in
        -p|--project-name)
            PROJECT_NAME="$2"
            shift 2
            ;;
        -u|--github-user)
            GITHUB_USER="$2"
            shift 2
            ;;
        -r|--repo-name)
            REPO_NAME="$2"
            shift 2
            ;;
        --skip-git)
            SKIP_GIT=true
            shift
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            error "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
done

get_input() {
    local prompt="$1"
    local default="$2"
    local result

    if [ -n "$default" ]; then
        read -p "$prompt [$default]: " result
        echo "${result:-$default}"
    else
        while [ -z "$result" ]; do
            read -p "$prompt: " result
        done
        echo "$result"
    fi
}

replace_in_file() {
    local file="$1"
    local find="$2"
    local replace="$3"

    if [ -f "$file" ]; then
        if [[ "$OSTYPE" == "darwin"* ]]; then
            sed -i '' "s|$find|$replace|g" "$file" 2>/dev/null || true
        else
            sed -i "s|$find|$replace|g" "$file" 2>/dev/null || true
        fi
    fi
}

replace_in_all_files() {
    local find="$1"
    local replace="$2"

    find . -type f \( -name "*.md" -o -name "*.go" -o -name "*.json" -o -name "*.yml" -o -name "*.yaml" -o -name "*.ts" -o -name "*.tsx" -o -name "*.graphql" \) \
        -not -path "./node_modules/*" \
        -not -path "./web/node_modules/*" \
        -not -path "./.git/*" \
        -exec sh -c '
            if [ -f "$1" ]; then
                if [[ "$OSTYPE" == "darwin"* ]]; then
                    sed -i "" "s|'"$find"'|'"$replace"'|g" "$1" 2>/dev/null || true
                else
                    sed -i "s|'"$find"'|'"$replace"'|g" "$1" 2>/dev/null || true
                fi
            fi
        ' _ {} \;
}

# Main script
echo ""
echo -e "${MAGENTA}========================================${NC}"
echo -e "${MAGENTA}  GRGN Stack Template Initialization   ${NC}"
echo -e "${MAGENTA}========================================${NC}"
echo ""

# Collect project information
if [ -z "$PROJECT_NAME" ]; then
    PROJECT_NAME=$(get_input "Enter your project name (e.g., MyApp)")
fi

if [ -z "$GITHUB_USER" ]; then
    GITHUB_USER=$(get_input "Enter your GitHub username")
fi

if [ -z "$REPO_NAME" ]; then
    default_repo=$(echo "$PROJECT_NAME" | tr '[:upper:]' '[:lower:]' | tr ' ' '-')
    REPO_NAME=$(get_input "Enter repository name" "$default_repo")
fi

ENV_PREFIX=$(echo "${PROJECT_NAME}_" | tr '[:lower:]' '[:upper:]' | tr ' ' '_' | tr '-' '_')

echo ""
info "Configuration Summary:"
echo "  Project Name:    $PROJECT_NAME"
echo "  GitHub User:     $GITHUB_USER"
echo "  Repository:      $REPO_NAME"
echo "  Env Prefix:      $ENV_PREFIX"
echo ""

read -p "Proceed with these settings? (Y/n) " confirm
if [ "$confirm" = "n" ] || [ "$confirm" = "N" ]; then
    warn "Aborted by user"
    exit 1
fi

echo ""

# Step 1: Update Go module paths
info "Updating Go module paths..."
replace_in_all_files "github.com/yourusername/grgn-stack" "github.com/$GITHUB_USER/$REPO_NAME"
success "Go module paths updated"

# Step 2: Update environment variable prefix
info "Updating environment variable prefix..."
replace_in_all_files "GRGN_STACK_" "$ENV_PREFIX"
success "Environment prefix updated to $ENV_PREFIX"

# Step 3: Update project name in documentation
info "Updating project name in documentation..."
replace_in_all_files "GRGN Stack" "$PROJECT_NAME"
replace_in_all_files "grgn-stack" "$REPO_NAME"
success "Project name updated"

# Step 4: Update README badges and links
info "Updating README badges..."
replace_in_file "README.md" "YOUR_USERNAME" "$GITHUB_USER"
replace_in_file "README.md" "YOUR_REPO" "$REPO_NAME"
success "README updated"

# Step 5: Update other documentation
info "Updating documentation links..."
for file in CONTRIBUTING.md TESTING-CI.md CI-CD.md CONFIG.md; do
    if [ -f "$file" ]; then
        replace_in_file "$file" "YOUR_USERNAME" "$GITHUB_USER"
        replace_in_file "$file" "YOUR_REPO" "$REPO_NAME"
    fi
done
success "Documentation updated"

# Step 6: Update LICENSE
info "Updating LICENSE..."
current_year=$(date +%Y)
license_owner=$(get_input "Enter name/organization for LICENSE" "$GITHUB_USER")
replace_in_file "LICENSE" "[YOUR NAME OR ORGANIZATION]" "$license_owner"
replace_in_file "LICENSE" "[YEAR]" "$current_year"
success "LICENSE updated"

# Step 7: Create .env files
info "Creating environment files..."
if [ ! -f ".env" ]; then
    cp ".env.example" ".env"
    if [[ "$OSTYPE" == "darwin"* ]]; then
        sed -i '' "s|GRGN_STACK_|${ENV_PREFIX}|g" ".env"
    else
        sed -i "s|GRGN_STACK_|${ENV_PREFIX}|g" ".env"
    fi
    success "Created .env from .env.example"
fi

if [ ! -f "web/.env" ] && [ -f "web/.env.example" ]; then
    cp "web/.env.example" "web/.env"
    success "Created web/.env from web/.env.example"
fi

# Step 8: Initialize git (optional)
if [ "$SKIP_GIT" = false ]; then
    info "Initializing Git repository..."

    # Remove existing .git if present (template origin)
    if [ -d ".git" ]; then
        rm -rf ".git"
    fi

    git init > /dev/null
    git add .
    git commit -m "Initial commit from GRGN Stack template" > /dev/null

    success "Git repository initialized"
    info "To add remote: git remote add origin https://github.com/$GITHUB_USER/$REPO_NAME.git"
fi

# Step 9: Clean up template files
info "Cleaning up template-specific files..."

read -p "Remove template initialization files? (y/N) " remove_files
if [ "$remove_files" = "y" ] || [ "$remove_files" = "Y" ]; then
    for file in init-template.ps1 init-template.sh TEMPLATE-SETUP.md USING-TEMPLATE.md; do
        if [ -f "$file" ]; then
            rm "$file"
            info "Removed $file"
        fi
    done

    # Amend the commit to remove template files
    if [ "$SKIP_GIT" = false ]; then
        git add .
        git commit --amend -m "Initial commit from GRGN Stack template" --no-edit > /dev/null
    fi
fi

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  Initialization Complete!             ${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo "Next steps:"
echo "  1. Review and edit .env and web/.env files"
echo "  2. Install dependencies: npm install && cd web && npm install"
echo "  3. Start development: docker-compose -f docker-compose.yml -f docker-compose.dev.yml up"
echo ""
echo "Documentation:"
echo "  - README.md         - Project overview"
echo "  - QUICK-REFERENCE.md - Command cheat sheet"
echo "  - ARCHITECTURE.md   - System architecture"
echo ""
success "Happy building with $PROJECT_NAME!"
echo ""
