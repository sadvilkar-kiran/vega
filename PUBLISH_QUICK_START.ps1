# Quick Start Script for Publishing Vega Framework to GitHub
# Run this script to prepare and publish your framework

Write-Host "🚀 Vega Framework - GitHub Publishing Helper" -ForegroundColor Cyan
Write-Host "============================================`n" -ForegroundColor Cyan

# Step 1: Clean up
Write-Host "Step 1: Cleaning up test files..." -ForegroundColor Yellow
Remove-Item -ErrorAction SilentlyContinue test_*.go.bak
Remove-Item -ErrorAction SilentlyContinue test_*.bat
Remove-Item -ErrorAction SilentlyContinue *.exe
Remove-Item -ErrorAction SilentlyContinue FIX_GIT_PATH.md
Remove-Item -ErrorAction SilentlyContinue TEST_*.md
Remove-Item -ErrorAction SilentlyContinue QUICK_START.md
Write-Host "✅ Cleanup complete`n" -ForegroundColor Green

# Step 2: Check if git is initialized
Write-Host "Step 2: Checking Git status..." -ForegroundColor Yellow
if (Test-Path ".git") {
    Write-Host "✅ Git repository already initialized" -ForegroundColor Green
} else {
    Write-Host "⚠️  Git not initialized. Run: git init" -ForegroundColor Yellow
}
Write-Host ""

# Step 3: Show current status
Write-Host "Step 3: Current Git status:" -ForegroundColor Yellow
git status --short
Write-Host ""

# Step 4: Instructions
Write-Host "📋 Next Steps:" -ForegroundColor Cyan
Write-Host "1. Review the changes above" -ForegroundColor White
Write-Host "2. Create GitHub repository at: https://github.com/new" -ForegroundColor White
Write-Host "   - Repository name: vega" -ForegroundColor White
Write-Host "   - Description: High-performance Go web framework" -ForegroundColor White
Write-Host "   - Visibility: Public" -ForegroundColor White
Write-Host "3. Run these commands:" -ForegroundColor White
Write-Host "   git add ." -ForegroundColor Gray
Write-Host "   git commit -m 'Initial commit: Vega Framework v1.0.0'" -ForegroundColor Gray
Write-Host "   git remote add origin https://github.com/sadvilkar-kiran/vega.git" -ForegroundColor Gray
Write-Host "   git branch -M main" -ForegroundColor Gray
Write-Host "   git push -u origin main" -ForegroundColor Gray
Write-Host "4. Create release tag:" -ForegroundColor White
Write-Host "   git tag -a v1.0.0 -m 'Vega Framework v1.0.0'" -ForegroundColor Gray
Write-Host "   git push origin v1.0.0" -ForegroundColor Gray
Write-Host ""
Write-Host "📖 For detailed instructions, see: PUBLISH_TO_GITHUB.md" -ForegroundColor Cyan
Write-Host ""

