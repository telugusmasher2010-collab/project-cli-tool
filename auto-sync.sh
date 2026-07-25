#!/bin/bash
# proj-init Auto-Sync Watcher (for git-bash / WSL / Linux / Mac)
# Watches for changes and syncs to GitHub every 30 seconds
# Automatically changes to the script's directory

cd "$(dirname "$0")" || exit 1

echo "============================================"
echo " proj-init Auto-Sync Watcher"
echo " Close this terminal to stop watching"
echo "============================================"
echo ""

# Check if gh is authenticated
if ! gh auth status >/dev/null 2>&1; then
    echo "[!] Not logged into GitHub. Run: gh auth login"
    exit 1
fi

echo "[i] Watching for changes every 30 seconds..."
echo ""

while true; do
    clear 2>/dev/null || cls 2>/dev/null
    echo "============================================"
    echo " proj-init Auto-Sync - Running"
    echo " $(date)"
    echo "============================================"
    echo ""
    
    # Stage all changes
    git add -A 2>&1
    
    # Check if there's anything to commit
    if ! git diff --cached --quiet 2>/dev/null; then
        echo "[*] Changes detected! Committing..."
        git commit -m "auto-sync $(date '+%Y-%m-%d %H:%M:%S')"
        echo ""
        echo "[*] Pulling latest from GitHub..."
        git pull --rebase --autostash
        echo ""
        echo "[*] Pushing to GitHub..."
        git push
        echo ""
        echo "[✓] Synced successfully at $(date '+%H:%M:%S')"
    else
        echo "[.] No local changes."
        echo ""
        echo "[*] Checking for remote changes..."
        git pull --rebase --autostash
    fi
    
    echo ""
    echo "Next check in 30 seconds..."
    echo "Press Ctrl+C to stop."
    
    sleep 30
done
