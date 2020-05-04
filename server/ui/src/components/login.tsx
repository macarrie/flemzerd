import React from "react";

import API from "../utils/api";
import Auth from "../auth";

import { FaUser } from "react-icons/fa";
import { FaLock } from "react-icons/fa";

type State = {
    BadCredentials: boolean,
    ServerUnavailable: boolean,
    FormData: any,
}

class Login extends React.Component<any, State> {
    state: State = {
        BadCredentials: false,
        ServerUnavailable: false,
        FormData: {
            username: "",
            password: "",
        }
    };

    constructor(props: any) {
        super(props);

        this.handleChange = this.handleChange.bind(this);
        this.handleSubmit = this.handleSubmit.bind(this);
    }

    handleChange(event: any) {
        event.persist();

        let target: string = event.target.name;
        let value: string;
        if (event.target.type === "checkbox") {
            value = event.target.checked;
        } else {
            value = event.target.value;
        }
        let stateTmp = this.state;
        stateTmp.FormData[target] = value;
        this.setState(stateTmp);
    }

    handleSubmit(event: any) {
        event.preventDefault();
        console.log(this.state);

        API.Auth.login(this.state.FormData.username, this.state.FormData.password).then(response => {
            console.log("Auth response: ", response);
            this.setState({
                BadCredentials: false,
                ServerUnavailable: false,
            });
            Auth.setToken(response.data.token, this.state.FormData.rememberme);
            window.location.reload();
        }).catch(error => {
            console.log("Auth error: ", error.response);
            if (error.response.status === 401) {
                this.setState({BadCredentials: true});
                return;
            }
            this.setState({ServerUnavailable: true});
        })
    }

    render() {
        return (
            <div className="">
                <div className="container">
                    <div className="columns is-centered is-mobile">
                        <div className={"column is-narrow"}>
                            <a className="logo" href="/dashboard">
                                <span className="navbar-brand">flemzer</span>
                            </a>
                        </div>
                    </div>
                    <div className="columns is-centered is-mobile">
                        <form className="column is-one-third-desktop is-four-fifths-mobile"
                              onSubmit={this.handleSubmit}>
                            <div className="field">
                                <label className="label"></label>
                                <div className={"control has-icons-right"}>
                                    <input className="input"
                                           name="username"
                                           placeholder="Username"
                                           onChange={this.handleChange}
                                           value={this.state.FormData.username}
                                           required
                                           type="text"/>
                                    <span className="icon is-small is-right">
                                        <FaUser/>
                                    </span>
                                </div>
                            </div>

                            <div className="field">
                                <label className="label" uk-icon="icon: lock"></label>
                                <div className={"control has-icons-right"}>
                                    <input className="input"
                                           name="password"
                                           placeholder="Password"
                                           onChange={this.handleChange}
                                           value={this.state.FormData.password}
                                           required
                                           type="password"/>
                                    <span className="icon is-small is-right">
                                        <FaLock/>
                                    </span>
                                </div>
                            </div>

                            <div className="field">
                                <label className="checkbox">
                                    <input name="rememberme"
                                           onChange={this.handleChange}
                                           value={this.state.FormData.rememberme}
                                           type="checkbox"/>
                                       <span></span>
                                       Remember me
                                </label>
                            </div>

                            <div className="field">
                                <div className="control">
                                    <input className="button is-link is-fullwidth"
                                           value="Log in"
                                           type="submit"/>
                                </div>
                            </div>
                        </form>
                    </div>
                    {this.state.BadCredentials && (
                        <>
                            <div className="columns is-centered">
                                <div className="column is-one-third">
                                    <div className="message is-danger">
                                        <div className={"message-body"}>
                                            Bad credentials
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </>
                    )}

                    {this.state.ServerUnavailable && (
                        <>
                        <div className="columns is-centered">
                            <div className="column is-one-third">
                                <div className="message is-danger">
                                    <div className={"message-body"}>
                                        Server unavailable
                                    </div>
                                </div>
                            </div>
                        </div>
                        </>
                        )}
                </div>
            </div>
        );
    }
}

export default Login;
