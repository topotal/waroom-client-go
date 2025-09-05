package main

import (
	"context"
	"fmt"
	"log"

	"waroom-client-go/waroom"

	"github.com/google/uuid"
)

// https://api.staging.waroom.com/api/v0

func main() {
	ctx := context.Background()

	client, err := waroom.NewClient(ctx,
		waroom.WithBaseURL("https://api.staging.waroom.com/api/v0"),
		waroom.WithUserAgent("waroom-sdk-go/0.1"),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Client initialized")

	// 例: 1ページだけ取得
	uuid, _ := uuid.Parse("d32fc8dd-e424-4974-b27e-c55bb323a6cb")
	res, err := client.GetIncident(ctx, uuid)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}
