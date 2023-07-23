package pkg

import (
	"encoding/json"
	"fmt"
)

type StatsGather struct {
	Stats
	IndexesStats IndexesStats
}

func NewStatsGather() *StatsGather {
	return &StatsGather{
		IndexesStats: IndexesStats{Indexes: map[string]IndexStats{}},
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
	stats := s.computeIndexesStats(node)
	if s.ExecutionTime == 0.0 {
		return stats
	}

	for indexName, _ := range stats.Indexes {
		index := stats.Indexes[indexName]
		index.Percentage = (index.TotalTime / s.ExecutionTime) * 100
		stats.Indexes[indexName] = index
	}

	return stats
}

func (s *StatsGather) computeIndexesStats(node Node) IndexesStats {
	if node[INDEX_NAME] != nil {
		indexName := node[INDEX_NAME].(string)

		indexes := s.IndexesStats.Indexes[indexName]
		indexes.Nodes = append(indexes.Nodes, IndexNode{
			Id:            node[NODE_ID].(string),
			Type:          node[NODE_TYPE].(string),
			ExclusiveTime: node[EXCLUSIVE_DURATION].(float64),
			Condition:     node[INDEX_CONDITION].(string),
		})
		indexes.TotalTime += node[EXCLUSIVE_DURATION].(float64)

		s.IndexesStats.Indexes[indexName] = indexes
	}

	if node[PLANS_PROP] != nil {
		for _, subNode := range node[PLANS_PROP].([]interface{}) {
			s.ComputeIndexesStats(subNode.(Node))
		}
	}

	return s.IndexesStats
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
