package main

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"

	"github.com/topotal/waroom-client-go/waroom"
)

func main() {
	getServiceLabels()
	createServiceLabel("fuga", "ffF000")
}

func getServiceLabels() {
	ctx := context.Background()
	c := newClient()

	svc := waroom.NewLabelService(c)
	param := waroom.ListServiceLabelsParams{
		Name: "ないとぅー",
	}

	//TODO: pagination
	res, err := svc.ListServiceLabels(ctx, param)
	if err != nil {
		log.Fatal(err)
	}

	for _, label := range res {
		fmt.Println(label.UUID)
		fmt.Println(label.Name)
		fmt.Println(label.Color)
	}
}

func createServiceLabel(name, color string) {
	ctx := context.Background()
	c := newClient()

	svc := waroom.NewLabelService(c)
	req := waroom.Label{
		Name: name,
	}
	req.Color.SetTo(color)

	param := waroom.CreateServiceLabelParams{
		Name: "ないとぅー",
	}

	res, err := svc.CreateServiceLabel(ctx, req, param)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res.UUID)
	fmt.Println(res.Name)
	fmt.Println(res.Color)
}

func deleteServiceLabel(service_name string, id uuid.UUID) {
	ctx := context.Background()
	c := newClient()

	svc := waroom.NewLabelService(c)
	param := waroom.DeleteServiceLabelParams{
		Name: service_name,
		UUID: id,
	}
	if err := svc.DeleteServiceLabel(ctx, param); err != nil {
		log.Fatal(err)
	}
}

func newClient() *waroom.Client {
	ctx := context.Background()

	client, err := waroom.NewClient(ctx,
		waroom.WithBaseURL("https://api.staging.waroom.com/api/v0"),
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
