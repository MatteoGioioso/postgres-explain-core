package pkg

import (
	"fmt"
)

type PlanRow struct {
	Level         int     `json:"level"`
	Nodes         string  `json:"nodes"`
	Inclusive     float64 `json:"inclusive"`
	Loops         float64 `json:"loops"`
	Rows          float64 `json:"rows"`
	Exclusive     float64 `json:"exclusive"`
	Rows_x        float64 `json:"rows_x"`
	ExecutionTime float64 `json:"execution_time"`
}

type Summary struct {
	planTable []PlanRow
}

func NewSummary() *Summary {
	return &Summary{
		planTable: make([]PlanRow, 0),
	}
}

func (s *Summary) Do(node Node, stats Stats) []PlanRow {
	s.recurseNode(node, stats, 0)

	return s.planTable
}

func (s *Summary) recurseNode(node Node, stats Stats, level int) {
	s.planTable = append(s.planTable, PlanRow{
		Level:         level,
		Nodes:         s.formatNode(node, level),
		Inclusive:     node[PropsExported.ACTUAL_TOTAL_TIME_PROP].(float64),
		Loops:         node[PropsExported.ACTUAL_LOOPS_PROP].(float64),
		Rows:          node[PropsExported.ACTUAL_ROWS_PROP].(float64),
		Exclusive:     node[PropsExported.ACTUAL_DURATION_PROP].(float64),
		Rows_x:        node[PropsExported.PLANNER_ESTIMATE_FACTOR].(float64),
		ExecutionTime: stats.ExecutionTime,
	})

	if node[PropsExported.PLANS_PROP] != nil {
		for _, subNode := range node[PropsExported.PLANS_PROP].([]interface{}) {
			s.recurseNode(subNode.(Node), stats, level+1)
		}
	}
}

func (s *Summary) formatNode(node Node, level int) string {
	return fmt.Sprintf(`${"-".repeat(level)}> ${node[props.NODE_TYPE_PROP]} 
      (cost=${node[props.STARTUP_COST]}...${node[props.TOTAL_COST_PROP]} rows=${node[props.PLAN_ROWS_PROP]} width=${node[props.PLAN_WIDTH]}) 
      (actual time=cost=${node[props.ACTUAL_STARTUP_TIME_PROP]}...${node[props.ACTUAL_TOTAL_TIME_PROP]} rows=${node[props.ACTUAL_ROWS_PROP]} loops=${node[props.ACTUAL_LOOPS_PROP]})
      Buffers shared hits: ${node['Shared Hit Blocks']}, read: ${node['Shared Read Blocks']}`)
}
