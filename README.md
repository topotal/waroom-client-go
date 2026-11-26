# waroom-client-go

waroom-client-go is a Go client library for the Waroom API.

> [!WARNING]
> The API of this library is not stable yet.

## Usage

```go
import (
    "context"
    "fmt"

    "github.com/topotal/waroom-client-go/waroom"
)

func main() {
    ctx := context.Background()

    client, err := waroom.NewClient(ctx)
    if err != nil {
        fmt.Println("Error creating waroom client:", err)
        return
    }

    res, err := client.Services.List(ctx, waroom.ListServicesParams{})
    if err != nil {
        fmt.Println("Error listing services:", err)
        return
    }

    fmt.Printf("%+v\n", res)
}
```

`NewClient` reads the API key from `WithAPIKey` or the `WAROOM_API_KEY` environment variable, and sends requests to the production API (`https://api.app.waroom.com`, `waroom.DefaultBaseURL`). Use `waroom.WithBaseURL` to target another host.

For debugging, `waroom.WithRequestDump()` and `waroom.WithResponseDump()` write the HTTP traffic to stderr. Credential headers such as `Authorization` are redacted in the request dump.

## Available services

Access each service through a field on `*waroom.Client` (e.g. `client.Incidents.List(ctx, ...)`).

| Field | Methods |
| --- | --- |
| `Incidents` | `List`, `Get`, `Create`, `Update`, `Delete`, `UpdateStatus`, `UpdateSeverity`, `UpdateCommander`, `UpdateResponders`, `CreateMetrics` |
| `Services` | `List`, `UpdateIncidentService`, `GetArchitectureContext`, `UpdateArchitectureContext`, `GetNotificationSlackChannel`, `PutNotificationSlackChannel`, `DeleteNotificationSlackChannel`, `GetIncidentDocumentTemplate`, `UpdateIncidentDocumentTemplate`, `DeleteIncidentDocumentTemplate` |
| `Labels` | `List`, `Create`, `Update`, `Delete`, `UpdateIncidentLabels` |
| `ActionItems` | `List`, `Create`, `Get`, `Delete` |
| `RunbookRules` | `List`, `Create`, `Update`, `Delete` |
| `WebhookConfigurations` | `List`, `Get` |
| `Postmortems` | `List`, `Create`, `Get`, `Update`, `Delete`, `GetTemplate`, `UpdateTemplate` |
| `PostmortemTemplates` | `List`, `Create`, `Get`, `Update`, `Delete` |
| `Runbooks` | `List`, `Create`, `Get`, `Update`, `Delete` |
| `IncidentDescriptionTemplate` | `Get`, `Update`, `Delete` |
| `Users` | `List` |

For convenience, `Incidents` also exposes wrappers that delegate to the canonical services: `ListActionItems` / `CreateActionItem` / `GetActionItem` / `DeleteActionItem` (→ `ActionItems`), `UpdateLabels` (→ `Labels.UpdateIncidentLabels`), and `UpdateService` (→ `Services.UpdateIncidentService`).

## Development

### Prerequisites

- Task
    - https://taskfile.dev/docs/installation
- GoReleaser
    - https://goreleaser.com/install/


### Release

Releases are automated by [tagpr](https://github.com/Songmu/tagpr) (`.github/workflows/release.yaml`). Don't push tags or edit `CHANGELOG.md` by hand.

1. On every push to `main`, tagpr opens or updates a release pull request (`tagpr` label) that bumps the version and adds the release's entry to `CHANGELOG.md`.
2. The release notes are GitHub's generated notes from the merged pull requests, grouped by label as configured in `.github/release.yml` (`enhancement` → Features, `bug` → Fixes, ...). Label your pull requests accordingly.
3. The next version is a patch bump by default. Add the `tagpr:minor` or `tagpr:major` label to the release pull request for a bigger bump.
4. Merging the release pull request tags `main`, and GoReleaser publishes the GitHub Release.
