# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [Unreleased]

First release of the Go client library for the Waroom API.

### Features

- Client built on the published Waroom public API spec (OpenAPI 3.1), with each resource exposed as a service field (e.g. `client.Incidents.List(ctx, ...)`).
- Supported services: `Incidents`, `Services`, `Labels`, `ActionItems`, `Runbooks`, `RunbookRules`, `Postmortems`, `PostmortemTemplates`, `IncidentDescriptionTemplate`, `WebhookConfigurations`, `Users`.
- API key authentication via the `WithAPIKey` option or the `WAROOM_API_KEY` environment variable.
- `APIError` type for server-side API errors, carrying the status code and error messages.
- Client options: `WithBaseURL`, `WithHTTPClient`, `WithUserAgent`, `WithRequestDump`, `WithResponseDump`.

[Unreleased]: https://github.com/topotal/waroom-client-go/compare/v0.0.0-alpha...HEAD
