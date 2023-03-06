import React, { memo } from 'react'
import { Handle, Position } from 'reactflow'
import { Box } from '@cloudscape-design/components'
import { PlanRow } from '../types'
// @ts-ignore
import Highlight from 'react-highlight'
import { betterNumbers, getCellWarningColor } from '../utils'

const getScopes = (data: PlanRow): JSX.Element => {
  return (
    <Box variant="p">
      <p style={{margin: 0}}>{data.node.scope && `on`} <code>{data.node.scope}</code></p>
      <p style={{margin: 0}}>{data.node.index && `by`} <code> {data.node.index}</code></p>
      <p style={{margin: 0}}>{data.node.filters && `filter by`} <code> {data.node.filters}</code></p>
    </Box>
  )
}

// @ts-ignore
export const NodeWidget = memo(({ data, isConnectable }: { data: PlanRow }) => {
  return (
    <>
      <Handle
        type="target"
        position={Position.Left}
        isConnectable={false}
        style={{backgroundColor: 'transparent', color: 'transparent'}}
      />

      <div style={{
        border: '1px solid black',
        borderRadius: '10px',
        padding: '10px',
      }}>
        <Box variant="h3">
          {data.node.operation}
        </Box>
        {getScopes(data)}
        <Box variant="p">
          Time: {betterNumbers(data.exclusive)} ms
        </Box>
        <div
          style={{
            backgroundColor: getCellWarningColor(data.exclusive, data.execution_time),
            width: '100%',
            height: '10px',
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

