// Command genwaroom generates the public waroom package facade from the
// ogen-generated client in internal/ogen, in two stages:
//
//	genwaroom analyze  internal/ogen → api/sdk-model.yaml (semantic model)
//	genwaroom render   api/sdk-model.yaml → waroom/*_gen.go
//
// The analyze stage reads the compiled ogen package via go/packages (the
// source of truth for type names and signatures) and the OpenAPI spec only
// for the operationId → tag mapping used to group operations into
// services. Its output, api/sdk-model.yaml, is a committed and purely
// derived artifact: reviewing its diff shows how spec/ogen changes affect
// the public API surface. The render stage contains all Go source
// assembly; it reads nothing but the model.
//
// The generator relies on the following structural conventions of ogen
// output. If an operation violates one of them, generation fails listing
// the offending operationIds instead of emitting wrong code:
//
//   - The Invoker interface lists one method per operation, with a doc
//     comment containing "<Name> invokes <operationId> operation." and a
//     trailing "<METHOD> <path>" line.
//   - Each operation returns a sealed sum type <Op>Res with exactly one
//     success variant suffixed OK, Created, or NoContent.
//   - Error variants are suffixed BadRequest, Unauthorized,
//     PaymentRequired, NotFound, or UnprocessableEntity, and their body is
//     one of: `Errors []string`, `Error OptString`, or `Error string`.
//   - A list success variant has exactly two fields: a slice of items and
//     a ResponseMetadata struct with CurrentPage/NextPage/PrevPage/
//     TotalPages/TotalCount.
//   - A request wrapper <Op>Req with a single struct-typed field named
//     <Op>Req<Field> is unwrapped so methods take the inner type.
//
// Exceptions are handled via .genwaroom.yaml (include/skip lists, group
// and rename overrides); conventions cover everything else.
package main
