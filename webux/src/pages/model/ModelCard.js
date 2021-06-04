import { useContext, useEffect, useState } from 'react';
import { useHistory, useLocation } from "react-router-dom";
import {
    Button,
    Card,
    Form,
    Grid,
    Icon,
} from "tabler-react";
import GlobalContext from '../../context/globalContext';
import HelpPopup from "../../components/help/HelpPopup";
import HelpButtonPopup from "../../components/help/HelpButtonPopup";
import * as Utils from '../../utils/utils';
import "./ModelCard.css";
import * as Notification from "../../components/notification/Notification";

const ModelCard = (props) => {
    const globalContext = useContext(GlobalContext)
    const [model, setModel] = useState(props.data);
    const [errors, setErrors] = useState({});
    const [backendError, setBackendError] = useState("");
    const history = useHistory();
    const queryParams = new URLSearchParams(useLocation().search);
    const fromIntro = (queryParams.has("intro"));

    useEffect(() => {
        //console.log("ModelForm useEffect data: ", props.data)

        props.data.capabilityModel = JSON.stringify(props.data.capabilityModel, null, 2);
        setModel(props.data);
        setErrors({});
        setBackendError("");
    }, [props.data]);

    const changeHandler = (event) => {
        let updatedModel = {
            ...model,
            [event.target.name]: event.target.value
        }
        setModel(updatedModel);
        //console.log("value changed target: ", event.target, " value: ", event.target.value, ", updatedModel: ", updatedModel);
    }

    const changeIDHandler = (event) => {
        let updatedModel = {
            ...model,
            [event.target.name]: event.target.value.toLowerCase()
        }
        setModel(updatedModel);
    }

    const onSubmit = async (event) => {
        event.preventDefault();

        // validate form
        let newErrors = {};
        let hasErrors = false;
        if (model.id.trim() === "" || !model.id.match(/^[0-9a-z]+$/)) {
            newErrors.id = true;
            hasErrors = true;
        }
        if (model.name.trim() === "") {
            newErrors.name = true;
            hasErrors = true;
        }
        if (model.capabilityModel.trim() === "") {
            newErrors.capabilityModel = true;
            hasErrors = true;
        }
        setErrors(newErrors);

        if (!hasErrors) {
            //console.log("Form submitted mode: ", props.mode, " model: ", model);
            try {
                const updatedModel = {
                    ...model,
                    capabilityModel: JSON.parse(model.capabilityModel)
                };
                if (props.mode === "add") {
                    await globalContext.addModel(updatedModel);
                    Notification.addNotification("success", "Success", `Device model '${updatedModel.name}' is added.`);
                } else {
                    await globalContext.updateModel(updatedModel);
                    Notification.addNotification("success", "Success", `Device model '${updatedModel.name}' is updated.`);
                }

                if (fromIntro) {
                    history.push("/");
                } else {
                    history.push(`/model/${model.id}`);
                }
            } catch (ex) {
                setBackendError(Utils.getErrorMessage(ex, "error saving model"));
            }
        }
    }

    const deleteHandler = async (event) => {
        event.preventDefault();
        //console.log("Delete Model", model.id);

        try {
            await globalContext.deleteModel(model.id);
            Notification.addNotification("success", "Success", `Device model '${model.name}' is deleted.`);
            history.push("/model");
        } catch (ex) {
            setBackendError(Utils.getErrorMessage(ex, "error deleting model"));
        }
    };

    const title = props.mode ? (props.mode === "add") ? "Add new device model" : "Edit model - " + props.data.name : "";
    //console.log("ModelForm - mode: ", props.mode, " backendError: ", props.backendError, " data: ", props.data);

    return <>
        <div>{props.backendError}</div>
        <form onSubmit={onSubmit}>
            <Card>
                <Card.Header>
                    <Card.Title>{title}</Card.Title>
                    <Card.Options>
                        <span title="Save this application">
                            <Button
                                color="primary"
                                size="sm"
                                icon="save"
                                className="ml-2"
                                onClick={onSubmit}
                            >Save</Button>
                        </span>
                        {props.mode === "add" &&
                            <span title="Cancel all changes">
                                <Button color="primary" size="sm" className="ml-2"
                                    onClick={() => history.push("/model")} >Cancel</Button>
                            </span>
                        }
                    </Card.Options>
                </Card.Header>
                {
                    props.mode !== "add" &&
                    <Card.Header className="simToolBar">
                        <span title="Delete this device model">
                            <Button color="danger" size="sm" outline icon="trash-2" type="button" onClick={deleteHandler}>Delete</Button></span>
                    </Card.Header>
                }
                <Card.Body>
                    {backendError && backendError.length > 0 && <div className="alert alert-danger">
                        <Icon prefix="fe" name="alert-triangle" />{" "}
                        {backendError}
                    </div>}
                    <p>
                        Device models describe the behavior of a device using <a href="https://github.com/Azure/opendigitaltwins-dtdl/blob/master/DTDL/v2/dtdlv2.md">DTDL</a> JSON Model.
                        Simulated devices can be spun up based on these device models.
                        Same model can be used across multiple applications.
                        Make sure that this model exists in the IoT Central application before starting the simulation.
                    </p>
                    <Form.Group
                        isRequired
                        label="Model ID"
                    >
                        <Grid.Row gutters="xs">
                            <Grid.Col>
                                <Form.Input
                                    name="id"
                                    value={model.id}
                                    onChange={changeIDHandler}
                                    disabled={props.mode !== "add"}
                                    invalid={errors.id ? true : false}
                                    feedback="Model ID is required (only alphanumeric characters are allowed)"
                                />
                            </Grid.Col>
                            <Grid.Col
                                auto
                                className="align-self-center"
                            >
                                <HelpPopup content={<>
                                    <p>Unique ID for the model. Only lowecase alphanumeric characters are allowed. E.g.: <strong>mymodel</strong></p>
                                    <p>Devices in a simulation are named SimID-AppID-ModelID-###.</p>Choose this ID wisely.</>} />
                            </Grid.Col>
                        </Grid.Row>
                    </Form.Group>
                    <Form.Group
                        isRequired
                        label="Model Name"
                    >
                        <Grid.Row gutters="xs">
                            <Grid.Col>
                                <Form.Input
                                    name="name"
                                    value={model.name}
                                    required
                                    onChange={changeHandler}
                                    invalid={errors.name ? true : false}
                                    feedback="Model Name is required"
                                />
                            </Grid.Col>
                            <Grid.Col
                                auto
                                className="align-self-center"
                            >
                                <HelpPopup content={<>Descriptive name of the model.</>} />
                            </Grid.Col>
                        </Grid.Row>
                    </Form.Group>
                    <div className="modelLearnMore">
                        <HelpButtonPopup
                            content={<>
                                <h6>DTDL support in Starling</h6>
                                <ol className="helpButton">
                                    <li><strong>Supported data types:</strong> boolean, date, datetime, double, duration, float, geopoint, integer, long, string, time.</li>
                                    <li><strong>Unsupported data types:</strong> enum, map, object, vector, array, hashmaps, event, state.</li>
                                    <li><strong>Interfaces</strong> are supported. <strong>Components</strong> are not supported.</li>
                                    <li><strong>Direct methods</strong> are acknowledged. They currently do not return any data.</li>
                                    <li><strong>C2D commands</strong> are not <i>completed</i> or return any data as response.</li>
                                    <li><strong>Regular devices</strong> are supported. <strong>Gateways and Edge devices</strong> are not supported.</li>
                                </ol>
                            </>}
                            label={<>DTDL Support</>}
                        />
                    </div>
                    <Form.Group
                        isRequired
                        label="Capability Model"
                    >
                        <Grid.Row gutters="xs">
                            <Grid.Col>
                                <Form.Textarea
                                    name="capabilityModel"
                                    value={model.capabilityModel}
                                    onChange={changeHandler}
                                    className="capabilityModel"
                                    invalid={errors.capabilityModel ? true : false}
                                    feedback="Capability Model is required"
                                    rows={15}
                                />
                            </Grid.Col>
                            <Grid.Col
                                auto
                                className="align-self-center"
                            >
                                <HelpPopup content={<>
                                    <p>Enter a <a
                                        href="https://github.com/Azure/opendigitaltwins-dtdl/blob/master/DTDL/v2/dtdlv2.md">DTDL</a> JSON Model.</p>
                                    You can get this from IoT Central application <strong> Device Templates</strong> <Icon prefix="fe" name="arrow-right" /> {" "}
                                    <strong>Select your model</strong> <Icon prefix="fe" name="arrow-right" /> {" "}
                                    <strong>Export</strong>.</>} />
                            </Grid.Col>
                        </Grid.Row>
                    </Form.Group>
                </Card.Body>
            </Card>
        </form>
    </>
}

export default ModelCard;
