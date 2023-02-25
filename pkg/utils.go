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
