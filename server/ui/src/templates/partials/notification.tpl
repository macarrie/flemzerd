{{ define "notification" }}
{{ $readClass := "" }}
{{ $readFilterClass := "" }}
{{  if .Read }}
    {{ $readClass = "notification_read" }}
    {{ $readFilterClass = "read" }}
{{ end }}
<li class="{{ $readFilterClass }}">
    <table
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
            <td class="uk-width-expand">
                {{ if eq (getNotificationType .) "new_episode" }}
                <span>
                    New episode aired: <a href="/episodes/{{ .Episode.ID }}"> {{ .Episode.TvShow.GetTitle }} S{{ printf "%02d" .Episode.Season }}E{{ printf "%02d" .Episode.Number }}: {{ .Episode.Title }}</a>
                </span>
                {{ end }}

                {{ if eq (getNotificationType .) "new_movie" }}
                <span>
                    New movie: <a href="/movies/{{ .Movie.ID }}">{{ .Movie.GetTitle }}</a>
                </span>
                {{ end }}

                {{ if eq (getNotificationType .) "download_start" }}
                <span>
                    Download started: 
                    {{ if .Movie.ID }}
                    <a href="/movies/{{ .Movie.ID }}">{{ .Movie.GetTitle }}</a>
                    {{ end }}
                    {{ if .Episode.ID }}
                    <a href="/episodes/{{ .Episode.ID }}"> {{ .Episode.TvShow.GetTitle }} S{{ printf "%02d" .Episode.Season }}E{{ printf "%02d" .Episode.Number }}: {{ .Episode.Title }}</a>
                    {{ end }}
                </span>
                {{ end }}

                {{ if eq (getNotificationType .) "download_success" }}
                <span>
                    Download success: 
                    {{ if .Movie.ID }}
                    <a href="/movies/{{ .Movie.ID }}">{{ .Movie.GetTitle }}</a>
                    {{ end }}
                    {{ if .Episode.ID }}
                    <a href="/episodes/{{ .Episode.ID }}"> {{ .Episode.TvShow.GetTitle }} S{{ printf "%02d" .Episode.Season }}E{{ printf "%02d" .Episode.Number }}: {{ .Episode.Title }}</a>
                    {{ end }}
                </span>
                {{ end }}

                {{ if eq (getNotificationType .) "download_failure" }}
                <span>
                    Download failure: 
                    {{ if .Movie.ID }}
                    <a href="/movies/{{ .Movie.ID }}">{{ .Movie.GetTitle }}</a>
                    {{ end }}
                    {{ if .Episode.ID }}
                    <a href="/episodes/{{ .Episode.ID }}"> {{ .Episode.TvShow.GetTitle }} S{{ printf "%02d" .Episode.Season }}E{{ printf "%02d" .Episode.Number }}: {{ .Episode.Title }}</a>
                    {{ end }}
                </span>
                {{ end }}

                {{ if eq (getNotificationType .) "text" }}
                <span>
                    {{ .Title }}
                </span>
                {{ end }}

                {{ if eq (getNotificationType .) "no_torrents" }}
                <span>
                    No torrents found: 
                    {{ if .Movie.ID }}
                    <a href="/movies/{{ .Movie.ID }}">{{ .Movie.GetTitle }}</a>
                    {{ end }}
                    {{ if .Episode.ID }}
                    <a href="/episodes/{{ .Episode.ID }}"> {{ .Episode.TvShow.GetTitle }} S{{ printf "%02d" .Episode.Season }}E{{ printf "%02d" .Episode.Number }}: {{ .Episode.Title }}</a>
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
            <td class="notification_action_button uk-text-success"
                data-controller="notifications"
                data-notifications-id="{{ .ID }}"
                data-action="click->notifications#markAsRead">
                <span uk-icon="icon: check"></span>
            </td>
            {{ end }}
        </tr>
        <tr class="notification-detailed-content-{{ .ID }}" hidden>
            <td></td>
            <td colspan="2">
                {{ if eq (getNotificationType .) "new_episode" }}
                <span>
                    <blockquote class="uk-text-small uk-padding-left-small">{{ .Episode.Overview }}</blockquote>
                </span>
                {{ end }}
                {{ if eq (getNotificationType .) "new_movie" }}
                <span>
                    <blockquote class="uk-text-small uk-padding-left-small">{{ .Movie.Overview }}</blockquote>
                </span>
                {{ end }}
                {{ if eq (getNotificationType .) "download_start" }}
                <span>
                    {{ if .Movie.ID }}
                    <blockquote class="uk-text-normal uk-padding-left-small">{{ .Movie.Overview }}</blockquote>
                    {{ end }}
                    {{ if .Episode.ID }}
                    <blockquote class="uk-text-normal uk-padding-left-small">{{ .Episode.Overview }}</blockquote>
                    {{ end }}
                </span>
                {{ end }}
                {{ if eq (getNotificationType .) "download_success" }}
                <span>
                    {{ if .Movie.ID }}
                    <blockquote class="uk-text-normal uk-padding-left-small">{{ .Movie.Overview }}</blockquote>
                    {{ end }}
                    {{ if .Episode.ID }}
                    <blockquote class="uk-text-normal uk-padding-left-small">{{ .Episode.Overview }}</blockquote>
                    {{ end }}
                </span>
                {{ end }}
                {{ if eq (getNotificationType .) "download_failure" }}
                <span>
                    {{ if .Movie.ID }}
                    <blockquote class="uk-text-normal uk-padding-left-small">{{ .Movie.Overview }}</blockquote>
                    {{ end }}
                    {{ if .Episode.ID }}
                    <blockquote class="uk-text-normal uk-padding-left-small">{{ .Episode.Overview }}</blockquote>
                    {{ end }}
                </span>
                {{ end }}
                {{ if eq (getNotificationType .) "text" }}
                <span>
                    {{ .Content }}
                </span>
                {{ end }}
                {{ if eq (getNotificationType .) "no_torrents" }}
                <span>
                    No torrents have been found for media. Try adding new indexers or wait for torrents to be available (for recently released movies or episodes)
                </span>
                {{ end }}
            </td>
            <td></td>
        </tr>
    </table>
</li>
{{ end }}
