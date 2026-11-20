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

1. Move the entries under `## [Unreleased]` in `CHANGELOG.md` into a new `## [vX.Y.Z] - YYYY-MM-DD` section, update the compare links at the bottom, and merge it to `main`.
2. Tag the merge commit and push the tag:
   ```sh
   git tag vX.Y.Z origin/main
   git push origin vX.Y.Z
   ```
3. `.github/workflows/release.yaml` runs GoReleaser and publishes a GitHub Release whose notes are that version's section of `CHANGELOG.md` (tags with a suffix such as `-alpha` become prereleases). The workflow fails if `CHANGELOG.md` has no section for the tag.
