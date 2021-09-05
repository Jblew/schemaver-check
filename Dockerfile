FROM golang:1.17-alpine
WORKDIR /src
ADD . /src
# RUN ./test.sh
RUN GOBIN=/bin/ CGO_ENABLED=0 go install

CMD ["/bin/sh", "-c", "\
    SCHEMAVERCHECK_SKIP_REMOTE_CHECK=1 \
    SCHEMAVERCHECK_SCHEMA_PATH=/src/mock/schema.json \
    /bin/schemaver-check --data-file /src/mock/data_valid.json --definition-name \"ChartSpec\" \
    "]