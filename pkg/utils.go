package pkg

import (
	"encoding/json"
	"fmt"
)

func GetStatsFromPlans(plans string) (StatsFromPlan, error) {
	type Plans []StatsFromPlan

	p := Plans{}
	if err := json.Unmarshal([]byte(plans), &p); err != nil {
		return StatsFromPlan{}, fmt.Errorf("could not unmarshal plan: %v", err)
	}

	return p[0], nil
}

func GetRootNodeFromPlans(plans string) (Node, error) {
	p := Plans{}
	if err := json.Unmarshal([]byte(plans), &p); err != nil {
		return nil, fmt.Errorf("could not unmarshal plan: %v", err)
	}

	return p[0].Plan, nil
}

// ConvertTextPlanToJson
// TODO: https://github.com/dalibo/pev2/blob/652bbc9041fde4f1df7df03e2500c2beaea5c3f5/src/services/plan-service.ts#L493
func ConvertTextPlanToJson(text string) (Node, error) {
	return nil, nil
}