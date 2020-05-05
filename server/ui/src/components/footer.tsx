import React from "react";

import {RiQuestionLine} from "react-icons/ri";
import {RiGithubLine} from "react-icons/ri";

class Footer extends React.Component {
    render() {
        return (
            <nav className="navbar is-light is-transparent footer">
                <div className="container">
                    <div className="navbar-brand">
                        <a className="navbar-item logo" href="/dashboard">
                            <span className="navbar-brand">flemzer</span>
                        </a>
                    </div>
                    <div id="navbarLinks" className={"navbar-menu"}>
                        <div className={"navbar-end"}>
                            <a href="https://github.com/macarrie/flemzerd/wiki" className={"navbar-item"}>
                                <div>
                                    <span className="icon">
                                        <RiQuestionLine/>
                                    </span>
                                    <span className={"header-label"}> Documentation </span>
                                </div>
                            </a>
                            <a href="https://github.com/macarrie/flemzerd" className={"navbar-item"}>
                                <div>
                                    <span className="icon">
                                        <RiGithubLine/>
                                    </span>
                                    <span className={"header-label"}> Github </span>
                                </div>
                            </a>
                        </div>
                    </div>
                </div>
            </nav>
        );
    }
}

export default Footer;
