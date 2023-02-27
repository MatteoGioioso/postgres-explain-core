export interface NodeSummary {
  operation: string
  level: number
  costs: string
  buffers: string
  relation: string
}

export interface PlanRow {
  level: number
  node: NodeSummary
  inclusive: number
  loops: number
  rows: number
  exclusive: number
  rows_x: number
  execution_time: number
}

export interface Stats {
  planning_time: number
  execution_time: number
}

export interface SummaryTableProps {
  summary: PlanRow[]
  stats: Stats
}