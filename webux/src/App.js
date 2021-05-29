import React from 'react';
import { BrowserRouter as Router, Route, Switch } from "react-router-dom";
import { GlobalContextProvider } from './context/globalContext';
import './App.css';

import HomePage from './pages/HomePage';
import AppLandingPage from './pages/app/AppLandingPage';
import AppPage from './pages/app/AppPage';
import ModelLandingPage from './pages/model/ModelLandingPage';
import ModelPage from './pages/model/ModelPage';
import SimLandingPage from './pages/sim/SimLandingPage';
import SimPage from './pages/sim/SimPage';
import SettingsPage from './pages/settings/SettingsPage';
import Error404Page from './pages/error/Error404Page';
import "tabler-react/dist/Tabler.css";

function App() {
  return (
    <React.StrictMode>
      <Router>
        <GlobalContextProvider>
          <Switch>
            <Route exact path="/" component={HomePage} />
            <Route exact path="/app" component={AppLandingPage} />
            <Route exact path="/app/:id" component={AppPage} />
            <Route exact path="/model" component={ModelLandingPage} />
            <Route exact path="/model/:id" component={ModelPage} />
            <Route exact path="/sim" component={SimLandingPage} />
            <Route exact path="/sim/:id" component={SimPage} />
            <Route exact path="/settings" component={SettingsPage} />
            <Route component={Error404Page} />
          </Switch>
        </GlobalContextProvider>
      </Router>
    </React.StrictMode>
  );
}

export default App;
