package pkg

import (
	"encoding/json"
	"fmt"
)

func GetStatsFromPlans(plans string) (Stats, error) {
	type Plans []StatsFromPlan

	p := Plans{}
	if err := json.Unmarshal([]byte(plans), &p); err != nil {
		return Stats{}, fmt.Errorf("could not unmarshal plan: %v", err)
	}

	s := Stats{}

	if p[0].Plan.ExecutionTime != 0 {
		s.PlanningTime = p[0].Plan.PlanningTime
		s.ExecutionTime = p[0].Plan.ExecutionTime
	} else {
		s.PlanningTime = p[0].PlanningTime
		s.ExecutionTime = p[0].ExecutionTime
	}

	return s, nil
}

func GetRootNodeFromPlans(plans string) (Node, error) {
	p := Plans{}
	if err := json.Unmarshal([]byte(plans), &p); err != nil {
		return nil, fmt.Errorf("could not unmarshal plan: %v", err)
	}

	return p[0].Plan, nil
}
