import { useContext, useEffect } from 'react';
import {
    Page
} from "tabler-react";
import "tabler-react/dist/Tabler.css";
import { Redirect } from "react-router-dom";
import GlobalContext from '../../context/globalContext';
import SiteWrapper from '../../components/site/SiteWrapper';
import NoDataFoundCard from '../../components/nodata/NoDataFoundCard';

const ModelLandingPage = () => {
    const globalContext = useContext(GlobalContext)

    // Called on mount to ensure reference data is loaded if coming from shortcut
    useEffect(() => {
        if (!globalContext.initialized) {
            globalContext.initializeData();
        }
    }, [globalContext])

    const modelCount = globalContext.models ? globalContext.models.length : 0;
    const redir = modelCount > 0 ? `/model/${globalContext.models[0].id}` : "";
    //console.log("modelCount: ", modelCount);

    return modelCount > 0 ? <Redirect to={redir} /> :
        <SiteWrapper>
            <Page.Content title="">
                <NoDataFoundCard
                    message="Device models"
                    description="Device models will show up here. Simulated devices can be spun up based on these device models. Add a device model here."
                    actionName="Add Device Model"
                    actionUrl="/model/add?new"
                    actionIcon="plus"
                    noDataImage="/images/emptyModels.svg"
                />
            </Page.Content>
        </SiteWrapper>
        ;
}

export default ModelLandingPage;