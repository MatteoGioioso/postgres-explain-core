package pkg

import (
	"strconv"
	"strings"
)

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

func (c *Comparator) compareNodes() {
	for _, optNode := range c.planOptimized.Summary {
		for _, prevNode := range c.planPrev.Summary {
			if JaccardSimilarity(optNode, prevNode) > 0.5 {

			} else {

			}
		}
	}
}

func JaccardSimilarity(optimized, prev PlanRow) float64 {
	groupOptimized := make([]token, 0)

	for _, s := range strings.Split(optimized.Operation, " ") {
		groupOptimized = append(groupOptimized, token{
			value:  s,
			weight: defaultWeight,
		})
	}
	groupOptimized = append(groupOptimized, token{
		value:  optimized.Scopes.Table,
		weight: tableWeight,
	})
	groupOptimized = append(groupOptimized, token{
		value:  strconv.Itoa(optimized.Level),
		weight: levelWeight,
	})

	groupPrev := make([]token, 0)
	for _, s := range strings.Split(prev.Operation, " ") {
		groupPrev = append(groupPrev, token{
			value:  s,
			weight: defaultWeight,
		})
	}
	groupPrev = append(groupPrev, token{
		value:  prev.Scopes.Table,
		weight: tableWeight,
	})
	groupPrev = append(groupPrev, token{
		value:  strconv.Itoa(prev.Level),
		weight: levelWeight,
	})

	intersection := make(map[token]bool)
	union := make(map[token]bool)

	for _, t := range groupOptimized {
		union[t] = true
	}

	for _, t := range groupPrev {
		union[t] = true
		if found, _ := optimizedContains(groupOptimized, t.value); found {
			intersection[t] = true
		}
	}

	intersectionWeight := 0.0
	unionWeight := 0.0
	for t := range union {
		unionWeight += t.weight
	}

	for t := range intersection {
		intersectionWeight += t.weight
	}

	jaccard := intersectionWeight / unionWeight
	return jaccard
}

func optimizedContains(set1 []token, token string) (bool, int) {
	for i, t := range set1 {
		if t.value == token {
			return true, i
		}
	}
	return false, -1
}

type token struct {
	value  string
	weight float64
}

const (
	defaultWeight = 1
	levelWeight   = 2
	tableWeight   = 5
)
