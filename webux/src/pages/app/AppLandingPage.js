import { useContext } from 'react';
import {
    Page
} from "tabler-react";
import "tabler-react/dist/Tabler.css";
import { Redirect } from "react-router-dom";
import GlobalContext from '../../context/globalContext';
import SiteWrapper from '../../components/site/SiteWrapper';
import NoDataFoundCard from '../../components/nodata/NoDataFoundCard';

const AppLandingPage = () => {
    const globalContext = useContext(GlobalContext)
    const appCount = globalContext.apps ? globalContext.apps.length : 0;
    const redir = appCount > 0 ? `/app/${globalContext.apps[0].id}` : "";

    const actions = [
        {
            actionName: "Add Application",
            actionUrl: "/app/add?new",
            actionIcon: "plus"
        }
    ];
    return appCount > 0 ? <Redirect to={redir} /> :
        <SiteWrapper>
            <Page.Content title="">
                <NoDataFoundCard
                    message="IoT Central Applications"
                    description="IoT Central Applications will show up here. Simulated devices are created in these IoT Central applications."
                    actions={actions}
                    noDataImage="/images/emptyApps.svg"
                />
            </Page.Content>
        </SiteWrapper>
        ;
}

export default AppLandingPage;