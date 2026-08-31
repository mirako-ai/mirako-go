# Repository Guidelines

## About testing code
`mirako-go` ships generated client stubs based on the Mirako's live Openapi spec.
- Do not add endpoint-specific tests for generated schemas, paths, payloads, or types.
- Test only clearly necessary, handwritten SDK behavior, preferably at the cross-cutting client level.
- Regeneration must pass `go build ./...` and `go test ./...` before release.
