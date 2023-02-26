package pkg

// https://github.com/dalibo/pev2/blob/652bbc9041fde4f1df7df03e2500c2beaea5c3f5/src/services/plan-service.ts#L108

import (
	"math"
	"strings"
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

func (ps *PlanEnricher) AnalyzePlan(rootNode Node) {
	ps.processNode(rootNode)
	rootNode[MAXIMUM_ROWS_PROP] = ps.maxRows
	rootNode[MAXIMUM_COSTS_PROP] = ps.maxCost
	rootNode[MAXIMUM_DURATION_PROP] = ps.maxDuration

	ps.findOutlierNodes(rootNode)
}

func (ps *PlanEnricher) isCTE(node Node) bool {
	return node[PARENT_RELATIONSHIP] == "InitPlan" && strings.HasPrefix(node[SUBPLAN_NAME].(string), "CTE")
}

func (ps *PlanEnricher) processNode(node Node) {
	ps.calculatePlannerEstimate(node)

	for key, value := range node {
		ps.calculateMaximums(node, key, value)

		if key == PLANS_PROP {
			for _, value := range value.([]interface{}) {
				if !ps.isCTE(node) {

				}
				if ps.isCTE(node) {

				}

				ps.processNode(value.(Node))
			}
		}
	}

	ps.calculateActuals(node)
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
