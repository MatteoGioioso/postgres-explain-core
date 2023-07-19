package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
)

type cte struct {
	id    string
	level int
}

type Summary struct {
	planTable []PlanRow
	ctes      map[string]cte
}

func NewSummary() *Summary {
	return &Summary{
		planTable: make([]PlanRow, 0),
		ctes:      map[string]cte{},
	}
}

func (s *Summary) Do(node Node, stats Stats) []PlanRow {
	s.recurseNode(node, stats, 0, "")
	s.recurseCTEsNodes(node[CTES].(map[string]Node), stats)
	return s.planTable
}

func (s *Summary) recurseNode(node Node, stats Stats, level int, parentId string) {
	id := uuid.New().String()

	row := PlanRow{
		NodeId:       id,
		NodeParentId: parentId,
		Level:        level,
		Operation:    node[NODE_TYPE_PROP].(string),
		Scopes:       s.scopes(node),
		Inclusive:    node[ACTUAL_TOTAL_TIME_PROP].(float64),
		Loops:        node[ACTUAL_LOOPS_PROP].(float64),
		Exclusive:    node[EXCLUSIVE_DURATION].(float64),
		Rows: Rows{
			Total:               node[ACTUAL_ROWS_PROP].(float64),
			PlannedRows:         node[PLAN_ROWS_PROP].(float64),
			Removed:             node[ROWS_REMOVED_BY_FILTER].(float64),
			Filters:             node[FILTER].(string),
			EstimationFactor:    node[PLANNER_ESTIMATE_FACTOR].(float64),
			EstimationDirection: node[PLANNER_ESTIMATE_DIRECTION].(string),
		},
		ExecutionTime: stats.ExecutionTime,
		Buffers: Buffers{
			EffectiveBlocksRead:    getEffectiveBlocksRead(node),
			EffectiveBlocksWritten: getEffectiveBlocksWritten(node),
		},
		Costs: Costs{
			StartupCost: node[STARTUP_COST].(float64),
			TotalCost:   node[TOTAL_COST_PROP].(float64),
			PlanWidth:   node[PLAN_WIDTH].(float64),
		},
	}

	if node[SHARED_HIT_BLOCKS] != nil {
		row.Buffers.Reads = node[SHARED_READ_BLOCKS].(float64)
		row.Buffers.Written = node[SHARED_WRITTEN_BLOCKS].(float64)
		row.Buffers.Hits = node[SHARED_HIT_BLOCKS].(float64)

		if node[EXCLUSIVE+SHARED_READ_BLOCKS] != nil {
			row.Buffers.ExclusiveReads = node[EXCLUSIVE+SHARED_READ_BLOCKS].(float64)
			row.Buffers.ExclusiveWritten = node[EXCLUSIVE+SHARED_WRITTEN_BLOCKS].(float64)
			row.Buffers.ExclusiveHits = node[EXCLUSIVE+SHARED_HIT_BLOCKS].(float64)
		}
	}

	if node[TEMP_READ_BLOCKS] != nil {
		row.Buffers.TempReads = node[TEMP_READ_BLOCKS].(float64)
		row.Buffers.TempWritten = node[TEMP_WRITTEN_BLOCKS].(float64)

		if node[EXCLUSIVE+TEMP_READ_BLOCKS] != nil {
			row.Buffers.ExclusiveTempReads = node[EXCLUSIVE+TEMP_READ_BLOCKS].(float64)
			row.Buffers.ExclusiveTempWritten = node[EXCLUSIVE+TEMP_WRITTEN_BLOCKS].(float64)
		}
	}

	if node[CTE_SUBPLAN_OF] != nil {
		row.SubPlanOf = node[CTE_SUBPLAN_OF].(string)
	}

	s.planTable = append(s.planTable, row)

	// If the node is a CTE assign it to the CTEs map, the map will later be used to get the parentId in case the node
	// is part of a CTE
	if node[NODE_TYPE_PROP].(string) == CTE_SCAN {
		s.ctes[node[CTE_NAME].(string)] = cte{
			id:    id,
			level: level,
		}
	}

	if node[PLANS_PROP] != nil {
		subNodes := node[PLANS_PROP].([]interface{})
		for _, subNode := range subNodes {
			// CTE will be recurse in a second moment in the recurseCTEsNodes method
			if subNode.(Node)[IS_CTE_ROOT] != "true" {
				s.recurseNode(subNode.(Node), stats, level+1, id)
			}
		}
	}
}

func (s *Summary) scopes(node Node) NodeScopes {
	operation := node[NODE_TYPE_PROP].(string)
	op, ok := operationsMap[operation]
	if !ok {
		op = operationsMap["Default"]
	}

	return NodeScopes{
		Table:     convertPropToString(node[op.RelationName]),
		Filters:   convertPropToString(node[op.Filter]),
		Index:     convertPropToString(node[op.Index]),
		Key:       convertPropToString(node[op.Key]),
		Method:    convertPropToString(node[op.Method]),
		Condition: convertPropToString(node[op.Condition]),
	}
}

func (s *Summary) recurseCTEsNodes(ctesNodes map[string]Node, stats Stats) {
	for cteName, node := range ctesNodes {
		cte := s.ctes[cteName]
		delete(node, IS_CTE_ROOT)
		s.recurseNode(node, stats, cte.level+1, cte.id)
	}
}

func convertPropToString(prop interface{}) string {
	switch r := prop.(type) {
	case string:
		return r
	case []interface{}: // When Sorting we can have an array of sorting keys
		marshal, err := json.MarshalIndent(r, "", "    ")
		if err != nil {
			panic(fmt.Errorf("could not marshal node operation scope into []string: %v", err))
		}
		return string(marshal)
	default:
		return ""
	}
}

func getEffectiveBlocksRead(node Node) float64 {
	if node[EXCLUSIVE+LOCAL_READ_BLOCKS] != nil {
		return node[EXCLUSIVE+LOCAL_READ_BLOCKS].(float64)
	} else if node[EXCLUSIVE+TEMP_READ_BLOCKS] != nil {
		return node[EXCLUSIVE+TEMP_READ_BLOCKS].(float64)
	} else {
		return node[EXCLUSIVE+SHARED_READ_BLOCKS].(float64)
	}
}

func getEffectiveBlocksWritten(node Node) float64 {
	if node[EXCLUSIVE+LOCAL_WRITTEN_BLOCKS] != nil {
		return node[EXCLUSIVE+LOCAL_WRITTEN_BLOCKS].(float64)
	} else if node[EXCLUSIVE+TEMP_WRITTEN_BLOCKS] != nil {
		return node[EXCLUSIVE+TEMP_WRITTEN_BLOCKS].(float64)
	} else {
		return node[EXCLUSIVE+SHARED_WRITTEN_BLOCKS].(float64)
	}
}
