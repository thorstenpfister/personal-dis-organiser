#!/usr/bin/env bash

# Script to safely add Terry Pratchett quotes to user configuration
# Only adds if not already present to prevent duplicates

CONFIG_FILE="$HOME/.config/personal-disorganizer/config.json"
QUOTE_FILE="quotes/pratchett.json"

# Create config file if it doesn't exist
if [ ! -f "$CONFIG_FILE" ]; then
    echo "Creating default config file..."
    mkdir -p "$(dirname "$CONFIG_FILE")"
    cat > "$CONFIG_FILE" << 'EOF'
{
  "calendar_urls": [],
  "data_file": "data.json",
  "quote_files": [],
  "refresh_interval": 300,
  "date_format": "2006-01-02",
  "time_format": "15:04",
  "theme": "dracula"
}
EOF
fi

# Check if pratchett.json is already in the quote_files array
if grep -q '"quotes/pratchett.json"' "$CONFIG_FILE"; then
    echo "Terry Pratchett quotes already configured in config.json"
    exit 0
fi

# Check if quote_files array is empty and add the quote file (handle both formatted and compact JSON)
if grep -q '"quote_files":\[\]' "$CONFIG_FILE" || grep -q '"quote_files": \[\]' "$CONFIG_FILE"; then
    # Array is empty, add the first quote file
    sed -i.bak 's/"quote_files":\[\]/"quote_files":["quotes\/pratchett.json"]/' "$CONFIG_FILE"
    sed -i.bak 's/"quote_files": \[\]/"quote_files": ["quotes\/pratchett.json"]/' "$CONFIG_FILE"
    echo "Added Terry Pratchett quotes to empty quote_files array"
elif grep -q '"quote_files":\[' "$CONFIG_FILE" || grep -q '"quote_files": \[' "$CONFIG_FILE"; then
    # Array has existing items, add to the end
    sed -i.bak 's/"quote_files":\[\([^]]*\)\]/"quote_files":[\1,"quotes\/pratchett.json"]/' "$CONFIG_FILE"
    sed -i.bak 's/"quote_files": \[\([^]]*\)\]/"quote_files": [\1, "quotes\/pratchett.json"]/' "$CONFIG_FILE"
    echo "Added Terry Pratchett quotes to existing quote_files array"
else
    echo "Warning: Could not find quote_files array in config.json"
    echo "Please manually add \"quotes/pratchett.json\" to the quote_files array in $CONFIG_FILE"
    exit 1
fi

# Clean up backup file
rm -f "$CONFIG_FILE.bak"

echo "Configuration updated successfully!"