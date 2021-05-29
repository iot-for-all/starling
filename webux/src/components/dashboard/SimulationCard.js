import { useContext } from 'react';
import { Link, useHistory } from 'react-router-dom';
import {
    Button,
    Card,
    Header,
    Icon,
    Progress,
    Table,
    Text
} from "tabler-react";
import "tabler-react/dist/Tabler.css";
import "./Dashboard.css";
import SimulationDeviceStatsTableRow from './SimulationDeviceStatsTableRow';
import moment from "moment";
import GlobalContext from '../../context/globalContext';
import * as Notification from "../../components/notification/Notification";

function formatNumber(num) {
    return num.toString().replace(/(\d)(?=(\d{3})+(?!\d))/g, '$1,')
}

const SimulationCard = (props) => {
    const globalContext = useContext(GlobalContext)
    const history = useHistory();
    const { sim } = props;

    const totalProvisionedDevices = sim.devices.reduce((currentNumber, device) => {
        return currentNumber + device.provisionedCount;
    }, 0);
    const totalSimulatedDevices = sim.devices.reduce((currentNumber, device) => {
        return currentNumber + device.simulatedCount;
    }, 0);
    const totalConnectedDevices = sim.devices.reduce((currentNumber, device) => {
        return currentNumber + device.connectedCount;
    }, 0);
    const percentConnected = (totalSimulatedDevices > 0) ? Math.floor(100 * totalConnectedDevices / totalSimulatedDevices) : 0;
    let progressColor = "gray";
    if (sim.status === "running") {
        if (percentConnected === 100) {
            progressColor = "green";
        } else if (percentConnected >= 70 && percentConnected < 100) {
            progressColor = "orange";
        } else if (percentConnected < 70) {
            progressColor = "red";
        }
    }

    let statusColor = "green";
    let canStart = false;
    let canStop = false;
    let simStatusName = "";
    let simStatusBadge = "";
    if (sim.status === "ready") {
        statusColor = "secondary";
        canStart = true;
        simStatusName = "Ready";
        const msg = "No devices are connected as simulation is not running.";
        simStatusBadge = <div className="float-left">
            <span className="text-dark" title={msg}><Icon prefix="fe" name="check" /></span>
        </div>;
    } else if (sim.status === "provisioning") {
        statusColor = "primary";
        simStatusName = "Devices Provisioning";
        const msg = "No devices are connected as simulation is not running.";
        simStatusBadge = <div className="float-left">
            <span className="text-dark" title={msg}><Icon prefix="fe" name="check" /></span>
        </div>;
    } else if (sim.status === "running") {
        simStatusName = "Running";
        if (totalSimulatedDevices !== totalConnectedDevices) {
            statusColor = "danger";
            const msg = formatNumber(totalSimulatedDevices - totalConnectedDevices) + " devices are not connected.";
            simStatusBadge = <div className="float-left">
                <span className="text-danger" title={msg}><Icon prefix="fe" name="alert-triangle" /></span>
            </div>;
        } else {
            const msg = "All " + formatNumber(totalSimulatedDevices) + " devices are connected.";
            simStatusBadge = <div className="float-left">
                <span className="text-success" title={msg}><Icon prefix="fe" name="check" /></span>
            </div>;
        }
        canStop = true;
    } else if (sim.status === "deleting") {
        statusColor = "danger";
        simStatusName = "Deleting Simulation";
        const msg = "Deleting devices.";
        simStatusBadge = <div className="float-left">
            <span className="text-dark" title={msg}><Icon prefix="fe" name="check" /></span>
        </div>;
    }
    const textStatusColor = "text-" + statusColor;

    const deviceRows = sim.devices.map((device) => {
        return <SimulationDeviceStatsTableRow
            key={device.modelId}
            status={sim.status}
            model={device.modelId}
            provisionedCount={device.provisionedCount}
            simulatedCount={device.simulatedCount}
            connectedCount={device.connectedCount}
        />;
    });

    const dt = new Date(Date.parse(sim.lastUpdatedTime));
    const lastUpdateStr = moment(dt).format("L LTS");

    const startHandler = async () => {
        await globalContext.startSimulation(sim.id);
        Notification.addNotification("success", "Success", `Simulation '${sim.name}' started.`);
    };

    const stopHandler = async () => {
        await globalContext.stopSimulation(sim.id);
        Notification.addNotification("success", "Success", `Simulation '${sim.name}' stopped.`);
    };

    const exportHandler = async () => {
        await globalContext.exportSimulation(sim.id);
    };

    const deleteHandler = async () => {
        Notification.addNotification("info", "Started", `Deleting simulation '${sim.name}'.`);
        await globalContext.deleteSimulation(sim.id);
    };

    const simBusy = (sim.status !== "ready");

    const app = globalContext.getApplication(sim.targetId);
    const targetName = (app) ? app.name : sim.targetId;

    return (
        <div className="dashboardItem">
            <Card statusColor={statusColor} className={"dashboardCard"}>
                <Card.Header>
                    <Card.Title>{sim.name}</Card.Title>
                    <Card.Options>
                        {canStart && <span title="Start this simulation">
                            <Button color="primary" size="sm" icon="play"
                                className="ml-2" onClick={startHandler}>Start</Button></span>}
                        {canStop && <span title="Stop this simulation">
                            <Button color="primary" size="sm" icon="square"
                                className="ml-2" onClick={stopHandler}>Stop</Button></span>}
                    </Card.Options>
                </Card.Header>
                <Card.Header className="simToolBar">
                    <span title="Reconfigure this simulation">
                        <Button color="primary" size="sm" outline icon="edit-2" onClick={() => history.push(`/sim/${sim.id}`)}></Button></span>
                    <span title="Provision devices for this simulation">
                        <Button color="primary" size="sm" outline icon="grid" disabled={simBusy} className="ml-2" onClick={() => history.push(`/sim/${sim.id}?provision`)}></Button></span>
                    <span title="Export this simulation as a shell script">
                        <Button color="primary" size="sm" outline icon="share" className="ml-2" onClick={exportHandler}></Button></span>
                    <span title="Delete this simulation">
                        <Button color="danger" size="sm" outline icon="trash-2" disabled={simBusy} className="ml-2" onClick={deleteHandler}></Button></span>
                </Card.Header>
                <Card.Body>
                    <div className="simProgressCardContainer">
                        <div className="simProgressCardItem">
                            <Table
                                cards={true}
                                responsive={true}
                                className="simFieldsTable">
                                <Table.Header>
                                    <Table.Row>
                                        <Table.ColHeader colSpan={2}>
                                            Configuration
                                    </Table.ColHeader>
                                    </Table.Row>
                                </Table.Header>
                                <Table.Body>
                                    <Table.Row>
                                        <Table.Col>ID</Table.Col>
                                        <Table.Col>
                                            <Link to={`/sim/${props.sim.id}`} >{props.sim.id}</Link>
                                        </Table.Col>
                                    </Table.Row>
                                    <Table.Row>
                                        <Table.Col>Target</Table.Col>
                                        <Table.Col>
                                            <Link to={`/app/${props.sim.targetId}`} >{targetName}</Link>
                                        </Table.Col>
                                    </Table.Row>
                                    <Table.Row>
                                        <Table.Col>Provisioned Devices</Table.Col>
                                        <Table.Col>{formatNumber(totalProvisionedDevices)}</Table.Col>
                                    </Table.Row>
                                    <Table.Row>
                                        <Table.Col>Telemetry Interval</Table.Col>
                                        <Table.Col>{formatNumber(sim.telemetryInterval)} {(sim.telemetryInterval > 1) ? "secs" : "sec"}</Table.Col>
                                    </Table.Row>
                                    <Table.Row>
                                        <Table.Col>Telemetry Batch Size</Table.Col>
                                        <Table.Col>{formatNumber(sim.telemetryBatchSize)}</Table.Col>
                                    </Table.Row>
                                    <Table.Row>
                                        <Table.Col>Telemetry Format</Table.Col>
                                        <Table.Col>{formatNumber(sim.telemetryFormat)}</Table.Col>
                                    </Table.Row>
                                    <Table.Row>
                                        <Table.Col>Reported Prop. Interval</Table.Col>
                                        <Table.Col>{formatNumber(sim.reportedPropertyInterval)} {(sim.reportedPropertyInterval > 1) ? "secs" : "sec"}</Table.Col>
                                    </Table.Row>
                                    <Table.Row>
                                        <Table.Col>Wave Groups</Table.Col>
                                        <Table.Col>{formatNumber(sim.waveGroupCount)}</Table.Col>
                                    </Table.Row>
                                    <Table.Row>
                                        <Table.Col>Wave Group Interval</Table.Col>
                                        <Table.Col>{formatNumber(sim.waveGroupInterval)} {(sim.waveGroupInterval > 1) ? "secs" : "sec"}</Table.Col>
                                    </Table.Row>
                                    <Table.Row>
                                        <Table.Col></Table.Col>
                                        <Table.Col></Table.Col>
                                    </Table.Row>
                                </Table.Body>
                            </Table>
                        </div>
                        <div>
                            <div className="simProgressCardItem">
                                <Card>
                                    <Card.Body className="text-center">
                                        <Header size={5}>Status</Header>
                                        <Header size={2} className={textStatusColor}>{simStatusName}</Header>
                                        <div className="float-right">
                                            <Text size="sm" className="simProgressTime">Since {lastUpdateStr}</Text>
                                        </div>
                                    </Card.Body>
                                </Card>
                            </div>
                            <div className="simProgressCardItem">
                                <Card>
                                    <Card.Body className="text-center">
                                        <Header size={5}>Connected Devices</Header>
                                        <Header size={2} className="text-dark">{formatNumber(totalConnectedDevices)}</Header>

                                        <Progress size="sm">
                                            <Progress.Bar color={progressColor} width={percentConnected} />
                                        </Progress>
                                        {simStatusBadge}
                                        <div className="float-right">
                                            <Text size="sm" muted>{percentConnected}%</Text>
                                        </div>
                                    </Card.Body>
                                </Card>
                            </div>
                        </div>
                    </div>
                    <Table
                        cards={true}
                        striped={true}
                        responsive={true}
                        className="table-vcenter"
                    >
                        <Table.Header>
                            <Table.Row>
                                <Table.ColHeader>Model</Table.ColHeader>
                                <Table.ColHeader>Provisioned</Table.ColHeader>
                                <Table.ColHeader>Simulated</Table.ColHeader>
                                <Table.ColHeader>Connected</Table.ColHeader>
                            </Table.Row>
                        </Table.Header>
                        <Table.Body>
                            {deviceRows}
                        </Table.Body>
                    </Table>
                </Card.Body>
            </Card>
        </div >
    );
}

export default SimulationCard;