# Push Promptito to GitHub
# Run this script from the promptito directory

# 1. Initialize git (if not already initialized)
if [ ! -d .git ]; then
    echo "Initializing git..."
    git init
fi

# 2. Add all files
echo "Adding files..."
git add .

# 3. Create initial commit
echo "Creating commit..."
git commit -m "feat: initial Promptito release

📁 PROMPT + RAPIDITO = The fastest way to manage AI prompts

Features:
- Zero-dependency Go binary
- Searchable prompt library
- REST API for AI agents
- Standards-compliant metadata
- Self-hosted, no tracking

Security:
- Read-only server design
- Input validation
- Path traversal protection
- DoS protection built-in"

# 4. Add remote (if not already set)
if ! git remote -v | grep -q origin; then
    echo "Adding remote..."
    git remote add origin https://github.com/jacotoledo/Promptito.git
fi

# 5. Push to GitHub
echo "Pushing to GitHub..."
git branch -M main
git push -u origin main

echo ""
echo "Done! Your repo is live at: https://github.com/jacotoledo/Promptito"
