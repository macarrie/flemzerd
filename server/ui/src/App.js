import React from "react";
import {BrowserRouter as Router, Route} from "react-router-dom";
//TODO: Load scss instead of css
import "./css/style.scss";

import Header from './components/header';
import Dashboard from './components/dashboard';
import TvShows from './components/tvshows/index';
import Movies from './components/movies/index';
import Status from './components/status';
import Settings from './components/settings/index';
import Notifications from './components/notifications/index';


function AppRouter() {
    return (
        <Router>
            <div>
                <Header />

                <Route path="/dashboard" exact  component={Dashboard} />
                <Route path="/tvshows"          component={TvShows} />
                <Route path="/movies"           component={Movies} />
                <Route path="/status"           component={Status} />
                <Route path="/settings"         component={Settings} />

                <Route path="/notifications"         component={Notifications} />
            </div>
        </Router>
    );
}

export default AppRouter;
