import React from "react";
import { NavLink } from "react-router-dom";

import API from '../utils/api';

class Header extends React.Component {
    notifications_refresh_interval;
    state = {
        unread_notifications_counter: 0,
    };

    constructor(props) {
        super(props);

        this.countNotifications = this.countNotifications.bind(this);
    }

    componentDidMount() {
        this.countNotifications();

        this.notifications_refresh_interval = setInterval(this.countNotifications().bind(this), 30000);
    }

    componentWillUnmount() {
        clearInterval(this.notifications_refresh_interval);
    }

    countNotifications() {
        console.log("Get notifications");
        API.Notifications.unread().then(response => {
            console.log("Notifications: ", response.data);
            this.setState({unread_notifications: response.data.length});
        }).catch(error => {
            console.log("Get notifications error: ", error);
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
                            <li><NavLink activeClassName="active" to="/dashboard">   <span className="uk-icon uk-margin-small-right" data-uk-icon="ratio: 0.8; icon: home"></span>           Dashboard   </NavLink></li>
                            <li><NavLink activeClassName="active" to="/tvshows">     <span className="uk-icon uk-margin-small-right" data-uk-icon="ratio: 0.8; icon: laptop"></span>         TVShows     </NavLink></li>
                            <li><NavLink activeClassName="active" to="/movies">      <span className="uk-icon uk-margin-small-right" data-uk-icon="ratio: 0.8; icon: video-camera"></span>   Movies      </NavLink></li>
                            <li><NavLink activeClassName="active" to="/status">      <span className="uk-icon uk-margin-small-right" data-uk-icon="ratio: 0.8; icon: heart"></span>          Status      </NavLink></li>
                            <li className="notifications_counter configuration_errors_notification">
                                <NavLink activeClassName="active" to="/settings">
                                    <span className="uk-icon uk-margin-small-right" data-uk-icon="ratio: 0.8; icon: settings"></span>
                                    Settings
                                    <span className="uk-badge uk-border-pill uk-text-small">
                                        {this.state.unread_notifications_counter}
                                    </span>
                                </NavLink>
                            </li>
                        </ul>
                    </div>
                    <div className="uk-navbar-right">
                        <ul className="uk-navbar-nav">
                            <li>
                                <NavLink type="button"
                                    className="notifications_counter"
                                    to="/notifications">
                                    <span data-uk-icon="bell"></span>
                                    <span className="uk-badge uk-border-pill uk-text-small"
                                        data-target="notifications-watcher.counter">
                                        0
                                    </span>
                                </NavLink>
                            </li>
                            <li>
                                <button className="uk-icon-link uk-hidden@m" data-uk-icon="menu" data-uk-toggle="target: #offcanvas-nav-primary"></button>
                            </li>
                        </ul>
                    </div>
                </div>

                <div id="offcanvas-nav-primary" data-uk-offcanvas="mode: push; overlay: true; flip: true">
                    <div className="uk-offcanvas-bar uk-flex uk-flex-column">

                        <ul className="uk-nav uk-nav-default">
                            <li><NavLink activeClassName="uk-active" to="/">        <span className="uk-icon uk-margin-small-right" data-uk-icon="home"></span>         Dashboard   </NavLink></li>
                            <li><NavLink activeClassName="uk-active" to="/tvshows"> <span className="uk-icon uk-margin-small-right" data-uk-icon="laptop"></span>       TVShows     </NavLink></li>
                            <li><NavLink activeClassName="uk-active" to="/movies">  <span className="uk-icon uk-margin-small-right" data-uk-icon="video-camera"></span> Movies      </NavLink></li>
                            <li><NavLink activeClassName="uk-active" to="/status">  <span className="uk-icon uk-margin-small-right" data-uk-icon="heart"></span>        Status      </NavLink></li>
                            <li><NavLink activeClassName="uk-active" to="/settings"><span className="uk-icon uk-margin-small-right" data-uk-icon="settings"></span>     Settings    </NavLink></li>
                        </ul>

                    </div>
                </div>
            </nav>
        );
    }
}

export default Header;
