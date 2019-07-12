import React from "react";
import {BrowserRouter as Router, Redirect, Route} from "react-router-dom";

import "./css/style.scss";

import Auth from './auth';

import Header from './components/header';
import Footer from './components/footer';

import Login from './components/login';

import Dashboard from './components/dashboard';
import TvShows from './components/tvshows/index';
import Movies from './components/movies/index';
import Status from './components/status';
import Settings from './components/settings/index';
import Notifications from './components/notifications';


function Root() {
    return (
        <Redirect to="/dashboard" />
    );
}

function AppRouter() {
    if (!Auth.IsLoggedIn()) {
        return (
            <Login/>
        );
    }

    return (
        <Router>
            <div>
                <Header />

                <Route path="/"          exact  component={Root} />
                <Route path="/dashboard" exact  component={Dashboard} />
                <Route path="/tvshows"          component={TvShows} />
                <Route path="/movies"           component={Movies} />
                <Route path="/status"           component={Status} />
                <Route path="/settings"         component={Settings} />

                <Route path="/notifications"         component={Notifications} />

                <Footer/>
            </div>
        </Router>
    );
}

export default AppRouter;
