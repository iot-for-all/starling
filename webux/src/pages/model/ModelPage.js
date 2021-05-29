import { useContext, useEffect, useState } from 'react';
import {
    Icon,
    List,
    Page
} from "tabler-react";
import "tabler-react/dist/Tabler.css";
import { Link, useHistory, useLocation, useParams } from "react-router-dom";
import GlobalContext from '../../context/globalContext';
import Toolbar from '../../components/toolbar/Toolbar';
import SiteWrapper from '../../components/site/SiteWrapper';
import ListDetails from '../../components/listdetails/ListDetails';
import ModelCard from './ModelCard';
import * as Utils from "../../utils/utils";

const ModelPage = () => {
    const globalContext = useContext(GlobalContext)

    const params = useParams();
    const queryParams = new URLSearchParams(useLocation().search);
    const [model, setModel] = useState();
    const [backendError, setBackendError] = useState("");
    const id = params.id;
    const history = useHistory();
    const pageMode = (queryParams.has("new")) ? "add" : "edit";

    // Called on mount to ensure reference data is loaded if coming from shortcut
    useEffect(() => {
        if (!globalContext.initialized) {
            globalContext.initializeData();
        }

        //console.log("in global useeffect, id: ", id);
    }, [globalContext])

    // Only called if a new mount or id has changed
    useEffect(() => {
        if (!globalContext.initialized) {
            globalContext.initializeData();
        }
        //console.log("ModelPage - useeffect, id: ", id);
        if (pageMode === 'add') {
            setModel({ id: '', name: '', capabilityModel: [] });
        } else {
            if (!id) {
                setModel(globalContext.models[0]);
            } else {
                const remoteModel = globalContext.getModel(id);
                if (!remoteModel) {
                    // modelId passed in url is invalid or does not exist
                    history.replace("/model");
                }
                setModel({ ...remoteModel });
            }
        }
        setBackendError("");

        // ignore global context dependency error
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [globalContext.models, id])

    const modelCount = Utils.formatCount(globalContext.models, "device model");
    //console.log("ModelPage - id: ", id, ", pageMode: ", pageMode, " backendError: ", backendError, " model:", model);

    const modelsList = model ? globalContext.models.map((element) => {
        return (
            <List.GroupItem
                className="d-flex align-items-center"
                to={"/model/" + element.id}
                RootComponent={Link}
                key={element.id}
                active={element.id === model?.id}
                action
            >
                {element.name}
            </List.GroupItem>
        );
    }) : "";

    const modelCard = model ? <ModelCard
        data={model}
        mode={pageMode}
        backendError={backendError}
    /> : "";

    return (
        <SiteWrapper>
            <Page.Content title="Device Models">
                <Toolbar
                    countMessage={modelCount}
                >
                    <span title="Add a Device Model">
                        <Link
                            to="/model/add?new"
                            className="btn btn-sm btn-primary"
                        >
                            <Icon prefix="fe" name="plus" />
                            Add New
                        </Link>
                    </span>
                </Toolbar>
                <ListDetails
                    listTitle={"Models"}
                    list={modelsList}
                    detailsTitle={"Edit Device Model"}
                    detailsForm={modelCard}
                    backendError={backendError}
                />
            </Page.Content>
        </SiteWrapper>
    );
}

export default ModelPage;