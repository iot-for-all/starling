import { useContext, useState } from 'react';
import { useHistory } from "react-router-dom";
import {
    Button,
    Card,
    Form,
    Grid,
    Icon,
} from "tabler-react";
import GlobalContext from '../../context/globalContext';
import * as Utils from '../../utils/utils';
import HelpPopup from "../../components/help/HelpPopup";
import * as Notification from "../../components/notification/Notification";

const ImportCard = (props) => {
    const globalContext = useContext(GlobalContext)
    const [appId, setAppId] = useState("");
    const [errors, setErrors] = useState({});
    const [backendError, setBackendError] = useState("");
    const history = useHistory();
    
    const changeHandler = (event) => {
        setAppId(event.target.value);
        //console.log("value changed target: ", event.target, " value: ", event.target.value, ", updatedSim: ", updatedSim);
    }

    const onSubmit = async (event) => {
        event.preventDefault();

        // validate form
        let newErrors = {};
        let hasErrors = false;
        
        if (appId.trim() === "") {
            newErrors.appId = true;
            hasErrors = true;
        }
        setErrors(newErrors);

        if (!hasErrors) {
            //console.log("Form submitted mode: ", props.mode, " sim: ", sim);
            try {
                Notification.addNotification("info", "Started", `Importing devices models from '${appId}'.`);
                await globalContext.importModels(appId);
                Notification.addNotification("success", "Success", `Devices models are imported from '${appId}'.`);
    
                history.push("/model");
            } catch (ex) {
                setBackendError(Utils.getErrorMessage(ex, "error importing device models"));
            }
        }
    }

    const appsArr = globalContext.apps ? Array.from(globalContext.apps) : [];
    const appsList = appsArr.map((app) => {
        return (
            <option key={app.id} value={app.id}>{app.name}</option>
        );
    });
    return <>
        <div>{props.backendError}</div>
        <form onSubmit={onSubmit} method="post">
            <Card>
                <Card.Header>
                    <Card.Title>Import Device Models</Card.Title>
                    <Card.Options>
                        <span title="Import Device Models">
                            <Button color="primary" size="sm" icon="download" className="ml-2" type="submit">Import</Button>
                        </span>
                        <span title="Cancel all changes">
                            <Button color="primary" size="sm" className="ml-2" outline type="button"
                                onClick={() => history.push("/model")} >Cancel</Button>
                        </span>
                    </Card.Options>
                </Card.Header>
                <Card.Body>
                    {backendError && backendError.length > 0 && <div className="alert alert-danger">
                        <Icon prefix="fe" name="alert-triangle" />{" "}
                        {backendError}
                    </div>}
                    <p>
                        Import all published device models from the application below.
                    </p>
                    <Form.FieldSet>
                        <Grid.Row>
                            <Grid.Col colSpan="2">
                                <Form.Group
                                    isRequired
                                    label="Application to import device models"
                                >
                                    <Grid.Row gutters="xs">
                                        <Grid.Col>
                                            <Form.Select
                                                name="appId"
                                                value={appId}
                                                required
                                                onChange={changeHandler}
                                                invalid={errors.appId ? true : false}
                                                feedback="Application is required"
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
                                            <HelpPopup content={<>The IoT Central application from which the device models are imported.</>} />
                                        </Grid.Col>
                                    </Grid.Row>
                                </Form.Group>

                            </Grid.Col>
                        </Grid.Row>
                    </Form.FieldSet>
                </Card.Body>
            </Card>
        </form>
    </>
}

export default ImportCard;
