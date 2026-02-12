#!/usr/bin/env sh

source ./auth.sh

echo -e "\nPublic"
curl -s http://localhost:8080/hello-public

echo -e "\nToken"
curl -s http://localhost:8080/hello-token \
-H "Authorization: Bearer $ACCESS_TOKEN" | jq '.'

echo -e "\nAdmin"
curl -s http://localhost:8080/hello-admin \
-H "Authorization: Bearer $ACCESS_TOKEN" | jq '.'
