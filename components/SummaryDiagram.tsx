import React, {useMemo, useEffect} from 'react'
import ReactFlow, {
    Controls,
    Edge,
    MarkerType,
    MiniMap,
    Node,
    Panel,
    Position,
    useEdgesState,
    useNodesState,
    useReactFlow,
    ReactFlowProvider
} from 'reactflow'
import {SummaryTableProps} from './interfaces'
import {PlanRow} from './types'
import 'reactflow/dist/style.css'
import {NodeWidget} from './diagram/NodeWidget'
import {EdgeWidget} from './diagram/EdgeWidget'
import {Box} from "@cloudscape-design/components";

const elk = new ELK();

const getLayoutedElements = (nodes, edges, options = {}) => {
    const graph = {
        id: 'root',
        layoutOptions: options,
        children: nodes.map((node) => ({
            ...node,
            targetPosition: 'top',
            sourcePosition: 'bottom',
            width: 250,
            height: 80,
        })),
        edges: edges,
    };

    return elk
        .layout(graph)
        .then((layoutedGraph) => ({
            nodes: layoutedGraph.children.map((node) => ({
                ...node,
                position: {x: node.x, y: node.y},
            })),

            edges: layoutedGraph.edges,
        }))
        .catch(console.error);
};

export const Diagram = ({summary, stats}: SummaryTableProps) => {
    const { fitView } = useReactFlow()
    const [nodes, setNodes, onNodesChange] = useNodesState([])
    const [edges, setEdges, onEdgesChange] = useEdgesState([])
    const nodeTypes = useMemo(() => ({special: NodeWidget}), [])
    const edgeTypes = useMemo(() => ({special: EdgeWidget}), [])

    // Elk has a *huge* amount of options to configure. To see everything you can
    // tweak check out:
    //
    // - https://www.eclipse.org/elk/reference/algorithms.html
    // - https://www.eclipse.org/elk/reference/options.html
    const elkOptions = {
        'elk.algorithm': 'layered',
        'elk.spacing.nodeNode': 80,
        'elk.layered.spacing.nodeNodeBetweenLayers': 80,
    };

    const calculateNodes = () => {
        const initialNodes = []
        const initialEdges = []

        for (let i = 0; i < summary.length; i++) {
            const row: PlanRow = summary[i]

            const node: Node = {
                id: row.node_id,
                data: {
                    ...row,
                },
                targetPosition: Position.Top,
                sourcePosition: Position.Bottom,
                type: 'special',
                draggable: true,
            }

            initialNodes.push(node)

            if (row.node_parent_id === "") continue

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

            initialEdges.push(edge)
        }

        return {
            initialNodes, initialEdges
        }
    }

    useEffect(() => {
        const {initialNodes, initialEdges} = calculateNodes()

        const opts = {'elk.direction': 'DOWN', ...elkOptions};

        getLayoutedElements(initialNodes, initialEdges, opts).then(({nodes: layoutedNodes, edges: layoutedEdges}) => {
            setNodes(layoutedNodes);
            setEdges(layoutedEdges);

            window.requestAnimationFrame(() => fitView());
        });
    }, [])

    return (
        <div style={{height: '1000px'}}>
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
                <MiniMap/>
            </ReactFlow>
        </div>
    )
}

export const SummaryDiagram = ({summary, stats}: SummaryTableProps) => {
    return (
        <ReactFlowProvider>
            <Diagram summary={summary} stats={stats}/>
        </ReactFlowProvider>
    )
}