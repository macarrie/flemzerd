import { Controller } from "stimulus"

export default class extends Controller {
    static get targets() {
        return ["loading", "start", "done", "validation", "usercode", "verificationurl"];
    }

    connect() {
        this.loadingTarget.classList.remove("uk-hidden");
        this.startTarget.classList.add("uk-hidden");
        this.doneTarget.classList.add("uk-hidden");
        this.validationTarget.classList.add("uk-hidden");

        fetch('/api/v1/modules/watchlists/trakt/token', {
            method: "GET",
        }).then(response => {
            return response.json()
        }).then(response => {
            if (response["access_token"] != "") {
                this.trakt_auth = true;
                this.trakt_token = response;
            } else {
                this.startTarget.classList.remove("uk-hidden");
            }

            this.adjustTraktAuthStatus();
        });
    }

    start_auth() {
        this.startTarget.classList.add("uk-hidden");
        this.loadingTarget.classList.remove("uk-hidden");

        fetch('/api/v1/modules/watchlists/trakt/auth', {
            method: "GET",
        }).then(response => {
            if (response.status == 200) {
                var waitForDeviceCode = setInterval(() => {
                    fetch('/api/v1/modules/watchlists/trakt/devicecode', {
                        method: "GET",
                    }).then(response => {
                        return response.json()
                    }).then(response => {
                        if (response["device_code"] != "") {
                            this.loadingTarget.classList.add("uk-hidden");
                            this.usercodeTarget.innerHTML = response["user_code"];
                            this.verificationurlTarget.innerHTML = response["verification_url"];
                            this.verificationurlTarget.href = response["verification_url"];
                            this.validationTarget.classList.remove("uk-hidden");
                            this.trakt_device_code = response;

                            clearInterval(waitForDeviceCode);
                            var waitForTraktToken = setInterval(() => {
                                fetch('/api/v1/modules/watchlists/trakt/token', {
                                    method: "GET",
                                }).then(response => {
                                    return response.json()
                                }).then(response => {
                                    if (response["access_token"] != "") {
                                        this.trakt_auth = true;
                                        this.trakt_token = response;
                                    }
                                });

                                fetch('/api/v1/modules/watchlists/trakt/auth_errors', {
                                    method: "GET",
                                }).then(response => {
                                    return response.json()
                                }).then(response => {
                                    this.trakt_auth_errors = response;
                                });

                                if (this.trakt_auth_errors != null && this.trakt_auth_errors.length != 0) {
                                    this.trakt_auth = false;
                                    clearInterval(waitForTraktToken);
                                    return;
                                }

                                if (this.trakt_token != null && Object.keys(this.trakt_token).length !== 0) {
                                    if (this.trakt_token.access_token != "") {
                                        clearInterval(waitForTraktToken);
                                    }
                                }

                                this.adjustTraktAuthStatus();
                            }, 3000);

                            setTimeout(function() {
                                clearInterval(waitForTraktToken);
                            }, this.trakt_device_code.expires_in * 1000);
                        }
                    }).catch(response => { 
                        UIkit.notification({
                            message: 'Movie restore failed',
                            status: 'danger',
                            pos: 'top-right',
                            timeout: 5000
                        });
                    });
                }, 2000);
                setTimeout(function() {
                    clearInterval(waitForDeviceCode);
                }, 10000);
                return;
            } else if (response.status == 204) {
                console.log("Trakt auth already done");
            } else {
                console.log("Trakt auth error");
            }
            return;
        }).catch(response => { 
            UIkit.notification({
                message: 'Movie restore failed',
                status: 'danger',
                pos: 'top-right',
                timeout: 5000
            });
            console.log("Request error: ", response) ;
        });
    }

    adjustTraktAuthStatus() {
        if (this.trakt_auth == true) {
            this.loadingTarget.classList.add("uk-hidden");
            this.doneTarget.classList.remove("uk-hidden");
            this.startTarget.classList.add("uk-hidden");
            this.validationTarget.classList.add("uk-hidden");
        } else {
            this.loadingTarget.classList.add("uk-hidden");
        }
    }
}
