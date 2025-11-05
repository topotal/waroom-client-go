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

    svc := waroom.NewServiceService(client)
    res, err := svc.ListServices(ctx, &waroom.ListServiceParams{})
    if err != nil {
        fmt.Println("Error listing services:", err)
        return
    }

    fmt.Printf("%+v\n", res)
}
```

## Development

### Prerequisites

- Task
    - https://taskfile.dev/docs/installation
- GoReleaser
    - https://goreleaser.com/install/

