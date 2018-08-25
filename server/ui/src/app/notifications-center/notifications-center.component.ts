import { Component, OnInit, OnDestroy } from '@angular/core';
import { Observable, of, Subject } from 'rxjs';
import { forkJoin } from "rxjs";
import { takeUntil  } from 'rxjs/operators'; // for rxjs ^5.5.0 lettable operators

import { Notification } from '../notification';

import { NotificationsService } from '../notifications.service';
import { UtilsService } from '../utils.service';
import { ConfigService } from '../config.service';

@Component({
    selector: 'app-notifications-center',
    templateUrl: './notifications-center.component.html',
    styleUrls: ['./notifications-center.component.scss']
})
export class NotificationsCenterComponent implements OnInit, OnDestroy {
    private subs :Subject<any> = new Subject();

    notifications :Notification[];
    show_read_notifications :boolean;

    constructor(
        private notificationsService :NotificationsService,
        private utils :UtilsService,
    ) {
        notificationsService.notifications.pipe(takeUntil(this.subs)).subscribe(notifs => {
            this.notifications = notifs;
        });
    }

    ngOnInit() {
        this.show_read_notifications = true;
        this.notificationsService.getNotifications();
    }

    ngOnDestroy() {
        this.subs.next();
        this.subs.complete();
    }

    deleteAll() {
        this.notificationsService.deleteNotifications();
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
