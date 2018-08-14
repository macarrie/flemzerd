import { Component, OnInit, OnDestroy } from '@angular/core';
import { Observable, of, Subject } from 'rxjs';
import { HttpClient } from "@angular/common/http";
import { takeUntil  } from 'rxjs/operators'; // for rxjs ^5.5.0 lettable operators

import { ConfigService } from '../config.service';

@Component({
    selector: 'app-settings',
    templateUrl: './settings.component.html',
    styleUrls: ['./settings.component.scss']
})
export class SettingsComponent implements OnInit, OnDestroy {
    private subs :Subject<any> = new Subject();

    config :any;

    trakt_auth :boolean = false;
    trakt_device_code :any;
    trakt_token :any;
    trakt_auth_errors :any;

    telegram_chat_id :number = 0;
    telegram_auth_code :number = 0;

    constructor(
        private configService :ConfigService,
        private http :HttpClient
    ) {
        configService.config.pipe(takeUntil(this.subs)).subscribe(cfg => {
            this.config = cfg;
        });
    }

    ngOnInit() {
        this.configService.getConfig();
        this.GetTraktToken();
        this.GetTelegramChatID();
    }

    ngOnDestroy() {
        this.subs.next();
        this.subs.complete();
    }

    keys(obj :any) {
        if (obj == undefined) {
            return [];
        }
        return Object.keys(obj);
    }

    startTraktAuth() {
        this.http.get("/api/v1/modules/watchlists/trakt/auth", {observe: 'response'}).subscribe(response => {
            if (response.status == 200) {
                var waitForDeviceCode = setInterval(() => {
                    this.http.get("/api/v1/modules/watchlists/trakt/devicecode", {observe: 'response'}).subscribe(response => {
                        if (response.body["device_code"] != "") {
                            this.trakt_device_code = response.body;
                            clearInterval(waitForDeviceCode);
                            var waitForTraktToken = setInterval(() => {
                                this.GetTraktToken();
                                this.checkTraktAuthErrors();

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
                            }, 3000);
                            setTimeout(function() {
                                clearInterval(waitForTraktToken);
                            }, this.trakt_device_code.expires_in * 1000);
                        }
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
        });
    };

    checkTraktAuthErrors() {
        this.http.get("/api/v1/modules/watchlists/trakt/auth_errors").subscribe(response => {
            this.trakt_auth_errors = response;
        });
    };

    GetTraktToken() {
        this.http.get("/api/v1/modules/watchlists/trakt/token").subscribe(response => {
            if (response["access_token"] != "") {
                this.trakt_auth = true;
                this.trakt_token = response;
            }
        });
    };

    GetTelegramChatID() {
        this.http.get("/api/v1/modules/notifiers/telegram/chatid").subscribe(response => {
            if (response["chat_id"] != 0) {
                this.telegram_chat_id = response["chat_id"];
            }
        });
    };

    startTelegramAuth() {
        this.http.get("/api/v1/modules/notifiers/telegram/auth", {observe: 'response'}).subscribe(response => {
            if (response.status == 200) {
                var waitForAuthCode = setInterval(() => {
                    this.http.get("/api/v1/modules/notifiers/telegram/auth_code", {observe: 'response'}).subscribe(response => {
                        if (response.body["auth_code"] != "") {
                            this.telegram_auth_code = response.body["auth_code"];
                            clearInterval(waitForAuthCode);
                            var waitForChatID = setInterval(() => {
                                this.GetTelegramChatID();

                                if (this.telegram_chat_id != 0) {
                                    clearInterval(waitForChatID);
                                }
                            }, 3000);
                            setTimeout(function() {
                                clearInterval(waitForChatID);
                            }, 120 * 1000);
                        }
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
        });
    };
}
