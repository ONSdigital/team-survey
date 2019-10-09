#!/usr/bin/env bash
cd team-survey/web/test;
cypress run --env HOST="${HOST}",API_SERVER="${API_SERVER}",ADMIN_USERNAME="${ADMIN_USERNAME}",ADMIN_PASSWORD="${ADMIN_PASSWORD}"