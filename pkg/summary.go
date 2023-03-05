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
	// CTE will be recurse in a second moment in the recurseCTEsNodes method
	if node[IS_CTE_ROOT] == "true" {
		return
	}

	id := uuid.New().String()

	row := PlanRow{
		NodeId:       id,
		NodeParentId: parentId,
		Level:        level,
		Node:         s.getNode(node, level),
		Inclusive:    node[ACTUAL_TOTAL_TIME_PROP].(float64),
		Loops:        node[ACTUAL_LOOPS_PROP].(float64),
		Exclusive:    node[EXCLUSIVE_DURATION].(float64),
		Rows: Rows{
			Total:               node[ACTUAL_ROWS_PROP].(float64),
			Removed:             node[ROWS_REMOVED_BY_FILTER].(float64),
			Filters:             node[FILTER].(string),
			EstimationFactor:    node[PLANNER_ESTIMATE_FACTOR].(float64),
			EstimationDirection: node[PLANNER_ESTIMATE_DIRECTION].(string),
		},
		ExecutionTime: stats.ExecutionTime,
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
		for _, subNode := range node[PLANS_PROP].([]interface{}) {
			s.recurseNode(subNode.(Node), stats, level+1, id)
		}
	}
}

func (s *Summary) getNode(node Node, level int) NodeSummary {
	costsVals := []interface{}{
		node[STARTUP_COST],
		node[TOTAL_COST_PROP],
		node[PLAN_ROWS_PROP],
		node[PLAN_WIDTH],
		node[ACTUAL_STARTUP_TIME_PROP],
		node[ACTUAL_TOTAL_TIME_PROP],
		node[ACTUAL_ROWS_PROP],
		node[ACTUAL_LOOPS_PROP],
	}

	var rel string
	if node[RELATION_NAME_PROP] != nil {
		rel = node[RELATION_NAME_PROP].(string)
	}

	operation := node[NODE_TYPE_PROP].(string)
	return NodeSummary{
		Operation: operation,
		Scope:     s.GetOperationScope(operation, node),
		Relation:  rel,
		Level:     level,
		Costs: fmt.Sprintf(
			"(cost=%v...%v rows=%v width=%v) (actual time=%v...%v rows=%v loops=%v)",
			costsVals...,
		),
		Buffers: fmt.Sprintf(
			"Buffers shared hits: %v, read: %v, written: %v",
			node[SHARED_HIT_BLOCKS],
			node[SHARED_READ_BLOCKS],
			node[SHARED_WRITTEN_BLOCKS],
		),
	}
}

var operationsMap = map[string]Operation{
	SEQUENTIAL_SCAN: {
		Scope: RELATION_NAME,
	},
	HASH_JOIN: {
		Scope: HASH_CONDITION_PROP,
	},
	SORT: {
		Scope: SORT_KEY,
	},
}

func (s *Summary) GetOperationScope(op string, node Node) string {
	if operation, ok := operationsMap[op]; ok {
		switch r := node[operation.Scope].(type) {
		case string:
			return r
		case []interface{}: // When Sorting we can have an array of sorting keys
			marshal, err := json.MarshalIndent(r, "", "    ")
			if err != nil {
				panic(fmt.Errorf("could not marshal node operation scope into []string: %v", err))
			}
			return string(marshal)
		}
	}

	return "-"
}

func (s *Summary) recurseCTEsNodes(ctesNodes map[string]Node, stats Stats) {
	for cteName, node := range ctesNodes {
		cte := s.ctes[cteName]
		delete(node, IS_CTE_ROOT)
		s.recurseNode(node, stats, cte.level+1, cte.id)
	}
}
