package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/topotal/waroom-client-go/waroom"
)

func main() {
	listEscalationPolicies()
	// getEscalationPolicySchedule("escalation-policy-uuid")
}

func listEscalationPolicies() {
	ctx := context.Background()
	c := newClient()

	//TODO: pagination
	res, err := c.OnCallEscalationPolicies.List(ctx, waroom.ListOnCallEscalationPoliciesParams{})
	if err != nil {
		log.Fatal(err)
	}

	for _, policy := range res {
		fmt.Println(policy.UUID, policy.Name, policy.Timezone, policy.TeamCallMode)
		for _, step := range policy.Steps {
			fmt.Printf("  step %d (%s): %s\n", step.Position, step.StepType, step.Team.Name)
		}
	}
}

func createEscalationPolicy(teamUUID string) {
	ctx := context.Background()
	c := newClient()

	req := waroom.CreateOnCallEscalationPolicyReq{
		Name:     "Primary escalation",
		Timezone: waroom.NewOptString("Asia/Tokyo"),
		Steps: []waroom.CreateOnCallEscalationPolicyStep{
			{TeamUUID: teamUUID},
		},
		TimeSlots: []waroom.CreateOnCallEscalationPolicyTimeSlot{
			{DayOfWeek: 1, StartTime: "09:00", EndTime: "18:00"},
		},
	}
	req.TeamCallMode.SetTo("on_duty_only")

	policy, err := c.OnCallEscalationPolicies.Create(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(policy.UUID)
}

func getEscalationPolicySchedule(id string) {
	ctx := context.Background()
	c := newClient()

	params := waroom.GetOnCallEscalationPolicyScheduleParams{UUID: id}
	params.StartDate.SetTo(time.Now())

	schedule, err := c.OnCallEscalationPolicies.GetSchedule(ctx, params)
	if err != nil {
		log.Fatal(err)
	}

	for _, day := range schedule.Days {
		for _, entry := range day.Entries {
			nickname := "(nobody)"
			if user, ok := entry.OnDutyUser.Get(); ok {
				nickname = user.Nickname
			}
			fmt.Printf("%s %s-%s step %d: %s\n",
				day.Date.Format(time.DateOnly),
				entry.StartAt.Format("15:04"),
				entry.EndAt.Format("15:04"),
				entry.EscalationPolicyStep.Position,
				nickname,
			)
		}
	}
}

func newClient() *waroom.Client {
	ctx := context.Background()

	client, err := waroom.NewClient(ctx,
		waroom.WithBaseURL("https://api.staging.waroom.com"),
		waroom.WithUserAgent("waroom-sdk-go/0.0.1"),
	)
	if err != nil {
		log.Fatal(err)
	}

	return client
}
