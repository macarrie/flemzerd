import React from "react";
import { NavLink } from "react-router-dom";


function Header() {
    return (
        <nav className="uk-navbar-container uk-navbar-transparent topbar uk-margin-small-bottom" data-uk-sticky>
            <div className="uk-container uk-navbar uk-flex uk-flex-middle">
                <div className="uk-navbar-left">
                    <a className="uk-navbar-item uk-logo" href="/dashboard">
                        <span className="navbar-brand" href="#">flemzer</span>
                    </a>
                    <ul className="uk-navbar-nav uk-visible@m">
                        <li><NavLink activeClassName="uk-active" to="/">            <span className="uk-icon uk-margin-small-right" data-uk-icon="ratio: 0.8; icon: home"></span>           Dashboard   </NavLink></li>
                        <li><NavLink activeClassName="uk-active" to="/tvshows">     <span className="uk-icon uk-margin-small-right" data-uk-icon="ratio: 0.8; icon: laptop"></span>         TVShows     </NavLink></li>
                        <li><NavLink activeClassName="uk-active" to="/movies">      <span className="uk-icon uk-margin-small-right" data-uk-icon="ratio: 0.8; icon: video-camera"></span>   Movies      </NavLink></li>
                        <li><NavLink activeClassName="uk-active" to="/status">      <span className="uk-icon uk-margin-small-right" data-uk-icon="ratio: 0.8; icon: heart"></span>          Status      </NavLink></li>
                        <li className="notifications_counter configuration_errors_notification"
                            data-controller="configuration-errors"
                            data-configuration-errors-refresh-interval="120">
                            <NavLink activeClassName="uk-active" to="/settings">
                                <span className="uk-icon uk-margin-small-right" data-uk-icon="ratio: 0.8; icon: settings"></span>
                                Settings
                                <span className="uk-badge uk-border-pill uk-text-small"
                                    data-target="configuration-errors.counter">
                                    0
                                </span>
                            </NavLink>
                        </li>
                    </ul>
                </div>
                <div className="uk-navbar-right">
                    <ul className="uk-navbar-nav">
                        <li data-controller="notifications-watcher"
                            data-notifications-watcher-refresh-interval="15">
                            <button type="button"
                                className="notifications_counter"
                                href="/notifications">
                                <span data-uk-icon="bell"></span>
                                <span className="uk-badge uk-border-pill uk-text-small"
                                    data-target="notifications-watcher.counter">
                                    0
                                </span>
                            </button>
                            <div className="uk-width-large" data-uk-drop>
                                <div className="uk-box-shadow-medium"
                                    data-target="notifications-watcher.popup">Loading notifications</div>
                            </div>
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

export default Header;
