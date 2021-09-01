There is a need for small cli utility that would:

1. Verify JSON against schema (for demo and testing sets)
2. Verify schema version against schema version verification endpoint
3. Should be easy to install on docker images
4. Should work on alpine
5. Candidate lang: golang
   — easy to build and install library from github
   - no need to deploy binaries anywhere (just go install github...)
   - can be easily build on multistage build — see above
6. The utility:
   - verify-data
   - assert-schema-version-compatibility
