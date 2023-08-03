package pkg

type Node = map[string]interface{}

// StatsFromPlan Statistic can be found in different forms
type StatsFromPlan struct {
	Plan struct {
		ExecutionTime float64 `json:"Execution Time"`
		PlanningTime  float64 `json:"Planning Time"`
	} `json:"plan"`
	ExecutionTime float64 `json:"Execution Time"`
	PlanningTime  float64 `json:"Planning Time"`
}

type Stats struct {
	ExecutionTime    float64 `json:"execution_time"`
	PlanningTime     float64 `json:"planning_time"`
	MaxRows          float64 `json:"max_rows"`
	MaxDuration      float64 `json:"max_duration"`
	MaxCost          float64 `json:"max_cost"`
	MaxBlocksRead    float64 `json:"max_blocks_read"`
	MaxBlocksWritten float64 `json:"max_blocks_written"`
}

type Plans []struct {
	Plan map[string]interface{} `json:"Plan"`
}

type IndexesStats struct {
	Indexes []IndexStats `json:"indexes"`
}

type Explained struct {
	Summary      []PlanRow    `json:"summary"`
	Stats        Stats        `json:"stats"`
	IndexesStats IndexesStats `json:"indexes_stats"`
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
	TotalPerNode        float64 `json:"total_per_node"`
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
	Launched float64  `json:"launched"`
	Planned  float64  `json:"planned"`
	List     []Worker `json:"list"`
}

type PlanRow struct {
	NodeId                     string     `json:"node_id"`
	NodeParentId               string     `json:"node_parent_id"`
	Operation                  string     `json:"operation"`
	Level                      int        `json:"level"`
	Branch                     string     `json:"branch"`
	Scopes                     NodeScopes `json:"scopes"`
	Inclusive                  float64    `json:"inclusive"`
	Loops                      float64    `json:"loops"`
	Rows                       Rows       `json:"rows"`
	Costs                      Costs      `json:"costs"`
	Exclusive                  float64    `json:"exclusive"`
	ExecutionTime              float64    `json:"execution_time"`
	Buffers                    Buffers    `json:"buffers"`
	SubPlanOf                  string     `json:"sub_plan_of"`
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

type IndexStats struct {
	Nodes      []IndexNode `json:"indexes"`
	TotalTime  float64     `json:"total_time"`
	Percentage float64     `json:"percentage"`
	Name       string      `json:"name"`
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

const (
	Timing   = Kind("timing")
	Quantity = Kind("quantity")
	DiskSize = Kind("disk_size")
	Blocks   = Kind("blocks")
)
