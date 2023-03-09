import React, {useMemo, useEffect} from 'react'
import ReactFlow, {Controls, Edge, MarkerType, MiniMap, Node, Panel, Position, useEdgesState, useNodesState} from 'reactflow'
import {SummaryTableProps} from './interfaces'
import {PlanRow} from './types'
import 'reactflow/dist/style.css'
import {NodeWidget} from './diagram/NodeWidget'
import {EdgeWidget} from './diagram/EdgeWidget'
import {Box} from "@cloudscape-design/components";

export const SummaryDiagram = ({summary, stats}: SummaryTableProps) => {
    const [nodes, setNodes, onNodesChange] = useNodesState([])
    const [edges, setEdges, onEdgesChange] = useEdgesState([])
    const nodeTypes = useMemo(() => ({special: NodeWidget}), [])
    const edgeTypes = useMemo(() => ({special: EdgeWidget}), [])
    const calculatePosition = (row: PlanRow, i: number) => {
        return {x: row.position.x_factor, y: row.position.y_factor}
    }

    const calculateNodes = () => {
        const initialNodes = []
        const initialEdges = []
        for (let i = 0; i < summary.length; i++) {
            const row: PlanRow = summary[i]

            const node: Node = {
                id: row.node_id,
                position: {x: row.position.x_factor * 500, y: row.position.y_factor * 250},
                data: {
                    ...row,
                },
                targetPosition: Position.Left,
                sourcePosition: Position.Right,
                type: 'special',
                draggable: true,
            }

            const edge: Edge = {
                id: `${row.node_id}-${row.node_parent_id}`,
                source: row.node_id,
                target: row.node_parent_id,
                style: {strokeWidth: Math.log((row.rows.total / stats.max_rows) * 100) * 10},
                data: {
                    rows: row.rows.total,
                },
                type: 'special',
            }

            initialNodes.push(node)
            initialEdges.push(edge)
        }

        return {
            initialNodes, initialEdges,
        }
    }

    useEffect(() => {
        const {initialNodes, initialEdges} = calculateNodes()
        setNodes(initialNodes)
        setEdges(initialEdges)
    }, [])

    return (
        <div style={{height: '600px'}}>
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
                <Panel position="bottom-center" style={{zIndex: 0}}>
                    <i
                        style={{
                            marginTop: "7px",
                            width: "500px",
                            background: "#dedcdc",
                            height: "2px",
                            float: "left",
                        }}
                    />
                    <i
                        style={{
                            width: 0,
                            height: 0,
                            borderTop: "8px solid transparent",
                            borderBottom: "8px solid transparent",
                            borderLeft: "10px solid #dedcdc",
                            float: "right",
                        }}
                    />
                    <Box variant="p" color="text-label">Flow</Box>
                </Panel>
                <MiniMap/>
            </ReactFlow>
        </div>
    )
}