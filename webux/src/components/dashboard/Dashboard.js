import { useContext, useEffect, useState } from 'react';
import { Link } from "react-router-dom";
import {
    Icon,
    Page
} from "tabler-react";
import "tabler-react/dist/Tabler.css";
import SimulationCard from './SimulationCard';
import Intro from '../intro/Intro';
import Toolbar from "../toolbar/Toolbar";
import GlobalContext from '../../context/globalContext';
import * as Utils from '../../utils/utils';
import "./Dashboard.css";

const Dashboard = () => {
    const globalContext = useContext(GlobalContext)
    const [backendError, setBackendError] = useState()

    useEffect(() => {
        let timer = setInterval(() => {
            if (!globalContext.initialized) {
                globalContext.initializeData();
            }

            const refreshDashboard = async () => {
                try {
                    await globalContext.listSimulations();
                    setBackendError("");
                } catch (err) {
                    let msg = Utils.getErrorMessage(err, "error listing simulations")
                    if (msg === "Network Error") {
                        msg += ". Make sure that the Starling server is running.";
                    }
                    setBackendError(msg);
                }
            };
            refreshDashboard();
        }, 5000);
        return () => {
            clearInterval(timer);
        }

        // ignore global context dependency error
        // eslint-disable-next-line react-hooks/exhaustive-deps
    });

    // Called on mount to ensure reference data is loaded if coming from shortcut
    useEffect(() => {
        try {
            if (!globalContext.initialized) {
                globalContext.initializeData();
            }
            if (!globalContext.config) {
                setBackendError("Make sure that the Starling server is running.");
            } else {
                setBackendError("");
            }
        } catch (err) {
            let msg = Utils.getErrorMessage(err, "error listing simulations")
            if (msg === "Network Error") {
                msg += ". Make sure that the Starling server is running.";
            }
            setBackendError(msg);
        }
    }, [globalContext]);

    const sims = globalContext.simulations;
    const cards = sims.map(sim =>
        <SimulationCard
            key={sim.id}
            sim={sim} />
    );

    if (globalContext.initialized && sims.length === 0) {
        return (
            <Page.Content title="Welcome to Starling!">
                <Intro />
            </Page.Content>
        );
    }

    const simCount = Utils.formatCount(sims, "simulation");
    return (
        <Page.Content title="Dashboard">
            {backendError && backendError.length > 0 &&
                <div className="alert alert-danger">
                    <Icon prefix="fe" name="alert-triangle" />{" "}
                    {backendError}
                </div>
            }
            {
                backendError && backendError.length === 0 &&
                <Toolbar countMessage={simCount}>
                    <span title="Create a new simulation">
                        <Link
                            to="/sim/add?new"
                            className="btn btn-sm btn-primary"
                        >
                            <Icon prefix="fe" name="plus" />
                    New Simulation
                    </Link>
                    </span>
                </Toolbar>
            }
            {
                sims && sims.length > 0 &&
                <div className="dashboard">
                    {cards}
                </div>
            }
        </Page.Content>
    );
};

export default Dashboard;