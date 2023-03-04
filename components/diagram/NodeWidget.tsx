import React, { memo } from 'react'
import { Handle, Position } from 'reactflow'
import { Cards, Link } from '@cloudscape-design/components'

// @ts-ignore
export const NodeWidget = memo(({ data, isConnectable }) => {
  return (
    <>
      <Handle
        type="target"
        position={Position.Left}
        isConnectable={false}
      />
      <Cards
        cardDefinition={{
          header: item => (
            <Link fontSize="heading-m">{data.label}</Link>
          ),
          sections: [
            {
              id: 'description',
              header: 'Description',
              content: item => item.description,
            },
            {
              id: 'type',
              header: 'Type',
              content: item => item.type,
            },
            {
              id: 'size',
              header: 'Size',
              content: item => item.size,
            },
          ],
        }}
        cardsPerRow={[
          { cards: 1 },
          { minWidth: 500, cards: 2 },
        ]}
        items={[
          {
            name: 'Item 1',
            alt: 'First',
            description: 'This is the first item',
            type: '1A',
            size: 'Small',
          },

        ]}
        loadingText="Loading resources"
      />
      <Handle
        type="source"
        position={Position.Right}
        isConnectable={false}
      />
    </>
  )
})

