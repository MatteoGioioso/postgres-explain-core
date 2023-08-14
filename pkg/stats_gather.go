package pkg

import (
	"encoding/json"
	"fmt"
	"sort"
)

type StatsGather struct {
	Stats
	indexesStats map[string]IndexStats
	tablesStats  map[string]TableStats
	nodesStats   map[string]NodeStats
	jit          *JIT
	triggers     []struct {
		Name  string  `json:"Trigger Name"`
		Time  float64 `json:"Time"`
		Calls string  `json:"Calls"`
	}
}

func NewStatsGather() *StatsGather {
	return &StatsGather{
		indexesStats: make(map[string]IndexStats),
		tablesStats:  make(map[string]TableStats),
		nodesStats:   make(map[string]NodeStats),
	}
}

func (s *StatsGather) GetStatsFromPlans(plans string) error {
	var p []StatsFromPlan
	if err := json.Unmarshal([]byte(plans), &p); err != nil {
		return fmt.Errorf("could not unmarshal plan: %v", err)
	}

	if p[0].Plan.ExecutionTime != 0 {
		s.PlanningTime = p[0].Plan.PlanningTime
		s.ExecutionTime = p[0].Plan.ExecutionTime
	} else {
		s.PlanningTime = p[0].PlanningTime
		s.ExecutionTime = p[0].ExecutionTime
	}

	s.jit = p[0].JIT
	s.triggers = p[0].Triggers

	return nil
}

func (s *StatsGather) ComputeStats(node Node) Stats {
	s.calculateMaximums(node)
	s.findOutlierNodes(node)

	return Stats{
		ExecutionTime:    s.ExecutionTime,
		PlanningTime:     s.PlanningTime,
		MaxRows:          s.MaxRows,
		MaxDuration:      s.MaxDuration,
		MaxCost:          s.MaxCost,
		MaxBlocksRead:    getMaxBlocksRead(node),
		MaxBlocksWritten: getMaxBlocksWritten(node),
		MaxBlocksHit:     getMaxBlocksHits(node),
	}
}

func (s *StatsGather) ComputeIndexesStats(node Node) IndexesStats {
	s.computeIndexesStats(node)

	// For only EXPLAIN plans 'Execution Time" is missing
	if s.ExecutionTime != 0.0 {
		for indexName, index := range s.indexesStats {
			index.Percentage = (index.TotalTime / s.ExecutionTime) * 100
			s.indexesStats[indexName] = index
		}
	}

	indexesSlice := make([]IndexStats, 0)
	for indexName, index := range s.indexesStats {
		index.Name = indexName
		indexesSlice = append(indexesSlice, index)
	}

	sort.Slice(indexesSlice, func(i, j int) bool {
		return indexesSlice[i].TotalTime > indexesSlice[j].TotalTime
	})

	return IndexesStats{
		Indexes: indexesSlice,
	}
}

func (s *StatsGather) ComputeTablesStats(node Node) TablesStats {
	s.computeTablesStats(node)

	// For only EXPLAIN plans 'Execution Time" is missing
	if s.ExecutionTime != 0.0 {
		for tableName, table := range s.tablesStats {
			table.Percentage = (table.TotalTime / s.ExecutionTime) * 100
			s.tablesStats[tableName] = table
		}
	}

	tablesSlice := make([]TableStats, 0)
	for tableName, table := range s.tablesStats {
		table.Name = tableName
		tablesSlice = append(tablesSlice, table)
	}

	sort.Slice(tablesSlice, func(i, j int) bool {
		return tablesSlice[i].TotalTime > tablesSlice[j].TotalTime
	})

	return TablesStats{
		Tables: tablesSlice,
	}
}

func (s *StatsGather) ComputeNodesStats(node Node) NodesStats {
	s.computeNodesStats(node)

	// For only EXPLAIN plans 'Execution Time" is missing
	if s.ExecutionTime != 0.0 {
		for nName, n := range s.nodesStats {
			n.Percentage = (n.TotalTime / s.ExecutionTime) * 100
			s.nodesStats[nName] = n
		}
	}

	nodesSlice := make([]NodeStats, 0)
	for nName, n := range s.nodesStats {
		n.Name = nName
		nodesSlice = append(nodesSlice, n)
	}

	sort.Slice(nodesSlice, func(i, j int) bool {
		return nodesSlice[i].TotalTime > nodesSlice[j].TotalTime
	})

	return NodesStats{
		Nodes: nodesSlice,
	}
}

func (s *StatsGather) ComputeJITStats() *JIT {
	return s.jit
}

