import { Component } from '@angular/core';

import { NotificationsService } from './notifications.service';
import { ConfigService } from './config.service';

@Component({
    selector: 'app-root',
    templateUrl: './app.component.html',
    styleUrls: ['./app.component.scss']
})
export class AppComponent {
    title = 'flemzer';

    notificationsCounter :number;
    notificationsRefresh :any;

    configRefresh :any;
    config :any;

    constructor(
        private notificationsService :NotificationsService,
        private configService :ConfigService
    ) {
        notificationsService.unreadNotificationsCount.subscribe(count => {
            this.notificationsCounter = count;
        });
        configService.config.subscribe(cfg => {
            this.config = cfg;
        });
    }

    ngOnInit() {
        this.notificationsCounter = 0;
        this.notificationsService.getNotifications();
        this.notificationsRefresh = setInterval(() => {
            this.notificationsService.getNotifications();
        }, 30000);

        this.configService.getConfig();
        this.configRefresh = setInterval(() => {
            this.configService.getConfig();
        }, 60000);
    }
}
