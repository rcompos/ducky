#!/bin/bash
HUB_HOST="blackduck.eng.netapp.com"
HUB_AUTH_TOKEN=${BD_TOKEN}

if [ -z "$HUB_AUTH_TOKEN" ]; then
  echo Set BlackDuck Hub token as envvar BD_TOKEN
fi

REQUEST_URL="https://"$HUB_HOST"/api/tokens/authenticate"

curl --request POST --url $REQUEST_URL --header 'authorization: token '$HUB_AUTH_TOKEN --header 'cache-control: no-cache'
