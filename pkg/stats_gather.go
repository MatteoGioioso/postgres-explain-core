package pkg

import (
	"encoding/json"
	"fmt"
	"sort"
)

type StatsGather struct {
	Stats
	indexesStats map[string]IndexStats
}

func NewStatsGather() *StatsGather {
	return &StatsGather{
		indexesStats: make(map[string]IndexStats),
	}
}

func (s *StatsGather) GetStatsFromPlans(plans string) error {
	type Plans []StatsFromPlan

	p := Plans{}
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

func (s *StatsGather) computeIndexesStats(node Node) {
	if node[INDEX_NAME] != nil {
		indexName := node[INDEX_NAME].(string)

		indexes := s.indexesStats[indexName]
		indexNode := IndexNode{
			Id:            node[NODE_ID].(string),
			Type:          node[NODE_TYPE].(string),
			ExclusiveTime: node[EXCLUSIVE_DURATION].(float64),
		}

		if node[INDEX_CONDITION] != nil {
			indexNode.Condition = node[INDEX_CONDITION].(string)
		}

		indexes.Nodes = append(indexes.Nodes, indexNode)
		indexes.TotalTime += node[EXCLUSIVE_DURATION].(float64)

		s.indexesStats[indexName] = indexes
	}

	if node[PLANS_PROP] != nil {
		for _, subNode := range node[PLANS_PROP].([]interface{}) {
			s.ComputeIndexesStats(subNode.(Node))
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

	if key == TOTAL_COST_PROP && s.MaxCost < valueFloat {
		s.MaxCost = valueFloat
	}

	if key == EXCLUSIVE_DURATION && s.MaxDuration < valueFloat {
		s.MaxDuration = valueFloat
	}
}
