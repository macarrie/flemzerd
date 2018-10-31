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

        fetch('/api/v1/modules/notifiers/telegram/chatid', {
            method: "GET",
        }).then(response => {
            return response.json()
        }).then(response => {
            if (response["chat_id"] != "") {
                this.telegram_auth = true;
                this.telegram_chat_id = response["chat_id"];
            } else {
                this.startTarget.classList.remove("uk-hidden");
            }

            this.adjustTelegramAuthStatus();
        });
    }

    start_auth() {
        this.startTarget.classList.add("uk-hidden");
        this.loadingTarget.classList.remove("uk-hidden");

        fetch('/api/v1/modules/notifiers/telegram/auth', {
            method: "GET",
        }).then(response => {
            if (response.status == 200) {
                var waitForAuthCode = setInterval(() => {
                    fetch('/api/v1/modules/notifiers/telegram/auth_code', {
                        method: "GET",
                    }).then(response => {
                        return response.json()
                    }).then(response => {
                        if (response["auth_code"] != "") {
                            this.loadingTarget.classList.add("uk-hidden");
                            this.usercodeTarget.innerHTML = response["auth_code"];
                            this.validationTarget.classList.remove("uk-hidden");
                            this.telegram_auth_code = response["auth_code"];

                            clearInterval(waitForAuthCode);
                            var waitForChatId = setInterval(() => {
                                fetch('/api/v1/modules/notifiers/telegram/chatid', {
                                    method: "GET",
                                }).then(response => {
                                    return response.json()
                                }).then(response => {
                                    if (response["chat_id"] != "") {
                                        this.telegram_auth = true;
                                        this.telegram_chat_id = response;
                                    }
                                });

                                if (this.telegram_chat_id != 0) {
                                    this.telegram_auth = false;
                                    clearInterval(waitForChatId);
                                    return;
                                }

                                if (this.telegram_auth != null && Object.keys(this.trakt_token).length !== 0) {
                                    if (this.telegram_chat_id != "") {
                                        clearInterval(waitForAuthCode);
                                    }
                                }

                                this.adjustTelegramAuthStatus();
                            }, 3000);

                            setTimeout(function() {
                                clearInterval(waitForAuthCode);
                            }, 120 * 1000);
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
                    clearInterval(waitForAuthCode);
                }, 10000);
                return;
            } else if (response.status == 204) {
                console.log("Telegram auth already done");
            } else {
                console.log("Telegram auth error");
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

    adjustTelegramAuthStatus() {
        if (this.telegram_auth == true) {
            this.loadingTarget.classList.add("uk-hidden");
            this.doneTarget.classList.remove("uk-hidden");
            this.startTarget.classList.add("uk-hidden");
            this.validationTarget.classList.add("uk-hidden");
        } else {
            this.loadingTarget.classList.add("uk-hidden");
        }
    }
}
