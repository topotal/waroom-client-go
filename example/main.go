package main

import (
	"context"
	"fmt"
	"log"

	"github.com/topotal/waroom-client-go/waroom"

	"github.com/google/uuid"
)

func main() {
	// getIncidentUUID()
	listIncidents()
}

func getIncidentUUID() {
	ctx := context.Background()

	client := NewClient()

	// 例: 1ページだけ取得
	uuid, _ := uuid.Parse("d32fc8dd-e424-4974-b27e-c55bb323a6cb")
	res, err := client.GetIncident(ctx, uuid)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}

func listIncidents() {
	ctx := context.Background()
	client := NewClient()

	res, err := client.ListIncidents(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}

func NewClient() *waroom.Client {
	ctx := context.Background()

	client, err := waroom.NewClient(ctx,
		waroom.WithBaseURL("http://api.staging.waroom.com/api/v0"),
		waroom.WithUserAgent("waroom-sdk-go/0.0.1"),
		waroom.WithDump(true),
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Print("Client initialized")
	return client
}
