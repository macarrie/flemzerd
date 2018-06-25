import { Component, OnInit } from '@angular/core';
import { Observable, of } from 'rxjs';
import { HttpClient } from "@angular/common/http";

import { ConfigService } from '../config.service';

@Component({
    selector: 'app-settings',
    templateUrl: './settings.component.html',
    styleUrls: ['./settings.component.scss']
})
export class SettingsComponent implements OnInit {
    config :any;
    trakt_auth :boolean = false;
    trakt_device_code :any;
    trakt_token :any;
    trakt_auth_errors :any;

    constructor(
        private configService :ConfigService,
        private http :HttpClient
    ) {}

    ngOnInit() {
        this.getConfig();
        this.GetTraktToken();
    }

    getConfig() {
        this.configService.getConfig().subscribe(config => {
            this.config = config;
        })
    }

    keys(obj :any) {
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
}
