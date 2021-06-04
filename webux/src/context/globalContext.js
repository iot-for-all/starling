import React, { useState } from 'react';
import axios from 'axios';

const GlobalContext = React.createContext({
    initialized: false,
    apps: [],
    models: [],
    simulations: [],
    config: {},
    metricsStatus: {},
    initializeData: () => { },
    getApplication: (appId) => { },
    listApplications: () => { },
    addApplication: (payload) => { },
    updateApplication: (payload) => { },
    deleteApplication: (appId) => { },
    getModel: (modelId) => { },
    listModels: () => { },
    addModel: (payload) => { },
    updateModel: (payload) => { },
    deleteModel: (modelId) => { },
    importModels: (appId) => { },
    getSimulation: (simId) => { },
    listSimulations: () => { },
    addSimulation: (payload) => { },
    updateSimulation: (payload) => { },
    deleteSimulation: (simId) => { },
    startSimulation: (simId) => { },
    stopSimulation: (simId) => { },
    exportSimulation: (simId) => { },
    provisionSimulationDevices: (simId, payload) => { },
    getConfig: () => { },
    updateConfig: (payload) => { },
    refreshMetricsStatus: () => {},
});

export const GlobalContextProvider = (props) => {
    const [initialized, setInitialized] = useState(false);
    const [apps, setApps] = useState([]);
    const [models, setModels] = useState([]);
    const [simulations, setSimulations] = useState([]);
    const [config, setConfig] = useState();
    const [metricsStatus, setMetricsStatus] = useState();
    const BASE_URL = (process.env.NODE_ENV !== "production") ? "http://localhost:6001/webapi" : "/webapi";

    const initializeData = async () => {
        //try {
            const remoteApps = await axios.get(`${BASE_URL}/target`);
            const remoteModels = await axios.get(`${BASE_URL}/model`);
            const remoteSimulations = await axios.get(`${BASE_URL}/simulation`);
            const remoteConfig = await axios.get(`${BASE_URL}/config`);
            const remoteMetricsStatus = await axios.get(`${BASE_URL}/config/metricsStatus`);
            setApps(remoteApps.data);
            setModels(remoteModels.data);
            setSimulations(remoteSimulations.data);
            setConfig(remoteConfig.data);
            setMetricsStatus(remoteMetricsStatus.data);
            setInitialized(true);
        /*} catch (err) {
            console.log(err);
        }*/
    }

    const getApplication = (appId) => {
        return apps.find((x) => x.id === appId);
    }

    const listApplications = async () => {
        const res = await axios.get(`${BASE_URL}/target`);
        setApps(res.data);
        return res.data;
    }

    const addApplication = async (payload) => {
        await axios.post(`${BASE_URL}/target`, payload);
        await listApplications();
        await listModels();
    }

    const updateApplication = async (payload) => {
        await axios.put(`${BASE_URL}/target`, payload);
        await listApplications();
    }

    const deleteApplication = async (appId) => {
        await axios.delete(`${BASE_URL}/target/${appId}`);
        await listApplications();
    }

    const getModel = (modelId) => {
        return models.find((x) => x.id === modelId);
        //const res = await axios.get(`${BASE_URL}/model/${modelId}`);
        //return res.data;
    }

    const listModels = async () => {
        const res = await axios.get(`${BASE_URL}/model`);
        setModels(res.data);
        return res.data;
    }

    const addModel = async (payload) => {
        await axios.post(`${BASE_URL}/model`, payload);
        await listModels();
        await listSimulations();
    }

    const updateModel = async (payload) => {
        await axios.put(`${BASE_URL}/model`, payload);
        await listModels();
        await listSimulations();
    }

    const deleteModel = async (modelId) => {
        await axios.delete(`${BASE_URL}/model/${modelId}`);
        await listModels();
        await listSimulations();
    }

    const importModels = async (appId) => {
        await axios.post(`${BASE_URL}/target/${appId}/import`);
        await listModels();
        await listSimulations();
    }

    const getSimulation = (simId) => {
        return simulations.find((x) => x.id === simId);
    }

    const listSimulations = async () => {
        const res = await axios.get(`${BASE_URL}/simulation`);
        setSimulations(res.data);
        return res.data;
    }

    const addSimulation = async (payload) => {
        await axios.post(`${BASE_URL}/simulation`, payload);
        await listSimulations();
    }

    const updateSimulation = async (payload) => {
        await axios.put(`${BASE_URL}/simulation`, payload);
        await listSimulations();
    }

    const deleteSimulation = async (simId) => {
        await axios.delete(`${BASE_URL}/simulation/${simId}`);

        await listSimulations();
    }

    const startSimulation = async (simId) => {
        await axios.post(`${BASE_URL}/simulation/${simId}/start`);
        await listSimulations();
    }

    const stopSimulation = async (simId) => {
        await axios.post(`${BASE_URL}/simulation/${simId}/stop`);
        await listSimulations();
    }

    const exportSimulation = async (simId) => {
        const res = await axios.get(`${BASE_URL}/simulation/${simId}/export`);
        const downloadUrl = window.URL.createObjectURL(new Blob([res.data]));
        const link = document.createElement('a');
        link.href = downloadUrl;
        link.setAttribute('download', `loadData-${simId}.sh`);
        document.body.appendChild(link);
        link.click();
        link.remove();
    }

    const provisionSimulationDevices = async (simId, payload) => {
        await axios.post(`${BASE_URL}/simulation/${simId}/provision`, payload);
        await listSimulations();
    }

    const getConfig = () => {
        return config;
    }

    const updateConfig = async (payload) => {
        const remoteConfig = await axios.put(`${BASE_URL}/config`, payload);
        setConfig(remoteConfig.data);
    }

    const refreshMetricsStatus = async () => {
        const res = await axios.get(`${BASE_URL}/config/metricsStatus`);
        setMetricsStatus(res.data);
        return res.data;
    }

    return (
        <GlobalContext.Provider
            value={{
                initialized: initialized,
                apps: apps,
                models: models,
                simulations: simulations,
                config: config,
                metricsStatus: metricsStatus,
                initializeData: initializeData,
                getApplication: getApplication,
                listApplications: listApplications,
                addApplication: addApplication,
                updateApplication: updateApplication,
                deleteApplication: deleteApplication,
                getModel: getModel,
                listModels: listModels,
                addModel: addModel,
                updateModel: updateModel,
                deleteModel: deleteModel,
                importModels: importModels,
                getSimulation: getSimulation,
                listSimulations: listSimulations,
                addSimulation: addSimulation,
                updateSimulation: updateSimulation,
                deleteSimulation: deleteSimulation,
                startSimulation: startSimulation,
                stopSimulation: stopSimulation,
                exportSimulation: exportSimulation,
                provisionSimulationDevices: provisionSimulationDevices,
                getConfig: getConfig,
                updateConfig: updateConfig,
                refreshMetricsStatus: refreshMetricsStatus,
            }}
        >
            {props.children}
        </GlobalContext.Provider>
    );
};

export default GlobalContext;
