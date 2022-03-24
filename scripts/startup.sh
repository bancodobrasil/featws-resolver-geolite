#!/bin/sh
set -e

if [ -f /opt/Geolite2-City.mmdb ]; then
    export DATABASE_GEOLITE2=/opt/Geolite2-City.mmdb
fi

echo "Starting Resolver..."
resolver serve \
    --log-level=$LOG_LEVEL \
    --server-port=$SERVER_PORT \
    --geo-database-file="$DATABASE_GEOLITE2" \
    --cities-database-file="/opt/city-state.csv"
    
    