package pkg

import (
	"encoding/json"
	"fmt"
	"strings"
)

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
