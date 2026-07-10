package main

import (
	"context"
	"fmt"
	"log"

	"github.com/topotal/waroom-client-go/waroom"

	"github.com/google/uuid"
)

func main() {
	listIncidents()
}

func getIncidentUUID() {
	ctx := context.Background()

	client := newClient()
	svc := waroom.NewIncidentService(client)

	// 例: 1ページだけ取得
	uuid, _ := uuid.Parse("d32fc8dd-e424-4974-b27e-c55bb323a6cb")
	res, err := svc.GetIncident(ctx, uuid)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}

func listIncidents() {
	client := newClient()
	svc := waroom.NewIncidentService(client)
	param := waroom.ListIncidentsParams{}

	//TODO: pagination
	ctx := context.Background()
	res, err := svc.ListIncidents(ctx, param)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}

func newClient() *waroom.Client {
	ctx := context.Background()

	client, err := waroom.NewClient(ctx,
		waroom.WithBaseURL("https://api.staging.waroom.com"),
		waroom.WithUserAgent("waroom-sdk-go/0.0.1"),
		waroom.WithRequestDump(),
	)

	if err != nil {
		log.Fatal(err)
	}

	return client
}
