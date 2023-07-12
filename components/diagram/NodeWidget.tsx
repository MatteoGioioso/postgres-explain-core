import React, {memo} from 'react'
import {Handle, NodeToolbar, Position} from 'reactflow'
import {Box, Button, Modal, SpaceBetween} from '@cloudscape-design/components'
import {PlanRow} from '../types'
// @ts-ignore
import Highlight from 'react-highlight'
import {betterNumbers, getCellWarningColor, truncateText} from '../utils'

const getScopes = (data: PlanRow): JSX.Element => {
    return (
        <Box variant="p">
            {data.node.scope && (<p style={{margin: 0}}><b>on </b><code>{truncateText(data.node.scope, 25)}</code></p>)}
            {data.node.index && (<p style={{margin: 0}}><b>by </b><code>{truncateText(data.node.index, 25)}</code></p>)}
            {data.node.filters && (<p style={{margin: 0}}><b>filter by </b><code>{truncateText(data.node.filters, 25)}</code></p>)}
        </Box>
    )
}

// @ts-ignore
export const NodeWidget = memo(({data, isConnectable}: { data: PlanRow }) => {
    const [visible, setVisible] = React.useState(false);
    const [visibleExpand, setVisibleExpand] = React.useState(false);

    return (
        <>
            {/*<NodeToolbar isVisible={visibleExpand} position={Position.Top} offset={0}>*/}
            {/*    <button*/}
            {/*        onClick={e => setVisible(true)}*/}
            {/*        style={{*/}
            {/*            height: "20px",*/}
            {/*            width: "30px",*/}
            {/*            backgroundColor: "white",*/}
            {/*            border: "solid 1px",*/}
            {/*            borderRadius: "5px",*/}
            {/*        }}*/}
            {/*    ><Box fontSize="body-s">+</Box></button>*/}
            {/*</NodeToolbar>*/}
            <Handle
                type="target"
                position={Position.Top}
                isConnectable={false}
                style={{backgroundColor: 'transparent', color: 'transparent'}}
            />

            <div
                onMouseOver={e => setVisibleExpand(true)}
                onMouseLeave={e => setTimeout(() => {
                    setVisibleExpand(false)
                }, 2000)}
                style={{
                    border: '1px solid black',
                    borderRadius: '10px',
                    padding: '10px',
                    backgroundColor: 'white',
                    boxShadow: '4px 8px 5px -4px #C9C9C9',
                    width: '250px',
                    height: '80px'
                }}
            >
                <Box variant="h3">
                    {data.node.operation}
                </Box>
                {/*{getScopes(data)}*/}
                {/*<Box variant="p">*/}
                {/*    Time: {betterNumbers(data.exclusive)} ms*/}
                {/*</Box>*/}
                <div
                    style={{
                        backgroundColor: getCellWarningColor(data.exclusive, data.execution_time),
                        width: '100%',
                        height: '10px',
                    }}
                ></div>
                <Box variant="p">
                    Rows returned: {betterNumbers(data.rows.total)}
                </Box>
            </div>
            <Handle
                type="source"
                position={Position.Bottom}
                isConnectable={false}
            />

            {/*<Modal*/}
            {/*    onDismiss={() => setVisible(false)}*/}
            {/*    visible={visible}*/}
            {/*    size="large"*/}
            {/*    closeAriaLabel="Close modal"*/}
            {/*    header={data.node.operation}*/}
            {/*>*/}
            {/*   <SpaceBetween size="xl">*/}
            {/*       <Box variant="p">*/}
            {/*           {data.node.scope && (<p><b>on</b><Highlight>{data.node.scope}</Highlight></p>)}*/}
            {/*           {data.node.index && (<p><b>by</b><Highlight>{data.node.index}</Highlight></p>)}*/}
            {/*           {data.node.filters && (<p><b>filter by</b><Highlight>{data.node.filters}</Highlight></p>)}*/}
            {/*       </Box>*/}
            {/*       <Box>*/}
            {/*           <p><b>Time spent: </b>{data.exclusive}ms</p>*/}
            {/*           <p><b>Time impact: </b>{(Math.floor((data.exclusive / data.execution_time) * 100))}%</p>*/}
            {/*           <div*/}
            {/*               style={{*/}
            {/*                   backgroundColor: getCellWarningColor(data.exclusive, data.execution_time),*/}
            {/*                   width: '100%',*/}
            {/*                   height: '10px',*/}
            {/*               }}*/}
            {/*           ></div>*/}
            {/*       </Box>*/}

            {/*       <Box>*/}
            {/*           <Box variant="h3">Disk</Box>*/}
            {/*           <p style={{marginBottom: 0}}><b>Buffers hit:</b> {data.buffers.hits} blocks</p>*/}
            {/*           <p style={{marginBottom: 0}}><b>Buffers read:</b> {data.buffers.reads} blocks</p>*/}
            {/*           <p style={{marginBottom: 0}}><b>Buffers written:</b> {data.buffers.written} blocks</p>*/}
            {/*       </Box>*/}
            {/*   </SpaceBetween>*/}
            {/*</Modal>*/}
        </>
    )
})

