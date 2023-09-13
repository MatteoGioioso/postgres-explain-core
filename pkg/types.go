package pkg

type Node = map[string]interface{}

type JIT struct {
	Functions int `json:"Functions"`
	Options   struct {
		Inlining     bool `json:"Inlining"`
		Optimization bool `json:"Optimization"`
		Expressions  bool `json:"Expressions"`
		Deforming    bool `json:"Deforming"`
	} `json:"Options"`
	Timing struct {
		Generation   float64 `json:"Generation"`
		Inlining     float64 `json:"Inlining"`
		Optimization float64 `json:"Optimization"`
		Emission     float64 `json:"Emission"`
		Total        float64 `json:"Total"`
	} `json:"Timing"`
}

// StatsFromPlan Statistic can be found in different forms
type StatsFromPlan struct {
	Plan struct {
		ExecutionTime float64 `json:"Execution Time"`
		PlanningTime  float64 `json:"Planning Time"`
		JIT           *JIT    `json:"JIT,omitempty"`
	} `json:"plan"`
	ExecutionTime float64 `json:"Execution Time"`
	PlanningTime  float64 `json:"Planning Time"`
	JIT           *JIT    `json:"JIT,omitempty"`
	Triggers      []struct {
		Name  string  `json:"Trigger Name"`
		Time  float64 `json:"Time"`
		Calls string  `json:"Calls"`
	} `json:"Triggers,omitempty"`
}

type Stats struct {
	ExecutionTime    float64 `json:"execution_time"`
	PlanningTime     float64 `json:"planning_time"`
	MaxRows          float64 `json:"max_rows"`
	MaxDuration      float64 `json:"max_duration"`
	MaxCost          float64 `json:"max_cost"`
	MaxBlocksRead    float64 `json:"max_blocks_read"`
	MaxBlocksWritten float64 `json:"max_blocks_written"`
	MaxBlocksHit     float64 `json:"max_blocks_hit"`
}

type Plans []struct {
	Plan map[string]interface{} `json:"Plan"`
}

type IndexesStats struct {
	Indexes []IndexStats `json:"stats"`
}

type TablesStats struct {
	Tables []TableStats `json:"stats"`
}

type NodesStats struct {
	Nodes []NodeStats `json:"stats"`
}

type Explained struct {
	Summary       []PlanRow    `json:"summary"`
	Stats         Stats        `json:"stats"`
	IndexesStats  IndexesStats `json:"indexes_stats"`
	TablesStats   TablesStats  `json:"tables_stats"`
	NodesStats    NodesStats   `json:"nodes_stats"`
	JITStats      *JIT         `json:"jit_stats"`
	TriggersStats *Triggers    `json:"triggers_stats"`
}

type NodeScopes struct {
	Table     string `json:"table"`
	Filters   string `json:"filters"`
	Index     string `json:"index"`
	Key       string `json:"key"`
	Condition string `json:"condition"`
}

type Costs struct {
	StartupCost float64 `json:"startup_cost"`
	TotalCost   float64 `json:"total_cost"`
	PlanWidth   float64 `json:"plan_width"`
}

type Rows struct {
	Total               float64 `json:"total"`
	TotalAvg            float64 `json:"total_avg"`
	PlannedRows         float64 `json:"planned_rows"`
	Removed             float64 `json:"removed"`
	EstimationFactor    float64 `json:"estimation_factor"`
	EstimationDirection string  `json:"estimation_direction"`
}

