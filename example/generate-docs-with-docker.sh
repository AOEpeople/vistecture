#!/bin/sh

docker run -v $(pwd):/workspace aoepeople/vistecture:latest  sh -c "cd /workspace && ./generate-documentation.sh"