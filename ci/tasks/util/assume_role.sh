#!/bin/bash

set -euo pipefail

# Extract credentials from the passed in role/credentials.json file to
# be able to assume the deployment role. If this file doesn't exist, then
# we assume there's no role to assume!
if [ -f role/credentials.json ]; then
    echo "Assuming deployment role"
    export AWS_ACCESS_KEY_ID=$(jq -r .Credentials.AccessKeyId role/credentials.json)
    export AWS_SECRET_ACCESS_KEY=$(jq -r .Credentials.SecretAccessKey role/credentials.json)
    export AWS_SESSION_TOKEN=$(jq -r .Credentials.SessionToken role/credentials.json)
fi
