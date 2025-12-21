#!/bin/bash
# Test script to verify cache behavior

set -e

BLOG_DIR="/Users/tien-duy.nguyen/projects/oss/git/tduyng-webs/apps/tduyng.github.io"
TEST_FILE="$BLOG_DIR/content/blog/2020-04-07-challenges-ruby-01/index.md"

echo "=== Cache Test Script ==="
echo ""

# Build 1: Clean build
echo "Build 1: Clean build (no cache)"
cd "$BLOG_DIR"
rm -rf public
../gozzi/gozzi build > /dev/null 2>&1

# Build 2: Rebuild without changes (should have high cache hit)
echo "Build 2: Rebuild without changes (should be 100% cache hits)"
rm -rf public
time ../gozzi/gozzi build > /dev/null 2>&1

# Build 3: Touch a file (simulate content change)
echo ""
echo "Build 3: Modify a blog post content"
# Backup original
cp "$TEST_FILE" "$TEST_FILE.backup"
# Add a comment to the file (doesn't change rendered output much)
echo "" >> "$TEST_FILE"
echo "<!-- Test comment at $(date) -->" >> "$TEST_FILE"
time ../gozzi/gozzi build > /dev/null 2>&1

# Restore original
mv "$TEST_FILE.backup" "$TEST_FILE"

echo ""
echo "=== Test Complete ==="
echo "Note: Cache stats only shown in 'serve' mode."
echo "To see detailed cache stats, run: ../gozzi/gozzi serve"
