#!/bin/bash

# This script generates a new migration file with sequential numbering.

# Define the migration folder path
MIGRATION_DIR="migrations"

# Check if the migrations folder exists, if not, create it
if [ ! -d "$MIGRATION_DIR" ]; then
  echo "Creating migrations folder..."
  mkdir -p "$MIGRATION_DIR"
fi

# Get the last migration number
LAST_MIGRATION=$(ls -1 "$MIGRATION_DIR" | grep -E '^[0-9]{6}_.+\.up\.sql$' | sort -n | tail -n 1)

# Extract the latest migration number
if [ -z "$LAST_MIGRATION" ]; then
  MIGRATION_NUMBER="000001"
else
  MIGRATION_NUMBER=$(printf "%06d" $((${LAST_MIGRATION%%_*} + 1)))
fi

# Ask for the migration name
read -p "Enter migration name (e.g., cart): " MIGRATION_NAME

# Create the .up.sql and .down.sql files
UP_FILE="$MIGRATION_DIR/${MIGRATION_NUMBER}_${MIGRATION_NAME}.up.sql"
DOWN_FILE="$MIGRATION_DIR/${MIGRATION_NUMBER}_${MIGRATION_NAME}.down.sql"

# Generate the migration files with placeholders
echo "-- SQL for the 'up' migration" > "$UP_FILE"
echo "-- Add your 'up' migration SQL here" >> "$UP_FILE"

echo "-- SQL for the 'down' migration" > "$DOWN_FILE"
echo "-- Add your 'down' migration SQL here" >> "$DOWN_FILE"

# Output the paths of the generated files
echo "Migration files created:"
echo "  $UP_FILE"
echo "  $DOWN_FILE"

# Done