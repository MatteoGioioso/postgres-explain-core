package pkg

import (
	"fmt"
	"strings"
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
		NodeShortId:  0,
		NodeParentId: parentId,
		Level:        level,
		Operation:    s.getFullOperationName(node),
		Scopes:       s.scopes(node),
		Loops:        ConvertToFloat64(node[ACTUAL_LOOPS]),
		Inclusive:    ConvertToFloat64(node[ACTUAL_TOTAL_TIME]),
		Exclusive:    ConvertToFloat64(node[EXCLUSIVE_DURATION]),
		Timings: Timings{
			Inclusive:     ConvertToFloat64(node[ACTUAL_TOTAL_TIME]),
			Exclusive:     ConvertToFloat64(node[EXCLUSIVE_DURATION]),
			ExecutionTime: stats.ExecutionTime,
		},
		Rows: Rows{
			Total:               node[ACTUAL_ROWS+REVISED].(float64),
			TotalAvg:            ConvertToFloat64(node[ACTUAL_ROWS]),
			PlannedRows:         node[PLAN_ROWS].(float64),
			Removed:             getRowsRemovedByFilter(node),
			EstimationFactor:    node[PLANNER_ESTIMATE_FACTOR].(float64),
			EstimationDirection: node[PLANNER_ESTIMATE_DIRECTION].(string),
		},
		ExecutionTime: stats.ExecutionTime,
		Costs: Costs{
			StartupCost: node[STARTUP_COST].(float64),
			TotalCost:   node[EXCLUSIVE+TOTAL_COST].(float64),
			PlanWidth:   node[PLAN_WIDTH].(float64),
		},
		Workers:                    Workers{},
		DoesContainBuffers:         node[DOES_CONTAIN_BUFFERS].(bool),
		NodeTypeSpecificProperties: make([]Property, 0),
	}

	operation, ok := operationsMap[node[NODE_TYPE].(string)]
	if !ok {
		operation = operationsMap["Default"]
	}
	if operation.getSpecificProperties != nil {
		row.NodeTypeSpecificProperties = operation.getSpecificProperties(node)
	}
	if operation.getWorkers != nil {
		row.Workers.List = operation.getWorkers(node)
	}

	if node[WORKERS_PLANNED_BY_GATHER] != nil {
		row.Workers.Planned = node[WORKERS_PLANNED_BY_GATHER].(float64)
		row.Workers.Launched = ConvertToFloat64(node[WORKERS_LAUNCHED])
	}

	if node[DOES_CONTAIN_BUFFERS].(bool) {
		row.Buffers = Buffers{}
		row.Buffers.EffectiveBlocksRead = getEffectiveBlocksRead(node)
		row.Buffers.EffectiveBlocksWritten = getEffectiveBlocksWritten(node)
		row.Buffers.EffectiveBlocksHits = getEffectiveBlocksHits(node)
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

		row.Buffers.LocalReads = node[LOCAL_READ_BLOCKS].(float64)
		row.Buffers.LocalWritten = node[LOCAL_WRITTEN_BLOCKS].(float64)
		row.Buffers.LocalHits = node[LOCAL_HIT_BLOCKS].(float64)
		row.Buffers.LocalDirtied = node[LOCAL_DIRTIED_BLOCKS].(float64)

		if node[EXCLUSIVE+LOCAL_READ_BLOCKS] != nil {
			row.Buffers.ExclusiveLocalReads = node[EXCLUSIVE+LOCAL_READ_BLOCKS].(float64)
			row.Buffers.ExclusiveLocalWritten = node[EXCLUSIVE+LOCAL_WRITTEN_BLOCKS].(float64)
			row.Buffers.ExclusiveLocalHits = node[EXCLUSIVE+LOCAL_HIT_BLOCKS].(float64)
			row.Buffers.ExclusiveLocalDirtied = node[EXCLUSIVE+LOCAL_DIRTIED_BLOCKS].(float64)
		}
	}

	if node[CTE_SUBPLAN_OF] != nil {
		row.CteSubPlanOf = node[CTE_SUBPLAN_OF].(string)
		row.ParentPlanId = s.ctes[node[CTE_SUBPLAN_OF].(string)].id
	}
	if isSubPlan(node) {
		row.SubPlanOf = node[SUBPLAN_NAME].(string)
	}

	row.NodeFingerprint = s.computeNodeFingerprint(row)

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
	builder := strings.Builder{}
	if node[PARALLEL_AWARE] != nil {
		if node[PARALLEL_AWARE].(bool) {
			builder.WriteString("Parallel ")
		}
	}

	if node[JOIN_TYPE] != nil {
		builder.WriteString(fmt.Sprintf("%v ", node[JOIN_TYPE].(string)))
	}

	builder.WriteString(node[NODE_TYPE].(string))
	return builder.String()
}

func (s *Summary) scopes(node Node) NodeScopes {
	operation := node[NODE_TYPE].(string)
	op, ok := operationsMap[operation]
	if !ok {
		op = operationsMap["Default"]
	}

	return NodeScopes{
		Table:     ConvertScopeToString(node[op.RelationName]),
		Filters:   ConvertScopeToString(node[op.Filter]),
		Index:     ConvertScopeToString(node[op.Index]),
		Key:       ConvertScopeToString(node[op.Key]),
		Condition: ConvertScopeToString(node[op.Condition]),
	}
}

func (s *Summary) recurseCTEsNodes(ctesNodes map[string]Node, stats Stats) {
	for cteName, node := range ctesNodes {
		cte := s.ctes[cteName]
		delete(node, IS_CTE_ROOT)
		s.recurseNode(node, stats, cte.level+1, cte.id)
	}
}

func (s *Summary) computeNodeFingerprint(row PlanRow) string {
	return ""
}
