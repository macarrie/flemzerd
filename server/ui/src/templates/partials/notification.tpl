{{ define "notification" }}
{{ $readClass := "" }}
{{ $readFilterClass := "" }}
{{  if .Read }}
    {{ $readClass = "notification_read" }}
    {{ $readFilterClass = "read" }}
{{ end }}
<li class="{{ $readFilterClass }}">
    <table *ngIf="notification"
        class="uk-table-middle notification-table {{ $readClass }}">
        <tr uk-toggle="target: .notification-detailed-content-{{ .ID }}">
            {{ $statusClass := "" }}
            {{ if or (eq (getNotificationType .) "new_episode") (eq (getNotificationType .) "new_movie") }}
                {{ $statusClass = "info" }}
            {{ end }}
            {{ if eq (getNotificationType .) "download_start" }}
                {{ $statusClass = "info" }}
            {{ end }}
            {{ if eq (getNotificationType .) "download_success" }}
                {{ $statusClass = "ok" }}
            {{ end }}
            {{ if eq (getNotificationType .) "no_torrents" }}
                {{ $statusClass = "warning" }}
            {{ end }}
            {{ if eq (getNotificationType .) "download_failure" }}
                {{ $statusClass = "critical" }}
            {{ end }}
            {{ if eq (getNotificationType .) "text" }}
                {{ $statusClass = "critical" }}
            {{ end }}
            {{ if eq (getNotificationType .) "unknown" }}
                {{ $statusClass = "unknown" }}
            {{ end }}
            <td class="uk-width-auto uk-padding-medium-right notif-icon {{ $statusClass }}"><div></div></td>
            <!--<td class="notification_collapse"-->
                <!--(click)="toggleContent = !toggleContent">-->
                <!--<a class="btn btn-sm p-0 text-muted">-->
                    <!--<i class="material-icons md-18" *ngIf="!toggleContent">expand_more</i>-->
                    <!--<i class="material-icons md-18" *ngIf="toggleContent">expand_less</i>-->
                <!--</a>-->
            <!--</td>-->
            <td class="uk-width-expand">
                {{ if eq (getNotificationType .) "new_episode" }}
                <span *ngIf="notification.Type == constants.NOTIFICATION_NEW_EPISODE">
                    New episode aired: <a href="/episodes/{{ .Episode.ID }}"> {{ getShowTitle .Episode.TvShow }} S{{ printf "%02d" .Episode.Season }}E{{ printf "%02d" .Episode.Number }}: {{ .Episode.Title }}</a>
                </span>
                {{ end }}

                {{ if eq (getNotificationType .) "new_movie" }}
                <span *ngIf="notification.Type == constants.NOTIFICATION_NEW_MOVIE">
                    New movie: <a href="/movies/{{ .Movie.ID }}">{{ getMovieTitle .Movie }}</a>
                </span>
                {{ end }}

                {{ if eq (getNotificationType .) "download_start" }}
                <span *ngIf="notification.Type == constants.NOTIFICATION_DOWNLOAD_START">
                    Download started: 
                    {{ if .Movie.ID }}
                    <a href="/movies/{{ .Movie.ID }}" *ngIf="notification.Movie.ID != 0">{{ getMovieTitle .Movie }}</a>
                    {{ end }}
                    {{ if .Episode.ID }}
                    <a href="/episodes/{{ .Episode.ID }}"> {{ getShowTitle .Episode.TvShow }} S{{ printf "%02d" .Episode.Season }}E{{ printf "%02d" .Episode.Number }}: {{ .Episode.Title }}</a>
                    {{ end }}
                </span>
                {{ end }}

                {{ if eq (getNotificationType .) "download_success" }}
                <span *ngIf="notification.Type == constants.NOTIFICATION_DOWNLOAD_SUCCESS">
                    Download success: 
                    {{ if .Movie.ID }}
                    <a href="/movies/{{ .Movie.ID }}" *ngIf="notification.Movie.ID != 0">{{ getMovieTitle .Movie }}</a>
                    {{ end }}
                    {{ if .Episode.ID }}
                    <a href="/episodes/{{ .Episode.ID }}"> {{ getShowTitle .Episode.TvShow }} S{{ printf "%02d" .Episode.Season }}E{{ printf "%02d" .Episode.Number }}: {{ .Episode.Title }}</a>
                    {{ end }}
                </span>
                {{ end }}

                {{ if eq (getNotificationType .) "download_failure" }}
                <span *ngIf="notification.Type == constants.NOTIFICATION_DOWNLOAD_FAILURE">
                    Download failure: 
                    {{ if .Movie.ID }}
                    <a href="/movies/{{ .Movie.ID }}" *ngIf="notification.Movie.ID != 0">{{ getMovieTitle .Movie }}</a>
                    {{ end }}
                    {{ if .Episode.ID }}
                    <a href="/episodes/{{ .Episode.ID }}"> {{ getShowTitle .Episode.TvShow }} S{{ printf "%02d" .Episode.Season }}E{{ printf "%02d" .Episode.Number }}: {{ .Episode.Title }}</a>
                    {{ end }}
                </span>
                {{ end }}

                {{ if eq (getNotificationType .) "text" }}
                <span *ngIf="notification.Type == constants.NOTIFICATION_TEXT">
                    {{ .Title }}
                </span>
                {{ end }}

                {{ if eq (getNotificationType .) "no_torrents" }}
                <span *ngIf="notification.Type == constants.NOTIFICATION_NO_TORRENT">
                    No torrents found: 
                    {{ if .Movie.ID }}
                    <a href="/movies/{{ .Movie.ID }}" *ngIf="notification.Movie.ID != 0">{{ getMovieTitle .Movie }}</a>
                    {{ end }}
                    {{ if .Episode.ID }}
                    <a href="/episodes/{{ .Episode.ID }}"> {{ getShowTitle .Episode.TvShow }} S{{ printf "%02d" .Episode.Season }}E{{ printf "%02d" .Episode.Number }}: {{ .Episode.Title }}</a>
                    {{ end }}
                </span>
                {{ end }}

                {{ if eq (getNotificationType .) "unknown" }}
                (Unknown notification type) {{ .Title }}
                {{ end }}
            </td>
            <td class="uk-table-shrink uk-text-nowrap uk-text-muted uk-text-right">
                <small>
                    <!--( utils.formatDate(.CreatedAt, 'ddd DD MMM hh:mm') )-->
                    {{ .CreatedAt.Format "02/01/2006 15:04 " }}
                </small>
            </td>
            {{ if not .Read }}
            <td *ngIf="!notification.Read"
                class="notification_action_button uk-text-success"
                data-controller="notifications"
                data-notifications-id="{{ .ID }}"
                data-action="click->notifications#markAsRead"
                (click)="changeReadState(notification, true)">
                <!--<a class="btn btn-sm p-0">-->
                    <!--<i class="material-icons md-18">check</i>-->
                <!--</a>-->
                <span uk-icon="icon: check"></span>
            </td>
            {{ end }}
        </tr>
        <tr class="notification-detailed-content-{{ .ID }}" hidden>
            <td></td>
            <td colspan="2">
                {{ if eq (getNotificationType .) "new_episode" }}
                <span *ngIf="notification.Type == constants.NOTIFICATION_NEW_EPISODE">
                    <blockquote class="uk-text-small uk-padding-left-small">{{ .Episode.Overview }}</blockquote>
                </span>
                {{ end }}
                {{ if eq (getNotificationType .) "new_movie" }}
                <span *ngIf="notification.Type == constants.NOTIFICATION_NEW_MOVIE">
                    <blockquote class="uk-text-small uk-padding-left-small">{{ .Movie.Overview }}</blockquote>
                </span>
                {{ end }}
                {{ if eq (getNotificationType .) "download_start" }}
                <span *ngIf="notification.Type == constants.NOTIFICATION_DOWNLOAD_START">
                    {{ if .Movie.ID }}
                    <blockquote class="uk-text-normal uk-padding-left-small" *ngIf="notification.Movie.ID != 0">{{ .Movie.Overview }}</blockquote>
                    {{ end }}
                    {{ if .Episode.ID }}
                    <blockquote class="uk-text-normal uk-padding-left-small" *ngIf="notification.Episode.ID != 0">{{ .Episode.Overview }}</blockquote>
                    {{ end }}
                </span>
                {{ end }}
                {{ if eq (getNotificationType .) "download_success" }}
                <span *ngIf="notification.Type == constants.NOTIFICATION_DOWNLOAD_SUCCESS">
                    {{ if .Movie.ID }}
                    <blockquote class="uk-text-normal uk-padding-left-small" *ngIf="notification.Movie.ID != 0">{{ .Movie.Overview }}</blockquote>
                    {{ end }}
                    {{ if .Episode.ID }}
                    <blockquote class="uk-text-normal uk-padding-left-small" *ngIf="notification.Episode.ID != 0">{{ .Episode.Overview }}</blockquote>
                    {{ end }}
                </span>
                {{ end }}
                {{ if eq (getNotificationType .) "download_failure" }}
                <span *ngIf="notification.Type == constants.NOTIFICATION_DOWNLOAD_FAILURE">
                    {{ if .Movie.ID }}
                    <blockquote class="uk-text-normal uk-padding-left-small" *ngIf="notification.Movie.ID != 0">{{ .Movie.Overview }}</blockquote>
                    {{ end }}
                    {{ if .Episode.ID }}
                    <blockquote class="uk-text-normal uk-padding-left-small" *ngIf="notification.Episode.ID != 0">{{ .Episode.Overview }}</blockquote>
                    {{ end }}
                </span>
                {{ end }}
                {{ if eq (getNotificationType .) "text" }}
                <span *ngIf="notification.Type == constants.NOTIFICATION_TEXT">
                    {{ .Content }}
                </span>
                {{ end }}
                {{ if eq (getNotificationType .) "no_torrents" }}
                <span *ngIf="notification.Type == constants.NOTIFICATION_NO_TORRENT">
                    No torrents have been found for media. Try adding new indexers or wait for torrents to be available (for recently released movies or episodes)
                </span>
                {{ end }}
            </td>
            <td></td>
        </tr>
    </table>
</li>
{{ end }}
