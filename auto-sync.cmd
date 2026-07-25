@echo off
cd /d "%~dp0"
echo ============================================
echo  proj-init Auto-Sync Watcher
echo  Watches for changes and syncs to GitHub
echo  Close this window to stop watching
echo ============================================
echo.

REM Check if we have a token
gh auth status >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo [!] Not logged into GitHub. Run: gh auth login
    pause
    exit /b 1
)

echo [i] Watching for changes every 30 seconds...
echo.

:loop
cls
echo ============================================
echo  proj-init Auto-Sync - Running
echo  %date% %time%
echo ============================================
echo.

REM Check for changes
git add -A 2>&1

REM Check if there's anything to commit
git diff --cached --quiet 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo [*] Changes detected! Committing...
    git commit -m "auto-sync %date% %time%"
    echo.
    echo [*] Pulling latest from GitHub...
    git pull --rebase --autostash
    echo.
    echo [*] Pushing to GitHub...
    git push
    echo.
    echo [✓] Synced successfully at %time%
) else (
    echo [.] No changes detected.
    echo.
    echo [*] Checking for remote changes...
    git pull --rebase --autostash
)

echo.
echo Next check in 30 seconds...
echo Close this window to stop.

timeout /t 30 /nobreak >nul
goto loop
