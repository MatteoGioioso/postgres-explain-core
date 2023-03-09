import React, {FC} from 'react'
import {EdgeProps, getBezierPath, EdgeLabelRenderer, MarkerType} from 'reactflow'
import {Box} from '@cloudscape-design/components'
import {betterNumbers} from '../utils'

export const EdgeWidget: FC<EdgeProps> = ({
                                              id,
                                              sourceX,
                                              sourceY,
                                              targetX,
                                              targetY,
                                              sourcePosition,
                                              targetPosition,
                                              markerEnd,
                                              data,
                                              style
                                          }) => {
    const [edgePath, labelX, labelY] = getBezierPath({
        sourceX,
        sourceY,
        sourcePosition,
        targetX,
        targetY,
        targetPosition,
    })

    return (
        <>
            <path id={id} className="react-flow__edge-path" d={edgePath} markerEnd={markerEnd} style={style}/>
            <EdgeLabelRenderer>
                <div
                    className="nodrag nopan"
                    style={{
                        position: 'absolute',
                        transform: `translate(-50%, -50%) translate(${labelX}px,${labelY}px)`,
                        border: '1px solid black',
                        borderRadius: '5px',
                        padding: '10px',
                        backgroundColor: 'white',
                    }}
                >
                    <Box>
                        <p style={{margin: 0}}>Rows: <b>{betterNumbers(data.rows)}</b></p>
                    </Box>
                </div>
            </EdgeLabelRenderer>
        </>
    )
}