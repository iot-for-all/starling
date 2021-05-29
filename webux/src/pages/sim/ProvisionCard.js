import { useContext, useEffect, useState } from 'react';
import { useHistory } from "react-router-dom";
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
import * as Utils from '../../utils/utils';
import HelpPopup from "../../components/help/HelpPopup";
import "./SimCard.css";
import * as Notification from "../../components/notification/Notification";

const ProvisionCard = (props) => {
    const globalContext = useContext(GlobalContext)
    const [sim, setSim] = useState(props.data);
    const [errors, setErrors] = useState({});
    const [backendError, setBackendError] = useState("");
    const history = useHistory();

    useEffect(() => {
        setSim(props.data);
        setErrors({});
        setBackendError("");
    }, [props.data]);

    const deviceChangeHandler = (event, index) => {
        let updatedSim = {
            ...sim,
        }
        updatedSim.devices[index].provisionedCount = event.target.value;
        setSim(updatedSim);
        //console.log("value changed target: ", event.target, " index: ", index, " value: ", event.target.value, ", updatedSim: ", updatedSim);
    }

    const onSubmit = async (event) => {
        event.preventDefault();

        // validate form
        let newErrors = {};
        let hasErrors = false;
        setErrors(newErrors);

        if (!hasErrors) {
            //console.log("Form submitted mode: ", props.mode, " sim: ", sim);
            try {
                const updatedSim = {
                    ...sim
                };
                // convert strings to numbers
                for (let i = 0; i < updatedSim.devices.length; i++) {
                    updatedSim.devices[i].provisionedCount = +updatedSim.devices[i].provisionedCount;
                }

                await globalContext.provisionSimulationDevices(updatedSim.id, updatedSim.devices);
                Notification.addNotification("info", "Started", `Provisioning devices for simulation '${updatedSim.name}'.`);

                history.push(`/sim/${sim.id}`);
            } catch (ex) {
                setBackendError(Utils.getErrorMessage(ex, "error provisioning devices"));
            }
        }
    }

    const title = "Provision devices - " + sim.name;
    const deviceRows = sim && sim.devices.map((device, index) => {
        //const fieldName = device.id + "SimulatedCount";
        const fieldName = `devices[${index}]`;
        //console.log("setting device: ", device, " fieldName:", fieldName);

        return <Table.Row key={device.id}>
            <Table.Col>
                {device.modelId}
            </Table.Col>
            <Table.Col>
                <Form.Input
                    name={`${fieldName}.provisionedCount`}
                    value={sim.devices[index].provisionedCount}
                    required
                    type="number"
                    onChange={(event) => { deviceChangeHandler(event, index) }}
                    invalid={errors.devices && errors.devices[index].provisionedCount ? true : false}
                    feedback="Provisioned devices is required"
                />
            </Table.Col>
            <Table.Col>{ }</Table.Col>
        </Table.Row>;
    });

    return <>
        <div>{props.backendError}</div>
        <form onSubmit={onSubmit} method="post">
            <Card>
                <Card.Header>
                    <Card.Title>{title}</Card.Title>
                    <Card.Options>
                        <span title="Provision devices for this simulation">
                            <Button color="primary" size="sm" icon="grid" className="ml-2" type="submit">Provision</Button>
                        </span>
                        <span title="Cancel all changes">
                            <Button color="primary" size="sm" className="ml-2" outline type="button"
                                onClick={() => history.push("/sim/"+sim.id)} >Cancel</Button>
                        </span>
                    </Card.Options>
                </Card.Header>
                <Card.Body>
                    {backendError && backendError.length > 0 && <div className="alert alert-danger">
                        <Icon prefix="fe" name="alert-triangle" />{" "}
                        {backendError}
                    </div>}
                    <p>
                        Provision devices for this simulation. Devices can either be added or removed by increasing
                        or decreasing the numbers below.
                    </p>
                    <Form.FieldSet>
                        <Grid.Row>
                            <Grid.Col colSpan="2">
                                <h4>Provision Devices Count</h4>
                                <Text className="small">Devices are provisioned in IoT Central for the given model below.
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
                                                        Device Count
                                                    </Grid.Col>
                                                    <Grid.Col
                                                        auto
                                                        className="align-self-center"
                                                    >
                                                        <HelpPopup content={<>Number of devices to be provisioned for this simulation.</>} placement="bottom" />
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

export default ProvisionCard;
