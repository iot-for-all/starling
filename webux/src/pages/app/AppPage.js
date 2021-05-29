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
import AppCard from './AppCard';
import * as Utils from "../../utils/utils";

const AppPage = () => {
    const globalContext = useContext(GlobalContext)

    const params = useParams();
    const queryParams = new URLSearchParams(useLocation().search);
    const [app, setApp] = useState();
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

        if (pageMode === 'add') {
            setApp({ id: '', name: '', provisioningUrl: 'global.azure-devices-provisioning.net', idScope: '', masterKey: '', appUrl: '', appToken: '', importModels: true });
        } else {
            if (!id) {
                setApp(globalContext.apps[0]);
            } else {
                const remoteApp = globalContext.getApplication(id);
                if (!remoteApp) {
                    // appId passed in url is invalid or does not exist
                    history.replace("/app");
                }
                setApp({ ...remoteApp });
            }
        }
        setBackendError("");

        // ignore global context dependency error
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [globalContext.apps, id])

    const appCount = Utils.formatCount(globalContext.apps, "application");
    const appsList = app ? globalContext.apps.map((element) => {
        return (
            <List.GroupItem
                className="d-flex align-items-center"
                to={"/app/" + element.id}
                RootComponent={Link}
                key={element.id}
                active={element.id === app?.id}
                action
            >
                {element.name}
            </List.GroupItem>
        );
    }) : "";

    const appCard = app ? <AppCard
        data={app}
        mode={pageMode}
        backendError={backendError}
    /> : "";

    return (
        <SiteWrapper>
            <Page.Content title="IoT Central Applications">
                <Toolbar
                    countMessage={appCount}
                >
                    <span title="Add an IoT Central Application">
                        <Link
                            to="/app/add?new"
                            className="btn btn-sm btn-primary"
                        >
                            <Icon prefix="fe" name="plus" />
                            Add New
                        </Link>
                    </span>
                </Toolbar>
                <ListDetails
                    listTitle={"Applications"}
                    list={appsList}
                    detailsTitle={"Edit Application"}
                    detailsForm={appCard}
                    backendError={backendError}
                />
            </Page.Content>
        </SiteWrapper>
    );
}

export default AppPage;