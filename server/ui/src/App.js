import React from "react";
import {BrowserRouter as Router, Route} from "react-router-dom";
//TODO: Load scss instead of css
import "./css/style.scss";

import Header from './components/header';
import Dashboard from './components/dashboard';
import TvShows from './components/tvshows/index';
import EpisodeDetails from './components/tvshows/episode_details';
import Movies from './components/movies/index';
import Status from './components/status';
import Settings from './components/settings/index';


function AppRouter() {
    return (
        <Router>
            <div>
                <Header />

                <Route path="/" exact           component={Dashboard} />
                <Route path="/tvshows"          component={TvShows} />
                <Route path="/episodes/:id"     component={EpisodeDetails} />
                <Route path="/movies"           component={Movies} />
                <Route path="/status"           component={Status} />
                <Route path="/settings"         component={Settings} />
            </div>
        </Router>
    );
}

export default AppRouter;
