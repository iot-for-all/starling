import { useContext, useEffect, useState } from 'react';
import { Link, useHistory, useLocation } from "react-router-dom";
import {
    Button,
    Card,
    Form,
    Grid,
    Icon,
    Table,
    Text,
} from "tabler-react";
import GlobalContext from '../../context/globalContext';
import HelpPopup from "../../components/help/HelpPopup";
import "./SimCard.css";
import * as Utils from '../../utils/utils';
import * as Notification from "../../components/notification/Notification";

const SimCard = (props) => {
    const globalContext = useContext(GlobalContext)
    const [sim, setSim] = useState(props.data);
    const [errors, setErrors] = useState({});
    const [backendError, setBackendError] = useState("");
    const history = useHistory();
    const queryParams = new URLSearchParams(useLocation().search);
    const fromIntro = (queryParams.has("intro"));

    useEffect(() => {
        setSim(props.data);
        setErrors({});
        setBackendError("");
    }, [props.data]);

    useEffect(() => {
        //console.log("useEffect sim updated");

        if (sim && sim.id && sim.status !== "ready") {
            //console.log("useEffect sim refreshed");
            let timer = setInterval(() => {
                globalContext.listSimulations();
                const remoteSim = globalContext.getSimulation(sim.id);
                setSim(remoteSim);
            }, 5000);
            return () => {
                clearInterval(timer);
            }
        }

        // ignore global context dependency error
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [sim]);

    const changeHandler = (event) => {
        let updatedSim = {
            ...sim,
            [event.target.name]: event.target.value
        }
        setSim(updatedSim);
        //console.log("value changed target: ", event.target, " value: ", event.target.value, ", updatedSim: ", updatedSim);
    }

    const changeNumberHandler = (event) => {
        //console.log("value changed target: ", event.target, " value: ", event.target.value, " valid: ", event.target.validity.valid);
        if (event.target.value.match(/^[0-9]+$/) || event.target.value.trim().length === 0) {
            let val = 0;
            if (event.target.value.trim().length > 0) {
                val = +event.target.value;
            }
            let updatedSim = {
                ...sim,
                [event.target.name]: val
            }
            setSim(updatedSim);
        }
    }

    const deviceChangeHandler = (event, index) => {
        if (event.target.value.match(/^[0-9]+$/) || event.target.value.trim().length === 0) {
            let val = 0;
            if (event.target.value.trim().length > 0) {
                val = +event.target.value;
            }
            let updatedSim = {
                ...sim,
            }
            updatedSim.devices[index].simulatedCount = val;
            setSim(updatedSim);
        }
    }

    const onSubmit = async (event) => {
        event.preventDefault();

        // validate form
        let newErrors = {};
        let hasErrors = false;
        if (sim.name.trim() === "") {
            newErrors.name = true;
            hasErrors = true;
        }
        if (sim.targetId.trim() === "") {
            newErrors.targetId = true;
            hasErrors = true;
        }
        if (sim.waveGroupCount <= 0) {
            newErrors.waveGroupCount = true;
            hasErrors = true;
        }
        if (sim.waveGroupInterval <= 0) {
            newErrors.waveGroupInterval = true;
            hasErrors = true;
        }
        if (sim.telemetryBatchSize <= 0) {
            newErrors.telemetryBatchSize = true;
            hasErrors = true;
        }
        if (sim.telemetryInterval <= 0) {
            newErrors.telemetryInterval = true;
            hasErrors = true;
        }
        if (sim.reportedPropertyInterval <= 0) {
            newErrors.reportedPropertyInterval = true;
            hasErrors = true;
        }
        if (sim.disconnectBehavior.trim() === "") {
            newErrors.disconnectBehavior = true;
            hasErrors = true;
        }
        if (sim.telemetryFormat.trim() === "") {
            newErrors.telemetryFormat = true;
            hasErrors = true;
        }
        for (let i = 0; i < sim.devices.length; i++) {
            if (sim.devices[i].simulatedCount < 0) {
                const fieldName = `devices[${i}].simulatedCount`;
                newErrors[fieldName] = true;
                hasErrors = true;
            }
        }

        setErrors(newErrors);

        if (!hasErrors) {
            //console.log("Form submitted mode: ", props.mode, " sim: ", sim);
            try {
                const updatedSim = {
                    ...sim
                };
                // convert strings to numbers
                updatedSim.telemetryInterval = +updatedSim.telemetryInterval;
                updatedSim.telemetryBatchSize = +updatedSim.telemetryBatchSize;
                updatedSim.reportedPropertyInterval = +updatedSim.reportedPropertyInterval;
                updatedSim.waveGroupCount = +updatedSim.waveGroupCount;
                updatedSim.waveGroupInterval = +updatedSim.waveGroupInterval;
                for (let i = 0; i < updatedSim.devices.length; i++) {
                    updatedSim.devices[i].simulatedCount = +updatedSim.devices[i].simulatedCount;
                }

                let simId = sim.id;
                if (props.mode === "add") {
                    const addedSim = await globalContext.addSimulation(updatedSim);
                    simId = addedSim.id;
                    Notification.addNotification("success", "Success", `Simulation '${updatedSim.name}' is added.`);
                } else {
                    await globalContext.updateSimulation(updatedSim);
                    Notification.addNotification("success", "Success", `Simulation '${updatedSim.name}' is updated.`);
                }

                if (fromIntro) {
                    history.push("/");
                } else {
                    history.push(`/sim/${simId}`);
                }
            } catch (ex) {
                setBackendError(Utils.getErrorMessage(ex, "error saving simulation"));
            }
        }
    }

    const startHandler = async (event) => {
        event.preventDefault();
        await globalContext.startSimulation(sim.id);
        Notification.addNotification("success", "Success", `Simulation '${sim.name}' started.`);
    };

    const stopHandler = async (event) => {
        event.preventDefault();
        await globalContext.stopSimulation(sim.id);
        Notification.addNotification("success", "Success", `Simulation '${sim.name}' stopped.`);
    };

    const exportHandler = async (event) => {
        event.preventDefault();
        await globalContext.exportSimulation(sim.id);
    };

    const deleteHandler = async (event) => {
        event.preventDefault();
        //console.log("Delete Simulation", sim.id);

        try {
            Notification.addNotification("info", "Started", `Deleting simulation '${sim.name}'.`);
            await globalContext.deleteSimulation(sim.id);
            history.push(`/sim/${sim.id}`);
        } catch (ex) {
            setBackendError(Utils.getErrorMessage(ex, "error deleting simulation"));
        }
    };

    const title = props.mode ? (props.mode === "add") ? "Add new simulation" : "Edit simulation - " + props.data.name : "";
    //console.log("SimCard - mode: ", props.mode, " backendError: ", props.backendError, " data: ", props.data);

    const appsArr = globalContext.apps ? Array.from(globalContext.apps) : [];
    const appsList = appsArr.map((app) => {
        return (
            <option key={app.id} value={app.id}>{app.name}</option>
        );
    });

    const deviceRows = sim && sim.devices.map((device, index) => {
        //const fieldName = device.id + "SimulatedCount";
        const fieldName = `devices[${index}]`;
        //console.log("setting device: ", device, " fieldName:", fieldName);
        let statusBadge = "";
        if (device.simulatedCount === device.connectedCount) {
            const msg = "All " + Utils.formatNumber(device.connectedCount) + " devices are connected";
            statusBadge = <span className="status-icon bg-success" title={msg} />;
        } else {
            const msg = Utils.formatNumber(device.simulatedCount - device.connectedCount) + " devices are not connected.";
            statusBadge = <span className="status-icon bg-danger" title={msg} />;

        }
        if (sim.status !== "running") {
            statusBadge = <span className="status-icon bg-gray" title="No devices are connected." />;
        }
        const simBusy = (sim.status !== "ready");
        const model = globalContext.getModel(device.modelId);
        const modelName = (model) ? model.name : props.model;

        return <Table.Row key={device.id}>
            <Table.Col>
                <Link to={`/model/${device.modelId}`}> {modelName} </Link>
            </Table.Col>
            <Table.Col>{Utils.formatNumber(device.provisionedCount)}</Table.Col>
            <Table.Col>
                <Form.Input
                    name={`${fieldName}.simulatedCount`}
                    value={sim.devices[index].simulatedCount}
                    required
                    type="text"
                    pattern="[0-9]*"
                    onChange={(event) => { deviceChangeHandler(event, index) }}
                    disabled={simBusy}
                    invalid={errors.devices && errors.devices[index].simulatedCount ? true : false}
                    feedback="Simulated devices is required"
                />
            </Table.Col>
            <Table.Col>{statusBadge} {Utils.formatNumber(device.connectedCount)}</Table.Col>
            <Table.Col>{ }</Table.Col>
        </Table.Row>;
    });

    const simBusy = (sim.status !== "ready");

    const totalSimulatedDevices = sim.devices.reduce((currentNumber, device) => {
        return currentNumber + device.simulatedCount;
    }, 0);
    const totalConnectedDevices = sim.devices.reduce((currentNumber, device) => {
        return currentNumber + device.connectedCount;
    }, 0);
    let simStatusName = "Ready";
    let statusColor = "green";
    if (sim.status === "ready") {
        statusColor = "secondary";
        simStatusName = "Ready";
    } else if (sim.status === "provisioning") {
        statusColor = "primary";
        simStatusName = "Devices Provisioning";
    } else if (sim.status === "running") {
        simStatusName = "Running";
        if (totalSimulatedDevices !== totalConnectedDevices) {
            statusColor = "danger";
        }
    } else if (sim.status === "deleting") {
        simStatusName = "Deleting";
        statusColor = "danger";
    }
    const textStatusColor = "text-" + statusColor;
    return <>
        <div>{props.backendError}</div>
        <form onSubmit={onSubmit}>
            <Card>
                <Card.Header>
                    <Card.Title>{title}</Card.Title>
                    <Card.Options>
                        <span title="Save this simulation">
                            <Button color="primary" size="sm" icon="save" disabled={simBusy} className="ml-2" type="submit">Save</Button>
                        </span>
                    </Card.Options>
                </Card.Header>
                {
                    props.mode !== "add" &&
                    <Card.Header className="simToolBar">
                        {
                            sim.status === "ready" &&
                            <span title="Start this simulation">
                                <Button color="primary" size="sm" icon="play" type="button" onClick={startHandler}>Start</Button></span>
                        }
                        {
                            sim.status === "running" &&
                            <span title="Stop this simulation">
                                <Button color="primary" size="sm" icon="square" type="button" onClick={stopHandler}>Stop</Button></span>
                        }
                        <span title="Provision devices for this simulation">
                            <Button color="primary" size="sm" outline icon="grid" disabled={simBusy} type="button" className="ml-2" onClick={() => history.push(`/sim/${sim.id}?provision`)}>Provision</Button></span>
                        <span title="Export this simulation as a shell script">
                            <Button color="primary" size="sm" outline icon="share" disabled={props.mode === "add"} type="button" className="ml-2" onClick={exportHandler}>Export</Button></span>
                        <span title="Delete this simulation">
                            <Button color="danger" size="sm" outline icon="trash-2" disabled={simBusy} className="ml-2" type="button" onClick={deleteHandler}>Delete</Button></span>
                        <div className="card-options mr-2">
                            <div className="text-dark">Status:{" "}</div>
                            <div className={textStatusColor}>{simStatusName}</div>
                        </div>
                    </Card.Header>
                }
                <Card.Body>
                    {backendError && backendError.length > 0 && <div className="alert alert-danger">
                        <Icon prefix="fe" name="alert-triangle" />{" "}
                        {backendError}
                    </div>}
                    <p>
                        Configure a simulation using the parameters below.
                        Simulations cannot be updated when they are running.
                        Devices are automatically provisioned when a simulation is started and these devices are reused in future executions of this simulation.
                        However, you can explicitly provision or de-provision ahead of time.
                    </p>
                    <Grid.Row>
                        <Grid.Col>
                            {simBusy && <Text className="text-default small"><Icon prefix="fe" name="alert-triangle" />{" "} Simulation cannot be updated when it is busy.</Text>}
                        </Grid.Col>
                        <Grid.Col>
                            <div className="simLearnMore">
                                <Text className="text-default"><Icon prefix="fe" name="help-circle" />{" "} <a href="https://github.com/iot-for-all/starling/blob/main/docs/configure.md#configuring-simulation" target="_blank" rel="noreferrer">Help me configure this simulation</a></Text>
                            </div>
                        </Grid.Col>
                    </Grid.Row>
                    <Form.FieldSet>
                        <Grid.Row>
                            <Grid.Col>
                                <h4>Base Configuration</h4>
                                <Form.Group
                                    isRequired
                                    label="Simulation Name"
                                >
                                    <Grid.Row gutters="xs">
                                        <Grid.Col>
                                            <Form.Input
                                                name="name"
                                                value={sim.name}
                                                required
                                                onChange={changeHandler}
                                                disabled={simBusy}
                                                invalid={errors.name ? true : false}
                                                feedback="Simulation Name is required"
                                            />
                                        </Grid.Col>
                                        <Grid.Col
                                            auto
                                            className="align-self-center"
                                        >
                                            <HelpPopup content={<>Name of the Simulation that is shown everywhere.</>} />
                                        </Grid.Col>
                                    </Grid.Row>
                                </Form.Group>
                                <Form.Group
                                    isRequired
                                    label="Target Application"
                                >
                                    <Grid.Row gutters="xs">
                                        <Grid.Col>
                                            <Form.Select
                                                name="targetId"
                                                value={sim.targetId}
                                                required
                                                onChange={changeHandler}
                                                disabled={props.mode !== "add"}
                                                invalid={errors.targetId ? true : false}
                                                feedback="Target Application is required"
                                            >
                                                <option value="">
                                                    Select one
                                        </option>
                                                {appsList}
                                            </Form.Select>
                                        </Grid.Col>
                                        <Grid.Col
                                            auto
                                            className="align-self-center"
                                        >
                                            <HelpPopup content={<>The IoT Central application in which the simulated devices are created.</>} />
                                        </Grid.Col>
                                    </Grid.Row>
                                </Form.Group>
                                <Form.Group
                                    isRequired
                                    label="Device Disconnect Behavior"
                                >
                                    <Grid.Row gutters="xs">
                                        <Grid.Col>
                                            <Form.Select
                                                name="disconnectBehavior"
                                                value={sim.disconnectBehavior}
                                                required
                                                onChange={changeHandler}
                                                disabled={simBusy}
                                                invalid={errors.disconnectBehavior ? true : false}
                                                feedback="Device Disconnect Behavior is required"
                                            >
                                                <option value="never">
                                                    Never Disconnect
                                        </option>
                                                <option value="telemetry">
                                                    After Sending Telemetry
                                        </option>
                                            </Form.Select>
                                        </Grid.Col>
                                        <Grid.Col
                                            auto
                                            className="align-self-center"
                                        >
                                            <HelpPopup content={<><p>Select <strong>Never Disconnect</strong> if you want your devices to be connected always. Use this for a typical IoT device.</p>
                                                Select <strong>After Sending Telemetry</strong> if you want your devices to be disconnected after sending each telemetry batch. Use this for simulating occasionally connected devices.</>} />
                                        </Grid.Col>
                                    </Grid.Row>
                                </Form.Group>
                                <Form.Group
                                    isRequired
                                    label="Telemetry Format"
                                >
                                    <Grid.Row gutters="xs">
                                        <Grid.Col>
                                            <Form.Select
                                                name="telemetryFormat"
                                                value={sim.telemetryFormat}
                                                required
                                                onChange={changeHandler}
                                                disabled={simBusy}
                                                invalid={errors.telemetryFormat ? true : false}
                                                feedback="Telemetry Format is required"
                                            >
                                                <option value="default">Default (JSON)</option>
                                                <option value="opcua">OPCUA (JSON)</option>
                                            </Form.Select>
                                        </Grid.Col>
                                        <Grid.Col
                                            auto
                                            className="align-self-center"
                                        >
                                            <HelpPopup content={<><p>Select <strong>Default</strong> for a typical IoT device.</p>
                                                Select <strong>OPCUA</strong> for sending data in OPCUA format. This is used for Industrial IoT (IIoT) devices.</>} />
                                        </Grid.Col>
                                    </Grid.Row>
                                </Form.Group>
                            </Grid.Col>
                            <Grid.Col>
                                <h4>Data Rates</h4>
                                <Form.Group
                                    isRequired
                                    label="Telemetry Interval (secs)"
                                >
                                    <Grid.Row gutters="xs">
                                        <Grid.Col>
                                            <Form.Input
                                                name="telemetryInterval"
                                                value={sim.telemetryInterval}
                                                required
                                                type="text"
                                                pattern="[0-9]*"
                                                onChange={changeNumberHandler}
                                                disabled={simBusy}
                                                invalid={errors.telemetryInterval ? true : false}
                                                feedback="Telemetry Interval (whole number) is required"
                                            />
                                        </Grid.Col>
                                        <Grid.Col
                                            auto
                                            className="align-self-center"
                                        >
                                            <HelpPopup content={<>Enter how often telemetry is sent from the device to IoT Central.</>} />
                                        </Grid.Col>
                                    </Grid.Row>
                                </Form.Group>
                                <Form.Group
                                    isRequired
                                    label="Telemetry Batch Size"
                                >
                                    <Grid.Row gutters="xs">
                                        <Grid.Col>
                                            <Form.Input
                                                name="telemetryBatchSize"
                                                value={sim.telemetryBatchSize}
                                                required
                                                type="text"
                                                pattern="[0-9]*"
                                                onChange={changeNumberHandler}
                                                disabled={simBusy}
                                                invalid={errors.telemetryBatchSize ? true : false}
                                                feedback="Telemetry Batch Size (whole number) is required"
                                            />
                                        </Grid.Col>
                                        <Grid.Col
                                            auto
                                            className="align-self-center"
                                        >
                                            <HelpPopup content={<>Number of telemetry messages to be sent everytime 'Telemetry Batch Size' seconds.</>} />
                                        </Grid.Col>
                                    </Grid.Row>
                                </Form.Group>
                                <Form.Group
                                    isRequired
                                    label="Reported Property Interval (secs)"
                                >
                                    <Grid.Row gutters="xs">
                                        <Grid.Col>
                                            <Form.Input
                                                name="reportedPropertyInterval"
                                                value={sim.reportedPropertyInterval}
                                                required
                                                type="text"
                                                pattern="[0-9]*"
                                                onChange={changeNumberHandler}
                                                disabled={simBusy}
                                                invalid={errors.reportedPropertyInterval ? true : false}
                                                feedback="Reported Property Interval (whole number) is required"
                                            />
                                        </Grid.Col>
                                        <Grid.Col
                                            auto
                                            className="align-self-center"
                                        >
                                            <HelpPopup content={<>Enter how often reported properties are sent from the device to IoT Central.</>} />
                                        </Grid.Col>
                                    </Grid.Row>
                                </Form.Group>
                                <Form.Group
                                    isRequired
                                    label="Wave Group Count"
                                >
                                    <Grid.Row gutters="xs">
                                        <Grid.Col>
                                            <Form.Input
                                                name="waveGroupCount"
                                                value={sim.waveGroupCount}
                                                required
                                                type="text"
                                                pattern="[0-9]*"
                                                onChange={changeNumberHandler}
                                                disabled={simBusy}
                                                invalid={errors.waveGroupCount ? true : false}
                                                feedback="Wave Group Count (whole number) is required"
                                            />
                                        </Grid.Col>
                                        <Grid.Col
                                            auto
                                            className="align-self-center"
                                        >
                                            <HelpPopup content={<><p>Devices can be divided into several number of waves. Telemetry is sent sequentially from devices in one wave after another.</p>
                                                <p>For a simple simulation, set it to <strong>1</strong>.</p>
                                                Say, if you want to divide your fleet of devices into 2 waves, 30 seconds apart, set <strong>Wave Group Count</strong> as <strong>2</strong> and <strong>Wave Group Interval</strong> as <strong>30</strong>.
                                            </>} />
                                        </Grid.Col>
                                    </Grid.Row>
                                </Form.Group>
                                <Form.Group
                                    isRequired
                                    label="Wave Group Interval (secs)"
                                >
                                    <Grid.Row gutters="xs">
                                        <Grid.Col>
                                            <Form.Input
                                                name="waveGroupInterval"
                                                value={sim.waveGroupInterval}
                                                required
                                                type="text"
                                                pattern="[0-9]*"
                                                onChange={changeNumberHandler}
                                                disabled={simBusy}
                                                invalid={errors.waveGroupInterval ? true : false}
                                                feedback="Wave Group Interval (whole number) is required"
                                            />
                                        </Grid.Col>
                                        <Grid.Col
                                            auto
                                            className="align-self-center"
                                        >
                                            <HelpPopup content={<><p>Interval between each wave group.</p>
                                                <p>If the <strong>Wave Group Count</strong> is set as <strong>1</strong>, you can leave <strong>Wave Group Interval</strong> as <strong>1</strong>.</p>
                                                Say, if you want to divide your fleet of devices into 2 waves, 30 seconds apart, set <strong>Wave Group Count</strong> as <strong>2</strong> and <strong>Wave Group Interval</strong> as <strong>30</strong>.
                                            </>} />
                                        </Grid.Col>
                                    </Grid.Row>
                                </Form.Group>
                            </Grid.Col>
                        </Grid.Row>
                        <Grid.Row>
                            <Grid.Col colSpan="2">
                                <h4>Simulated Devices</h4>
                                <Text className="small">Devices are automatically provisioned in IoT Central when the simulation is started.
                                These devices are reused whenever this simulation is executed.
                                You can add/delete provisioned devices using the Provision button above.
                                All provisioned devices will be deleted when simulation is deleted.
                                    <p></p>
                                    <Text className="text-default"><strong><Icon prefix="fe" name="alert-triangle" />{" "} Provisioned devices in IoT Central will be billed to your Azure subscription.</strong></Text>
                                </Text>

                                <Table
                                    cards={true}
                                    striped={true}
                                    responsive={true}
                                    className="table-vcenter fillerTable"
                                >
                                    <Table.Header>
                                        <Table.Row>
                                            <Table.ColHeader>
                                                <Grid.Row gutters="xs">
                                                    <Grid.Col>
                                                        Model
                                                    </Grid.Col>
                                                    <Grid.Col
                                                        auto
                                                        className="align-self-center"
                                                    >
                                                        <HelpPopup content={<>Devices of this type are simulated.</>} placement="bottom" />
                                                    </Grid.Col>
                                                </Grid.Row>
                                            </Table.ColHeader>
                                            <Table.ColHeader>
                                                <Grid.Row gutters="xs">
                                                    <Grid.Col>
                                                        Provisioned
                                                    </Grid.Col>
                                                    <Grid.Col
                                                        auto
                                                        className="align-self-center"
                                                    >
                                                        <HelpPopup content={<>Number of devices currently provisioned for this simulation.</>} placement="bottom" />
                                                    </Grid.Col>
                                                </Grid.Row>
                                            </Table.ColHeader>
                                            <Table.ColHeader>
                                                <Grid.Row gutters="xs">
                                                    <Grid.Col>
                                                        Simulated
                                                    </Grid.Col>
                                                    <Grid.Col
                                                        auto
                                                        className="align-self-center"
                                                    >
                                                        <HelpPopup content={<>Number of devices requested to be simulated.</>} placement="bottom" />
                                                    </Grid.Col>
                                                </Grid.Row>
                                            </Table.ColHeader>
                                            <Table.ColHeader>
                                                <Grid.Row gutters="xs">
                                                    <Grid.Col>
                                                        Connected
                                                    </Grid.Col>
                                                    <Grid.Col
                                                        auto
                                                        className="align-self-center"
                                                    >
                                                        <HelpPopup content={<>Number of devices currently connected for this simulation.</>} placement="bottom" />
                                                    </Grid.Col>
                                                </Grid.Row>
                                            </Table.ColHeader>
                                            <Table.ColHeader></Table.ColHeader>
                                        </Table.Row>
                                    </Table.Header>
                                    <Table.Body>
                                        {deviceRows}
                                    </Table.Body>
                                </Table>
                            </Grid.Col>
                        </Grid.Row>
                    </Form.FieldSet>
                </Card.Body>
            </Card>
        </form>
    </>
}

export default SimCard;
