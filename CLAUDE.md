# Claude Code Context

This file provides context for Claude Code (and other LLM assistants) when working with this repository.

## Project Context Reference

**IMPORTANT**: Before working on this codebase, read [llms.txt](llms.txt) for comprehensive project context.

The [llms.txt](llms.txt) file contains essential information for LLMs:
- Project description and purpose (Go Client SDK for Waroom API)
- Official documentation references (Waroom API, OpenAPI Specification, Ogen)
- Implementation guidelines and best practices
- Links to example files and resources

## Instructions for LLM Assistants

1. **Always read llms.txt first** when starting work on this repository
2. Follow the Ogen documentation and OpenAPI Specification linked in llms.txt
3. Refer to the example directory for usage patterns
4. Ensure all implementations comply with the Waroom API specification

## Quick Start

When implementing features or fixing bugs:
- Check [llms.txt](llms.txt) for relevant documentation links
- Review existing code in `waroom/` for patterns
- Test using examples in `example/` directory
- Use `task test` to run tests

## Bug Reports and Behavior Investigation

waroom-client-go is the entry point for bug reports and API behavior questions across the Waroom ecosystem. [waroomctl](https://github.com/topotal/waroomctl) and [terraform-provider-waroom](https://github.com/topotal/terraform-provider-waroom) search existing issues here and in [waroom-backend](https://github.com/topotal/waroom-backend), then file new ones here instead of going to waroom-backend or [waroom](https://github.com/topotal/waroom) directly. Do the same when a problem comes up in this repository.

When working on such an issue, triage it here first and decide whether an upstream issue is needed:

1. **Reproduce and locate the cause**, using `waroom.WithRequestDump()` / `waroom.WithResponseDump()` to see the actual traffic:
   - A bug in the `waroom/` wrappers or in how the client is generated → fix it here. No upstream issue.
   - `api/openapi.yaml` is behind the published spec (`https://static.waroom.com/openapi/v0.yaml`) → sync the spec and regenerate. No upstream issue.
   - The published spec is wrong, or the API behaves differently from it or inconsistently across endpoints → file in **waroom-backend**.
   - A product or specification decision rather than an implementation bug (what the API should offer, behavior shared with the web app) → file in **waroom**.
2. **Search upstream before filing**: `gh search issues --repo topotal/waroom-backend --repo topotal/waroom "<keywords>"`, including closed issues. Link an existing one instead of filing a duplicate.
3. **Link both ways.** The upstream issue links back to the waroom-client-go issue, and the waroom-client-go issue records the upstream link and stays open until the client side follows the fix (spec sync, regenerate, release). Write upstream issues in the language their existing issues use.
4. Waroom is multi-tenant: use placeholders for tenant data (incident titles, service names, user names) in every issue, and never paste API keys.
