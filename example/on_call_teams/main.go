package main

import (
	"context"
	"fmt"
	"log"

	"github.com/topotal/waroom-client-go/waroom"
)

func main() {
	listOnCallTeams()
	// createOnCallTeam("SRE", "FF4030", "auth0|xxxxxxxx")
}

func listOnCallTeams() {
	ctx := context.Background()
	c := newClient()

	param := waroom.ListOnCallTeamsParams{
		Page:    waroom.NewOptInt(1),
		PerPage: waroom.NewOptInt(100),
	}

	teams, err := c.OnCallTeams.List(ctx, param)
	if err != nil {
		log.Fatal(err)
	}

	for _, team := range teams {
		fmt.Println(team.UUID, team.Name, team.Color)
		if onDuty, ok := team.OnDutyParticipant.Get(); ok {
			fmt.Println("  on duty:", onDuty.User.Nickname)
		}
	}
}

func createOnCallTeam(name, color, userSub string) {
	ctx := context.Background()
	c := newClient()

	body := waroom.CreateOnCallTeamReq{
		Name:  name,
		Color: color,
		Rotations: []waroom.CreateOnCallTeamRotationReq{
			{
				Timezone:         "Asia/Tokyo",
				StartedAt:        "2026-01-05T10:00:00+09:00",
				ShiftHandoffTime: "10:00",
				ShiftHandoffWday: waroom.NewOptNilInt(1),
				IntervalDays:     waroom.NewOptInt(7),
				Participants: []waroom.CreateOnCallTeamParticipantReq{
					{UserSub: userSub, Position: 0},
				},
			},
		},
	}

	team, err := c.OnCallTeams.Create(ctx, body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(team.UUID, team.Name, team.Color)
}

func updateOnCallTeam(id string) {
	ctx := context.Background()
	c := newClient()

	body := waroom.UpdateOnCallTeamReq{
		Name: waroom.NewOptString("updated-team"),
	}

	team, err := c.OnCallTeams.Update(ctx, body, waroom.UpdateOnCallTeamParams{UUID: id})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(team.UUID, team.Name, team.Color)
}

func deleteOnCallTeam(id string) {
	ctx := context.Background()
	c := newClient()

	if err := c.OnCallTeams.Delete(ctx, waroom.DeleteOnCallTeamParams{UUID: id}); err != nil {
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
