import { useContext } from 'react';
import {
    Page
} from "tabler-react";
import "tabler-react/dist/Tabler.css";
import { Redirect } from "react-router-dom";
import GlobalContext from '../../context/globalContext';
import SiteWrapper from '../../components/site/SiteWrapper';
import NoDataFoundCard from '../../components/nodata/NoDataFoundCard';

const SimLandingPage = () => {
    const globalContext = useContext(GlobalContext)
    const simCount = globalContext.simulations ? globalContext.simulations.length : 0;
    const redir = simCount > 0 ? `/sim/${globalContext.simulations[0].id}` : "";

    let actionName = "Add Simulation";
    let actionUrl = "/sim/add?new";
    let description = "";
    if (!globalContext.models || globalContext.models.length === 0) {
        actionName = "Add Device Model";
        actionUrl = "/model/add?new";
        description = "You need to add a Device Model before you create a new Simulation.";
    } else if (!globalContext.apps || globalContext.apps.length === 0) {
        actionName = "Add Application";
        actionUrl = "/app/add?new";
        description = "You need to add an IoT Central Application before you create a new Simulation.";
    }

    const actions = [
        {
            actionName: actionName,
            actionUrl: actionUrl,
            actionIcon: "plus"
        }
    ];
    return simCount > 0 ? <Redirect to={redir} /> :
        <SiteWrapper>
            <Page.Content title="">
                <NoDataFoundCard
                    message="Simulations"
                    description="Simulations will show up here. Multiple simulations can be concurrently executed against IoT Central applications."
                    description2={description}
                    actions={actions}
                    noDataImage="/images/emptySimulations.svg"
                />
            </Page.Content>
        </SiteWrapper>
        ;
}

export default SimLandingPage;