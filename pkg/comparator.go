package pkg

type Comparator struct {
	planPrev      Explained
	planOptimized Explained
}

func NewComparator(planPrev Explained, planOptimized Explained) *Comparator {
	return &Comparator{planPrev: planPrev, planOptimized: planOptimized}
}

func (c *Comparator) Compare() (Comparison, error) {
	return Comparison{
		GeneralStats: c.compareGeneralStats(),
	}, nil
}

func (c *Comparator) compareGeneralStats() ComparisonGeneralStats {
	return ComparisonGeneralStats{
		ExecutionTime:    c.getPropComparison(c.planPrev.Stats.ExecutionTime, c.planOptimized.Stats.ExecutionTime, false),
		PlanningTime:     c.getPropComparison(c.planPrev.Stats.PlanningTime, c.planOptimized.Stats.PlanningTime, false),
		MaxDuration:      c.getPropComparison(c.planPrev.Stats.MaxDuration, c.planOptimized.Stats.MaxDuration, false),
		MaxCost:          c.getPropComparison(c.planPrev.Stats.MaxCost, c.planOptimized.Stats.MaxCost, true),
		MaxBlocksRead:    c.getPropComparison(c.planPrev.Stats.MaxBlocksRead, c.planOptimized.Stats.MaxBlocksRead, false),
		MaxBlocksWritten: c.getPropComparison(c.planPrev.Stats.MaxBlocksWritten, c.planOptimized.Stats.MaxBlocksWritten, false),
		MaxBlocksHit:     c.getPropComparison(c.planPrev.Stats.MaxBlocksHit, c.planOptimized.Stats.MaxBlocksHit, true),
	}
}

func (c *Comparator) getPropComparison(prev, optimized float64, isInverted bool) PropComparison {
	comparison := PropComparison{
		Previous:  prev,
		Optimized: optimized,
	}
	if isInverted {
		comparison.HasImproved = prev > optimized
	} else {
		comparison.HasImproved = prev < optimized
	}

	return comparison
}
