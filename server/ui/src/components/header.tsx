import React from "react";
import {NavLink} from "react-router-dom";

import API from '../utils/api';
import Const from "../const";
import Auth from "../auth";

import { RiHome2Line } from "react-icons/ri";
import { RiTvLine } from "react-icons/ri";
import { RiFilmLine } from "react-icons/ri";
import { RiHospitalLine } from "react-icons/ri";
import { RiSettings3Line } from "react-icons/ri";
import { RiNotification3Line } from "react-icons/ri";

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
                Auth.logout();
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
                Auth.logout();
            }
            this.setState({load_error: true});
        });
    }

    getSettingsClassnames() :String {
        if (this.state.config_errors_count > 0) {
            return "notifications_counter configuration_errors_notification";
        }

        return "";
    }

    render() {
        return (
            <nav className="navbar is-light is-transparent is-fixed-top">
                <div className="container">
                    <div className="navbar-brand">
                        <a className="navbar-item logo" href="/dashboard">
                            <span className="navbar-brand">flemzer</span>
                        </a>
                        <a role="button" className="navbar-burger burger" aria-label="menu" aria-expanded="false"
                           data-target="navbarLinks">
                            <span aria-hidden="true"></span>
                            <span aria-hidden="true"></span>
                            <span aria-hidden="true"></span>
                        </a>
                    </div>
                    <div id="navbarLinks" className={"navbar-menu"}>
                        <div className={"navbar-start"}>
                            <NavLink activeClassName="active" to="/dashboard" className={"navbar-item"}>
                                <div>
                                    <span className="icon">
                                        <RiHome2Line/>
                                    </span>
                                    <span className={"header-label"}> Dashboard </span>
                                </div>
                            </NavLink>
                            <NavLink activeClassName="active" to="/tvshows" className={"navbar-item"}>
                                <div>
                                    <span className="icon">
                                        <RiTvLine/>
                                    </span>
                                    <span className={"header-label"}> TV Shows </span>
                                </div>
                            </NavLink>
                            <NavLink activeClassName="active" to="/movies" className={"navbar-item has-text-grey"}>
                                <div>
                                    <span className="icon">
                                        <RiFilmLine/>
                                    </span>
                                    <span className={"header-label"}> Movies </span>
                                </div>
                            </NavLink>
                            <NavLink activeClassName="active" to="/status" className={"navbar-item has-text-grey"}>
                                <div>
                                    <span className="icon">
                                        <RiHospitalLine />
                                    </span>
                                    <span className={"header-label"}> Status </span>
                                </div>
                            </NavLink>
                            <NavLink activeClassName="active" to="/settings" className={`navbar-item has-text-grey ${this.getSettingsClassnames()}`}>
                                <div>
                                    <span className="icon">
                                        <RiSettings3Line />
                                    </span>
                                    <span className={"header-label"}> Settings </span>
                                    {this.state.config_errors_count > 0 && (
                                        <span
                                            className={`tag ${this.state.config_errors_count > 0 ? "alerted" : ""}`}>
                                                {this.state.config_errors_count}
                                        </span>
                                    )}
                                </div>
                            </NavLink>
                        </div>
                        <div className="navbar-end">
                            {this.state.load_error !== false && (
                                <div className="navbar-item">
                                    <i className="server-unavailable">
                                        Server unavailable
                                    </i>
                                </div>
                            )}
                            <NavLink type="button"
                                     activeClassName="active"
                                     className="notifications_counter navbar-item"
                                     to="/notifications">
                                <div>
                                    <span className="icon">
                                        <RiNotification3Line />
                                    </span>

                                    <span className={`tag ${this.state.unread_notifications_counter > 0 ? "alerted" : ""}`}>
                                            {this.state.unread_notifications_counter}
                                    </span>
                                </div>
                            </NavLink>
                        </div>
                    </div>
                </div>
            </nav>
        );
    }
}

export default Header;
