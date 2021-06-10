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
import * as Utils from '../../utils/utils';
import * as Notification from "../../components/notification/Notification";

const AppCard = (props) => {
    const globalContext = useContext(GlobalContext)
    const [app, setApp] = useState(props.data);
    const [errors, setErrors] = useState({});
    const [backendError, setBackendError] = useState("");
    const history = useHistory();
    const queryParams = new URLSearchParams(useLocation().search);
    const fromIntro = (queryParams.has("intro"));

    useEffect(() => {
        //console.log("ModelForm useEffect data: ", props.data)

        setApp(props.data);
        setErrors({});
        setBackendError("");
    }, [props.data]);

    const changeHandler = (event) => {
        let updatedApp = {
            ...app,
            [event.target.name]: event.target.value
        }
        setApp(updatedApp);
        //console.log("value changed target: ", event.target, " value: ", event.target.value, ", updatedApp: ", updatedApp);
    }

    const changeCheckHandler = (event) => {
        let updatedApp = {
            ...app,
            [event.target.name]: event.target.checked
        }
        setApp(updatedApp);
        //console.log("value changed target: ", event.target, " value: ", event.target.value, ", updatedApp: ", updatedApp);
    }
    
    const onSubmit = async (event) => {
        event.preventDefault();

        // validate form
        let newErrors = {};
        let hasErrors = false;
        if (app.name.trim() === "") {
            newErrors.name = true;
            hasErrors = true;
        }
        if (app.appUrl.trim() === "") {
            newErrors.appUrl = true;
            hasErrors = true;
        }
        if (app.appToken.trim() === "") {
            newErrors.appToken = true;
            hasErrors = true;
        }
        if (app.provisioningUrl.trim() === "") {
            newErrors.provisioningUrl = true;
            hasErrors = true;
        }
        if (app.idScope.trim() === "") {
            newErrors.idScope = true;
            hasErrors = true;
        }
        if (app.masterKey.trim() === "") {
            newErrors.masterKey = true;
            hasErrors = true;
        }

        setErrors(newErrors);

        if (!hasErrors) {
            //console.log("Form submitted mode: ", props.mode, " app: ", app);
            try {
                const updatedApp = {
                    ...app
                };
                let appId = app.id;
                if (props.mode === "add") {
                    const addedApp = await globalContext.addApplication(updatedApp);
                    appId = addedApp.id;
                    Notification.addNotification("success", "Success", `Application '${updatedApp.name}' is added.`);
                } else {
                    await globalContext.updateApplication(updatedApp);
                    Notification.addNotification("success", "Success", `Application '${updatedApp.name}' is updated.`);
                }

                if (fromIntro) {
                    history.push("/");
                } else {
                    history.push(`/app/${appId}`);
                }
            } catch (ex) {
                setBackendError(Utils.getErrorMessage(ex, "error saving application"));
            }
        }
    }

    const importHandler = async (event) => {
        event.preventDefault();
        try {
            Notification.addNotification("info", "Started", `Importing devices models from '${app.name}'.`);
            await globalContext.importModels(app.id);
            Notification.addNotification("success", "Success", `Devices models are imported from '${app.name}'.`);
        } catch (ex) {
            setBackendError(Utils.getErrorMessage(ex, "error importing models"));
        }
    };

    const deleteHandler = async (event) => {
        event.preventDefault();
        //console.log("Delete App", app.id);

        try {
            await globalContext.deleteApplication(app.id);
            Notification.addNotification("success", "Success", `Application '${app.name}' is deleted.`);
            history.push("/app");
        } catch (ex) {
            setBackendError(Utils.getErrorMessage(ex, "error deleting application"));
        }
    };

    const title = props.mode ? (props.mode === "add") ? "Add new application" : "Edit application - " + props.data.name : "";
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
                                    onClick={() => history.push("/app")} >Cancel</Button>
                            </span>
                        }
                    </Card.Options>
                </Card.Header>
                {props.mode !== "add" &&
                    <Card.Header className="simToolBar">
                        <span title="Import all device models from this application">
                            <Button color="primary" size="sm" outline icon="download" onClick={importHandler}>Import Device Models</Button></span>
                        <span title="Delete this application">
                            <Button color="danger" outline size="sm" icon="trash-2" className="ml-2" onClick={deleteHandler}>Delete</Button>
                        </span>
                    </Card.Header>
                }
                <Card.Body>
                    {backendError && backendError.length > 0 && <div className="alert alert-danger">
                        <Icon prefix="fe" name="alert-triangle" />{" "}
                        {backendError}
                    </div>}
                    <p>
                        Simulated devices are created in an IoT Central application. Create an IoT Central and
                        enter these details so that the devices can be provisioned using the credentials given below.
                        Multiple simulations can be created against an application and executed simultaneously.
                    </p>
                    <Form.Group
                        isRequired
                        label="Application Name"
                    >
                        <Grid.Row gutters="xs">
                            <Grid.Col>
                                <Form.Input
                                    name="name"
                                    value={app.name}
                                    required
                                    onChange={changeHandler}
                                    invalid={errors.name ? true : false}
                                    feedback="Application Name is required"
                                />
                            </Grid.Col>
                            <Grid.Col
                                auto
                                className="align-self-center"
                            >
                                <HelpPopup content={<>Descriptive name of the IoT Central application.</>} />
                            </Grid.Col>
                        </Grid.Row>
                    </Form.Group>
                    <Form.Group
                        isRequired
                        label="Application URL"
                    >
                        <Grid.Row gutters="xs">
                            <Grid.Col>
                                <Form.InputGroup>
                                    <Form.InputGroupPrepend>
                                        <Form.InputGroupText>
                                            https://
                                </Form.InputGroupText>
                                    </Form.InputGroupPrepend>
                                    <Form.Input
                                        name="appUrl"
                                        value={app.appUrl}
                                        required
                                        onChange={changeHandler}
                                        invalid={errors.appUrl ? true : false}
                                        feedback="Application URL is required"
                                    />
                                </Form.InputGroup>
                            </Grid.Col>
                            <Grid.Col
                                auto
                                className="align-self-center"
                            >
                                <HelpPopup content={<>
                                    <p>URL to access IoT Central application.</p>
                                    This is used by the IoT Central API to delete the devices when a simulation is deleted.
                                    </>} />
                            </Grid.Col>
                        </Grid.Row>
                    </Form.Group>
                    <Form.Group
                        isRequired
                        label="Application API Token"
                    >
                        <Grid.Row gutters="xs">
                            <Grid.Col>
                                <Form.Input
                                    name="appToken"
                                    value={app.appToken}
                                    required
                                    onChange={changeHandler}
                                    invalid={errors.name ? true : false}
                                    feedback="Application API Token is required"
                                />
                            </Grid.Col>
                            <Grid.Col
                                auto
                                className="align-self-center"
                            >
                                <HelpPopup content={<>
                                    <p><strong>API Token</strong> to access the IoT Central application.</p>
                                    <p>You can get this from IoT Central application <strong> Administration</strong> <Icon prefix="fe" name="arrow-right" /> {" "}
                                        <strong>API Token</strong> <Icon prefix="fe" name="arrow-right" /> {" "}
                                        <strong>Generate Token</strong>.</p>
                                    API Token is used to delete devices from IoT Central. Make sure that this token has delete devices permission.</>} />
                            </Grid.Col>
                        </Grid.Row>
                    </Form.Group>
                    <Form.Group
                        isRequired
                        label="Device Provisioning URL"
                    >
                        <Grid.Row gutters="xs">
                            <Grid.Col>
                                <Form.InputGroup>
                                    <Form.InputGroupPrepend>
                                        <Form.InputGroupText>
                                            https://
                                </Form.InputGroupText>
                                    </Form.InputGroupPrepend>
                                    <Form.Input
                                        name="provisioningUrl"
                                        value={app.provisioningUrl}
                                        required
                                        onChange={changeHandler}
                                        invalid={errors.provisioningUrl ? true : false}
                                        feedback="Device Provisioning URL is required"
                                    />
                                </Form.InputGroup>
                            </Grid.Col>
                            <Grid.Col
                                auto
                                className="align-self-center"
                            >
                                <HelpPopup content={<>
                                    <p>The endpoint for Device Provisioning Service (DPS).</p>
                                    <p>Typically it is https://<strong>global.azure-devices-provisioning.net</strong>.</p>
                                Rarely this might be different.
                                </>} />
                            </Grid.Col>
                        </Grid.Row>
                    </Form.Group>
                    <Form.Group
                        isRequired
                        label="ID Scope"
                    >
                        <Grid.Row gutters="xs">
                            <Grid.Col>
                                <Form.Input
                                    name="idScope"
                                    value={app.idScope}
                                    required
                                    onChange={changeHandler}
                                    invalid={errors.idScope ? true : false}
                                    feedback="ID Scope is required"
                                />
                            </Grid.Col>
                            <Grid.Col
                                auto
                                className="align-self-center"
                            >
                                <HelpPopup content={<>
                                    <p>The <strong>ID Scope</strong> for the IoT Central application.</p>
                                    <p>You can get it from IoT Central application <strong> Administration</strong> <Icon prefix="fe" name="arrow-right" /> {" "}
                                        <strong>Device Connection</strong> <Icon prefix="fe" name="arrow-right" /> {" "}
                                        <strong>ID Scope</strong>.</p>
                                    This ID Scope is used during provisioning of simulated devices.
                                    </>} />
                            </Grid.Col>
                        </Grid.Row>
                    </Form.Group>
                    <Form.Group
                        isRequired
                        label="Device Connection SAS Key"
                    >
                        <Grid.Row gutters="xs">
                            <Grid.Col>
                                <Form.Input
                                    name="masterKey"
                                    value={app.masterKey}
                                    required
                                    onChange={changeHandler}
                                    invalid={errors.masterKey ? true : false}
                                    feedback="Device Connection SAS Key is required"
                                />
                            </Grid.Col>
                            <Grid.Col
                                auto
                                className="align-self-center"
                            >
                                <HelpPopup content={<>
                                    <p>The <strong>Device Connection SAS Key</strong> for the IoT Central application.</p>
                                    <p>You can get it from IoT Central application <strong>Administration</strong> <Icon prefix="fe" name="arrow-right" /> {" "}
                                        <strong>Device Connection</strong> <Icon prefix="fe" name="arrow-right" /> {" "}
                                        <strong>SAS-IoT-Devices</strong> <Icon prefix="fe" name="arrow-right" /> {" "}
                                        <strong>Shared access signature (SAS)</strong> <Icon prefix="fe" name="arrow-right" /> {" "}
                                        <strong>Primary key</strong>.</p>
                                    <p>This Device Connection SAS Key is used during provisioning of simulated devices.</p>
                                    Currently, X509 certificates are not supported in Starling.
                                    </>} />
                            </Grid.Col>
                        </Grid.Row>
                    </Form.Group>
                    {props.mode === "add" &&
                        <Form.Group>
                            <Grid.Row gutters="xs">
                                <Grid.Col>
                                    <Form.Checkbox
                                        name="importModels"
                                        label="Automatically import all device models from this application"
                                        checked={app.importModels}
                                        onChange={changeCheckHandler}
                                    />
                                </Grid.Col>
                                <Grid.Col
                                    auto
                                    className="align-self-center"
                                >
                                    <HelpPopup content={<><p>After adding the application, import all device models from the application.</p>
                                                You can always import them after adding the application.</>} />
                                </Grid.Col>
                            </Grid.Row>
                        </Form.Group>
                    }
                </Card.Body>
            </Card>
        </form>
    </>
}

export default AppCard;
