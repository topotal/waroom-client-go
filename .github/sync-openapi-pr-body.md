Sync `api/openapi.yaml` with the published spec and regenerate `internal/ogen`.

| | |
| --- | --- |
| spec | ${OPENAPI_URL} |
| version | `${MAIN_VERSION}` → `${VERSION}` |
| last-modified | ${LAST_MODIFIED} |
| etag | `${ETAG}` |
| x-amz-version-id | `${S3_VERSION_ID}` |
| sha256 | `${SHA256}` |
| replaces | ${REPLACES} |

## Added / removed operationIds

```diff
${OPS}
```

## Checklist

- [ ] UnitTest passes (a failure means existing wrappers are affected by the spec change)
- [ ] File issues for added operations with the spec-coverage skill
