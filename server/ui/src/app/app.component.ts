import { Component } from '@angular/core';

import { NotificationsService } from './notifications.service';

@Component({
    selector: 'app-root',
    templateUrl: './app.component.html',
    styleUrls: ['./app.component.scss']
})
export class AppComponent {
    title = 'flemzer';
    notificationsCounter :number;
    notificationsRefresh :any;

    constructor(
        private notificationsService :NotificationsService
    ) {
        notificationsService.unreadNotificationsCount.subscribe(count => {
            this.notificationsCounter = count;
        });
    }

    ngOnInit() {
        this.notificationsCounter = 0;
        this.notificationsService.getNotifications();
        this.notificationsRefresh = setInterval(() => {
            this.notificationsService.getNotifications();
        }, 10000);
    }
}
