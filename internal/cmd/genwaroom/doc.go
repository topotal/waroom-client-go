// Command genwaroom generates both the low-level ogen client and the
// public waroom package facade, in two stages:
//
//	genwaroom analyze  spec → internal/ogen (via ogen's own generator)
//	                        → api/sdk-model.yaml (semantic model)
//	genwaroom render   api/sdk-model.yaml → waroom/*_gen.go
//
// The analyze stage runs ogen's generator library (gen.NewGenerator) and
// derives the model directly from ogen's intermediate representation
// (gen/ir) — the same data ogen writes internal/ogen from, so type names
// match by construction. The spec's operationId/tags come along on the
// same IR. Its output, api/sdk-model.yaml, is a committed and purely
// derived artifact: reviewing its diff shows how spec changes affect the
// public API surface. The render stage contains all Go source assembly;
// it reads nothing but the model.
//
// The generator relies on the following structural conventions. If an
// operation violates one of them, generation fails listing the offending
// operationIds instead of emitting wrong code:
//
//   - Each operation has exactly one 2xx response variant; error variant
//     names carry a suffix from the closed set BadRequest, Unauthorized,
//     PaymentRequired, NotFound, UnprocessableEntity matching their
//     status code.
//   - Error bodies are one of: `Errors []string`, `Error OptString`, or
//     `Error string` (extracted at runtime by newAPIError via getters).
//   - A list success variant has exactly two fields: a slice of items and
//     a ResponseMetadata struct with CurrentPage/NextPage/PrevPage/
//     TotalPages/TotalCount.
//   - A request wrapper <Op>Req with a single struct-typed field named
//     <Op>Req<Field> is unwrapped so methods take the inner type.
//   - Header/cookie parameters and optional request bodies are not
//     supported.
//
// Exceptions are handled via .genwaroom.yaml (include/skip lists, group
// and rename overrides); conventions cover everything else.
package main
