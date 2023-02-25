package pkg

import (
	"encoding/json"
	"fmt"
)

func GetStatsFromPlans(plans string) (Stats, error) {
	type Plans []Stats

	p := Plans{}
	if err := json.Unmarshal([]byte(plans), &p); err != nil {
		return Stats{}, fmt.Errorf("could not unmarshal plan: %v", err)
	}

	return p[0], nil
}

func GetRootNodeFromPlans(plans string) (Node, error) {
	p := Plans{}
	if err := json.Unmarshal([]byte(plans), &p); err != nil {
		fmt.Printf("%+v\n", err)
		return nil, fmt.Errorf("could not unmarshal plan: %v", err)
	}

	return p[0].Plan, nil
}
