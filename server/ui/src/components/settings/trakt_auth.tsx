import React from "react";

import API from "../../utils/api";

type State = {
    verification_url :string,
    user_code :string,
    loading :boolean,
    token :any,
    authenticated :boolean,
};

class TraktAuth extends React.Component<any, State> {
    state :State = {
        verification_url: "",
        user_code: "",
        loading: true,
        token: null,
        authenticated: false,

    }
    constructor(props :any) {
        super(props);

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
        let trakt_auth_errors :any = null;
        let trakt_device_code :any = null;

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

                                if (trakt_auth_errors != null && trakt_auth_errors!.length !== 0) {
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
                <span className="has-text-success">
                    <small><i>configured</i></small>
                </span>
            );
        } else if (this.state.loading) {
            return (
                <button
                    className="button is-loading is-small is-outlined is-info is-disabled">
                    Loading
                </button>
            );
        } else if (this.state.user_code !== "") {
            return "";
        }

        return (
            <button onClick={this.startAuth}
                    className="button is-small is-outlined is-info">
                Connect to trakt.tv
            </button>
        );
    }

    render() {
        return (
            <>
                <div className="column is-narrow">
                    {this.traktAuthStatusButton()}
                </div>
                {(this.state.user_code && !this.state.authenticated) && (
                    <div className={"column is-full"}>
                        <hr />
                        Enter the following validation code into <a target="_blank" rel="noopener noreferrer" href={this.state.verification_url}>{this.state.verification_url}</a>
                        <div className={"columns is-mobile is-centered"}>
                            <div className={"column is-narrow"}>
                                    <h5 className="title is-3 has-text-info trakt-auth-code">{this.state.user_code}</h5>
                            </div>
                        </div>
                    </div>
                )}
            </>
        );
    }
}

export default TraktAuth;
