# Schemaver-check

Features:

1. Uses single json schema file with multiple definitions (configurable via `SCHEMAVERCHECK_SCHEMA_FILE`)
2. Multiple definitions should be named as `#/definition/${DefinitionName}`
3. Verifies local data mocks against schema and specified definition
4. Calls remote endpoint to verify if specified definition is compatible with remote one
5. Has ability to skip remote check (via `SCHEMAVERCHECK_SKIP_REMOTE_CHECK`)
6. Retries the remote check 10 times before failing (useful for long container startup time)

Environment variables

1. `SCHEMAVERCHECK_SCHEMA_FILE=/path/to/schema.json`
2. `SCHEMAVERCHECK_ENDPOINT_URL_FORMAT=http://local-service:80/verify_json_schema?definitionName=%s`
3. `SCHEMAVERCHECK_SKIP_REMOTE_CHECK=1` — allows skipping remote check (useful for local testing)

## Remote endpoint spec

The remote check endpoint compatibility guide:

1. Accepts POST request with a GET/url parameter with DefinitionName
2. The body of the request consists of JSON Schema file
3. 200 — If the specified definition of sent schema is compatible with the one's of endpoint
4. 409 — If incomptible
5. 4xx — On request errors
6. 5xx — On server error
7. Error response should have `error` field (string)
