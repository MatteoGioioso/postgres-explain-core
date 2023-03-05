import React, { useMemo, useEffect } from 'react'
import ReactFlow, { Controls, Edge, Node, Position, useEdgesState, useNodesState } from 'reactflow'
import { SummaryTableProps } from './interfaces'
import { PlanRow } from './types'
import 'reactflow/dist/style.css'
import { NodeWidget } from './diagram/NodeWidget'
import { Box } from '@cloudscape-design/components'
import { EdgeWidget } from './diagram/EdgeWidget'

export const SummaryDiagram = ({ summary, stats }: SummaryTableProps) => {
  const [nodes, setNodes, onNodesChange] = useNodesState([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState([]);
  const nodeTypes = useMemo(() => ({ special: NodeWidget }), []);
  const edgeTypes = useMemo(() => ({ special: EdgeWidget }), []);

  const calculateNodes = () => {
    const initialNodes = []
    const initialEdges = []
    for (let i = 0; i < summary.length; i++) {
      const row: PlanRow = summary[i]
      const node: Node = {
        id: row.node_id,
        position: { x: 350 * (summary.length - i), y: 350 * (row.level + 1) },
        data: {
          ...row
        },
        targetPosition: Position.Left,
        sourcePosition: Position.Right,
        type: 'special',
        draggable: true
      }

      const edge: Edge = {
        id: `${row.node_id}-${row.node_parent_id}`,
        source: row.node_id,
        target: row.node_parent_id,
        // TODO: change stroke based on the amount of time or rows
        // style: { strokeWidth: row.rows.total /100 },
        data: {
          rows: row.rows.total
        },
        type: 'special',
      }

      initialNodes.push(node)
      initialEdges.push(edge)
    }

    return {
      initialNodes, initialEdges
    }
  }

  useEffect(() => {
    const {initialNodes, initialEdges} = calculateNodes()
    setNodes(initialNodes);
    setEdges(initialEdges);
  }, []);

  return (
    <div style={{ height: '600px' }}>
      <ReactFlow
        fitView
        nodes={nodes}
        edges={edges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        nodeTypes={nodeTypes}
        edgeTypes={edgeTypes}
      >
        <Controls/>
      </ReactFlow>
    </div>
  )
}