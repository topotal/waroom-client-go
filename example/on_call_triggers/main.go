package main

import (
	"context"
	"fmt"
	"log"

	"github.com/topotal/waroom-client-go/waroom"
)

func main() {
	svc := "awesome-service"

	listOnCallTriggers(svc)
	// createOnCallTrigger(svc, "critical incidents", "<escalation-policy-uuid>")
}

func listOnCallTriggers(sname string) {
	ctx := context.Background()
	c := newClient()

	//TODO: pagination
	triggers, err := c.OnCallTriggers.List(ctx, waroom.ListOnCallTriggersParams{
		ServiceName: waroom.NewOptString(sname),
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, trigger := range triggers {
		fmt.Println(trigger.UUID)
		fmt.Println(trigger.Name)
		fmt.Println(trigger.Enabled)
		fmt.Println(trigger.SeverityFilters)
	}
}

func createOnCallTrigger(sname, name, escalationPolicyUUID string) {
	ctx := context.Background()
	c := newClient()

	req := waroom.CreateOnCallTriggerReq{
		Name:                  name,
		ServiceName:           sname,
		EscalationPolicyUuids: []string{escalationPolicyUUID},
		RetryCount:            waroom.NewOptInt(1),
	}
	req.SeverityFilters = append(req.SeverityFilters, "critical", "high")

	trigger, err := c.OnCallTriggers.Create(ctx, req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(trigger.UUID)
	fmt.Println(trigger.Name)
}

func updateOnCallTrigger(id string) {
	ctx := context.Background()
	c := newClient()

	trigger, err := c.OnCallTriggers.Update(ctx, waroom.UpdateOnCallTriggerReq{
		Enabled: waroom.NewOptBool(false),
	}, waroom.UpdateOnCallTriggerParams{UUID: id})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(trigger)
}

func deleteOnCallTrigger(id string) {
	ctx := context.Background()
	c := newClient()

	if err := c.OnCallTriggers.Delete(ctx, waroom.DeleteOnCallTriggerParams{UUID: id}); err != nil {
		log.Fatal(err)
	}
}

func newClient() *waroom.Client {
	ctx := context.Background()

	client, err := waroom.NewClient(ctx,
		waroom.WithBaseURL("https://api.staging.waroom.com"),
		waroom.WithUserAgent("waroom-sdk-go/0.0.1"),
		waroom.WithRequestDump(),
		waroom.WithResponseDump(),
	)
	if err != nil {
		log.Fatal(err)
	}

	return client
}
