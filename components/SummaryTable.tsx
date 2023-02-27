import React from 'react'
import { Box, ColumnLayout, SpaceBetween, Table, TableProps } from '@cloudscape-design/components'
import { PlanRow, SummaryTableProps } from './interfaces'

function getCellWarningColor(reference: number, total: number): string {
  if (total === 0 || total === undefined) return '#fff'

  const percentage = (reference / total) * 100
  if (percentage <= 10) {
    return '#fff'
  }

  if (percentage > 10 && percentage < 50) {
    return '#fe8'
  }

  if (percentage > 50 && percentage < 90) {
    return '#e80'
  }

  if (percentage > 90) {
    return '#800'
  }

  return '#fff'
}

// Cloudscape does not allow for cell background to be set:
// https://github.com/cloudscape-design/components/discussions/718
const ComparatorCell = ({ prop, totalProp }: {prop: number, totalProp: number}) => {
  return (
    <>
      <div style={{
        fontSize: '20px',
        display: 'inline',
        color: getCellWarningColor(prop, totalProp),
      }}>&#9632;</div>
      <div style={{ color: '#2f2f2f', display: 'inline', marginLeft: '5px' }}>{prop}</div>
    </>
  )
}

const explainerColumns: Array<TableProps.ColumnDefinition<PlanRow>> = [
  {
    id: 'exclusive',
    header: 'Exclusive',
    cell: (e) => <ComparatorCell prop={e.exclusive} totalProp={e.execution_time}/>,
  },
  {
    id: 'inclusive',
    header: 'Inclusive',
    cell: (e) => <ComparatorCell prop={e.inclusive} totalProp={e.execution_time}/>,
  },
  {
    id: 'rows',
    header: 'Rows',
    cell: (e) => Math.floor(e.rows),
  },
  {
    id: 'rows_x',
    header: 'Rows_x',
    cell: (e) => Math.floor(e.rows_x),
  },
  {
    id: 'node',
    header: 'Node',
    cell: (e) => (
      <ColumnLayout columns={1} variant="text-grid">
        <SpaceBetween size="l" direction="horizontal">
          {"---".repeat(e.node.level) + "->"}
          <SpaceBetween size="xxxs">
            <div>
              <Box variant="awsui-key-label">{e.node.operation} {e.node.relation && `on`} {e.node.relation}</Box>
            </div>
            <div>
              <div>{e.node.costs}</div>
            </div>
            <div>
              <div>{e.node.buffers}</div>
            </div>
          </SpaceBetween>
        </SpaceBetween>
      </ColumnLayout>
    ),
  },
]

export const SummaryTable = ({summary, stats}: SummaryTableProps) => {
  return (
    <Table
      items={summary}
      columnDefinitions={explainerColumns}
      variant="embedded"
      footer={
        <Box textAlign="center">
          <div>Planning time
            time: {stats.planning_time}</div>
          <div>Total Execution
            time: <b>{stats.execution_time}</b></div>
        </Box>
      }
    />
  )
}