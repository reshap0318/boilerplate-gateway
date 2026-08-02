#!/bin/sh
set -e

PRIVATE_KEY="${JWT_PRIVATE_KEY_PATH:-storage/keys/private.pem}"
PUBLIC_KEY="${JWT_PUBLIC_KEY_PATH:-storage/keys/public.pem}"

if [ ! -f "$PRIVATE_KEY" ] || [ ! -f "$PUBLIC_KEY" ]; then
    echo "JWT keys not found, generating..."
    /genkey -f
fi

echo "Running database migrations..."
/migrate up

exec /server
