package main

import (
	"context"
	"log"

	"github.com/topotal/waroom-client-go/waroom"
)

func newClient() *waroom.WaroomClient {
	ctx := context.Background()

	client, err := waroom.NewWaroomClient(ctx,
		waroom.WithBaseURL("http://api.staging.waroom.com/api/v0"),
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
