import React from "react";
import {NavLink} from "react-router-dom";

import API from '../utils/api';
import Const from "../const";

type State = {
    unread_notifications_counter: number,
    config_errors_count: number,
    load_error: boolean,
};

class Header extends React.Component<any, State> {
    notifications_refresh_interval: number;
    config_check_refresh_interval: number;

    state: State = {
        unread_notifications_counter: 0,
        config_errors_count: 0,
        load_error: false,
    };

    constructor(props: any) {
        super(props);

        this.notifications_refresh_interval = 0;
        this.config_check_refresh_interval = 0;

        this.countNotifications = this.countNotifications.bind(this);
    }

    componentDidMount() {
        this.countNotifications();
        this.countConfigurationErrors();

        this.notifications_refresh_interval = window.setInterval(this.countNotifications.bind(this), Const.NOTIFICATIONS_REFRESH);
        this.config_check_refresh_interval = window.setInterval(this.countConfigurationErrors.bind(this), Const.DATA_REFRESH);
    }

    componentWillUnmount() {
        clearInterval(this.notifications_refresh_interval);
        clearInterval(this.config_check_refresh_interval);
    }

    countNotifications() {
        API.Notifications.unread().then(response => {
            this.setState({
                load_error: false,
                unread_notifications_counter: response.data.length
            });
        }).catch(error => {
            console.log("Get notifications error: ", error);
            if (error.response.status === 401) {
                window.location.reload();
            }
            this.setState({load_error: true});
        });
    }

    countConfigurationErrors() {
        API.Config.check().then(response => {
            this.setState({
                load_error: false,
                config_errors_count: response.data.length
            });
        }).catch(error => {
            console.log("Check configuration error: ", error);
            if (error.response.status === 401) {
                window.location.reload();
            }
            this.setState({load_error: true});
        });
    }

    render() {
        return (
            <nav className="uk-navbar-container uk-navbar-transparent topbar uk-margin-small-bottom" data-uk-sticky>
                <div className="uk-container uk-navbar uk-flex uk-flex-middle">
                    <div className="uk-navbar-left">
                        <a className="uk-navbar-item uk-logo" href="/dashboard">
                            <span className="navbar-brand">flemzer</span>
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
                                    {this.state.config_errors_count > 0 && (
                                        <span
                                            className={`uk-badge uk-border-pill uk-text-small ${this.state.config_errors_count > 0 ? "alerted" : ""}`}>
                                            {this.state.config_errors_count}
                                        </span>
                                    )}
                                </NavLink>
                            </li>
                        </ul>
                    </div>
                    <div className="uk-navbar-right">
                        <ul className="uk-navbar-nav">
                            {this.state.load_error !== false && (
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
