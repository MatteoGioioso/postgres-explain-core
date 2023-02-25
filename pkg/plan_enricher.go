package pkg

import (
	"math"
)

type PlanEnricher struct {
	maxRows     float64
	maxCost     float64
	maxDuration float64
}

func NewPlanEnricher() *PlanEnricher {
	return &PlanEnricher{
		maxRows:     0,
		maxCost:     0,
		maxDuration: 0,
	}
}

func (ps *PlanEnricher) AnalyzePlan(plan map[string]interface{}) {
	ps.processNode(plan)
	plan[MAXIMUM_ROWS_PROP] = ps.maxRows
	plan[MAXIMUM_COSTS_PROP] = ps.maxCost
	plan[MAXIMUM_DURATION_PROP] = ps.maxDuration

	ps.findOutlierNodes(plan)
}

func (ps *PlanEnricher) processNode(node Node) {
	ps.calculatePlannerEstimate(node)
	ps.calculateActuals(node)
	for key, value := range node {
		ps.calculateMaximums(node, key, value)

		if key == PLANS_PROP {
			for _, value := range value.([]interface{}) {
				ps.processNode(value.(Node))
			}
		}
	}
}

func (ps *PlanEnricher) calculateMaximums(node Node, key string, value interface{}) {
	var valueFloat float64
	switch value.(type) {
	case float64:
		valueFloat = value.(float64)
	default:
		return
	}

	if key == ACTUAL_ROWS_PROP && ps.maxRows < valueFloat {
		ps.maxRows = valueFloat
	}

	if key == ACTUAL_COST_PROP && ps.maxCost < valueFloat {
		ps.maxCost = valueFloat
	}

	if key == ACTUAL_DURATION_PROP && ps.maxDuration < valueFloat {
		ps.maxDuration = valueFloat
	}
}

func (ps *PlanEnricher) findOutlierNodes(node Node) {
	node[SLOWEST_NODE_PROP] = false
	node[LARGEST_NODE_PROP] = false
	node[COSTLIEST_NODE_PROP] = false

	if node[ACTUAL_COST_PROP] == ps.maxCost {
		node[COSTLIEST_NODE_PROP] = true
	}
	if node[ACTUAL_ROWS_PROP] == ps.maxRows {
		node[LARGEST_NODE_PROP] = true
	}
	if node[ACTUAL_DURATION_PROP] == ps.maxDuration {
		node[SLOWEST_NODE_PROP] = true
	}

	for key, value := range node {
		if key == PLANS_PROP {
			for _, subNode := range value.([]interface{}) {
				ps.findOutlierNodes(subNode.(Node))
			}
		}
	}
}

func (ps *PlanEnricher) calculatePlannerEstimate(node Node) {
	node[PLANNER_ESTIMATE_FACTOR] = node[ACTUAL_ROWS_PROP].(float64) / node[PLAN_ROWS_PROP].(float64)
	node[PLANNER_ESTIMATE_DIRECTION] = node[PLANNER_ESTIMATE_FACTOR].(float64) >= 1

	if node[PLANNER_ESTIMATE_FACTOR].(float64) < 1 {
		node[PLANNER_ESTIMATE_DIRECTION] = EstimateDirectionOver
		node[PLANNER_ESTIMATE_FACTOR] = node[PLAN_ROWS_PROP].(float64) / node[ACTUAL_ROWS_PROP].(float64)
	}

	if math.IsInf(node[PLANNER_ESTIMATE_FACTOR].(float64), 0) {
		node[PLANNER_ESTIMATE_FACTOR] = float64(0)
		return
	}
	if math.IsNaN(node[PLANNER_ESTIMATE_FACTOR].(float64)) {
		node[PLANNER_ESTIMATE_FACTOR] = float64(0)
	}
}

func (ps *PlanEnricher) calculateActuals(node Node) {
	node[ACTUAL_DURATION_PROP] = node[ACTUAL_TOTAL_TIME_PROP]
	node[ACTUAL_COST_PROP] = node[TOTAL_COST_PROP]

	if node["Plans"] == nil {
		return
	}
	plans := node["Plans"].([]interface{})

	for _, subPlan := range plans {
		sp := subPlan.(Node)
		if sp[NODE_TYPE_PROP].(string) != CTE_SCAN_PROP {
			node[ACTUAL_DURATION_PROP] = node[ACTUAL_DURATION_PROP].(float64) - sp[ACTUAL_TOTAL_TIME_PROP].(float64)
			node[ACTUAL_COST_PROP] = node[ACTUAL_COST_PROP].(float64) - sp[TOTAL_COST_PROP].(float64)
		}
	}

	if node[ACTUAL_COST_PROP].(float64) < 0 {
		node[ACTUAL_COST_PROP] = 0
	}

	node[ACTUAL_DURATION_PROP] = node[ACTUAL_DURATION_PROP].(float64) * node[ACTUAL_LOOPS_PROP].(float64)
}
