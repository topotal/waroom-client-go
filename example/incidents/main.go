package main

import (
	"context"
	"errors"
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

	c := newClient()

	uuid, _ := uuid.Parse("d32fc8dd-e424-4974-b27e-c55bb323a6cb")
	res, err := c.Incidents().GetIncident(ctx, uuid)
	if errors.Is(err, waroom.ErrNotFound) {
		fmt.Println("incident not found")
		return
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}

func listIncidents() {
	ctx := context.Background()
	c := newClient()

	// 1ページ分をページ情報付きで取得
	list, err := c.Incidents().ListIncidents(ctx, waroom.ListIncidentsOptions{
		SortOrder: waroom.NewOptListIncidentsSortOrder(waroom.ListIncidentsSortOrderDesc),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d incidents (page %d/%d)\n", list.TotalCount, list.CurrentPage, list.TotalPages)
	fmt.Println(list.Items)

	// 全件が欲しい場合はイテレータでページングを自動化
	for inc, err := range c.Incidents().ListIncidentsAll(ctx, waroom.ListIncidentsOptions{}) {
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(inc.Title)
	}
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
