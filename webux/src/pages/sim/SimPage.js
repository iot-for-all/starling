import { useContext, useEffect, useState } from 'react';
import {
    Icon,
    List,
    Page
} from "tabler-react";
import "tabler-react/dist/Tabler.css";
import { Link, Redirect, useHistory, useLocation, useParams } from "react-router-dom";
import GlobalContext from '../../context/globalContext';
import Toolbar from '../../components/toolbar/Toolbar';
import SiteWrapper from '../../components/site/SiteWrapper';
import ListDetails from '../../components/listdetails/ListDetails';
import SimCard from './SimCard';
import ProvisionCard from './ProvisionCard';
import * as Utils from "../../utils/utils";

const SimPage = () => {
    const globalContext = useContext(GlobalContext)

    const params = useParams();
    const queryParams = new URLSearchParams(useLocation().search);
    const [sim, setSim] = useState();
    const [backendError, setBackendError] = useState("");
    const id = params.id;
    const history = useHistory();
    const pageMode = (queryParams.has("new")) ? "add" : "edit";
    const isProvisionMode = (queryParams.has("provision"));

    // Called on mount to ensure reference data is loaded if coming from shortcut
    
    // Only called if a new mount or id has changed
    useEffect(() => {
        if (pageMode === 'add') {
            let devices = [];
            for (let i = 0; i < globalContext.models.length; i++) {
                devices.push({ id: globalContext.models[i].id, modelId: globalContext.models[i].id, provisionedCount: 0, simulatedCount: 0, connectedCount: 0 });
            }

            setSim({
                id: "",
                name: "",
                targetId: "",
                status: "ready",
                waveGroupCount: 1,
                waveGroupInterval: 1,
                telemetryBatchSize: 1,
                telemetryInterval: 120,
                reportedPropertyInterval: 14400,
                disconnectBehavior: "never",
                telemetryFormat: "default",
                devices: devices
            });
        } else {
            if (!id) {
                setSim(globalContext.simulations[0]);
            } else {
                const remoteSim = globalContext.getSimulation(id);
                if (!remoteSim) {
                    // simId passed in url is invalid or does not exist
                    history.replace("/sim");
                }

                setSim({ ...remoteSim });
            }
        }
        setBackendError("");

        // ignore global context dependency error
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [globalContext.simulations, globalContext.initialized, id])

    const simCount = Utils.formatCount(globalContext.simulations, "simulation");

    // simId passed in url is invalid or does not exist
    if (!globalContext.models || globalContext.models.length === 0
        || !globalContext.apps || globalContext.apps.length === 0
    ) {
        return <Redirect to={'/sim'} />
    }

    const simsList = sim ? globalContext.simulations.map((element) => {
        return (
            <List.GroupItem
                className="d-flex align-items-center"
                to={"/sim/" + element.id}
                RootComponent={Link}
                key={element.id}
                active={element.id === sim?.id}
                action
            >
                {element.name}
            </List.GroupItem>
        );
    }) : "";

    const simCard = sim ? <SimCard
        data={sim}
        mode={pageMode}
        backendError={backendError}
    /> : "";

    const provisionCard = sim ? <ProvisionCard
        data={sim}
        backendError={backendError}
    /> : "";

    return (
        <SiteWrapper>
            <Page.Content title="Simulations">
                <Toolbar
                    countMessage={simCount}
                >
                    <span title="Add a Simulation">
                        <Link
                            to="/sim/add?new"
                            className="btn btn-sm btn-primary"
                        >
                            <Icon prefix="fe" name="plus" />
                            Add New
                        </Link>
                    </span>
                </Toolbar>
                <ListDetails
                    listTitle={"Simulations"}
                    list={simsList}
                    detailsTitle={"Edit Application"}
                    detailsForm={isProvisionMode ? provisionCard : simCard}
                    backendError={backendError}
                />
            </Page.Content>
        </SiteWrapper>
    );
}

export default SimPage;