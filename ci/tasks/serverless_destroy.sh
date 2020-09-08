#!/bin/bash

set -euo pipefail

: ${WORKSPACE}
: ${AWS_REGION}
: ${S3_NAME}
: ${S3_KEY}

. ./spp-bw-etl/ci/tasks/util/assume_role.sh
. ./spp-bw-etl/ci/tasks/util/setup_terraform.sh

terraform destroy --auto-approve

echo "done"