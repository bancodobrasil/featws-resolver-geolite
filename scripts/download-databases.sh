#!/bin/sh
set -e

if test -z $1; then
  echo "MAXMIND_LICENSE_KEY is empty"
  echo "Skipping GeoLite2-City database download"
else
  echo "MAXMIND_LICENSE_KEY=$MAXMIND_LICENSE_KEY"
  curl "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=$MAXMIND_LICENSE_KEY&suffix=tar.gz" \
    --output GeoLite2-City.tar.gz

  mkdir temp
  tar -xvf GeoLite2-City.tar.gz -C temp
  mv temp/*/GeoLite2-City.mmdb Geolite2-City.mmdb
  rm -rf temp
  rm GeoLite2-City.tar.gz
  #city state csv for Brazil
  curl "https://raw.githubusercontent.com/chandez/Estados-Cidades-IBGE/master/sql/Municipios.sql" --output Municipios.sql
  awk -F ',' '{print "BR," $4 "," $5}' Municipios.sql | sed -e "s/''/#/g"  | sed -e "s/'//g" | sed -e "s/)//g" | sed -e "s/;//g" | sed -e s/", "/,/g | sed -e "s/#/'/g" > city-state.csv
  rm Municipios.sql
fi
