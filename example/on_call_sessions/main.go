package main

import (
	"context"
	"fmt"
	"log"

	"github.com/topotal/waroom-client-go/waroom"
)

func main() {
	listOnCallSessions("running")
	// getOnCallSession("session-uuid")
}

func listOnCallSessions(status string) {
	ctx := context.Background()
	c := newClient()

	param := waroom.ListOnCallSessionsParams{
		Status: waroom.NewOptOnCallSessionStatusFilter(waroom.OnCallSessionStatusFilter(status)),
	}

	//TODO: pagination
	res, err := c.OnCallSessions.List(ctx, param)
	if err != nil {
		log.Fatal(err)
	}

	for _, session := range res {
		fmt.Println(session.UUID)
		fmt.Println(session.Status)
		fmt.Println(session.Incident.UUID)
		fmt.Println(session.CreatedAt)
	}
}

func getOnCallSession(uuid string) {
	ctx := context.Background()
	c := newClient()

	session, err := c.OnCallSessions.Get(ctx, waroom.GetOnCallSessionParams{UUID: uuid})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(session.UUID)
	fmt.Println(session.Status)
	for _, activity := range session.Activities {
		fmt.Println(activity.OccurredAt, activity.ActivityType)
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
