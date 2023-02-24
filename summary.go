package postgres_explain

import "fmt"

type planRow struct {
	level         int
	nodes         string
	inclusive     float64
	loops         float64
	rows          float64
	exclusive     float64
	rows_x        float64
	executionTime float64
}

type Summary struct {
	planTable []planRow
}

func (s *Summary) RecurseNode(node Node, stats map[string]interface{}, level int) {
	s.planTable = append(s.planTable, planRow{
		level:         level,
		nodes:         s.formatNode(node, level),
		inclusive:     node[PropsExported.ACTUAL_TOTAL_TIME_PROP].(float64),
		loops:         node[PropsExported.ACTUAL_LOOPS_PROP].(float64),
		rows:          node[PropsExported.ACTUAL_ROWS_PROP].(float64),
		exclusive:     node[PropsExported.ACTUAL_DURATION_PROP].(float64),
		rows_x:        node[PropsExported.PLANNER_ESTIMATE_FACTOR].(float64),
		executionTime: stats[PropsExported.EXECUTION_TIME_PROP].(float64),
	})

	if node[PropsExported.PLANS_PROP] != nil {
		for _, subNode := range node[PropsExported.PLANS_PROP].([]interface{}) {
			s.RecurseNode(subNode.(Node), stats, level+1)
		}
	}
}

func (s *Summary) formatNode(node Node, level int) string {
	return fmt.Sprintf(`${"-".repeat(level)}> ${node[props.NODE_TYPE_PROP]} 
      (cost=${node[props.STARTUP_COST]}...${node[props.TOTAL_COST_PROP]} rows=${node[props.PLAN_ROWS_PROP]} width=${node[props.PLAN_WIDTH]}) 
      (actual time=cost=${node[props.ACTUAL_STARTUP_TIME_PROP]}...${node[props.ACTUAL_TOTAL_TIME_PROP]} rows=${node[props.ACTUAL_ROWS_PROP]} loops=${node[props.ACTUAL_LOOPS_PROP]})
      Buffers shared hits: ${node['Shared Hit Blocks']}, read: ${node['Shared Read Blocks']}`)
}
