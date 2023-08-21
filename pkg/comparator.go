package pkg

import (
	"fmt"
)

type Comparator struct {
	plan          ExplainedComparison
	planToCompare ExplainedComparison
}

func NewComparator(plan ExplainedComparison, planToCompare ExplainedComparison) *Comparator {
	return &Comparator{plan: plan, planToCompare: planToCompare}
}

func (c *Comparator) Compare() (Comparison, error) {
	return Comparison{
		GeneralStats: c.compareGeneralStats(),
	}, nil
}

func (c *Comparator) compareGeneralStats() ComparisonGeneralStats {
	return ComparisonGeneralStats{
		ExecutionTime:    getPropComparison(c.plan.Stats.ExecutionTime, c.planToCompare.Stats.ExecutionTime, false),
		PlanningTime:     getPropComparison(c.plan.Stats.PlanningTime, c.planToCompare.Stats.PlanningTime, false),
		MaxDuration:      getPropComparison(c.plan.Stats.MaxDuration, c.planToCompare.Stats.MaxDuration, false),
		MaxCost:          getPropComparison(c.plan.Stats.MaxCost, c.planToCompare.Stats.MaxCost, false),
		MaxBlocksRead:    getPropComparison(c.plan.Stats.MaxBlocksRead, c.planToCompare.Stats.MaxBlocksRead, true),
		MaxBlocksWritten: getPropComparison(c.plan.Stats.MaxBlocksWritten, c.planToCompare.Stats.MaxBlocksWritten, true),
		MaxBlocksHit:     getPropComparison(c.plan.Stats.MaxBlocksHit, c.planToCompare.Stats.MaxBlocksHit, false),
	}
}

type NodeComparator struct {
	node          PlanRow
	nodeToCompare PlanRow
}

func NewNodeComparator(node, nodeToCompare PlanRow) *NodeComparator {
	return &NodeComparator{
		node:          node,
		nodeToCompare: nodeToCompare,
	}
}

func (c *NodeComparator) Compare() (NodeComparison, error) {
	if c.node.Operation != c.nodeToCompare.Operation {
		return NodeComparison{}, fmt.Errorf("nodes have different operation name, they are not comparable")
	}

	if c.node.Scopes.Table != c.nodeToCompare.Scopes.Table {
		return NodeComparison{}, fmt.Errorf("nodes operation are acting on different table, they are not comparable")
	}

	return NodeComparison{
		Operation: c.node.Operation,
		Level:     0,
		Scopes: NodeScopesComparison{
			Filters: PropStringComparison{
				Original:  c.node.Scopes.Filters,
				ToCompare: c.nodeToCompare.Scopes.Filters,
			},
			Index: PropStringComparison{
				Original:  c.node.Scopes.Index,
				ToCompare: c.nodeToCompare.Scopes.Index,
			},
			Key: PropStringComparison{
				Original:  c.node.Scopes.Key,
				ToCompare: c.nodeToCompare.Scopes.Key,
			},
			Condition: PropStringComparison{
				Original:  c.node.Scopes.Condition,
				ToCompare: c.nodeToCompare.Scopes.Condition,
			},
		},
		Inclusive: getPropComparison(c.node.Inclusive, c.nodeToCompare.Inclusive, false),
		Loops:     getPropComparison(c.node.Loops, c.nodeToCompare.Loops, false),
		Rows: RowsComparison{
			Total:            getPropComparison(c.node.Rows.Total, c.nodeToCompare.Rows.Total, false),
			PlannedRows:      getPropComparison(c.node.Rows.PlannedRows, c.nodeToCompare.Rows.PlannedRows, false),
			Removed:          getPropComparison(c.node.Rows.Removed, c.nodeToCompare.Rows.Removed, false),
			EstimationFactor: getPropComparison(c.node.Rows.EstimationFactor, c.nodeToCompare.Rows.EstimationFactor, false),
		},
		Costs: CostsComparison{
			StartupCost: getPropComparison(c.node.Costs.StartupCost, c.nodeToCompare.Costs.StartupCost, false),
			TotalCost:   getPropComparison(c.node.Costs.TotalCost, c.nodeToCompare.Costs.TotalCost, false),
			PlanWidth:   getPropComparison(c.node.Costs.PlanWidth, c.nodeToCompare.Costs.PlanWidth, false),
		},
		Exclusive:     getPropComparison(c.node.Exclusive, c.nodeToCompare.Exclusive, false),
		ExecutionTime: getPropComparison(c.node.ExecutionTime, c.nodeToCompare.ExecutionTime, false),
		Buffers: BuffersComparison{
			EffectiveBlocksRead:    getPropComparison(c.node.Buffers.EffectiveBlocksRead, c.nodeToCompare.Buffers.EffectiveBlocksRead, false),
			EffectiveBlocksWritten: getPropComparison(c.node.Buffers.EffectiveBlocksWritten, c.nodeToCompare.Buffers.EffectiveBlocksWritten, false),
			EffectiveBlocksHits:    getPropComparison(c.node.Buffers.EffectiveBlocksHits, c.nodeToCompare.Buffers.EffectiveBlocksHits, false),
		},
	}, nil
}

func getPropComparison(current, toCompare float64, isInverted bool) PropComparison {
	comparison := PropComparison{
		Original:  current,
		ToCompare: toCompare,
	}
	if isInverted {
		comparison.HasImproved = current < toCompare
	} else {
		comparison.HasImproved = current > toCompare
	}

	if current != 0.0 || toCompare != 0.0 {
		comparison.PercentageImproved = ((toCompare - current) / ((current + toCompare) / 2)) * 100
	}

	return comparison
}
