import React from "react";

import API from "../../utils/api";

class TraktAuth extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            verification_url: "",
            user_code: "",
            loading: true,
            token: null,
            authenticated: false,
        };

        this.getToken = this.getToken.bind(this);
        this.traktAuthStatusButton = this.traktAuthStatusButton.bind(this);
        this.startAuth = this.startAuth.bind(this);
    }

    componentDidMount() {
        this.getToken();
    }

    getToken() {
        API.Modules.Watchlists.Trakt.get_token().then(response => {
            this.setState({
                token: response.data,
                loading: false,
                authenticated: response.data["access_token"] !== "",
            });
        });
    }

    startAuth() {
        let trakt_auth_errors = null;
        let trakt_device_code = null;

        this.setState({loading: true});

        API.Modules.Watchlists.Trakt.auth().then(response => {
            if (response.status === 200) {
                let waitForDeviceCode = setInterval(() => {
                    API.Modules.Watchlists.Trakt.get_device_code().then(response => {
                        if (response.data["device_code"] !== "") {
                            this.setState({
                                loading: false,
                                user_code: response.data["user_code"],
                                verification_url: response.data["verification_url"],
                            });
                            trakt_device_code = response.data;

                            clearInterval(waitForDeviceCode);
                            let waitForTraktToken = setInterval(() => {
                                API.Modules.Watchlists.Trakt.get_token().then(response => {
                                    if (response.data["access_token"] !== "") {
                                        this.setState({
                                            loading: false,
                                            token: response.data,
                                            authenticated: response.data["access_token"] !== "",
                                        });
                                    }
                                });

                                API.Modules.Watchlists.Trakt.get_auth_errors().then(response => {
                                    trakt_auth_errors = response.data;
                                });

                                if (trakt_auth_errors != null && trakt_auth_errors.length !== 0) {
                                    this.setState({
                                        loading: false,
                                    });

                                    clearInterval(waitForTraktToken);
                                    return;
                                }

                                if (this.state.token != null && Object.keys(this.state.token).length !== 0) {
                                    if (this.state.token.access_token !== "") {
                                        clearInterval(waitForTraktToken);
                                    }
                                }
                            }, 3000);

                            setTimeout(function () {
                                clearInterval(waitForTraktToken);
                            }, trakt_device_code.expires_in * 1000);
                        }
                    }).catch(error => {
                        console.log("Trakt auth error: ", error);
                    });
                }, 2000);
                setTimeout(function () {
                    clearInterval(waitForDeviceCode);
                }, 10000);
                return;
            } else if (response.status === 204) {
                console.log("Trakt auth already done");
            } else {
                console.log("Trakt auth error");
            }
        });
    }

    traktAuthStatusButton() {
        if (this.state.authenticated) {
            return (
                <span className="uk-text-success">
                    <small><i>configured</i></small>
                </span>
            );
        } else if (this.state.loading) {
            return (
                <div className="uk-flex uk-flex-middle">
                    <div className="uk-margin-small-right"
                         data-uk-spinner="ratio: 0.5"></div>
                    <div>
                        <small><i>Loading</i></small>
                    </div>
                </div>
            );
        } else if (this.state.user_code !== "") {
            return "";
        }

        return (
            <button
                onClick={this.startAuth}
                className="uk-button uk-button-default uk-text-danger uk-button-small uk-text-right">
                Connect to trakt.tv
            </button>
        );
    }

    render() {
        return (
            <>
                <div className="uk-float-right">
                    {this.traktAuthStatusButton()}
                </div>
                {(this.state.user_code && !this.state.authenticated) && (
                    <div>
                        <hr className="uk-margin-small-top"/>
                        Enter the following validation code into <a target="_blank"
                                                                    rel="noopener noreferrer"
                                                                    href={this.state.verification_url}>{this.state.verification_url}</a>
                        <h5 className="uk-h2 uk-text-center uk-text-primary">{this.state.user_code}</h5>
                    </div>
                )}
            </>
        );
    }
}

export default TraktAuth;
