package pkg

import (
	"encoding/json"
	"fmt"
	"strings"
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

func getMaxBlocksRead(rootNode Node) float64 {
	sum := 0.0
	if rootNode[SHARED_READ_BLOCKS] != nil {
		sum += rootNode[SHARED_READ_BLOCKS].(float64)
	}
	if rootNode[TEMP_READ_BLOCKS] != nil {
		sum += rootNode[TEMP_READ_BLOCKS].(float64)
	}
	if rootNode[LOCAL_READ_BLOCKS] != nil {
		sum += rootNode[LOCAL_READ_BLOCKS].(float64)
	}

	return sum
}

func getMaxBlocksWritten(rootNode Node) float64 {
	sum := 0.0
	if rootNode[SHARED_WRITTEN_BLOCKS] != nil {
		sum += rootNode[SHARED_WRITTEN_BLOCKS].(float64)
	}
	if rootNode[TEMP_WRITTEN_BLOCKS] != nil {
		sum += rootNode[TEMP_WRITTEN_BLOCKS].(float64)
	}
	if rootNode[LOCAL_WRITTEN_BLOCKS] != nil {
		sum += rootNode[LOCAL_WRITTEN_BLOCKS].(float64)
	}

	return sum
}

func IsCTE(node Node) bool {
	return node[PARENT_RELATIONSHIP] == "InitPlan" && strings.HasPrefix(node[SUBPLAN_NAME].(string), "CTE")
}
