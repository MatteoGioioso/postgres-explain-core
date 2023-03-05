import React from 'react'
import { Box, Button, ColumnLayout, Popover, SpaceBetween, Table, TableProps } from '@cloudscape-design/components'
import { PlanRow } from './types'
import { SummaryTableProps } from './interfaces'
// @ts-ignore
import Highlight from 'react-highlight'
import { betterNumbers, getCellWarningColor } from './utils'

const GenericDetailsPopover = (props: { name: string, content: any, children: string }) => {
  return (
    <Popover
      dismissAriaLabel="Close"
      header={props.name}
      content={props.content}
      triggerType="custom"
      dismissButton={false}
      position="top"
      size="large"
    >
      {props.children}
    </Popover>
  )
}

const getRowEstimateDirectionSymbol = (direction: string): string => {
  switch (direction) {
    case 'over':
      return '↑'
    case 'under':
      return '↓'
    default:
      return ''
  }
}

// Cloudscape does not allow for cell background to be set:
// https://github.com/cloudscape-design/components/discussions/718
const ComparatorCell = ({ prop, totalProp, name }: { prop: number, totalProp: number, name?: string }) => {
  return (
    <Popover
      dismissAriaLabel="Close"
      header={name}
      content={prop}
      triggerType="custom"
      dismissButton={false}
      position="top"
    >
      <div style={{
        fontSize: '20px',
        display: 'inline',
        color: getCellWarningColor(prop, totalProp),
      }}>&#9632;</div>
      <div style={{ color: '#2f2f2f', display: 'inline', marginLeft: '5px' }}>{betterNumbers(prop)}</div>
    </Popover>
  )
}

const explainerColumns: Array<TableProps.ColumnDefinition<PlanRow>> = [
  {
    id: 'exclusive',
    header: 'Time',
    cell: (e) => <ComparatorCell prop={e.exclusive} totalProp={e.execution_time} name={'Exclusive time'}/>,
  },
  {
    id: 'inclusive',
    header: 'Cumulative Time',
    cell: (e) => <ComparatorCell prop={e.inclusive} totalProp={e.execution_time} name={'Inclusive time'}/>,
  },
  {
    id: 'rows',
    header: 'Rows',
    cell: (e) => <GenericDetailsPopover content={e.rows.total}
                                        name="Rows">{betterNumbers(e.rows.total)}</GenericDetailsPopover>,
  },
  {
    id: 'rows-removed',
    header: 'Rows Removed',
    cell: (e) => (
      <>
        {
          e.rows.filters && (
            <>
              - {' '}
              <GenericDetailsPopover
                content={
                  <div>
                    <p>Filters: <Highlight>{e.rows.filters}</Highlight></p>
                    <p>Removed: {e.rows.removed}</p>
                  </div>
                }
                name="Rows removed by a filter"
              >
                {betterNumbers(e.rows.removed)}
              </GenericDetailsPopover>
            </>
          )
        }
      </>
    ),
  },
  {
    id: 'rows_x',
    header: (
      <>
        Rows E
        <Popover
          dismissAriaLabel="Close"
          header={'Rows ES'}
          content={'Rows estimate factor'}
          triggerType="custom"
          dismissButton={false}
          position="top"
        >
          <Button iconName="status-info" variant="icon"/>
        </Popover>
      </>

    ),
    cell: (e) => (
      <>
        {getRowEstimateDirectionSymbol(e.rows.estimation_direction) + ' '}
        <GenericDetailsPopover content={e.rows.estimation_factor}
                               name="Rows estimate factor">{betterNumbers(e.rows.estimation_factor)}</GenericDetailsPopover>
      </>
    ),
  },
  {
    id: 'loops',
    header: 'Loops',
    cell: (e) => <GenericDetailsPopover name={'Loops'}
                                        content={e.loops}>{betterNumbers(e.loops)}</GenericDetailsPopover>,
  },
  {
    id: 'node',
    header: 'Node',
    cell: (e) => (
      <ColumnLayout columns={1} variant="text-grid">
        <SpaceBetween size="l" direction="horizontal">
          {'└' + '──'.repeat(e.node.level) + '->'}
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
    minWidth: 1100,
  },
]

export const SummaryTable = ({ summary, stats }: SummaryTableProps) => {
  return (
    <Table
      items={summary}
      columnDefinitions={explainerColumns}
      variant="embedded"
      footer={
        <Box textAlign="left">
          <h3>Planning time
            time: {stats.planning_time} ms</h3>
          <h2>Execution
            time: <b>{stats.execution_time} ms</b></h2>
        </Box>
      }
    />
  )
}