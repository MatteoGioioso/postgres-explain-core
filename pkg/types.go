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
