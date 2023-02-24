package pkg

import "encoding/json"

const (
	// plan property keys
	NODE_TYPE_PROP           = "Node Type"
	ACTUAL_ROWS_PROP         = "Actual Rows"
	PLAN_ROWS_PROP           = "Plan Rows"
	ACTUAL_TOTAL_TIME_PROP   = "Actual Total Time"
	ACTUAL_LOOPS_PROP        = "Actual Loops"
	TOTAL_COST_PROP          = "Total Cost"
	PLANS_PROP               = "Plans"
	STARTUP_COST             = "Startup Cost"
	PLAN_WIDTH               = "Plan Width"
	ACTUAL_STARTUP_TIME_PROP = "Actual Startup Time"
	RELATION_NAME_PROP       = "Relation Name"
	SCHEMA_PROP              = "Schema"
	ALIAS_PROP               = "Alias"
	GROUP_KEY_PROP           = "Group Key"
	SORT_KEY_PROP            = "Sort Key"
	JOIN_TYPE_PROP           = "Join Type"
	INDEX_NAME_PROP          = "Index Name"
	HASH_CONDITION_PROP      = "Hash Cond"
	EXECUTION_TIME_PROP      = "Execution Time"

	// computed
	COMPUTED_TAGS_PROP = "*Tags"

	COSTLIEST_NODE_PROP = "*Costliest Node (by cost)"
	LARGEST_NODE_PROP   = "*Largest Node (by rows)"
	SLOWEST_NODE_PROP   = "*Slowest Node (by duration)"

	MAXIMUM_COSTS_PROP         = "*Most Expensive Node (cost)"
	MAXIMUM_ROWS_PROP          = "*Largest Node (rows)"
	MAXIMUM_DURATION_PROP      = "*Slowest Node (time)"
	ACTUAL_DURATION_PROP       = "*Actual Duration"
	ACTUAL_COST_PROP           = "*Actual Cost"
	PLANNER_ESTIMATE_FACTOR    = "*Planner Row Estimate Factor"
	PLANNER_ESTIMATE_DIRECTION = "*Planner Row Estimate Direction"

	CTE_SCAN_PROP = "CTE Scan"
	CTE_NAME_PROP = "CTE Name"

	ARRAY_INDEX_KEY = "arrayIndex"

	EstimateDirectionOver  = "over"
	EstimateDirectionUnder = "under"
)

type PlanProps struct {
	NODE_TYPE_PROP           string
	ACTUAL_ROWS_PROP         string
	PLAN_ROWS_PROP           string
	ACTUAL_TOTAL_TIME_PROP   string
	ACTUAL_LOOPS_PROP        string
	TOTAL_COST_PROP          string
	PLANS_PROP               string
	STARTUP_COST             string
	PLAN_WIDTH               string
	ACTUAL_STARTUP_TIME_PROP string
	RELATION_NAME_PROP       string
	SCHEMA_PROP              string
	ALIAS_PROP               string
	GROUP_KEY_PROP           string
	SORT_KEY_PROP            string
	JOIN_TYPE_PROP           string
	INDEX_NAME_PROP          string
	HASH_CONDITION_PROP      string
	EXECUTION_TIME_PROP      string

	// computed
	COMPUTED_TAGS_PROP string

	COSTLIEST_NODE_PROP string
	LARGEST_NODE_PROP   string
	SLOWEST_NODE_PROP   string

	MAXIMUM_COSTS_PROP         string
	MAXIMUM_ROWS_PROP          string
	MAXIMUM_DURATION_PROP      string
	ACTUAL_DURATION_PROP       string
	ACTUAL_COST_PROP           string
	PLANNER_ESTIMATE_FACTOR    string
	PLANNER_ESTIMATE_DIRECTION string

	CTE_SCAN_PROP string
	CTE_NAME_PROP string

	ARRAY_INDEX_KEY string
}

func (p PlanProps) ToJSON() []byte {
	marshal, err := json.Marshal(p)
	if err != nil {
		return nil
	}

	return marshal
}

var PropsExported = PlanProps{
	NODE_TYPE_PROP:           NODE_TYPE_PROP,
	EXECUTION_TIME_PROP:      EXECUTION_TIME_PROP,
	ACTUAL_ROWS_PROP:         ACTUAL_ROWS_PROP,
	PLAN_ROWS_PROP:           PLAN_ROWS_PROP,
	ACTUAL_TOTAL_TIME_PROP:   ACTUAL_TOTAL_TIME_PROP,
	ACTUAL_LOOPS_PROP:        ACTUAL_LOOPS_PROP,
	TOTAL_COST_PROP:          TOTAL_COST_PROP,
	PLANS_PROP:               PLANS_PROP,
	STARTUP_COST:             STARTUP_COST,
	PLAN_WIDTH:               PLAN_WIDTH,
	ACTUAL_STARTUP_TIME_PROP: ACTUAL_STARTUP_TIME_PROP,
	RELATION_NAME_PROP:       RELATION_NAME_PROP,
	SCHEMA_PROP:              SCHEMA_PROP,
	ALIAS_PROP:               ALIAS_PROP,
	GROUP_KEY_PROP:           GROUP_KEY_PROP,
	SORT_KEY_PROP:            SORT_KEY_PROP,
	JOIN_TYPE_PROP:           JOIN_TYPE_PROP,
	INDEX_NAME_PROP:          INDEX_NAME_PROP,
	HASH_CONDITION_PROP:      HASH_CONDITION_PROP,

	// computed
	COMPUTED_TAGS_PROP: COMPUTED_TAGS_PROP,

	COSTLIEST_NODE_PROP: COSTLIEST_NODE_PROP,
	LARGEST_NODE_PROP:   LARGEST_NODE_PROP,
	SLOWEST_NODE_PROP:   SLOWEST_NODE_PROP,

	MAXIMUM_COSTS_PROP:         MAXIMUM_COSTS_PROP,
	MAXIMUM_ROWS_PROP:          MAXIMUM_ROWS_PROP,
	MAXIMUM_DURATION_PROP:      MAXIMUM_DURATION_PROP,
	ACTUAL_DURATION_PROP:       ACTUAL_DURATION_PROP,
	ACTUAL_COST_PROP:           ACTUAL_COST_PROP,
	PLANNER_ESTIMATE_FACTOR:    PLANNER_ESTIMATE_FACTOR,
	PLANNER_ESTIMATE_DIRECTION: PLANNER_ESTIMATE_DIRECTION,

	CTE_SCAN_PROP: CTE_SCAN_PROP,
	CTE_NAME_PROP: CTE_NAME_PROP,

	ARRAY_INDEX_KEY: ARRAY_INDEX_KEY,
}
