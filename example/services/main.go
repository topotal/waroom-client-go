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
	listServices()
}

func getIncidentUUID() {
	ctx := context.Background()

	client := newClient()

	// 例: 1ページだけ取得
	uuid, _ := uuid.Parse("d32fc8dd-e424-4974-b27e-c55bb323a6cb")
	res, err := client.GetIncident(ctx, uuid)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}

func listServices() {
	ctx := context.Background()
	c := newClient()

	svc := waroom.NewServiceService(c)
	param := waroom.ListServicesParams{}

	//TODO: pagination
	res, err := svc.ListServices(ctx, param)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}

func newClient() *waroom.WaroomClient {
	ctx := context.Background()

	client, err := waroom.NewWaroomClient(ctx,
		waroom.WithBaseURL("http://api.staging.waroom.com/api/v0"),
		// waroom.WithBaseURL("http://127.0.0.1:4010"),
		waroom.WithUserAgent("waroom-sdk-go/0.0.1"),
		waroom.WithResponseDump(),
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Print("Client initialized")
	return client
}
