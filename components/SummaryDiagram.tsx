import React from 'react'
// @ts-ignore
import { SummaryTableProps } from './interfaces'
import createEngine, {
  DefaultLinkModel,
  DefaultNodeModel,
  DefaultPortModel,
  DiagramModel,
} from '@projectstorm/react-diagrams'

import { CanvasWidget } from '@projectstorm/react-canvas-core'
import { PlanRow } from './types'

export const SummaryDiagram = ({ summary, stats }: SummaryTableProps) => {
  const engine = createEngine()
  const model = new DiagramModel()

  const nodesMap: { [key: string]: DefaultPortModel } = {}

  for (let i = 0; i < summary.length; i++) {
    const row: PlanRow = summary[i]
    const node = new DefaultNodeModel({
      name: `${row.node.operation} on ${row.node.relation}`,
      color: 'rgb(0,192,255)',
    })
    node.setPosition(100 * i+1, 100)
    model.addNode(node)

    nodesMap[row.node_id] = node.addOutPort('Out')
  }

  for (let i = 0; i < summary.length; i++) {
    const row: PlanRow = summary[i]
    if (nodesMap[row.node_parent_id]) {

      const link =  nodesMap[row.node_id].link<DefaultLinkModel>(nodesMap[row.node_parent_id])
      model.addLink(link)
    }
  }

  engine.setModel(model)

  return (
    <div className="summary-diagram-container" style={{
      fontFamily: 'sans-serif',
      textAlign: 'center',
      height: 'calc(100vh - 100px)',
      width: '90%',
    }}>
      <CanvasWidget className="summary-diagram-canvas" engine={engine}/>
    </div>
  )
}