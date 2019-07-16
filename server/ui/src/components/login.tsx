import React from "react";

import API from "../utils/api";
import Auth from "../auth";

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
            <div className="uk-container uk-position-center">
                <div className="uk-flex uk-flex-middle uk-flex-center" data-uk-grid>
                    <div className="uk-width-1-1 uk-text-center">
                        <a className="uk-navbar-item uk-logo" href="/dashboard">
                            <span className="navbar-brand">flemzer</span>
                        </a>
                    </div>
                    <div className="uk-width-1-1@s uk-width-1-2@m">
                        <form className="uk-grid-small" data-uk-grid
                              onSubmit={this.handleSubmit}>
                            <div className="uk-inline uk-width-1-1">
                                <span className="uk-form-icon uk-form-icon-flip" uk-icon="icon: user"></span>
                                <input className="uk-input uk-width-1-1"
                                       name="username"
                                       placeholder="Username"
                                       onChange={this.handleChange}
                                       value={this.state.FormData.username}
                                       required
                                       type="text"/>
                            </div>

                            <div className="uk-inline uk-width-1-1">
                                <span className="uk-form-icon uk-form-icon-flip" uk-icon="icon: lock"></span>
                                <input className="uk-input uk-width-1-1"
                                       name="password"
                                       placeholder="Password"
                                       onChange={this.handleChange}
                                       value={this.state.FormData.password}
                                       required
                                       type="password"/>
                            </div>

                            <div className="uk-inline uk-width-1-1">
                                <label className="uk-flex uk-flex-middle">
                                    <input className="uk-checkbox uk-margin-small-right"
                                           name="rememberme"
                                           onChange={this.handleChange}
                                           value={this.state.FormData.RememberMe}
                                           type="checkbox"/>
                                    <span>
                                        Remember me
                                   </span>
                                </label>
                            </div>

                            <div className="uk-inline uk-width-1-1">
                                <input className="uk-button uk-button-primary uk-input"
                                       value="Log in"
                                       type="submit"/>
                            </div>
                        </form>
                    </div>
                    {this.state.BadCredentials && (
                        <>
                            <div className="uk-width-1-1">
                            </div>
                            <div className="uk-with-1-1@s uk-width-1-2@m uk-flex-center">
                                <div className="uk-alert-danger" data-uk-alert>
                                    Bad credentials
                                </div>
                            </div>
                        </>
                    )}

                    {this.state.ServerUnavailable && (
                        <>
                            <div className="uk-width-1-1">
                            </div>
                            <div className="uk-with-1-1@s uk-width-1-2@m uk-flex-center">
                                <div className="uk-alert-danger" data-uk-alert>
                                    Server unavailable
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
