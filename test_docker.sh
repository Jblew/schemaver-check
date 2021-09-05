ARTIFACT="schemaver-check"
docker build --progress=plain -t ${ARTIFACT} .
docker run -it ${ARTIFACT}