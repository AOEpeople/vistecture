FROM golang:alpine

ADD https://github.com/AOEpeople/vistecture/releases/download/v2.0.9/vistecture-linux /usr/local/bin/vistecture
RUN chmod +x /usr/local/bin/vistecture
COPY result /artefacts/result
COPY demoproject /artefacts/demoproject

WORKDIR /artefacts
CMD vistecture --config=demoproject/project.yml --skipValidation=1 serve --staticDocumentsFolder=result
EXPOSE 8080
