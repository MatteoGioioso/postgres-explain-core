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
	ExecutionTime float64 `json:"execution_time"`
	PlanningTime  float64 `json:"planning_time"`
	MaxRows       float64 `json:"max_rows"`
}

type Plans []struct {
	Plan map[string]interface{} `json:"Plan"`
}

type Explained struct {
	Summary []PlanRow `json:"summary"`
	Stats   Stats     `json:"stats"`
}

type Position struct {
	XFactor float64 `json:"x_factor"`
	YFactor float64 `json:"y_factor"`
}

type NodeSummary struct {
	Operation string `json:"operation"`
	Scope     string `json:"scope"`
	Level     int    `json:"level"`
	Costs     string `json:"costs"`
	Buffers   string `json:"buffers"`
	Relation  string `json:"relation"`
	Filters   string `json:"filters"`
	Index     string `json:"index"`
}

type Rows struct {
	Total               float64 `json:"total"`
	Removed             float64 `json:"removed"`
	Filters             string  `json:"filters"`
	EstimationFactor    float64 `json:"estimation_factor"`
	EstimationDirection string  `json:"estimation_direction"`
}

type Buffers struct {
	Reads   float64 `json:"reads"`
	Written float64 `json:"written"`
	Hits    float64 `json:"hits"`
}

type PlanRow struct {
	NodeId        string      `json:"node_id"`
	NodeParentId  string      `json:"node_parent_id"`
	Level         int         `json:"level"`
	Branch        string      `json:"branch"`
	Node          NodeSummary `json:"node"`
	Inclusive     float64     `json:"inclusive"`
	Loops         float64     `json:"loops"`
	Rows          Rows        `json:"rows"`
	Exclusive     float64     `json:"exclusive"`
	ExecutionTime float64     `json:"execution_time"`
	Buffers       Buffers     `json:"buffers"`
	SubPlanOf     string      `json:"sub_plan_of"`
	Position      Position    `json:"position"`
}

type Operation struct {
	Scope  string `json:"scope"`
	Index  string `json:"index"`
	Filter string `json:"filter"`
}