type Buffers struct {
	Reads   float64 `json:"reads"`
	Written float64 `json:"written"`
	Hits    float64 `json:"hits"`
	Dirtied float64 `json:"dirtied"`

	LocalReads   float64 `json:"local_reads"`
	LocalWritten float64 `json:"local_written"`
	LocalHits    float64 `json:"local_hits"`
	LocalDirtied float64 `json:"local_dirtied"`

	TempReads   float64 `json:"temp_reads"`
	TempWritten float64 `json:"temp_written"`

	ExclusiveReads   float64 `json:"exclusive_reads"`
	ExclusiveWritten float64 `json:"exclusive_written"`
	ExclusiveHits    float64 `json:"exclusive_hits"`
	ExclusiveDirtied float64 `json:"exclusive_dirtied"`

	ExclusiveTempReads   float64 `json:"exclusive_temp_reads"`
	ExclusiveTempWritten float64 `json:"exclusive_temp_written"`

	ExclusiveLocalReads   float64 `json:"exclusive_local_reads"`
	ExclusiveLocalWritten float64 `json:"exclusive_local_written"`
	ExclusiveLocalHits    float64 `json:"exclusive_local_hits"`
	ExclusiveLocalDirtied float64 `json:"exclusive_local_dirtied"`

	EffectiveBlocksRead    float64 `json:"effective_blocks_read"`
	EffectiveBlocksWritten float64 `json:"effective_blocks_written"`
	EffectiveBlocksHits    float64 `json:"effective_blocks_hits"`
}

type Worker struct {
	Number float64 `json:"number"`
	Loops  float64 `json:"loops"`
	Rows   float64 `json:"rows"`
	Time   float64 `json:"time"`
}

type Workers struct {
	Launched float64      `json:"launched"`
	Planned  float64      `json:"planned"`
	List     [][]Property `json:"list"`
}

type Timings struct {
	Inclusive     float64 `json:"inclusive"`
	Exclusive     float64 `json:"exclusive"`
	ExecutionTime float64 `json:"execution_time"`
}

type PlanRow struct {
	NodeId                     string     `json:"node_id"`
	NodeShortId                float64    `json:"node_short_id"`
	NodeFingerprint            string     `json:"node_fingerprint"`
	NodeParentId               string     `json:"node_parent_id"`
	Operation                  string     `json:"operation"`
	Level                      int        `json:"level"`
	Scopes                     NodeScopes `json:"scopes"`
	Inclusive                  float64    `json:"inclusive"`
	Timings                    Timings    `json:"timings"`
	Loops                      float64    `json:"loops"`
	Rows                       Rows       `json:"rows"`
	Costs                      Costs      `json:"costs"`
	Exclusive                  float64    `json:"exclusive"`
	ExecutionTime              float64    `json:"execution_time"`
	Buffers                    Buffers    `json:"buffers"`
	SubPlanOf                  string     `json:"sub_plan_of"`
	CteSubPlanOf               string     `json:"cte_sub_plan_of"`
	ParentPlanId               string     `json:"parent_plan_id"`
	DoesContainBuffers         bool       `json:"does_contain_buffers"`
	Workers                    Workers    `json:"workers"`
	NodeTypeSpecificProperties []Property `json:"node_type_specific_properties"`
}

type Operation struct {
	RelationName string `json:"relation_name"`
	Index        string `json:"index"`
	Filter       string `json:"filter"`
	Key          string `json:"key"`
	Condition    string `json:"condition"`

	getSpecificProperties func(node Node) []Property
	getWorkers            func(node Node) [][]Property
}

type Scope struct {
	Name    string `json:"name"`
	Prepend string `json:"prepend"`
}

type IndexNode struct {
	Id            string  `json:"id"`
	Type          string  `json:"type"`
	ExclusiveTime float64 `json:"exclusive_time"`
	Condition     string  `json:"condition"`
}

type TableNode struct {
	Id            string  `json:"id"`
	Type          string  `json:"type"`
	ExclusiveTime float64 `json:"exclusive_time"`
}

type IndexStats struct {
	Nodes      []IndexNode `json:"nodes"`
	TotalTime  float64     `json:"total_time"`
	Percentage float64     `json:"percentage"`
	Name       string      `json:"name"`
}

type TableStats struct {
	Nodes      []TableNode `json:"nodes"`
	TotalTime  float64     `json:"total_time"`
	Percentage float64     `json:"percentage"`
	Name       string      `json:"name"`
}

type NodeStats struct {
	Nodes      []NodeNode `json:"nodes"`
	TotalTime  float64    `json:"total_time"`
	Percentage float64    `json:"percentage"`
	Name       string     `json:"name"`
}

