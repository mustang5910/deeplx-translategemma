#!/bin/sh
set -e

# If no config file exists, copy the example config as a default
if [ ! -f /app/deeplx-api.yaml ]; then
    echo "No config file found at /app/deeplx-api.yaml, copying example config..."
    cp /app/deeplx-api.yaml.example /app/deeplx-api.yaml
fi

exec "$@"
