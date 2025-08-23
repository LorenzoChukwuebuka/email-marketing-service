#!/bin/sh
# healthcheck.sh
curl -X GET http://localhost:9002/api/v1/health || exit 1