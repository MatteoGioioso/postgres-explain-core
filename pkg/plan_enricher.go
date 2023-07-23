package pkg

import (
	"github.com/google/uuid"
	"math"
	"reflect"
	"strings"
)

type PlanEnricher struct {
	ctes            map[string]Node
	containsBuffers bool
}

func NewPlanEnricher() *PlanEnricher {
	return &PlanEnricher{
		ctes:            map[string]Node{},
		containsBuffers: false,
	}
}

func (ps *PlanEnricher) AnalyzePlan(rootNode Node) {
	ps.processNode(rootNode)
	rootNode[CTES] = ps.ctes
}

func (ps *PlanEnricher) checkBuffers(node Node) {
	if node[SHARED_HIT_BLOCKS] != nil {
		ps.containsBuffers = true
		node[DOES_CONTAIN_BUFFERS] = true
		return
	}

	ps.containsBuffers = false
	node[DOES_CONTAIN_BUFFERS] = false
}

func (ps *PlanEnricher) processNode(node Node) {
	node[NODE_ID] = uuid.New().String()

	ps.checkBuffers(node)
	ps.calculatePlannerEstimate(node)

	// Iterate over all the node's properties: "Startup Cost", "Planning Time", "Plans", ect...
	for name, value := range node {
		// If the key is "Plans", then iterated over all the sub nodes
		if name == PLANS_PROP {
			for _, child := range value.([]interface{}) {
				childNode := child.(Node)

				// Add workers planned info to parallel nodes (ie. Gather children)
				if !IsCTE(childNode) && childNode[PARENT_RELATIONSHIP] != "InitPlan" && childNode[PARENT_RELATIONSHIP] != "SubPlan" {
					if node[WORKERS_PLANNED] != nil {
						childNode[WORKERS_PLANNED_BY_GATHER] = node[WORKERS_PLANNED]
					} else {
						childNode[WORKERS_PLANNED_BY_GATHER] = node[WORKERS_PLANNED_BY_GATHER]
					}

					if node[WORKERS_LAUNCHED] != nil {
						childNode[WORKERS_LAUNCHED] = node[WORKERS_LAUNCHED]
					}
				}

				// Plans belonging to CTEs are not found as direct child of CTEs nodes,
				// { "Node Type": "CTE Scan" }
				// Instead they just appears as child nodes of root, thus they have to be
				// grouped and put back in the root node
				if IsCTE(childNode) {
					subPlanName := strings.ReplaceAll(childNode[SUBPLAN_NAME].(string), "CTE ", "")
					childNode[IS_CTE_ROOT] = "true"
					childNode[CTE_SUBPLAN_OF] = subPlanName
					ps.ctes[subPlanName] = childNode
				}
				if node[CTE_SUBPLAN_OF] != nil {
					childNode[CTE_SUBPLAN_OF] = node[CTE_SUBPLAN_OF]
				}

				ps.processNode(childNode)
			}
		}
	}

	ps.calculateActuals(node)
	ps.calculateExclusive(node)
}

func (ps *PlanEnricher) calculatePlannerEstimate(node Node) {
	if node[ACTUAL_ROWS] != nil && node[PLAN_ROWS] != nil {
		node[PLANNER_ESTIMATE_DIRECTION] = EstimateDirectionNone
		node[PLANNER_ESTIMATE_FACTOR] = node[PLAN_ROWS].(float64) / node[ACTUAL_ROWS].(float64)

		if node[ACTUAL_ROWS].(float64) < node[PLAN_ROWS].(float64) {
			node[PLANNER_ESTIMATE_DIRECTION] = EstimateDirectionOver
		}

		if node[ACTUAL_ROWS].(float64) > node[PLAN_ROWS].(float64) {
			node[PLANNER_ESTIMATE_DIRECTION] = EstimateDirectionUnder
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
	if node[ACTUAL_TOTAL_TIME] != nil {
		// since time is reported for an individual loop, actual duration must be adjusted by number of loops
		// number of workers is also taken into account
		workers := ps.getWorkers(node)
		node[ACTUAL_TOTAL_TIME] = (node[ACTUAL_TOTAL_TIME].(float64) * node[ACTUAL_LOOPS].(float64)) / workers
		if node[ACTUAL_STARTUP_TIME] != nil {
			node[ACTUAL_STARTUP_TIME] = (node[ACTUAL_STARTUP_TIME].(float64) * node[ACTUAL_LOOPS].(float64)) / workers
		} else {
			node[ACTUAL_STARTUP_TIME] = 0.0
		}
		node[EXCLUSIVE_DURATION] = node[ACTUAL_TOTAL_TIME]

		duration := (node[EXCLUSIVE_DURATION].(float64)) - ps.childrenDuration(node, 0)
		if duration > 0 {
			node[EXCLUSIVE_DURATION] = duration
		} else {
			node[EXCLUSIVE_DURATION] = 0.0
		}
	}
	node[ACTUAL_DURATION] = node[ACTUAL_TOTAL_TIME]
	node[ACTUAL_COST_PROP] = node[TOTAL_COST_PROP]

	if node[FILTER] == nil {
		node[FILTER] = ""
	}

	for _, name := range []string{
		ACTUAL_ROWS,
		PLAN_ROWS,
		ROWS_REMOVED_BY_FILTER,
		ROWS_REMOVED_BY_JOIN_FILTER,
	} {
		if node[name] != nil {
			loops := 1.0
			if node[ACTUAL_LOOPS] != nil {
				loops = node[ACTUAL_LOOPS].(float64)
			}

			if ps.getWorkers(node) > 1 {
				node[name+REVISED] = node[name].(float64)
			} else {
				// TODO it could be that the parser has a bug in which it will print a string
				if reflect.TypeOf(node[name]).Name() == "string" {
					node[name+REVISED] = 0.0
					continue
				}
				node[name+REVISED] = node[name].(float64) * loops
			}
		} else {
			node[name+REVISED] = 0.0
		}
	}
}

func (ps *PlanEnricher) getWorkers(node Node) float64 {
	workers := 1.0
	if node[WORKERS_PLANNED_BY_GATHER] != nil {
		workers = node[WORKERS_PLANNED_BY_GATHER].(float64) + 1.0
	}
	return workers
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
		if node[property] == nil {
			node[property] = 0.0
			continue
		}

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

		node[exclusiveProp] = node[property].(float64) - sum
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
