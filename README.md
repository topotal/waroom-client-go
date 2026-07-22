# waroom-client-go

waroom-client-go is a Go client library for the Waroom API.

> [!WARNING]
> The API of this library is not stable yet.

## Usage

```go
import (
    "context"
    "fmt"

    "github.com/waroom/waroom-client-go/waroom"
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