func (s *StatsGather) ComputeTriggersStats() []Trigger {
	if s.triggers != nil {
		triggers := make([]Trigger, 0)
		for _, trigger := range s.triggers {
			calls := ConvertStringToFloat64(trigger.Calls)
			triggers = append(triggers, Trigger{
				Name:    trigger.Name,
				Time:    trigger.Time,
				Calls:   calls,
				AvgTime: trigger.Time / calls,
			})
		}

		return triggers
	}

	return nil
}

func (s *StatsGather) computeIndexesStats(node Node) {
	if node[INDEX_NAME] != nil {
		indexName := node[INDEX_NAME].(string)

		indexes := s.indexesStats[indexName]
		indexNode := IndexNode{
			Id:            node[NODE_ID].(string),
			Type:          node[NODE_TYPE].(string),
			ExclusiveTime: ConvertToFloat64(node[EXCLUSIVE_DURATION]),
		}

		if node[INDEX_CONDITION] != nil {
			indexNode.Condition = node[INDEX_CONDITION].(string)
		}

		indexes.Nodes = append(indexes.Nodes, indexNode)
		indexes.TotalTime += ConvertToFloat64(node[EXCLUSIVE_DURATION])

		s.indexesStats[indexName] = indexes
	}

	if node[PLANS_PROP] != nil {
		for _, subNode := range node[PLANS_PROP].([]interface{}) {
			s.computeIndexesStats(subNode.(Node))
		}
	}
}

func (s *StatsGather) computeTablesStats(node Node) {
	if node[RELATION_NAME] != nil {
		tableName := node[RELATION_NAME].(string)

		tables := s.tablesStats[tableName]
		tableNode := TableNode{
			Id:            node[NODE_ID].(string),
			Type:          node[NODE_TYPE].(string),
			ExclusiveTime: ConvertToFloat64(node[EXCLUSIVE_DURATION]),
		}

		tables.Nodes = append(tables.Nodes, tableNode)
		tables.TotalTime += ConvertToFloat64(node[EXCLUSIVE_DURATION])

		s.tablesStats[tableName] = tables
	}

	if node[PLANS_PROP] != nil {
		for _, subNode := range node[PLANS_PROP].([]interface{}) {
			s.computeTablesStats(subNode.(Node))
		}
	}
}

func (s *StatsGather) computeNodesStats(node Node) {
	if node[NODE_TYPE] != nil {
		nodeType := node[NODE_TYPE].(string)

		nodeStats := s.nodesStats[nodeType]
		n := NodeNode{
			Id:            node[NODE_ID].(string),
			Type:          node[NODE_TYPE].(string),
			ExclusiveTime: ConvertToFloat64(node[EXCLUSIVE_DURATION]),
		}

		nodeStats.Nodes = append(nodeStats.Nodes, n)
		nodeStats.TotalTime += ConvertToFloat64(node[EXCLUSIVE_DURATION])

		s.nodesStats[nodeType] = nodeStats
	}

	if node[PLANS_PROP] != nil {
		for _, subNode := range node[PLANS_PROP].([]interface{}) {
			s.computeNodesStats(subNode.(Node))
		}
	}
}

func (s *StatsGather) findOutlierNodes(node Node) {
	node[SLOWEST_NODE_PROP] = false
	node[LARGEST_NODE_PROP] = false
	node[COSTLIEST_NODE_PROP] = false

	if node[ACTUAL_COST_PROP] == s.MaxCost {
		node[COSTLIEST_NODE_PROP] = true
	}
	if node[ACTUAL_ROWS] == s.MaxRows {
		node[LARGEST_NODE_PROP] = true
	}
	if node[ACTUAL_DURATION] == s.MaxDuration {
		node[SLOWEST_NODE_PROP] = true
	}

	if node[PLANS_PROP] != nil {
		for _, subNode := range node[PLANS_PROP].([]interface{}) {
			s.findOutlierNodes(subNode.(Node))
		}
	}
}

func (s *StatsGather) calculateMaximums(node Node) {
	for name, value := range node {
		s.getMaximum(name, value)
		if name == PLANS_PROP {
			for _, subNode := range value.([]interface{}) {
				sn := subNode.(Node)
				s.calculateMaximums(sn)
			}
		}
	}
}

func (s *StatsGather) getMaximum(key string, value interface{}) {
	var valueFloat float64
	switch value.(type) {
	case float64:
		valueFloat = value.(float64)
	default:
		return
	}

	if key == ACTUAL_ROWS+REVISED && s.MaxRows < valueFloat {
		s.MaxRows = valueFloat
	}

	if key == TOTAL_COST && s.MaxCost < valueFloat {
		s.MaxCost = valueFloat
	}

	if key == EXCLUSIVE_DURATION && s.MaxDuration < valueFloat {
		s.MaxDuration = valueFloat
	}
}
