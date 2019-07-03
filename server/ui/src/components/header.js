import React from "react";
import {NavLink} from "react-router-dom";

import API from '../utils/api';
import Const from "../const";

class Header extends React.Component {
    notifications_refresh_interval;
    config_check_refresh_interval;

    state = {
        unread_notifications_counter: 0,
        config_errors_count: 0,
        load_error: null,
    };

    constructor(props) {
        super(props);

        this.countNotifications = this.countNotifications.bind(this);
    }

    componentDidMount() {
        this.countNotifications();
        this.countConfigurationErrors();

        this.notifications_refresh_interval = setInterval(this.countNotifications.bind(this), Const.NOTIFICATIONS_REFRESH);
        this.config_check_refresh_interval = setInterval(this.countConfigurationErrors.bind(this), Const.DATA_REFRESH);
    }

    componentWillUnmount() {
        clearInterval(this.notifications_refresh_interval);
        clearInterval(this.config_check_refresh_interval);
    }

    countNotifications() {
        API.Notifications.unread().then(response => {
            this.setState({
                load_error: null,
                unread_notifications_counter: response.data.length
            });
        }).catch(error => {
            console.log("Get notifications error: ", error);
            this.setState({load_error: error});
        });
    }

    countConfigurationErrors() {
        API.Config.check().then(response => {
            this.setState({
                load_error: null,
                config_errors_count: response.data.length
            });
        }).catch(error => {
            console.log("Check configuration error: ", error);
            this.setState({load_error: error});
        });
    }

    render() {
        return (
            <nav className="uk-navbar-container uk-navbar-transparent topbar uk-margin-small-bottom" data-uk-sticky>
                <div className="uk-container uk-navbar uk-flex uk-flex-middle">
                    <div className="uk-navbar-left">
                        <a className="uk-navbar-item uk-logo" href="/dashboard">
                            <span className="navbar-brand" href="#">flemzer</span>
                        </a>
                        <ul className="uk-navbar-nav uk-visible@m">
                            <li><NavLink activeClassName="active" to="/dashboard"> <span
                                className="uk-icon uk-margin-small-right"
                                data-uk-icon="ratio: 0.8; icon: home"></span> Dashboard </NavLink></li>
                            <li><NavLink activeClassName="active" to="/tvshows"> <span
                                className="uk-icon uk-margin-small-right"
                                data-uk-icon="ratio: 0.8; icon: laptop"></span> TVShows </NavLink></li>
                            <li><NavLink activeClassName="active" to="/movies"> <span
                                className="uk-icon uk-margin-small-right"
                                data-uk-icon="ratio: 0.8; icon: video-camera"></span> Movies </NavLink></li>
                            <li><NavLink activeClassName="active" to="/status"> <span
                                className="uk-icon uk-margin-small-right"
                                data-uk-icon="ratio: 0.8; icon: heart"></span> Status </NavLink></li>
                            <li className="notifications_counter configuration_errors_notification">
                                <NavLink activeClassName="active" to="/settings">
                                    <span className="uk-icon uk-margin-small-right"
                                          data-uk-icon="ratio: 0.8; icon: settings"></span>
                                    Settings
                                    <span
                                        className={`uk-badge uk-border-pill uk-text-small ${this.state.config_errors_count > 0 ? "alerted" : ""}`}>
                                        {this.state.config_errors_count}
                                    </span>
                                </NavLink>
                            </li>
                        </ul>
                    </div>
                    <div className="uk-navbar-right">
                        <ul className="uk-navbar-nav">
                            {this.state.load_error != null && (
                                <li className="uk-padding-small">
                                    <i className="server-unavailable">
                                        Server unavailable
                                    </i>
                                </li>
                            )}
                            <li>
                                <NavLink type="button"
                                         className="notifications_counter"
                                         to="/notifications">
                                    <span data-uk-icon="bell"></span>
                                    <span
                                        className={`uk-badge uk-border-pill uk-text-small ${this.state.unread_notifications_counter > 0 ? "alerted" : ""}`}>
                                        {this.state.unread_notifications_counter}
                                    </span>
                                </NavLink>
                            </li>
                            <li>
                                <button className="uk-icon-link uk-hidden@m" data-uk-icon="menu"
                                        data-uk-toggle="target: #offcanvas-nav-primary"></button>
                            </li>
                        </ul>
                    </div>
                </div>

                <div id="offcanvas-nav-primary" data-uk-offcanvas="mode: push; overlay: true; flip: true">
                    <div className="uk-offcanvas-bar uk-flex uk-flex-column">

                        <ul className="uk-nav uk-nav-default">
                            <li><NavLink activeClassName="uk-active" to="/"> <span
                                className="uk-icon uk-margin-small-right" data-uk-icon="home"></span> Dashboard
                            </NavLink></li>
                            <li><NavLink activeClassName="uk-active" to="/tvshows"> <span
                                className="uk-icon uk-margin-small-right" data-uk-icon="laptop"></span> TVShows
                            </NavLink></li>
                            <li><NavLink activeClassName="uk-active" to="/movies"> <span
                                className="uk-icon uk-margin-small-right" data-uk-icon="video-camera"></span> Movies
                            </NavLink></li>
                            <li><NavLink activeClassName="uk-active" to="/status"> <span
                                className="uk-icon uk-margin-small-right" data-uk-icon="heart"></span> Status </NavLink>
                            </li>
                            <li><NavLink activeClassName="uk-active" to="/settings"><span
                                className="uk-icon uk-margin-small-right" data-uk-icon="settings"></span> Settings
                            </NavLink></li>
                        </ul>

                    </div>
                </div>
            </nav>
        );
    }
}

export default Header;
