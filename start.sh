#!/usr/bin/env sh

# Create client in hydra
HYDRA_CLIENT_JSON=$(cat <<-END
{
  "client_id": "integration-tests",
  "client_secret": "secret",
  "grant_types": ["client_credentials"],
  "scope": "scope1"
}
END
)

response=$(curl -d "${HYDRA_CLIENT_JSON}" -H 'Content-Type: application/json' -X POST http://127.0.0.1:4445/clients)

CLIENT_ID=$(echo "$response" | jq -r .client_id)
CLIENT_SECRET=$(echo "$response" | jq -r .client_secret)

# Show some info for using furby
echo "----------------------------------------"
echo "Your client_id is ${CLIENT_ID}"
echo "Your client_secret is ${CLIENT_SECRET}"
echo "----------------------------------------"


