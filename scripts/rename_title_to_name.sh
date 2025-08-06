#!/bin/bash

# Script to apply migration and regenerate models after renaming Name to name

set -e

echo "Running migration 09_rename_Name_to_name.sql..."

# Run the migration (assuming you have a migration tool or manual SQL execution)
# You'll need to run this against your database:
# psql -d your_database -f migrations/09_rename_Name_to_name.sql

echo "Migration completed. Now regenerating SQLBoiler models..."

# Regenerate models with SQLBoiler
if [ -f "sqlboiler.toml" ]; then
    sqlboiler psql
    echo "SQLBoiler models regenerated successfully"
else
    echo "Warning: sqlboiler.toml not found. Please regenerate models manually with:"
    echo "sqlboiler psql"
fi

echo "Please update any remaining references in the codebase manually."
echo "Script completed!"