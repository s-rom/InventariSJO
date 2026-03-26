#!/bin/sh
# Obtiene el cert de Let's Encrypt via DNS-01 (DuckDNS) y lo renueva cada 12h.
# No necesita puertos abiertos en el router.
set -e

while true; do
    echo "Running certbot..."
    certbot certonly \
        --non-interactive \
        --agree-tos \
        --email "$ACME_EMAIL" \
        --preferred-challenges dns \
        --authenticator dns-duckdns \
        --dns-duckdns-token "$DUCKDNS_TOKEN" \
        --dns-duckdns-propagation-seconds 60 \
        --keep-until-expiring \
        -d inventarisjo.duckdns.org
    echo "Certbot done. Sleeping 12h..."
    sleep 12h
done
