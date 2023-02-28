package pkg

type Node = map[string]interface{}

type StatsFromPlan struct {
	ExecutionTime float64 `json:"Execution Time"`
	PlanningTime  float64 `json:"Planning Time"`
}

type Stats struct {
	ExecutionTime float64 `json:"execution_time"`
	PlanningTime  float64 `json:"planning_time"`
}

type Plans []struct {
	Plan map[string]interface{} `json:"Plan"`
}

type Explained struct {
	Summary []PlanRow `json:"summary"`
	Stats   Stats     `json:"stats"`
}

type NodeSummary struct {
	Operation string `json:"operation"`
	Level     int    `json:"level"`
	Costs     string `json:"costs"`
	Buffers   string `json:"buffers"`
	Relation  string `json:"relation"`
}

type PlanRow struct {
	Level         int            `json:"level"`
	Node          NodeSummary    `json:"node"`
	Inclusive     float64        `json:"inclusive"`
	Loops         float64        `json:"loops"`
	Rows          float64        `json:"rows"`
	Exclusive     float64        `json:"exclusive"`
	Rows_x        EstimateFactor `json:"rows_x"`
	ExecutionTime float64        `json:"execution_time"`
	Reads         float64        `json:"reads"`
	Written       float64        `json:"written"`
}

type EstimateFactor struct {
	Value     float64 `json:"value"`
	Direction string  `json:"direction"`
}
