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
	TOTAL_RUNTIME            = "Total Runtime"

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

	PARENT_RELATIONSHIP = "Parent Relationship"
	SUBPLAN_NAME        = "Subplan Name"

	CTE_SCAN_PROP = "CTE Scan"
	CTE_NAME_PROP = "CTE Name"

	ARRAY_INDEX_KEY = "arrayIndex"

	RELATION_NAME         = "Relation Name"
	SCHEMA                = "Schema"
	ALIAS                 = "Alias"
	GROUP_KEY             = "Group Key"
	SORT_KEY              = "Sort Key"
	SORT_METHOD           = "Sort Method"
	SORT_SPACE_TYPE       = "Sort Space Type"
	SORT_SPACE_USED       = "Sort Space Used"
	JOIN_TYPE             = "Join Type"
	INDEX_NAME            = "Index Name"
	HASH_CONDITION        = "Hash Cond"
	PARALLEL_AWARE        = "Parallel Aware"
	WORKERS               = "Workers"
	WORKERS_PLANNED       = "Workers Planned"
	WORKERS_LAUNCHED      = "Workers Launched"
	SHARED_HIT_BLOCKS     = "Shared Hit Blocks"
	SHARED_READ_BLOCKS    = "Shared Read Blocks"
	SHARED_DIRTIED_BLOCKS = "Shared Dirtied Blocks"
	SHARED_WRITTEN_BLOCKS = "Shared Written Blocks"
	TEMP_READ_BLOCKS      = "Temp Read Blocks"
	TEMP_WRITTEN_BLOCKS   = "Temp Written Blocks"
	LOCAL_HIT_BLOCKS      = "Local Hit Blocks"
	LOCAL_READ_BLOCKS     = "Local Read Blocks"
	LOCAL_DIRTIED_BLOCKS  = "Local Dirtied Blocks"
	LOCAL_WRITTEN_BLOCKS  = "Local Written Blocks"
	IO_READ_TIME          = "I/O Read Time"
	IO_WRITE_TIME         = "I/O Write Time"
	OUTPUT                = "Output"
	HEAP_FETCHES          = "Heap Fetches"
	WAL_RECORDS           = "WAL Records"
	WAL_BYTES             = "WAL Bytes"
	WAL_FPI               = "WAL FPI"
	FULL_SORT_GROUPS      = "Full-sort Groups"
	PRE_SORTED_GROUPS     = "Pre-sorted Groups"
	PRESORTED_KEY         = "Presorted Key"

	// computed by pev
	NODE_ID                             = "nodeId"
	EXCLUSIVE_DURATION                  = "*Duration (exclusive)"
	EXCLUSIVE_COST                      = "*Cost (exclusive)"
	ACTUAL_ROWS_REVISED                 = "*Actual Rows Revised"
	PLAN_ROWS_REVISED                   = "*Plan Rows Revised"
	ROWS_REMOVED_BY_FILTER_REVISED      = "*Rows Removed by Filter"
	ROWS_REMOVED_BY_JOIN_FILTER_REVISED = "*Rows Removed by Join Filter"

	EXCLUSIVE_SHARED_HIT_BLOCKS     = "*Shared Hit Blocks (exclusive)"
	EXCLUSIVE_SHARED_READ_BLOCKS    = "*Shared Read Blocks (exclusive)"
	EXCLUSIVE_SHARED_DIRTIED_BLOCKS = "*Shared Dirtied Blocks (exclusive)"
	EXCLUSIVE_SHARED_WRITTEN_BLOCKS = "*Shared Written Blocks (exclusive)"
	EXCLUSIVE_TEMP_READ_BLOCKS      = "*Temp Read Blocks (exclusive)"
	EXCLUSIVE_TEMP_WRITTEN_BLOCKS   = "*Temp Written Blocks (exclusive)"
	EXCLUSIVE_LOCAL_HIT_BLOCKS      = "*Local Hit Blocks (exclusive)"
	EXCLUSIVE_LOCAL_READ_BLOCKS     = "*Local Read Blocks (exclusive)"
	EXCLUSIVE_LOCAL_DIRTIED_BLOCKS  = "*Local Dirtied Blocks (exclusive)"
	EXCLUSIVE_LOCAL_WRITTEN_BLOCKS  = "*Local Written Blocks (exclusive)"

	EXCLUSIVE_IO_READ_TIME  = "*I/O Read Time (exclusive)"
	EXCLUSIVE_IO_WRITE_TIME = "*I/O Write Time (exclusive)"
	AVERAGE_IO_READ_TIME    = "*I/O Read Speed (exclusive)"
	AVERAGE_IO_WRITE_TIME   = "*I/O Write Speed (exclusive)"

	WORKERS_PLANNED_BY_GATHER = "*Workers Planned By Gather"

	CTE_SCAN      = "CTE Scan"
	CTE_NAME      = "CTE Name"
	FUNCTION_NAME = "Function Name"

	PEV_PLAN_TAG = "plan_"

	EstimateDirectionOver  = "over"
	EstimateDirectionUnder = "under"
	EstimateDirectionNone  = "none"
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
