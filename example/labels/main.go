package main

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"

	"github.com/topotal/waroom-client-go/waroom"
)

func main() {

	svc := "awesome-service"

	getServiceLabels(svc)
	// createServiceLabel(svc, "test-06", "FF4030")
}

func getServiceLabels(sname string) {
	ctx := context.Background()
	c := newClient()

	// 1ページ分をページ情報付きで取得
	list, err := c.Labels().ListServiceLabels(ctx, sname, waroom.ListServiceLabelsOptions{
		PerPage: waroom.NewOptInt(50),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d labels (page %d/%d)\n", list.TotalCount, list.CurrentPage, list.TotalPages)

	for _, label := range list.Items {
		fmt.Println(label.UUID)
		fmt.Println(label.Name)
		fmt.Println(label.Color)
	}

	// 全件が欲しい場合はイテレータでページングを自動化
	for label, err := range c.Labels().ListServiceLabelsAll(ctx, sname, waroom.ListServiceLabelsOptions{}) {
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(label.Name)
	}
}

func createServiceLabel(sname, label, color string) {
	ctx := context.Background()
	c := newClient()

	res, err := c.Labels().CreateServiceLabel(ctx, sname, waroom.CreateServiceLabelRequest{
		Name:  label,
		Color: color,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res.UUID)
	fmt.Println(res.Name)
	fmt.Println(res.Color)
}

func updateServiceLabel(sname string, id uuid.UUID) {
	ctx := context.Background()
	c := newClient()

	label, err := c.Labels().UpdateServiceLabel(ctx, sname, id, waroom.UpdateServiceLabelRequest{
		Name:  waroom.NewOptString("updated-label"),
		Color: waroom.NewOptString("FF00FF"),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(label)
}

func deleteServiceLabel(sname string, id uuid.UUID) {
	ctx := context.Background()
	c := newClient()

	if err := c.Labels().DeleteServiceLabel(ctx, sname, id); err != nil {
		log.Fatal(err)
	}
}

func newClient() *waroom.Client {
	ctx := context.Background()

	client, err := waroom.NewClient(ctx,
		waroom.WithBaseURL("https://api.staging.waroom.com"),
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
