import React from "react";

import API from "../../utils/api";

type State = {
    chat_id: number,
    auth_code: number;
    loading: boolean,
    authenticated: boolean,
    config_error: boolean,
};

class TelegramAuth extends React.Component<any, State> {
    state: State = {
        chat_id: 0,
        auth_code: 0,
        loading: true,
        authenticated: false,
        config_error: false,
    };

    constructor(props: any) {
        super(props);

        this.getChatID = this.getChatID.bind(this);
        this.telegramAuthStatusButton = this.telegramAuthStatusButton.bind(this);
        this.startAuth = this.startAuth.bind(this);
    }

    componentDidMount() {
        this.getChatID();
    }

    getChatID() {
        API.Modules.Notifiers.Telegram.get_chat_id().then(response => {
            this.setState({
                chat_id: response.data["chat_id"],
                loading: false,
                authenticated: response.data["chat_id"] !== 0,
            });
        }).catch(error => {
            console.log("Telegram chat ID load error: ", error);
            if (error.response.status === 404) {
                this.setState({config_error: true});
            }
        });
    }

    startAuth() {
        this.setState({loading: true});

        API.Modules.Notifiers.Telegram.auth().then(response => {
            if (response.status === 200) {
                let waitForAuthCode = setInterval(() => {
                    API.Modules.Notifiers.Telegram.get_auth_code().then(response => {
                        if (response.data["auth_code"] !== 0) {
                            this.setState({
                                loading: false,
                                auth_code: response.data["auth_code"],
                            });

                            clearInterval(waitForAuthCode);
                            let waitForChatID = setInterval(() => {
                                this.getChatID();

                                if (this.state.chat_id !== 0) {
                                    clearInterval(waitForChatID);
                                    return;
                                }
                            }, 3000);

                            // Set 120s timeout for telegram auth
                            setTimeout(function () {
                                clearInterval(waitForChatID);
                            }, 120000);
                        }
                    }).catch(error => {
                        console.log("Telegram auth error: ", error);
                    });
                }, 2000);
                setTimeout(function () {
                    clearInterval(waitForAuthCode);
                }, 10000);
                return;
            } else if (response.status === 204) {
                console.log("Telegram auth already done");
            } else {
                console.log("Telegram auth error");
            }
        });
    }

    telegramAuthStatusButton() {
        if (this.state.config_error) {
            return (
                <span className="has-text-danger">
                    <small><i>configuration error</i></small>
                </span>
            );
        }

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
        } else if (this.state.auth_code !== 0) {
            return "";
        }

        return (
            <button
                onClick={this.startAuth}
                className="button is-small is-outlined is-info">
                Connect to Telegram
            </button>
        );
    }

    render() {
        return (
            <>
                <div className="column is-narrow">
                    {this.telegramAuthStatusButton()}
                </div>
                {!!(this.state.auth_code && !this.state.authenticated) && (
                    <div className={"column is-full"}>
                        <hr />
                        Send the following code to <code>@flemzerd_bot</code> Telegram bot to link flemzerd with Telegram
                        <div className={"columns is-mobile is-centered"}>
                            <div className={"column is-narrow"}>
                                <h5 className="title is-3 has-text-info trakt-auth-code">{this.state.auth_code}</h5>
                            </div>
                        </div>
                    </div>
                )}
            </>
        );
    }
}

export default TelegramAuth;
