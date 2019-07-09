import React from "react";
import {Link} from "react-router-dom";


class Footer extends React.Component {
    render() {
        return (
            <nav className="uk-navbar-container uk-navbar-transparent uk-position-bottom footer">
                <div className="uk-container uk-navbar uk-flex uk-flex-middle">
                    <div className="uk-navbar-left">
                        <Link className="uk-navbar-item uk-logo" to="/dashboard">
                            <span className="navbar-brand">flemzer</span>
                        </Link>
                    </div>
                    <div className="uk-navbar-right">
                        <ul className="uk-navbar-nav">
                            <li>
                                <a href="https://macarrie.github.io/flemzerd/">
                                    <span className="uk-icon uk-margin-small-right" data-uk-icon="question"></span>
                                    Documentation
                                </a>
                            </li>
                            <li>
                                <a href="https://github.com/macarrie/flemzerd">
                                    <span className="uk-icon uk-margin-small-right" data-uk-icon="github"></span>
                                    GitHub
                                </a>
                            </li>
                        </ul>
                    </div>
                </div>
            </nav>
        );
    }
}

export default Footer;
