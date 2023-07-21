package pkg

const (
	// plan property keys
	NODE_TYPE           = "Node Type"
	ACTUAL_ROWS         = "Actual Rows"
	PLAN_ROWS           = "Plan Rows"
	ACTUAL_TOTAL_TIME   = "Actual Total Time"
	ACTUAL_LOOPS        = "Actual Loops"
	TOTAL_COST_PROP     = "Total Cost"
	PLANS_PROP          = "Plans"
	STARTUP_COST        = "Startup Cost"
	PLAN_WIDTH          = "Plan Width"
	ACTUAL_STARTUP_TIME = "Actual Startup Time"
	HASH_CONDITION_PROP = "Hash Cond"
	INDEX_CONDITION     = "Index Cond"
	TOTAL_RUNTIME       = "Total Runtime"

	// computed
	COMPUTED_TAGS_PROP = "*Tags"

	COSTLIEST_NODE_PROP = "*Costliest Node (by cost)"
	LARGEST_NODE_PROP   = "*Largest Node (by rows)"
	SLOWEST_NODE_PROP   = "*Slowest Node (by duration)"

	MAXIMUM_COSTS_PROP         = "*Most Expensive Node (cost)"
	MAXIMUM_ROWS_PROP          = "*Largest Node (rows)"
	MAXIMUM_DURATION_PROP      = "*Slowest Node (time)"
	ACTUAL_DURATION            = "*Actual Duration"
	ACTUAL_COST_PROP           = "*Actual Cost"
	PLANNER_ESTIMATE_FACTOR    = "*Planner Row Estimate Factor"
	PLANNER_ESTIMATE_DIRECTION = "*Planner Row Estimate Direction"

	PARENT_RELATIONSHIP = "Parent Relationship"
	SUBPLAN_NAME        = "Subplan Name"

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
	NODE_ID                     = "nodeId"
	EXCLUSIVE_DURATION          = "*Duration (exclusive)"
	EXCLUSIVE_COST              = "*Cost (exclusive)"
	ACTUAL_ROWS_REVISED         = "*Actual Rows Revised"
	PLAN_ROWS_REVISED           = "*Plan Rows Revised"
	ROWS_REMOVED_BY_FILTER      = "Rows Removed by Filter"
	ROWS_REMOVED_BY_JOIN_FILTER = "Rows Removed by Join Filter"
	FILTER                      = "Filter"

	WORKERS_PLANNED_BY_GATHER = "*Workers Planned By Gather"

	CTE_SCAN = "CTE Scan"
	CTE_NAME = "CTE Name"

	CTES = "CTEs"

	IS_CTE_ROOT    = "*Is CTE Root"
	CTE_SUBPLAN_OF = "*CTE Subplan Of"
	FUNCTION_NAME  = "Function Name"

	PEV_PLAN_TAG = "plan_"

	EstimateDirectionOver  = "over"
	EstimateDirectionUnder = "under"
	EstimateDirectionNone  = "none"

	// Operations
	SEQUENTIAL_SCAN   = "Seq Scan"
	INDEX_SCAN        = "Index Scan"
	INDEX_ONLY_SCAN   = "Index Only Scan"
	BITMAP_INDEX_SCAN = "Bitmap Index Scan"
	BITMAP_HEAP_SCAN  = "Bitmap Heap Scan"
	HASH              = "Hash"
	HASH_JOIN         = "Hash Join"
	HASH_AGGREGATE    = "HashAggregate"
	SORT              = "Sort"
	FUNCTION_SCAN     = "Function Scan"

	GROUP_AGGREGATE = "GroupAggregate"

	// Others

	X_POSITION_FACTOR = "*X Position Factor"
	Y_POSITION_FACTOR = "*Y Position Factor"

	EXCLUSIVE            = "Exclusive "
	REVISED              = " Revised"
	DOES_CONTAIN_BUFFERS = "Does contain buffers"
)
