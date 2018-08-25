import { Injectable } from '@angular/core';
import { Observable, of } from 'rxjs';
import { HttpClient } from "@angular/common/http";
import { Subject }    from 'rxjs';

import { Notification } from './notification';
import { Const } from './const';

@Injectable({
    providedIn: 'root'
})
export class NotificationsService {
    constructor(private http :HttpClient) {}

    private notificationSource =  new Subject<Notification[]>();
    private unreadNotificationSource =  new Subject<number>();

    notifications = this.notificationSource.asObservable();
    unreadNotificationsCount = this.unreadNotificationSource.asObservable();

    getNotifications() {
        this.http.get<Notification[]>('/api/v1/notifications/all').subscribe(notifs => {
            this.notificationSource.next(notifs);
            this.getUnreadNotificationCount();
        });
    }

    getReadNotifications(): Observable<Notification[]> {
        return this.http.get<Notification[]>('/api/v1/notifications/read');
    }

    getUnreadNotificationCount() {
        this.http.get<Notification[]>('/api/v1/notifications/unread').subscribe(notifs => {
            this.unreadNotificationSource.next(notifs.length);
        });
    }

    changeReadState(n :Notification, read :boolean) {
        if (n.ID != 0) {
            return this.http.post<Notification>('/api/v1/notifications/read/' + n.ID, {
                "Read": read
            });
        }
    }

    deleteNotifications() {
        this.http.delete<Notification[]>('/api/v1/notifications/all').subscribe(notifs => {
            this.notificationSource.next(notifs);
            this.getUnreadNotificationCount();
        });
    }
}
