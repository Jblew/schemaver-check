ARTIFACT="jsonschema-compatibility-checker"
docker build --progress=plain -t ${ARTIFACT} .
docker run -it ${ARTIFACT}