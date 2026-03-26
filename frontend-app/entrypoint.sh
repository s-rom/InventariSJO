#!/bin/sh
set -e

LE_CERT=/etc/letsencrypt/live/inventarisjo.duckdns.org/fullchain.pem
LE_KEY=/etc/letsencrypt/live/inventarisjo.duckdns.org/privkey.pem
ACTIVE=/etc/ssl/active

mkdir -p "$ACTIVE"

if [ -f "$LE_CERT" ]; then
    echo "Using Let's Encrypt certificate."
    cp "$LE_CERT" "$ACTIVE/cert.pem"
    cp "$LE_KEY"  "$ACTIVE/key.pem"
else
    echo "Let's Encrypt cert not ready yet, using self-signed."
    echo "Restart this container after certbot finishes to activate the LE cert."
    if [ ! -f "$ACTIVE/cert.pem" ]; then
        openssl req -x509 -nodes -days 3650 -newkey rsa:2048 \
            -keyout "$ACTIVE/key.pem" \
            -out    "$ACTIVE/cert.pem" \
            -subj "/CN=inventari-local"
    fi
fi

exec nginx -g "daemon off;"
