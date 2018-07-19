import { Component, OnInit } from '@angular/core';
import { Observable } from 'rxjs';
import { forkJoin } from "rxjs";

import { Notification } from '../notification';

import { NotificationsService } from '../notifications.service';
import { UtilsService } from '../utils.service';

@Component({
    selector: 'app-notifications-center',
    templateUrl: './notifications-center.component.html',
    styleUrls: ['./notifications-center.component.scss']
})
export class NotificationsCenterComponent implements OnInit {
    notifications :Notification[];
    show_read_notifications :boolean;

    constructor(
        private notificationsService :NotificationsService,
        private utils :UtilsService
    ) {
        notificationsService.notifications.subscribe(notifs => {
            this.notifications = notifs;
        });
    }

    ngOnInit() {
        this.show_read_notifications = true;
        this.notificationsService.getNotifications();
    }

    markAllAsRead() {
        let requests = new Array<Observable<Notification>>();
        for (let notif of this.notifications) {
            if (!notif.Read) {
                requests.push(this.notificationsService.changeReadState(notif, true));
            }
        }

        forkJoin(requests).subscribe(response => {
            this.notificationsService.getNotifications();
        });
    }
}
