import { useContext, useEffect, useState } from 'react';
import { useHistory } from "react-router-dom";
import {
    Button,
    Card,
    Form,
    Grid,
    Icon,
    Page,
    Text,
} from "tabler-react";

import "tabler-react/dist/Tabler.css";
import GlobalContext from '../../context/globalContext';
import SiteWrapper from '../../components/site/SiteWrapper';
import HelpPopup from "../../components/help/HelpPopup";
import * as Utils from '../../utils/utils';
import "./SettingsPage.css";
import * as Notification from "../../components/notification/Notification";

const SettingsPage = () => {
    const globalContext = useContext(GlobalContext)
    const history = useHistory();
    const [config, setConfig] = useState();
    const [errors, setErrors] = useState({});
    const [backendError, setBackendError] = useState("");
    const [simWarning, setSimWarning] = useState(false);
    const [appWarning, setAppWarning] = useState(false);
    const [saved, setSaved] = useState(false);

    // Called on mount to ensure reference data is loaded if coming from shortcut
    useEffect(() => {
        if (globalContext.config) {
            // deep copy the config object so that the original config containing geopoint array
            // does not get changed as we are converting it into string locally
            let localConfig = {
                ...globalContext.config,
            };
            localConfig.data = {
                ...globalContext.config.data
            };
            localConfig.http = {
                ...globalContext.config.http
            };
            localConfig.logger = {
                ...globalContext.config.logger
            };
            localConfig.simulation = {
                ...globalContext.config.simulation
            };
            if (localConfig.simulation.geopointData) {
                // default JSON Stringify is spewing too many lines, so we are building a simple JSON strings (one line per geopoint) below
                //localConfig.Simulation.geopointData = JSON.stringify(localConfig.Simulation.geopointData, null, 2);
                let str = "[\n";
                for (let i = 0; i < localConfig.simulation.geopointData.length; i++) {
                    str += "  [" + localConfig.simulation.geopointData[i][0] + ", " + localConfig.simulation.geopointData[i][1] + ", " + localConfig.simulation.geopointData[i][2] + "]";
                    if (i < localConfig.simulation.geopointData.length - 1) {
                        str += ",\n";
                    } else {
                        str += "\n";
                    }
                }
                str += "]";
                localConfig.simulation.geopointData = str;
            }

            setConfig(localConfig);
            setBackendError("");
        } else {
            setBackendError("Failed to get settings. Make sure that the Starling server is running.");
        }
        // ignore global context dependency error
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [globalContext.config])

    // need to change this
    const changeSimHandler = (event) => {
        let updatedSim = {
            ...config.simulation,
            [event.target.name]: event.target.value
        };
        let updatedConfig = {
            ...config,
            simulation: updatedSim,
        };
        setConfig(updatedConfig);
        setSimWarning(true);
        setSaved(false);
    }

    const changeSimCheckHandler = (event) => {
        let updatedSim = {
            ...config.simulation,
            [event.target.name]: event.target.checked
        };
        let updatedConfig = {
            ...config,
            simulation: updatedSim,
        };
        setConfig(updatedConfig);
        setSimWarning(true);
        setSaved(false);
    }

    // need to change this
    const changeDataHandler = (event) => {
        let updatedData = {
            ...config.data,
            [event.target.name]: event.target.value
        };
        let updatedConfig = {
            ...config,
            data: updatedData,
        };
        setConfig(updatedConfig);
        setAppWarning(true);
        setSaved(false);
    }

    // need to change this
    const changeLogHandler = (event) => {
        let updatedLog = {
            ...config.logger,
            [event.target.name]: event.target.value
        };
        let updatedConfig = {
            ...config,
            logger: updatedLog,
        };
        setConfig(updatedConfig);
        setAppWarning(true);
        setSaved(false);
    }
    
    // need to change this
    const changeSimulationNumberHandler = (event) => {
        //console.log("value changed target: ", event.target, " value: ", event.target.value, " valid: ", event.target.validity.valid);
        if (event.target.value.match(/^[0-9]+$/) || event.target.value.trim().length === 0) {
            let val = 0;
            if (event.target.value.trim().length > 0) {
                val = +event.target.value;
            }
            let updatedSim = {
                ...config.simulation,
                [event.target.name]: val
            };
            let updatedConfig = {
                ...config,
                simulation: updatedSim,
            };
            setConfig(updatedConfig);
            setSimWarning(true);
            setSaved(false);
        }
    }

    // need to change this
    const changeHttpNumberHandler = (event) => {
        //console.log("value changed target: ", event.target, " value: ", event.target.value, " valid: ", event.target.validity.valid);
        if (event.target.value.match(/^[0-9]+$/) || event.target.value.trim().length === 0) {
            let val = 0;
            if (event.target.value.trim().length > 0) {
                val = +event.target.value;
            }
            let updatedHttp = {
                ...config.http,
                [event.target.name]: val
            };
            let updatedConfig = {
                ...config,
                http: updatedHttp,
            };
            setConfig(updatedConfig);
            setAppWarning(true);
            setSaved(false);
        }
    }

    const onSubmit = async (event) => {
        event.preventDefault();

        // validate form
        let newErrors = {};
        let hasErrors = false;
        if (config.simulation.connectionTimeout <= 0) {
            newErrors.connectionTimeout = true;
            hasErrors = true;
        }
        if (config.simulation.telemetryTimeout <= 0) {
            newErrors.telemetryTimeout = true;
            hasErrors = true;
        }
        if (config.simulation.twinUpdateTimeout <= 0) {
            newErrors.twinUpdateTimeout = true;
            hasErrors = true;
        }
        if (config.simulation.commandTimeout <= 0) {
            newErrors.commandTimeout = true;
            hasErrors = true;
        }
        if (config.simulation.registrationAttemptTimeout <= 0) {
            newErrors.registrationAttemptTimeout = true;
            hasErrors = true;
        }
        if (config.simulation.maxConcurrentConnections <= 0) {
            newErrors.maxConcurrentConnections = true;
            hasErrors = true;
        }
        if (config.simulation.maxConcurrentTwinUpdates <= 0) {
            newErrors.maxConcurrentTwinUpdates = true;
            hasErrors = true;
        }
        if (config.simulation.maxConcurrentRegistrations <= 0) {
            newErrors.maxConcurrentRegistrations = true;
            hasErrors = true;
        }
        if (config.simulation.maxConcurrentDeletes <= 0) {
            newErrors.maxConcurrentDeletes = true;
            hasErrors = true;
        }
        if (config.simulation.maxRegistrationAttempts <= 0) {
            newErrors.maxRegistrationAttempts = true;
            hasErrors = true;
        }
        if (config.simulation.geopointData.trim() === "") {
            newErrors.geopointData = true;
            hasErrors = true;
        }
        if (config.http.adminPort <= 0) {
            newErrors.adminPort = true;
            hasErrors = true;
        }
        if (config.http.metricsPort <= 0) {
            newErrors.maxConcurrentRegistrations = true;
            hasErrors = true;
        }
        if (config.logger.logsDir.trim().length === 0) {
            newErrors.logsDir = true;
            hasErrors = true;
        }
        if (config.data.path.trim().length === 0) {
            newErrors.logsDir = true;
            hasErrors = true;
        }

        setErrors(newErrors);

        if (!hasErrors) {
            //console.log("Form submitted mode: ", props.mode, " sim: ", sim);
            try {
                const updatedConfig = {
                    ...config
                };

                let foundError = false;
                try {
                    updatedConfig.simulation.geopointData = JSON.parse(updatedConfig.simulation.geopointData);
                } catch (ex2) {
                    setBackendError("Error parsing Geopoint Data: " + Utils.getErrorMessage(ex2, "error parsing JSON"));
                    foundError = true;
                }

                if (!foundError) {
                    await globalContext.updateConfig(updatedConfig);
                    let msg = "Configuration is updated.";
                    if (simWarning) {
                        msg += " Restart simulations."
                    }
                    if (appWarning) {
                        msg += " Restart Starling server."
                    }
                    Notification.addNotification("success", "Success", msg);

                    setSaved(true);
                    history.replace(`/settings`);
                }
            } catch (ex) {
                setBackendError(Utils.getErrorMessage(ex, "error saving configuration"));
            }
        }
    }

    return <SiteWrapper>
        <Page.Content title="Settings">
            {backendError && backendError.length > 0 && <div className="alert alert-danger">
                <Icon prefix="fe" name="alert-triangle" />{" "}
                {backendError}
            </div>}
            {simWarning && saved &&
                <div className="alert alert-warning">
                    <Icon prefix="fe" name="alert-triangle" />{" "}
                    Some of the simulation level settings are changed. Restart simulations to apply changes.
                </div>
            }
            {appWarning && saved &&
                <div className="alert alert-warning">
                    <Icon prefix="fe" name="alert-triangle" />{" "}
                    Some of the application level settings are changed. Restart Starling server to apply changes.
                </div>
            }
            {
                config &&
                <form onSubmit={onSubmit}>
                    <Card>
                        <Card.Header>
                            <Card.Title>Starling Configuration</Card.Title>
                            <Card.Options>
                                <span title="Save this configuration">
                                    <Button
                                        color="primary"
                                        size="sm"
                                        icon="save"
                                        className="ml-2"
                                        onClick={onSubmit}
                                    >Save</Button>
                                </span>
                            </Card.Options>
                        </Card.Header>
                        <Card.Body>
                            <p>Configure Starling server using the following settings.</p>
                            <Form.FieldSet>
                                <Grid.Row>
                                    <Grid.Col>
                                        <h6>
                                            <Text className="text-default small"><strong><Icon prefix="fe" name="alert-triangle" />{" "} Restart simulations to apply these changes.</strong></Text>
                                        </h6>
                                        <h4>Simulation Settings</h4>
                                        <Form.Group
                                            isRequired
                                            label="Connection Timeout (ms)"
                                        >
                                            <Grid.Row gutters="xs">
                                                <Grid.Col>
                                                    <Form.Input
                                                        name="connectionTimeout"
                                                        value={config.simulation.connectionTimeout}
                                                        required
                                                        type="text"
                                                        pattern="[0-9]*"
                                                        onChange={changeSimulationNumberHandler}
                                                        invalid={errors.connectionTimeout ? true : false}
                                                        feedback="Connection Timeout (whole number) is required"
                                                    />
                                                </Grid.Col>
                                                <Grid.Col
                                                    auto
                                                    className="align-self-center"
                                                >
                                                    <HelpPopup content={<>Timeout in milliseconds for device connection.</>} />
                                                </Grid.Col>
                                            </Grid.Row>
                                        </Form.Group>
                                        <Form.Group
                                            isRequired
                                            label="Telemetry Timeout (ms)"
                                        >
                                            <Grid.Row gutters="xs">
                                                <Grid.Col>
                                                    <Form.Input
                                                        name="telemetryTimeout"
                                                        value={config.simulation.telemetryTimeout}
                                                        required
                                                        type="text"
                                                        pattern="[0-9]*"
                                                        onChange={changeSimulationNumberHandler}
                                                        invalid={errors.telemetryTimeout ? true : false}
                                                        feedback="Telemetry Timeout (whole number) is required"
                                                    />
                                                </Grid.Col>
                                                <Grid.Col
                                                    auto
                                                    className="align-self-center"
                                                >
                                                    <HelpPopup content={<>Timeout in milliseconds for sending telemetry.</>} />
                                                </Grid.Col>
                                            </Grid.Row>
                                        </Form.Group>
                                        <Form.Group
                                            isRequired
                                            label="Twin Update Timeout (ms)"
                                        >
                                            <Grid.Row gutters="xs">
                                                <Grid.Col>
                                                    <Form.Input
                                                        name="twinUpdateTimeout"
                                                        value={config.simulation.twinUpdateTimeout}
                                                        required
                                                        type="text"
                                                        pattern="[0-9]*"
                                                        onChange={changeSimulationNumberHandler}
                                                        invalid={errors.twinUpdateTimeout ? true : false}
                                                        feedback="Twin Update Timeout (whole number) is required"
                                                    />
                                                </Grid.Col>
                                                <Grid.Col
                                                    auto
                                                    className="align-self-center"
                                                >
                                                    <HelpPopup content={<>Timeout in milliseconds for updating device twin.
                                                This is used for acknowledging desired property changes or sending reported properties.</>} />
                                                </Grid.Col>
                                            </Grid.Row>
                                        </Form.Group>
                                        <Form.Group
                                            isRequired
                                            label="Command Timeout (ms)"
                                        >
                                            <Grid.Row gutters="xs">
                                                <Grid.Col>
                                                    <Form.Input
                                                        name="commandTimeout"
                                                        value={config.simulation.commandTimeout}
                                                        required
                                                        type="text"
                                                        pattern="[0-9]*"
                                                        onChange={changeSimulationNumberHandler}
                                                        invalid={errors.commandTimeout ? true : false}
                                                        feedback="Command Timeout (whole number) is required"
                                                    />
                                                </Grid.Col>
                                                <Grid.Col
                                                    auto
                                                    className="align-self-center"
                                                >
                                                    <HelpPopup content={<>Timeout in milliseconds for acknowledging commands.</>} />
                                                </Grid.Col>
                                            </Grid.Row>
                                        </Form.Group>
                                        <Form.Group
                                            isRequired
                                            label="Provisioning Timeout (ms)"
                                        >
                                            <Grid.Row gutters="xs">
                                                <Grid.Col>
                                                    <Form.Input
                                                        name="registrationAttemptTimeout"
                                                        value={config.simulation.registrationAttemptTimeout}
                                                        required
                                                        type="text"
                                                        pattern="[0-9]*"
                                                        onChange={changeSimulationNumberHandler}
                                                        invalid={errors.registrationAttemptTimeout ? true : false}
                                                        feedback="Registration Timeout (whole number) is required"
                                                    />
                                                </Grid.Col>
                                                <Grid.Col
                                                    auto
                                                    className="align-self-center"
                                                >
                                                    <HelpPopup content={<>Timeout in milliseconds for provisioning a device in Device Provisioning Service (DPS).</>} />
                                                </Grid.Col>
                                            </Grid.Row>
                                        </Form.Group>
                                        <Form.Group
                                            isRequired
                                            label="Max. Concurrent Telemetry"
                                        >
                                            <Grid.Row gutters="xs">
                                                <Grid.Col>
                                                    <Form.Input
                                                        name="maxConcurrentConnections"
                                                        value={config.simulation.maxConcurrentConnections}
                                                        required
                                                        type="text"
                                                        pattern="[0-9]*"
                                                        onChange={changeSimulationNumberHandler}
                                                        invalid={errors.maxConcurrentConnections ? true : false}
                                                        feedback="Max. Concurrent Telemetry (whole number) is required"
                                                    />
                                                </Grid.Col>
                                                <Grid.Col
                                                    auto
                                                    className="align-self-center"
                                                >
                                                    <HelpPopup content={<>Maximum number of devices to send telemetry at a time.
                                                You can simulate large number of devices (say 1,000) at any time, but concurrently this many (say 100) devices will send telemetry at the same instant.</>} />
                                                </Grid.Col>
                                            </Grid.Row>
                                        </Form.Group>
                                        <Form.Group
                                            isRequired
                                            label="Max. Concurrent Twin Updates"
                                        >
                                            <Grid.Row gutters="xs">
                                                <Grid.Col>
                                                    <Form.Input
                                                        name="maxConcurrentConnections"
                                                        value={config.simulation.maxConcurrentTwinUpdates}
                                                        required
                                                        type="text"
                                                        pattern="[0-9]*"
                                                        onChange={changeSimulationNumberHandler}
                                                        invalid={errors.maxConcurrentTwinUpdates ? true : false}
                                                        feedback="Max. Concurrent Twin Updates (whole number) is required"
                                                    />
                                                </Grid.Col>
                                                <Grid.Col
                                                    auto
                                                    className="align-self-center"
                                                >
                                                    <HelpPopup content={<>Maximum number of devices to send reported properties at a time.
                                                You can simulate large number of devices (say 1,000) at any time, but concurrently this many (say 100) devices will send reported properties at the same instant.</>} />
                                                </Grid.Col>
                                            </Grid.Row>
                                        </Form.Group>
                                        <Form.Group
                                            isRequired
                                            label="Max. Concurrent Device Registrations"
                                        >
                                            <Grid.Row gutters="xs">
                                                <Grid.Col>
                                                    <Form.Input
                                                        name="maxConcurrentRegistrations"
                                                        value={config.simulation.maxConcurrentRegistrations}
                                                        required
                                                        type="text"
                                                        pattern="[0-9]*"
                                                        onChange={changeSimulationNumberHandler}
                                                        invalid={errors.maxConcurrentRegistrations ? true : false}
                                                        feedback="Max. Concurrent Device Registrations (whole number) is required"
                                                    />
                                                </Grid.Col>
                                                <Grid.Col
                                                    auto
                                                    className="align-self-center"
                                                >
                                                    <HelpPopup content={<>Maximum number of devices to be provisioned at a time.
                                                    Keep this number low as Device Provisioning Service (DPS) will throttle.</>} />
                                                </Grid.Col>
                                            </Grid.Row>
                                        </Form.Group>
                                        <Form.Group
                                            isRequired
                                            label="Max. Concurrent Device Deletes"
                                        >
                                            <Grid.Row gutters="xs">
                                                <Grid.Col>
                                                    <Form.Input
                                                        name="maxConcurrentDeletes"
                                                        value={config.simulation.maxConcurrentDeletes}
                                                        required
                                                        type="text"
                                                        pattern="[0-9]*"
                                                        onChange={changeSimulationNumberHandler}
                                                        invalid={errors.maxConcurrentDeletes ? true : false}
                                                        feedback="Max. Concurrent Device Deletes (whole number) is required"
                                                    />
                                                </Grid.Col>
                                                <Grid.Col
                                                    auto
                                                    className="align-self-center"
                                                >
                                                    <HelpPopup content={<>Maximum number of devices to be deleted at a time.
                                                    Devices are deleted when the simulation is deleted or the provisioned count in a simulation is decreased.</>} />
                                                </Grid.Col>
                                            </Grid.Row>
                                        </Form.Group>
                                        <Form.Group
                                            isRequired
                                            label="Max. Registration Attempts"
                                        >
                                            <Grid.Row gutters="xs">
                                                <Grid.Col>
                                                    <Form.Input
                                                        name="maxRegistrationAttempts"
                                                        value={config.simulation.maxRegistrationAttempts}
                                                        required
                                                        type="text"
                                                        pattern="[0-9]*"
                                                        onChange={changeSimulationNumberHandler}
                                                        invalid={errors.maxRegistrationAttempts ? true : false}
                                                        feedback="Max. Registration Attempts (whole number) is required"
                                                    />
                                                </Grid.Col>
                                                <Grid.Col
                                                    auto
                                                    className="align-self-center"
                                                >
                                                    <HelpPopup content={<>Maximum number of times a device can attempt for registration before erroring out.</>} />
                                                </Grid.Col>
                                            </Grid.Row>
                                        </Form.Group>
                                        <Form.Group>
                                            <Grid.Row gutters="xs">
                                                <Grid.Col>
                                                    <Form.Checkbox
                                                        name="enableTelemetry"
                                                        label="Send Telemetry"
                                                        checked={config.simulation.enableTelemetry}
                                                        onChange={changeSimCheckHandler}
                                                    />
                                                </Grid.Col>
                                                <Grid.Col
                                                    auto
                                                    className="align-self-center"
                                                >
                                                    <HelpPopup content={<><p>Should telemetry be sent from the device?</p>
                                                If this setting is turned off, no devices will send any telemetry.</>} />
                                                </Grid.Col>
                                            </Grid.Row>
                                        </Form.Group>
                                        <Form.Group>
                                            <Grid.Row gutters="xs">
                                                <Grid.Col>
                                                    <Form.Checkbox
                                                        name="enableReportedProps"
                                                        label="Send Reported Properties"
                                                        checked={config.simulation.enableReportedProps}
                                                        onChange={changeSimCheckHandler}
                                                    />
                                                </Grid.Col>
                                                <Grid.Col
                                                    auto
                                                    className="align-self-center"
                                                >
                                                    <HelpPopup content={<><p>Should reported properties be sent from the device?</p>
                                                If this setting is turned off, no devices will send any reported properties.</>} />
                                                </Grid.Col>
                                            </Grid.Row>
                                        </Form.Group>
                                        <Form.Group>
                                            <Grid.Row gutters="xs">
                                                <Grid.Col>
                                                    <Form.Checkbox
                                                        name="enableTwinUpdateAcks"
                                                        label="Send Desired Property Acknowledgements"
                                                        checked={config.simulation.enableTwinUpdateAcks}
                                                        onChange={changeSimCheckHandler}
                                                    />
                                                </Grid.Col>
                                                <Grid.Col
                                                    auto
                                                    className="align-self-center"
                                                >
                                                    <HelpPopup content={<><p>When the device receives a desired property update, should reported properties be sent from the device?</p>
                                                If this setting is turned off, no devices will send any reported property acknowledgements for desired property updates.</>} />
                                                </Grid.Col>
                                            </Grid.Row>
                                        </Form.Group>
                                        <Form.Group>
                                            <Grid.Row gutters="xs">
                                                <Grid.Col>
                                                    <Form.Checkbox
                                                        name="enableCommandAcks"
                                                        label="Send Command Acknowedgements"
                                                        checked={config.simulation.enableCommandAcks}
                                                        onChange={changeSimCheckHandler}
                                                    />
                                                </Grid.Col>
                                                <Grid.Col
                                                    auto
                                                    className="align-self-center"
                                                >
                                                    <HelpPopup content={<><p>Should commands be acknowledgements be sent from the device?</p>
                                                If this setting is turned off, no devices will send any acknowledgements for Direct methods or C2D commands.</>} />
                                                </Grid.Col>
                                            </Grid.Row>
                                        </Form.Group>
                                    </Grid.Col>
                                    <Grid.Col>
                                        <h6>
                                            <Text className="text-default small"><strong><Icon prefix="fe" name="alert-triangle" />{" "} Restart Starling to apply these changes.</strong></Text>
                                        </h6>
                                        <h4>HTTP Settings</h4>
                                        <Form.Group
                                            isRequired
                                            label="Administration Server Port Number"
                                        >
                                            <Grid.Row gutters="xs">
                                                <Grid.Col>
                                                    <Form.Input
                                                        name="adminPort"
                                                        value={config.http.adminPort}
                                                        required
                                                        type="text"
                                                        pattern="[0-9]*"
                                                        onChange={changeHttpNumberHandler}
                                                        invalid={errors.adminPort ? true : false}
                                                        feedback="Administration Server Port Number (whole number) is required"
                                                    />
                                                </Grid.Col>
                                                <Grid.Col
                                                    auto
                                                    className="align-self-center"
                                                >
                                                    <HelpPopup content={<>The port on which the Starling administration server is listening.</>} />
                                                </Grid.Col>
                                            </Grid.Row>
                                        </Form.Group>
                                        <Form.Group
                                            isRequired
                                            label="Metrics Port Number"
                                        >
                                            <Grid.Row gutters="xs">
                                                <Grid.Col>
                                                    <Form.Input
                                                        name="metricsPort"
                                                        value={config.http.metricsPort}
                                                        required
                                                        type="text"
                                                        pattern="[0-9]*"
                                                        onChange={changeHttpNumberHandler}
                                                        invalid={errors.metricsPort ? true : false}
                                                        feedback="Metrics Port Number (whole number) is required"
                                                    />
                                                </Grid.Col>
                                                <Grid.Col
                                                    auto
                                                    className="align-self-center"
                                                >
                                                    <HelpPopup content={<>The port on which the Starling publishes Prometheus metrics.</>} />
                                                </Grid.Col>
                                            </Grid.Row>
                                        </Form.Group>
                                        <Form.Group
                                            isRequired
                                            label="Prometheus Port Number"
                                        >
                                            <Grid.Row gutters="xs">
                                                <Grid.Col>
                                                    <Form.Input
                                                        name="prometheusPort"
                                                        value={config.http.prometheusPort}
                                                        required
                                                        type="text"
                                                        pattern="[0-9]*"
                                                        onChange={changeHttpNumberHandler}
                                                        invalid={errors.prometheusPort ? true : false}
                                                        feedback="Prometheus Port Number (whole number) is required"
                                                    />
                                                </Grid.Col>
                                                <Grid.Col
                                                    auto
                                                    className="align-self-center"
                                                >
                                                    <HelpPopup content={<>The port on which the Promethus server is listening.</>} />
                                                </Grid.Col>
                                            </Grid.Row>
                                        </Form.Group>
                                        <Form.Group
                                            isRequired
                                            label="Grafana Port Number"
                                        >
                                            <Grid.Row gutters="xs">
                                                <Grid.Col>
                                                    <Form.Input
                                                        name="grafanaPort"
                                                        value={config.http.grafanaPort}
                                                        required
                                                        type="text"
                                                        pattern="[0-9]*"
                                                        onChange={changeHttpNumberHandler}
                                                        invalid={errors.grafanaPort ? true : false}
                                                        feedback="Grafana Port Number (whole number) is required"
                                                    />
                                                </Grid.Col>
                                                <Grid.Col
                                                    auto
                                                    className="align-self-center"
                                                >
                                                    <HelpPopup content={<>The port on which the Grafana dashboard server is listening.</>} />
                                                </Grid.Col>
                                            </Grid.Row>
                                        </Form.Group>

                                        <h4>Data Settings</h4>
                                        <Form.Group
                                            isRequired
                                            label="Database Directory"
                                        >
                                            <Grid.Row gutters="xs">
                                                <Grid.Col>
                                                    <Form.Input
                                                        name="path"
                                                        value={config.data.path}
                                                        required
                                                        type="text"
                                                        pattern="[0-9]*"
                                                        onChange={changeDataHandler}
                                                        invalid={errors.adminPort ? true : false}
                                                        feedback="Database Directory is required"
                                                    />
                                                </Grid.Col>
                                                <Grid.Col
                                                    auto
                                                    className="align-self-center"
                                                >
                                                    <HelpPopup content={<><p>The directory in which Starling stores its database.</p>
                                                        <p>By default it stores data in a <strong>[Starling binary directory]/.db</strong> directory.</p>
                                                Use forward / instead of \ in paths. E.g.: d:/starling/data
                                                </>} />
                                                </Grid.Col>
                                            </Grid.Row>
                                        </Form.Group>

                                        <h4>Logger Settings</h4>
                                        <Form.Group
                                            isRequired
                                            label="Logs Directory"
                                        >
                                            <Grid.Row gutters="xs">
                                                <Grid.Col>
                                                    <Form.Input
                                                        name="logsDir"
                                                        value={config.logger.logsDir}
                                                        required
                                                        type="text"
                                                        pattern="[0-9]*"
                                                        onChange={changeLogHandler}
                                                        invalid={errors.adminPort ? true : false}
                                                        feedback="Logs Directory is required"
                                                    />
                                                </Grid.Col>
                                                <Grid.Col
                                                    auto
                                                    className="align-self-center"
                                                >
                                                    <HelpPopup content={<><p>The directory in which Starling stores log files.</p>
                                                        <p>Log files are automatically rotated every <strong>30 days</strong> or when it reaches <strong>10MB</strong>. Last <strong>3</strong> files are kept.</p>
                                                        <p>By default it stores data in a <strong>[Starling binary directory]/logs</strong> directory.</p>
                                                Use forward / instead of \ in paths. E.g.: d:/starling/logs.
                                                </>} />
                                                </Grid.Col>
                                            </Grid.Row>
                                        </Form.Group>
                                        <Form.Group
                                            isRequired
                                            label="Log Level"
                                        >
                                            <Grid.Row gutters="xs">
                                                <Grid.Col>
                                                    <Form.Select
                                                        name="logLevel"
                                                        value={config.logger.logLevel}
                                                        required
                                                        onChange={changeLogHandler}
                                                        invalid={errors.logLevel ? true : false}
                                                        feedback="Log Level is required"
                                                    >
                                                        <option value="panic">Panic</option>
                                                        <option value="fatal">Fatal</option>
                                                        <option value="error">Error</option>
                                                        <option value="warn">Warning</option>
                                                        <option value="info">Information</option>
                                                        <option value="debug">Debug</option>
                                                        <option value="trace">Trace</option>
                                                    </Form.Select>
                                                </Grid.Col>
                                                <Grid.Col
                                                    auto
                                                    className="align-self-center"
                                                >
                                                    <HelpPopup content={<><p>Starling produces runtime logs stored in its <strong>Log Directory</strong>. Select the appropriate log level.</p>
                                                Leave this at <strong>Debug</strong> level as it produces useful information.</>} />
                                                </Grid.Col>
                                            </Grid.Row>
                                        </Form.Group>

                                        <h4>Data Generation</h4>
                                        <Form.Group
                                            isRequired
                                            label="Geopoint Data"
                                        >
                                            <Grid.Row gutters="xs">
                                                <Grid.Col>
                                                    <Form.Textarea
                                                        name="geopointData"
                                                        value={config.simulation.geopointData}
                                                        required
                                                        onChange={changeSimHandler}
                                                        className="geopointData"
                                                        invalid={errors.geopointData ? true : false}
                                                        feedback="Geopoint Data is required"
                                                        rows={15}
                                                    />
                                                </Grid.Col>
                                                <Grid.Col
                                                    auto
                                                    className="align-self-center"
                                                >
                                                    <HelpPopup content={<><p>These Geopoints data is used during simulating geopoint data types. All the points are sequentially used while sending telemetry messages containing geopoints. </p>
                                                        <p>Format is a JSON array with three floating point numbers: <strong>[latitude, longitude, altitude]</strong>.</p>
                                                        Default data is for a road around part of Microsoft campus in Redmond, WA, USA.
                                                </>} />
                                                </Grid.Col>
                                            </Grid.Row>
                                        </Form.Group>
                                    </Grid.Col>
                                </Grid.Row>
                            </Form.FieldSet>
                        </Card.Body>
                    </Card>
                </form>
            }
        </Page.Content>
    </SiteWrapper>;
}

export default SettingsPage;