type NodeNode struct {
	Id            string  `json:"id"`
	Type          string  `json:"type"`
	ExclusiveTime float64 `json:"exclusive_time"`
}

type Property struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	ValueFloat  float64 `json:"float"`
	ValueString string  `json:"string"`
	Skip        bool    `json:"skip"`
	Kind        Kind    `json:"kind"`
}

type Kind = string

type ExplainedComparison struct {
	Explained
	Query string `json:"query"`
}

type ComparisonGeneralStats struct {
	ExecutionTime    PropComparison `json:"execution_time"`
	PlanningTime     PropComparison `json:"planning_time"`
	MaxRows          PropComparison `json:"max_rows"`
	MaxDuration      PropComparison `json:"max_duration"`
	MaxCost          PropComparison `json:"max_cost"`
	MaxBlocksRead    PropComparison `json:"max_blocks_read"`
	MaxBlocksWritten PropComparison `json:"max_blocks_written"`
	MaxBlocksHit     PropComparison `json:"max_blocks_hit"`
}

type NodeComparison struct {
	NodeId          string               `json:"node_id"`
	NodeIdToCompare string               `json:"node_id_to_compare"`
	Operation       string               `json:"operation"`
	Level           int                  `json:"level"`
	Warnings        []string             `json:"warnings"`
	Infos           []string             `json:"infos"`
	Scopes          NodeScopesComparison `json:"scopes"`
	Inclusive       PropComparison       `json:"inclusive_time"`
	Loops           PropComparison       `json:"loops"`
	Rows            RowsComparison       `json:"rows"`
	Costs           CostsComparison      `json:"costs"`
	Exclusive       PropComparison       `json:"exclusive_time"`
	ExecutionTime   PropComparison       `json:"execution_time"`
	Buffers         BuffersComparison    `json:"buffers"`
}

type NodeScopesComparison struct {
	Filters   PropStringComparison `json:"filters"`
	Index     PropStringComparison `json:"index"`
	Key       PropStringComparison `json:"key"`
	Condition PropStringComparison `json:"condition"`
}

type RowsComparison struct {
	Total            PropComparison `json:"total_rows"`
	PlannedRows      PropComparison `json:"planned_rows"`
	Removed          PropComparison `json:"removed_rows"`
	EstimationFactor PropComparison `json:"rows_estimation_factor"`
}

type BuffersComparison struct {
	EffectiveBlocksRead    PropComparison `json:"effective_blocks_read"`
	EffectiveBlocksWritten PropComparison `json:"effective_blocks_written"`
	EffectiveBlocksHits    PropComparison `json:"effective_blocks_hits"`
}

type CostsComparison struct {
	StartupCost PropComparison `json:"startup_cost"`
	TotalCost   PropComparison `json:"total_cost"`
	PlanWidth   PropComparison `json:"plan_width"`
}

type Comparison struct {
	GeneralStats ComparisonGeneralStats `json:"general_stats"`
}

type PropComparison struct {
	Original           float64 `json:"current"`
	ToCompare          float64 `json:"to_compare"`
	HasImproved        bool    `json:"has_improved"`
	PercentageImproved float64 `json:"percentage_improved"`
}

type PropStringComparison struct {
	Original  string `json:"original"`
	ToCompare string `json:"to_compare"`
	AreSame   bool   `json:"are_same"`
}

type Trigger struct {
	Name    string  `json:"name"`
	Time    float64 `json:"time"`
	Calls   float64 `json:"calls"`
	AvgTime float64 `json:"avg_time"`
}

type Triggers struct {
	MaxTime float64   `json:"max_time"`
	Items   []Trigger `json:"items"`
}

type ExplainedError struct {
	Error   string `json:"error"`
	Details string `json:"error_details"`
	Stack   string `json:"error_stack"`
}

type ExplainedResponse struct {
	Error     string `json:"error"`
	Explained string `json:"explained"`
}

type ComparisonResponse struct {
	Error      string `json:"error"`
	Comparison string `json:"comparison"`
}
