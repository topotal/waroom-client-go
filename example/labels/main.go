package main

import (
	"context"
	"fmt"
	"log"
	"os"

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

	svc := waroom.NewLabelService(c)
	param := waroom.ListServiceLabelsParams{
		Name: sname,
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

func createServiceLabel(sname, label, color string) {
	ctx := context.Background()
	c := newClient()

	svc := waroom.NewLabelService(c)

	param := waroom.CreateServiceLabelParams{
		Name: sname,
	}

	body := waroom.CreateServiceLabelReq{
		Name:  label,
		Color: color,
	}

	res, err := svc.CreateServiceLabel(ctx, body, param)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res.UUID)
	fmt.Println(res.Name)
	fmt.Println(res.Color)
}

func updateServiceLabel(service_name string, id uuid.UUID) {
	ctx := context.Background()
	c := newClient()

	svc := waroom.NewLabelService(c)
	param := waroom.UpdateServiceLabelParams{
		UUID: id,
		Name: service_name,
	}

	body := waroom.UpdateServiceLabelReq{
		Name:  waroom.NewOptString("updated-label"),
		Color: waroom.NewOptString("FF00FF"),
	}

	label, err := svc.UpdateServiceLabel(ctx, body, param)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(label)
}

func deleteServiceLabel(sname string, id uuid.UUID) {
	ctx := context.Background()
	c := newClient()

	svc := waroom.NewLabelService(c)
	param := waroom.DeleteServiceLabelParams{
		Name: sname,
		UUID: id,
	}
	if err := svc.DeleteServiceLabel(ctx, param); err != nil {
		log.Fatal(err)
	}
}

func newClient() *waroom.Client {
	ctx := context.Background()

	client, err := waroom.NewClient(ctx, os.Getenv("WAROOM_API_KEY"),
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
