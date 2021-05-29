import { useContext } from 'react';
import { Link } from "react-router-dom";
import {
    Table
} from "tabler-react";
import GlobalContext from '../../context/globalContext';

function formatNumber(num) {
    return num.toString().replace(/(\d)(?=(\d{3})+(?!\d))/g, '$1,')
}

const SimulationDeviceStatsTableRow = (props) => {
    const globalContext = useContext(GlobalContext)

    let statusBadge = "";
    if (props.simulatedCount === props.connectedCount) {
        const msg = "All " + formatNumber(props.connectedCount) + " devices are connected";
        statusBadge = <span className="status-icon bg-success" title={msg} />;
    } else {
        const msg = formatNumber(props.simulatedCount - props.connectedCount) + " devices are not connected.";
        statusBadge = <span className="status-icon bg-danger" title={msg} />;

    }
    if (props.status !== "running") {
        statusBadge = <span className="status-icon bg-gray" title="No devices are connected." />;
    }
    
    const model = globalContext.getModel(props.model);
    const modelName = (model) ? model.name : props.model;
    return (<Table.Row>
        <Table.Col>                
            <Link to={`/model/${props.model}`}> {modelName} </Link>
        </Table.Col>
        <Table.Col>{formatNumber(props.provisionedCount)}</Table.Col>
        <Table.Col>{formatNumber(props.simulatedCount)}</Table.Col>
        <Table.Col>{statusBadge} {formatNumber(props.connectedCount)}</Table.Col>
    </Table.Row>);
};

export default SimulationDeviceStatsTableRow;