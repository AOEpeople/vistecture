#!/bin/sh

./generate-docs-with-docker.sh

echo "Starting vistecture serve.. open http://localhost:8080/"
docker run -v $(pwd):/workspace -p 8080:8080 aoepeople/vistecture vistecture --config=/workspace/demoproject/project.yml serve --staticDocumentsFolder=/workspace/result
