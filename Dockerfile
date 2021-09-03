FROM golang:1.17-alpine
WORKDIR /src
ADD . /src
RUN GOBIN=/bin/ CGO_ENABLED=0 go install

CMD ["/bin/sh", "-c", "\
    JSCC_SKIP_SCHEMAVER_COMPATIBILITY_CHECK=1 \
    JSCC_SCHEMA_FILE_PATH=/src/mock/schema.json \
    /bin/jsonschema-compatibility-checker --data-file /src/mock/data_valid.json --definition-name \"ChartSpec\" \
    "]