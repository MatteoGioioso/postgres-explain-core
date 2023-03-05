import React, { memo } from 'react'
import { Handle, Position } from 'reactflow'
import { Box } from '@cloudscape-design/components'
import { PlanRow } from '../types'
// @ts-ignore
import Highlight from 'react-highlight'
import { betterNumbers, getCellWarningColor } from '../utils'

// @ts-ignore
export const NodeWidget = memo(({ data, isConnectable }: { data: PlanRow }) => {
  return (
    <>
      <Handle
        type="target"
        position={Position.Left}
        isConnectable={false}
      />

      <div style={{
        border: '1px solid black',
        borderRadius: '10px',
        padding: "10px",
      }}>
        <Box variant="h2">
          {data.node.operation}
        </Box>
        <Box variant="p">
          <Highlight style={{padding: 0}}>
            {data.node.scope}
          </Highlight>
        </Box>
        <Box variant="p">
          Time: {betterNumbers(data.exclusive)} ms
        </Box>
        <div
          style={{
            backgroundColor: getCellWarningColor(data.exclusive, data.execution_time),
            width: '100%',
            height: '10px'
          }}
        ></div>
      </div>
      <Handle
        type="source"
        position={Position.Right}
        isConnectable={false}
      />
    </>
  )
})

