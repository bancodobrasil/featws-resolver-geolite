#!/bin/bash
set -e
set -a

source .env

if [[ -z $DATABASE_GEOLITE2 ]]; then
  echo "environment variable DATABASE_GEOLITE2 not defined!"
  exit 1
fi

if [[ -z $DATABASE_CITYSTATE ]]; then
  echo "environment variable DATABASE_CITYSTATE not defined!"
  exit 1
fi

if [[ -z $FEATWS_GEOLITE_TOKEN ]]; then
  echo "environment variable FEATWS_GEOLITE_TOKEN not defined!"
  exit 1
fi

if [[ ! -f "$DATABASE_GEOLITE2" || ! -f "$DATABASE_CITYSTATE"  ]]; then
  echo "Downloading databases..."
  scripts/download-databases.sh $FEATWS_GEOLITE_TOKEN
  echo "Finished!"
fi

go build -o resolver && ./resolver serve --log-level debug