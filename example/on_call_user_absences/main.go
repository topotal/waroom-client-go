package main

import (
	"context"
	"fmt"
	"log"

	"github.com/topotal/waroom-client-go/waroom"
)

func main() {

	// The user's sub. Characters such as "|" are percent-encoded by the client.
	sub := "auth0|1234567890"

	listAbsences(sub)
	// createAbsence(sub, "2026-12-01T00:00:00+09:00", "2026-12-03T00:00:00+09:00", "Vacation")
}

func listAbsences(sub string) {
	ctx := context.Background()
	c := newClient()

	param := waroom.ListOnCallUserAbsencesParams{
		Sub: sub,
	}
	// Also return finished one-off absences.
	param.IncludeExpired.SetTo(true)

	//TODO: pagination
	res, err := c.OnCallUserAbsences.List(ctx, param)
	if err != nil {
		log.Fatal(err)
	}

	for _, absence := range res {
		fmt.Println(absence.UUID)
		fmt.Println(absence.StartAt, absence.EndAt)
		fmt.Println(absence.RepeatType)
		if reason, ok := absence.Reason.Get(); ok {
			fmt.Println(reason)
		}
	}
}

func createAbsence(sub, startAt, endAt, reason string) {
	ctx := context.Background()
	c := newClient()

	param := waroom.CreateOnCallUserAbsenceParams{
		Sub: sub,
	}

	body := waroom.CreateOnCallUserAbsenceReq{
		StartAt: startAt,
		EndAt:   endAt,
		Reason:  waroom.NewOptNilString(reason),
	}
	// Omit RepeatType for a one-off absence ("none").
	body.RepeatType.SetTo("weekly")

	res, err := c.OnCallUserAbsences.Create(ctx, body, param)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res.UUID)
	fmt.Println(res.StartAt, res.EndAt)
	fmt.Println(res.RepeatType)
}

func updateAbsence(sub, id, endAt string) {
	ctx := context.Background()
	c := newClient()

	param := waroom.UpdateOnCallUserAbsenceParams{
		Sub:  sub,
		UUID: id,
	}

	body := waroom.UpdateOnCallUserAbsenceReq{
		EndAt: waroom.NewOptString(endAt),
	}

	absence, err := c.OnCallUserAbsences.Update(ctx, body, param)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(absence)
}

func deleteAbsence(sub, id string) {
	ctx := context.Background()
	c := newClient()

	param := waroom.DeleteOnCallUserAbsenceParams{
		Sub:  sub,
		UUID: id,
	}
	if err := c.OnCallUserAbsences.Delete(ctx, param); err != nil {
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
