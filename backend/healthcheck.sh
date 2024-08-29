#!/bin/sh
# healthcheck.sh
curl -X GET http://localhost:9000/api/v1/health || exit 1