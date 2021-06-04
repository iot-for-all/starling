import { useContext, useEffect, useState } from 'react';
import { useHistory } from "react-router-dom";
import {
    Button,
    Card,
    Grid,
    Header,
    Icon,
    Page,
    Text
} from "tabler-react";
import "tabler-react/dist/Tabler.css";
import GlobalContext from '../../context/globalContext';
import SiteWrapper from '../../components/site/SiteWrapper';
import * as Utils from '../../utils/utils';
import "./MetricsPage.css"
import * as Notification from "../../components/notification/Notification";

const MetricsPage = () => {
    const globalContext = useContext(GlobalContext)
    const history = useHistory();
    const [status, setStatus] = useState({ grafanaStatus: false, prometheusStatus: false });
    const [backendError, setBackendError] = useState("");

    // Called on mount to ensure reference data is loaded if coming from shortcut
    useEffect(() => {
        if (!globalContext.initialized) {
            globalContext.initializeData();
        }

        if (globalContext.metricsStatus) {
            setStatus(globalContext.metricsStatus);
            setBackendError("");
        }
        // ignore global context dependency error
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [globalContext.metricsStatus])

    const onSubmit = async (event) => {
        event.preventDefault();

        // validate form
        //console.log("Form submitted mode: ", props.mode, " sim: ", sim);
        try {
            const stats = await globalContext.refreshMetricsStatus();
            Notification.addNotification("success", "Success", `Metrics status refreshed.`);
            setStatus(stats)
            history.replace(`/metrics`);
        } catch (ex) {
            setBackendError(Utils.getErrorMessage(ex, "error refreshing metrics"));
        }
    }

    const starlingStatus = <span className={backendError.length === 0 ? "text-green" : "text-danger"}>
        {backendError.length === 0 ? "Healthy" : "Not running"}
    </span>;
    const prometheusStatus = <span className={status.prometheusServer ? "text-green" : "text-danger"}>
        {status.prometheusServer ? "Healthy" : "Not running"}
    </span>;
    const grafanaStatus = <span className={status.grafanaServer ? "text-green" : "text-danger"}>
        {status.grafanaServer ? "Healthy" : "Not running"}
    </span>;

    let prometheusLink = "";
    let grafanaLink = "";

    const metricsUrl = (globalContext.config) ? "http://localhost:" + globalContext.config.http.metricsPort + "/metrics" : "#";
    const starlingMetricsLink = <div>
        <Icon prefix="fe" name="external-link" />{" "}<a href={metricsUrl} target="_blank" rel="noreferrer">Raw Metrics</a>
    </div>;

    if (status.prometheusServer) {
        let url = "#";
        if (globalContext.config) {
            url = "http://localhost:" + globalContext.config.http.prometheusPort;
        }
        prometheusLink = <div>
            <Icon prefix="fe" name="external-link" />{" "}<a href={url} target="_blank" rel="noreferrer">Dashboard</a>
        </div>;
    } else {
        prometheusLink = <div>
            <Icon prefix="fe" name="external-link" />{" "}<a href="https://github.com/iot-for-all/starling/blob/main/README.md" target="_blank" rel="noreferrer">Help me configure</a>
        </div>;
    }

    if (status.grafanaServer) {
        let url = "#";
        if (globalContext.config) {
            url = "http://localhost:" + globalContext.config.http.grafanaPort;
        }
        grafanaLink = <div>
            <Icon prefix="fe" name="external-link" />{" "}<a href={url} target="_blank" rel="noreferrer">Dashboard</a>
        </div>;
    } else {
        grafanaLink = <div>
            <Icon prefix="fe" name="external-link" />{" "}<a href="https://github.com/iot-for-all/starling/blob/main/README.md" target="_blank" rel="noreferrer">Help me configure</a>
        </div>;
    }

    return <SiteWrapper>
        <Page.Content title="Metrics">
            {backendError && backendError.length > 0 && <div className="alert alert-danger">
                <Icon prefix="fe" name="alert-triangle" />{" "}
                {backendError}
            </div>}
            <form onSubmit={onSubmit}>
                <Card>
                    <Card.Header>
                        <Card.Title>Data Flow</Card.Title>
                        <Card.Options>
                            <span title="Refresh Metrics Status">
                                <Button
                                    color="primary"
                                    size="sm"
                                    icon="refresh-ccw"
                                    className="ml-2"
                                    onClick={onSubmit}
                                >Refresh</Button>
                            </span>
                        </Card.Options>
                    </Card.Header>
                    <Card.Body>
                        <div className="simLearnMore">
                            <p>
                                <Text className="text-default"><Icon prefix="fe" name="help-circle" />{" "} <a href="https://github.com/iot-for-all/starling" target="_blank" rel="noreferrer">Help me configure metrics pipeline</a></Text>
                            </p>
                        </div>
                        <p>
                            Starling publishes metrics through its metrics endpoint.
                            You can configure Prometheus to scrape these metrics and store in its timeseries database.
                            Grafana can be configured to show dashboards based on the data stored in Prometheus.
                        </p>
                        <div className="serverCardContainer">
                            <div className="serverCardItem">
                                <Card className="serverCard">
                                    <Card.Body className="text-center">
                                        <Header size={2}>Starling</Header>
                                        <Header size={5}>Status: {starlingStatus}</Header>
                                        {starlingMetricsLink}
                                    </Card.Body>
                                </Card>
                            </div>
                            <div className="serverCardArrow">
                                <div>
                                    <Text className="arrowText">Prometheus scrapes data from Starling</Text>
                                </div>
                                <span class="arrow-right"></span>
                            </div>
                            <div className="serverCardItem">
                                <Card className="serverCard">
                                    <Card.Body className="text-center">
                                        <img className="card-img-top" src="./images/prometheusLogo.png" alt="Prometheus" />
                                        <Header size={5}>Status: {prometheusStatus}</Header>
                                        {prometheusLink}
                                    </Card.Body>
                                </Card>
                            </div>
                            <div className="serverCardArrow">
                                <div>
                                    <Text className="arrowText">Grafana queries data from Prometheus</Text>
                                </div>
                                <span class="arrow-right"></span>
                            </div>
                            <div className="serverCardItem">
                                <Card className="serverCard">
                                    <Card.Body className="text-center">
                                        <img className="card-img-top" src="./images/grafanaLogo.jpg" alt="Grafana" />
                                        <Header size={5}>Status: {grafanaStatus}</Header>
                                        {grafanaLink}
                                    </Card.Body>
                                </Card>
                            </div>
                        </div>
                    </Card.Body>
                </Card>
            </form>
        </Page.Content>
    </SiteWrapper>;
}

export default MetricsPage;