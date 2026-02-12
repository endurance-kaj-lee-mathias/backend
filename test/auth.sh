#!/usr/bin/env sh

ACCESS_TOKEN=$(curl -s -X POST "http://localhost:8180/realms/endurance/protocol/openid-connect/token" \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -d "grant_type=password" \
    -d "client_id=backend" \
    -d "client_secret=Hvkh9JcoLjjNFT8FIgd3qPpJZWT3E7v9" \
    -d "username=kaj" \
    -d "password=password" \
	-d "scope=openid" \
	| jq -r '.access_token'
)

if [ "$ACCESS_TOKEN" = "null" ] || [ -z "$ACCESS_TOKEN" ]; then
    echo "token is null"
    exit 1
fi

echo "got access token: ${ACCESS_TOKEN:0:5}..."
