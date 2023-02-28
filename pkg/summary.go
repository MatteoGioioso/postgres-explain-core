package pkg

import (
	"fmt"
)

type Summary struct {
	planTable []PlanRow
}

func NewSummary() *Summary {
	return &Summary{
		planTable: make([]PlanRow, 0),
	}
}

func (s *Summary) Do(node Node, stats StatsFromPlan) []PlanRow {
	s.recurseNode(node, stats, 0)

	return s.planTable
}

func (s *Summary) recurseNode(node Node, stats StatsFromPlan, level int) {
	s.planTable = append(s.planTable, PlanRow{
		Level:     level,
		Node:      s.getNode(node, level),
		Inclusive: node[ACTUAL_TOTAL_TIME_PROP].(float64),
		Loops:     node[ACTUAL_LOOPS_PROP].(float64),
		Rows:      node[ACTUAL_ROWS_PROP].(float64),
		Exclusive: node[ACTUAL_DURATION_PROP].(float64),
		Rows_x: EstimateFactor{
			Value:     node[PLANNER_ESTIMATE_FACTOR].(float64),
			Direction: node[PLANNER_ESTIMATE_DIRECTION].(string),
		},
		ExecutionTime: stats.ExecutionTime,
	})

	if node[PLANS_PROP] != nil {
		for _, subNode := range node[PLANS_PROP].([]interface{}) {
			s.recurseNode(subNode.(Node), stats, level+1)
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

	return NodeSummary{
		Operation: node[NODE_TYPE_PROP].(string),
		Relation:  rel,
		Level:     level,
		Costs: fmt.Sprintf(
			"(cost=%v...%v rows=%v width=%v) (actual time=cost=%v...%v rows=%v loops=%v)",
			costsVals...,
		),
		Buffers: fmt.Sprintf(
			"Buffers shared hits: %v, read: %v, written: %v",
			node["Shared Hit Blocks"],
			node["Shared Read Blocks"],
			node["Shared Written Blocks"],
		),
	}
}
