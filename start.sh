#!/usr/bin/env sh

# Check if hydra is ready
HYDRA_HEALTH_CHECK_URL="http://hydra:4444/health/ready"

for I in $(seq 9)
do
  RETURN_CODE=$(curl -s -X GET -w "%{http_code}" -o /dev/null $HYDRA_HEALTH_CHECK_URL)
  if [[ $RETURN_CODE -eq 200 ]]
  then
    echo "Hydra is ready"
    break
  else
    echo "Hydra is still not ready"
    sleep 3
  fi
done

# Create client in hydra
RANDOM_CLIENT_ID=$(tr -dc A-Za-z0-9 </dev/urandom | head -c 13 ; echo '')
HYDRA_CLIENT_JSON=$(cat <<-END
{
  "client_id": "${RANDOM_CLIENT_ID}",
  "client_secret": "secret",
  "grant_types": ["client_credentials"],
  "scope": "scope1"
}
END
)

response=$(curl -d "${HYDRA_CLIENT_JSON}" -H 'Content-Type: application/json' -X POST http://hydra:4445/clients)

CLIENT_ID=$(echo "$response" | jq -r .client_id)
CLIENT_SECRET=$(echo "$response" | jq -r .client_secret)

## Create config for furby
CONFIG_VAR=$(cat <<-END
---

auth:
    type: 'noop'

client_credentials:
    id: "${CLIENT_ID}"
    secret: "${CLIENT_SECRET}"
    scopes:
      - "scope1"
    url: "http://hydra:4444/oauth2/token"
END
)

echo "${CONFIG_VAR}" > /go/src/github.com/dpattmann/furby/configs/furby_temp_config.yaml

# Show some info for using furby
echo "----------------------------------------"
echo "Your client_id is ${CLIENT_ID}"
echo "Your client_secret is ${CLIENT_SECRET}"
echo "----------------------------------------"

# Start furby with a created config
furby -p /go/src/github.com/dpattmann/furby/configs/furby_temp_config.yaml

