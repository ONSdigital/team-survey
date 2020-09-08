# Ci-bootstrap

Resources for setting up roles in a new AWS account

## Pre-requisites

The terraform state for the role is stored in the state bucket
for the environment under the key `bootstrap/deploy/sor-role.tfstate`

To this end, the bucket must have already been created in the environent (see
[bootstrap instructions in bpm-ci](https://github.com/ONSdigital/bpm-ci/tree/master/bootstrap))

## Running

```bash
# (export creds for the envionment you're bootstrapping)

# MUST switch into the folder first or relative paths will fail!
cd ci-bootstrap

# Create the role
# Where ENVIRONMENT is sandbox, cicd, uat, staging or prod
ENVIRONMENT=<accountname> ./bin/deploy.sh

# Delete the role
# Where ENVIRONMENT is sandbox, cicd, uat, staging or prod
ENVIRONMENT=<accountname> ./bin/destroy.sh
```
