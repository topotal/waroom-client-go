package main

import (
	"context"
	"fmt"
	"log"

	"github.com/topotal/waroom-client-go/waroom"
)

func main() {
	// getIncidentUUID()
	listServices()
}

func listServices() {
	ctx := context.Background()
	c := newClient()

	param := waroom.ListServicesParams{}

	//TODO: pagination
	res, err := c.Services.List(ctx, param)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}

func newClient() *waroom.Client {
	ctx := context.Background()

	client, err := waroom.NewClient(ctx,
		waroom.WithBaseURL("http://api.staging.waroom.com"),
		// waroom.WithBaseURL("http://127.0.0.1:4010"),
		waroom.WithUserAgent("waroom-sdk-go/0.0.1"),
		waroom.WithRequestDump(),
		waroom.WithResponseDump(),
	)

	if err != nil {
		log.Fatal(err)
	}

	return client
}
