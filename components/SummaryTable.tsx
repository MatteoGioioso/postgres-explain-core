import React from 'react'
import { Box, ColumnLayout, Popover, SpaceBetween, Table, TableProps } from '@cloudscape-design/components'
import { PlanRow } from './types'
import { SummaryTableProps } from './interfaces'

const betterNumbers = (num: number): string => {
  const ONE_MILLION = 1000000
  const THOUSAND = 1000
  const HUNDRED = 100
  const TEN = 10
  const ONE = 1

  if (num >= ONE_MILLION) {
    return `${Math.floor(num / ONE_MILLION)} Mil`
  }

  if (num > THOUSAND) {
    return `${Math.floor(num / THOUSAND)} K`
  }

  if (num <= THOUSAND && num >= HUNDRED) {
    return `${Math.floor((num / THOUSAND) * 10) / 10} K`
  }

  if (num <= HUNDRED && num >= TEN) {
    return `${Math.floor((num / HUNDRED) * 100) / 100} K`
  }

  if (num <= TEN && num >= ONE) {
    return `${Math.floor(num * 100) / 100}`
  }

  if (num <= ONE) {
    return `${Math.floor(num * 1000) / 1000}`
  }

  return num.toString()
}

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


const GenericNumberDetailsPopover = (props: { name: string, number: number, children: string }) => {
  return (
    <Popover
      dismissAriaLabel="Close"
      header={props.name}
      content={props.number}
      triggerType="custom"
      dismissButton={false}
      position="top"
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
    header: 'Exclusive',
    cell: (e) => <ComparatorCell prop={e.exclusive} totalProp={e.execution_time} name={'Exclusive time'}/>,
  },
  {
    id: 'inclusive',
    header: 'Inclusive',
    cell: (e) => <ComparatorCell prop={e.inclusive} totalProp={e.execution_time} name={'Inclusive time'}/>,
  },
  {
    id: 'rows',
    header: 'Rows',
    cell: (e) => <GenericNumberDetailsPopover number={e.rows} name="Rows">{betterNumbers(e.rows)}</GenericNumberDetailsPopover>,
  },
  {
    id: 'rows_x',
    header: 'Rows estimate factor',
    cell: (e) => (
      <>
        {getRowEstimateDirectionSymbol(e.rows_x.direction) + " "}
        <GenericNumberDetailsPopover number={e.rows_x.value} name="Rows estimate factor">{betterNumbers(e.rows_x.value)}</GenericNumberDetailsPopover>
      </>
    ),
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