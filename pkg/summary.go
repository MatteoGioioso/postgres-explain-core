package pkg

import (
	"encoding/json"
	"fmt"
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
	id := node[NODE_ID].(string)

	row := PlanRow{
		NodeId:       id,
		NodeParentId: parentId,
		Level:        level,
		Operation:    s.getFullOperationName(node),
		Scopes:       s.scopes(node),
		Loops:        node[ACTUAL_LOOPS].(float64),
		Inclusive:    node[ACTUAL_TOTAL_TIME].(float64),
		Exclusive:    node[EXCLUSIVE_DURATION].(float64),
		Rows: Rows{
			Total:               node[ACTUAL_ROWS+REVISED].(float64),
			TotalPerNode:        node[ACTUAL_ROWS+REVISED].(float64),
			PlannedRows:         node[PLAN_ROWS].(float64),
			Removed:             getRowsRemovedByFilter(node),
			EstimationFactor:    node[PLANNER_ESTIMATE_FACTOR].(float64),
			EstimationDirection: node[PLANNER_ESTIMATE_DIRECTION].(string),
		},
		ExecutionTime: stats.ExecutionTime,
		Costs: Costs{
			StartupCost: node[STARTUP_COST].(float64),
			TotalCost:   node[TOTAL_COST_PROP].(float64),
			PlanWidth:   node[PLAN_WIDTH].(float64),
		},
		Workers:            Workers{},
		DoesContainBuffers: node[DOES_CONTAIN_BUFFERS].(bool),
	}

	if operation, ok := operationsMap[node[NODE_TYPE].(string)]; ok {
		if operation.getSpecificProperties != nil {
			row.NodeTypeSpecificProperties = operation.getSpecificProperties(node)
		}
	}

	if node[WORKERS_PLANNED_BY_GATHER] != nil {
		row.Workers.Planned = node[WORKERS_PLANNED_BY_GATHER].(float64)
		row.Workers.Launched = node[WORKERS_LAUNCHED].(float64)
	}

	if node[WORKERS] != nil {
		for _, worker := range node[WORKERS].([]interface{}) {
			w := worker.(map[string]interface{})

			if w[ACTUAL_ROWS] != nil {
				row.Workers.List = append(row.Workers.List, Worker{
					Number: w["Worker Number"].(float64),
					Loops:  w[ACTUAL_LOOPS].(float64),
					Rows:   w[ACTUAL_ROWS].(float64),
					Time:   w[ACTUAL_TOTAL_TIME].(float64),
				})
			}
		}
	}

	if node[DOES_CONTAIN_BUFFERS].(bool) {
		row.Buffers = Buffers{}
		row.Buffers.EffectiveBlocksRead = getEffectiveBlocksRead(node)
		row.Buffers.EffectiveBlocksWritten = getEffectiveBlocksWritten(node)

		row.Buffers.Reads = node[SHARED_READ_BLOCKS].(float64)
		row.Buffers.Written = node[SHARED_WRITTEN_BLOCKS].(float64)
		row.Buffers.Hits = node[SHARED_HIT_BLOCKS].(float64)
		row.Buffers.Dirtied = node[SHARED_DIRTIED_BLOCKS].(float64)

		if node[EXCLUSIVE+SHARED_READ_BLOCKS] != nil {
			row.Buffers.ExclusiveReads = node[EXCLUSIVE+SHARED_READ_BLOCKS].(float64)
			row.Buffers.ExclusiveWritten = node[EXCLUSIVE+SHARED_WRITTEN_BLOCKS].(float64)
			row.Buffers.ExclusiveHits = node[EXCLUSIVE+SHARED_HIT_BLOCKS].(float64)
			row.Buffers.ExclusiveDirtied = node[EXCLUSIVE+SHARED_DIRTIED_BLOCKS].(float64)
		}

		row.Buffers.TempReads = node[TEMP_READ_BLOCKS].(float64)
		row.Buffers.TempWritten = node[TEMP_WRITTEN_BLOCKS].(float64)

		if node[EXCLUSIVE+TEMP_READ_BLOCKS] != nil {
			row.Buffers.ExclusiveTempReads = node[EXCLUSIVE+TEMP_READ_BLOCKS].(float64)
			row.Buffers.ExclusiveTempWritten = node[EXCLUSIVE+TEMP_WRITTEN_BLOCKS].(float64)
		}

		row.Buffers.Reads = node[LOCAL_READ_BLOCKS].(float64)
		row.Buffers.Written = node[LOCAL_WRITTEN_BLOCKS].(float64)
		row.Buffers.Hits = node[LOCAL_HIT_BLOCKS].(float64)
		row.Buffers.Dirtied = node[LOCAL_DIRTIED_BLOCKS].(float64)

		if node[EXCLUSIVE+LOCAL_READ_BLOCKS] != nil {
			row.Buffers.ExclusiveReads = node[EXCLUSIVE+LOCAL_READ_BLOCKS].(float64)
			row.Buffers.ExclusiveWritten = node[EXCLUSIVE+LOCAL_WRITTEN_BLOCKS].(float64)
			row.Buffers.ExclusiveHits = node[EXCLUSIVE+LOCAL_HIT_BLOCKS].(float64)
			row.Buffers.ExclusiveDirtied = node[EXCLUSIVE+LOCAL_DIRTIED_BLOCKS].(float64)
		}
	}

	if node[CTE_SUBPLAN_OF] != nil {
		row.SubPlanOf = node[CTE_SUBPLAN_OF].(string)
		row.ParentPlanId = s.ctes[node[CTE_SUBPLAN_OF].(string)].id
	}

	s.planTable = append(s.planTable, row)

	// If the node is a CTE assign it to the CTEs map, the map will later be used to get the parentId in case the node
	// is part of a CTE
	if node[NODE_TYPE].(string) == CTE_SCAN {
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

func (s *Summary) getFullOperationName(node Node) string {
	if node[PARALLEL_AWARE] != nil {
		return fmt.Sprintf("Parallel %s", node[NODE_TYPE].(string))
	}

	return node[NODE_TYPE].(string)
}

func (s *Summary) scopes(node Node) NodeScopes {
	operation := node[NODE_TYPE].(string)
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

func getRowsRemovedByFilter(node Node) float64 {
	op := node[NODE_TYPE].(string)
	filter, ok := filtersMap[op]
	if !ok {
		return node[ROWS_REMOVED_BY_FILTER+REVISED].(float64)
	}

	return node[filter+REVISED].(float64)
}
