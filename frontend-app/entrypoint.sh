#!/bin/sh
set -e

CERT_DIR=/etc/ssl/selfsigned

if [ ! -f "$CERT_DIR/cert.pem" ]; then
    echo "Generating self-signed certificate..."
    mkdir -p "$CERT_DIR"
    openssl req -x509 -nodes -days 3650 -newkey rsa:2048 \
        -keyout "$CERT_DIR/key.pem" \
        -out    "$CERT_DIR/cert.pem" \
        -subj "/CN=inventari-local" \
        -addext "subjectAltName=IP:0.0.0.0,IP:127.0.0.1"
    echo "Certificate generated."
fi

exec nginx -g "daemon off;"
