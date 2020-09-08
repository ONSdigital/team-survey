#!/bin/bash

# Runs baseline terraform steps to set up the s3 backend
# and switch to the correct workspace

set -euo pipefail

: ${WORKSPACE}
: ${TERRAFORM_SOURCE}

cd ${TERRAFORM_SOURCE}

# If a backend.tf hasn't already been written then generate one so that
# terraform can correctly configure its state files.
if [ ! -f backend.tf ]; then
echo "Writing backend.tf"
echo "\
provider \"aws\" {
  region = \"eu-west-2\"
}

terraform {
  backend \"s3\" {
    bucket               = \"${S3_NAME}\"
    key                  = \"${S3_KEY}\"
    region               = \"${AWS_REGION}\"
    workspace_key_prefix = \"workspaces\"
  }
}
" > backend.tf
fi

# Initialise terraform and switch to the appropriate workspace
echo "Initialising terraform and select workspace"
terraform init
terraform workspace select ${WORKSPACE} ||
    {
        echo "Workspace '${WORKSPACE}' not found, creating"
        terraform workspace new ${WORKSPACE}
    }
