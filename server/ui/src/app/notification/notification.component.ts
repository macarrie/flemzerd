import { Component, OnInit, Input } from '@angular/core';

import { Movie } from '../movie';
import { TvShow } from '../tvshow';
import { Notification } from '../notification';
import { Const } from '../const';

import { NotificationsService } from '../notifications.service';
import { UtilsService } from '../utils.service';

@Component({
    selector: 'notification',
    templateUrl: './notification.component.html',
    styleUrls: ['./notification.component.scss']
})
export class NotificationComponent implements OnInit {
    @Input() notification :Notification;
    constants :any;

    constructor(
        private notificationsService :NotificationsService,
        private utils :UtilsService
    ) {}

    ngOnInit() {
        this.constants = Const;
    }

    changeReadState(notif :Notification, read :boolean) {
        this.notificationsService.changeReadState(notif, read).subscribe(notif => {
            this.notificationsService.getNotifications();
        });
    }

    getShowTitle(show :TvShow) :string {
        if (show.UseDefaultTitle) {
            return show.Name;
        }

        return show.OriginalName;
    }

    getMovieTitle(movie :Movie) :string {
        if (movie.UseDefaultTitle) {
            return movie.Title;
        }

        return movie.OriginalTitle;
    }

    getStatusPill(n :Notification) {
        let status :Number = 0;

        switch (n.Type) {
            case Const.NOTIFICATION_NEW_EPISODE:
                status = Const.INFO;
                break;
            case Const.NOTIFICATION_NEW_MOVIE:
                status = Const.INFO;
                break;
            case Const.NOTIFICATION_DOWNLOAD_START:
                status = Const.INFO;
                break;
            case Const.NOTIFICATION_DOWNLOAD_SUCCESS:
                status = Const.OK;
                break;
            case Const.NOTIFICATION_DOWNLOAD_FAILURE:
                status = Const.CRITICAL;
                break;
            case Const.NOTIFICATION_TEXT:
                status = Const.INFO;
                break;
            case Const.NOTIFICATION_NO_TORRENT:
                status = Const.WARNING;
                break;
            default:
                status = Const.UNKNOWN;
                break;
        }

        let classes :string = "notif-icon ";

        switch (status) {
            case Const.OK:
                classes += "ok";
                break;
            case Const.WARNING:
                classes += "warning";
                break;
            case Const.CRITICAL:
                classes += "critical";
                break;
            case Const.UNKNOWN:
                classes += "unknown";
                break;
            case Const.INFO:
                classes += "info";
                break;
        }

        return classes;
    }
}
