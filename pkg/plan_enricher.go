package pkg

// https://github.com/dalibo/pev2/blob/652bbc9041fde4f1df7df03e2500c2beaea5c3f5/src/services/plan-service.ts#L108

import (
	"math"
	"strings"
)

type PlanEnricher struct {
	maxRows          float64
	maxCost          float64
	maxDuration      float64
	maxBlocksWritten float64
	maxBlocksRead    float64
	ctes             map[string]Node
}

func NewPlanEnricher() *PlanEnricher {
	return &PlanEnricher{
		maxRows:     0,
		maxCost:     0,
		maxDuration: 0,
		ctes:        map[string]Node{},
	}
}

func (ps *PlanEnricher) AnalyzePlan(rootNode Node, stats *Stats) {
	ps.processNode(rootNode)

	ps.calculateMaximums(rootNode)
	ps.findOutlierNodes(rootNode)
	stats.MaxRows = ps.maxRows
	stats.MaxCost = ps.maxCost
	stats.MaxDuration = ps.maxDuration
	stats.MaxBlocksRead = getMaxBlocksRead(rootNode)
	stats.MaxBlocksWritten = getMaxBlocksWritten(rootNode)
	rootNode[CTES] = ps.ctes
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

func (ps *PlanEnricher) processNode(node Node) {
	ps.calculatePlannerEstimate(node)

	// Iterate over all the node's properties: "Startup Cost", "Planning Time", "Plans", ect...
	for name, value := range node {
		// If the key is "Plans", then iterated over all the sub nodes
		if name == PLANS_PROP {
			for _, subNode := range value.([]interface{}) {
				sn := subNode.(Node)

				if !IsCTE(sn) && sn[PARENT_RELATIONSHIP] != "InitPlan" && sn[PARENT_RELATIONSHIP] != "SubPlan" {
					if sn[WORKERS_PLANNED] != nil {
						sn[WORKERS_PLANNED_BY_GATHER] = sn[WORKERS_PLANNED]
					} else {
						sn[WORKERS_PLANNED_BY_GATHER] = sn[WORKERS_PLANNED_BY_GATHER]
					}
				}

				// Plans belonging to CTEs are not found as direct child of CTEs nodes,
				// { "Node Type": "CTE Scan" }
				// Instead they just appears as child nodes of root, thus they have to be
				// grouped and put back in the root node
				if IsCTE(sn) {
					subPlanName := strings.ReplaceAll(sn[SUBPLAN_NAME].(string), "CTE ", "")
					sn[IS_CTE_ROOT] = "true"
					sn[CTE_SUBPLAN_OF] = subPlanName
					ps.ctes[subPlanName] = sn
				}
				if node[CTE_SUBPLAN_OF] != nil {
					sn[CTE_SUBPLAN_OF] = node[CTE_SUBPLAN_OF]
				}

				ps.processNode(sn)
			}
		}
	}

	ps.calculateActuals(node)
	ps.calculateExclusive(node)
}

func (ps *PlanEnricher) calculateMaximums(node Node) {
	for name, value := range node {
		ps.getMaximum(name, value)
		if name == PLANS_PROP {
			for _, subNode := range value.([]interface{}) {
				sn := subNode.(Node)
				ps.calculateMaximums(sn)
			}
		}
	}
}

func (ps *PlanEnricher) getMaximum(key string, value interface{}) {
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

	if key == TOTAL_COST_PROP && ps.maxCost < valueFloat {
		ps.maxCost = valueFloat
	}

	if key == EXCLUSIVE_DURATION && ps.maxDuration < valueFloat {
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
	if node[ACTUAL_ROWS_PROP] != nil && node[PLAN_ROWS_PROP] != nil {
		node[PLANNER_ESTIMATE_DIRECTION] = EstimateDirectionNone
		node[PLANNER_ESTIMATE_FACTOR] = float64(0)

		if node[ACTUAL_ROWS_PROP].(float64) < node[PLAN_ROWS_PROP].(float64) {
			node[PLANNER_ESTIMATE_DIRECTION] = EstimateDirectionOver
			node[PLANNER_ESTIMATE_FACTOR] = node[PLAN_ROWS_PROP].(float64) / node[ACTUAL_ROWS_PROP].(float64)
		}

		if node[ACTUAL_ROWS_PROP].(float64) > node[PLAN_ROWS_PROP].(float64) {
			node[PLANNER_ESTIMATE_DIRECTION] = EstimateDirectionUnder
			node[PLANNER_ESTIMATE_FACTOR] = node[ACTUAL_ROWS_PROP].(float64) / node[PLAN_ROWS_PROP].(float64)
		}
	} else {
		node[PLANNER_ESTIMATE_DIRECTION] = EstimateDirectionNone
		node[PLANNER_ESTIMATE_FACTOR] = float64(0)
	}

	// There is the possibility that the calculation of PLANNER_ESTIMATE_FACTOR will yield Inf or NaN:
	//  var zero interface{}
	//	zero = 0.0
	//	inf := 1 / zero.(float64)
	// This will be equal to +Inf
	if math.IsInf(node[PLANNER_ESTIMATE_FACTOR].(float64), 0) {
		node[PLANNER_ESTIMATE_FACTOR] = float64(0)
		return
	}
	if math.IsNaN(node[PLANNER_ESTIMATE_FACTOR].(float64)) {
		node[PLANNER_ESTIMATE_FACTOR] = float64(0)
	}
}

func (ps *PlanEnricher) calculateActuals(node Node) {
	if node[ACTUAL_TOTAL_TIME_PROP] != nil {
		// since time is reported for an individual loop, actual duration must be adjusted by number of loops
		// number of workers is also taken into account
		workers := 1.0
		if node[WORKERS_PLANNED_BY_GATHER] != nil {
			workers = node[WORKERS_PLANNED_BY_GATHER].(float64) + 1.0
		}
		node[ACTUAL_TOTAL_TIME_PROP] = (node[ACTUAL_TOTAL_TIME_PROP].(float64) * node[ACTUAL_LOOPS_PROP].(float64)) / workers
		if node[ACTUAL_STARTUP_TIME_PROP] != nil {
			node[ACTUAL_STARTUP_TIME_PROP] = (node[ACTUAL_STARTUP_TIME_PROP].(float64) * node[ACTUAL_LOOPS_PROP].(float64)) / workers
		} else {
			node[ACTUAL_STARTUP_TIME_PROP] = 0.0
		}
		node[EXCLUSIVE_DURATION] = node[ACTUAL_TOTAL_TIME_PROP]

		duration := (node[EXCLUSIVE_DURATION].(float64)) - ps.childrenDuration(node, 0)
		if duration > 0 {
			node[EXCLUSIVE_DURATION] = duration
		} else {
			node[EXCLUSIVE_DURATION] = 0.0
		}
	}
	node[ACTUAL_DURATION_PROP] = node[ACTUAL_TOTAL_TIME_PROP]
	node[ACTUAL_COST_PROP] = node[TOTAL_COST_PROP]

	if node[FILTER] == nil {
		node[FILTER] = ""
	}

	for _, name := range []string{
		ACTUAL_ROWS_PROP,
		PLAN_ROWS_PROP,
		ROWS_REMOVED_BY_FILTER,
		ROWS_REMOVED_BY_JOIN_FILTER,
	} {
		if node[name] != nil {
			loops := 1.0
			if node[ACTUAL_LOOPS_PROP] != nil {
				loops = node[ACTUAL_LOOPS_PROP].(float64)
			}

			node[name] = node[name].(float64) * loops
		} else {
			node[name] = 0.0
		}
	}
}

// Any node reports total of what it used itself, plus all that its sub-nodes used
// https://www.depesz.com/2021/06/20/explaining-the-unexplainable-part-6-buffers/
func (ps *PlanEnricher) calculateExclusive(node Node) {
	properties := []string{
		SHARED_HIT_BLOCKS,
		SHARED_READ_BLOCKS,
		SHARED_DIRTIED_BLOCKS,
		SHARED_WRITTEN_BLOCKS,
		TEMP_READ_BLOCKS,
		TEMP_WRITTEN_BLOCKS,
		LOCAL_HIT_BLOCKS,
		LOCAL_READ_BLOCKS,
		LOCAL_DIRTIED_BLOCKS,
		LOCAL_WRITTEN_BLOCKS,
		IO_READ_TIME,
		IO_WRITE_TIME,
	}

	for _, property := range properties {
		sum := 0.0

		exclusiveProp := EXCLUSIVE + property
		node[exclusiveProp] = node[property]

		if node[PLANS_PROP] == nil {
			continue
		}

		for _, subNode := range node[PLANS_PROP].([]interface{}) {
			sn := subNode.(Node)
			if sn[property] != nil {
				sum += sn[property].(float64)
			}
		}

		if node[property] != nil {
			node[exclusiveProp] = node[property].(float64) - sum
		}
	}
}

func (ps *PlanEnricher) childrenDuration(node Node, duration float64) float64 {
	if node[PLANS_PROP] == nil {
		return duration
	}

	for _, subNode := range node[PLANS_PROP].([]interface{}) {
		sn := subNode.(Node)
		if sn[PARENT_RELATIONSHIP] != "InitPlan" {
			if sn[EXCLUSIVE_DURATION] == nil {
				duration += 0.0
			} else {
				duration += sn[EXCLUSIVE_DURATION].(float64)
			}
			duration = ps.childrenDuration(sn, duration)
		}
	}

	return duration
}
