#!/bin/sh
# healthcheck.sh
wget --no-verbose --tries=1 --spider http://localhost:9000/health || exit 1