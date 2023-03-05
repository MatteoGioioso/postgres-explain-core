import React, { FC } from 'react'
import { EdgeProps, getBezierPath, EdgeLabelRenderer } from 'reactflow'
import { Box, Container, Flashbar, Popover } from '@cloudscape-design/components'
import { betterNumbers } from '../utils'

export const EdgeWidget: FC<EdgeProps> = ({
                                            id,
                                            sourceX,
                                            sourceY,
                                            targetX,
                                            targetY,
                                            sourcePosition,
                                            targetPosition,
                                            data,
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
      <path id={id} className="react-flow__edge-path" d={edgePath}/>
      <EdgeLabelRenderer>
        <Box>
          <div
            style={{
              position: 'absolute',
              transform: `translate(-50%, -50%) translate(${labelX}px,${labelY}px)`,
              // background: '#ffcc00',
              padding: 10,
              borderRadius: 5,
              fontSize: 12,
              fontWeight: 700,
            }}
            className="nodrag nopan"
          >
            <Flashbar items={[
              {
                type: 'info',
                dismissible: false,
                content: (
                  <>
                    Rows returned:
                    <p style={{marginBottom: 0}}><b>{betterNumbers(data.rows)}</b></p>
                  </>
                ),
                id: 'message_1',
              },
            ]}/>
          </div>
        </Box>
      </EdgeLabelRenderer>
    </>
  )
}