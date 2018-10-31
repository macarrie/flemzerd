{{ define "content" }}
<div data-controller="refresh"
     data-refresh-interval="60">
</div>

{{ if gt (len .trackedShows) 0 }}
<!--Add filter only if some shows are present to avoid uikit js errors-->
<div class="uk-container" uk-filter="target: .item-filter">
{{ else }}
<div class="uk-container">
{{ end }}

    {{ if eq (len .trackedShows) 0 }}
    <div class="uk-text-center">
        <span>
            <div class="uk-h3 uk-text-muted">
                No tracked TV shows
            </div>
            <a class="uk-button uk-button-default uk-button-small"
               data-controller="request"
               data-action="click->request#send"
               data-request-url="/api/v1/modules/watchlists/refresh"
               data-request-method="POST">
               <div data-target="request.label">
                   <span class="uk-icon" uk-icon="icon: refresh; ratio: 0.75"></span>
                   Refresh
               </div>
               <div data-target="request.wait">
                   <div uk-spinner="ratio: 0.5"></div>
                   <i>Refreshing</i>
               </div>
            </a>
        </span>
    </div>
    {{ end }}

    {{ if .downloadingEpisodes }}
    <div class="uk-width-1-1">
        <span class="uk-h3">Downloading episodes</span>
        <hr />
    </div>
    <div class="uk-width-1-1">
        <table class="uk-table uk-table-small uk-table-divider">
            <thead>
                <th>Episode name</th>
                <th class="uk-text-nowrap">Download status</th>
                <th class="uk-text-center uk-table-shrink uk-text-nowrap">Failed torrents</th>
                <th></th>
            </thead>

            <tbody>
                {{ range .downloadingEpisodes }}
                <tr *ngFor="let episode of downloadingEpisodes">
                    <td>
                        <a href="/episodes/{{ .ID }}">
                            {{ getShowTitle .TvShow }} S{{ printf "%02d" .Season }}E{{ printf "%02d" .Number }} : {{ .Title }}
                        </a>
                    </td>
                    <td class="uk-text-center uk-table-shrink uk-text-nowrap">
                        {{ if .DownloadingItem.Downloading }}
                        <span class="uk-text-success">Started at {{ .DownloadingItem.CurrentTorrent.CreatedAt.Format "15:04 02/01/2006" }}</span>
                        {{ end }}
                        {{ if .DownloadingItem.Pending }}
                        <span class="uk-text-italic">Looking for torrents</span>
                        {{ end }}
                    </td>
                    {{ $className := "uk-text-success" }}
                    {{ if (gt (len .DownloadingItem.FailedTorrents) 0) }}
                    {{ $className = "uk-text-danger" }}
                    {{ end }}
                    <td class="uk-text-center uk-text-bold {{ $className }}">
                        {{ len .DownloadingItem.FailedTorrents }}
                    </td>
                    <td>
                        <a class="uk-icon uk-icon-link uk-text-danger" uk-icon="icon: close; ratio: 0.75" (click)="stopDownload(episode)"
                            data-controller="episode"
                            data-action="click->episode#abortDownload"
                            data-episode-id="{{ .ID }}">
                        </a>
                    </td>
                </tr>
                {{ end }}
            </tbody>
        </table>
    </div>
    {{ end }}

    {{ if gt (len .trackedShows) 0 }}
    <div class="uk-grid" uk-grid>
        <div class="uk-width-expand">
            <span class="uk-h3">TV Shows</span>
        </div>
        <ul class="uk-subnav uk-subnav-pill">
            <li uk-filter-control class="active">
                <a href="">
                    All
                </a>
            </li>
            <li uk-filter-control="filter: .ended-show">
                <a href="">
                    Ended
                </a>
            </li>
            <li uk-filter-control="filter: .removed-show">
                <a href="">
                    Removed
                </a>
            </li>
            <li data-controller="request"
                data-request-url="/api/v1/actions/poll"
                data-request-method="POST">
                <a data-action="click->request#send">
                    <div data-target="request.label">
                        <span class="uk-icon" uk-icon="icon: refresh; ratio: 0.75"></span>
                        Check now
                    </div>
                    <div data-target="request.wait">
                        <div uk-spinner="ratio: 0.5"></div>
                        <i>Checking</i>
                    </div>
                </a>
            </li>
            <li data-controller="request"
                data-request-url="/api/v1/modules/watchlists/refresh"
                data-request-method="POST">
                <a data-action="click->request#send">
                    <div data-target="request.label">
                        <span class="uk-icon" uk-icon="icon: refresh; ratio: 0.75"></span>
                        Refresh
                    </div>
                    <div data-target="request.wait">
                        <div uk-spinner="ratio: 0.5"></div>
                        <i>Refreshing</i>
                    </div>
                </a>
            </li>
        </ul>
    </div>
    <hr />

    <div class="uk-grid uk-grid-small uk-child-width-1-6 item-filter" uk-grid>
        {{ range .trackedShows }}
        {{ $endedClass := "" }}
        {{ if or .DeletedAt (eq .Status 3) }}
        {{ $endedClass = "ended-show" }}
        {{ end }}
        <div class="{{ $endedClass }}" *ngFor="let tvshow of trackedShows">
            {{ template "tvshow_miniature" . }}
        </div>
        {{ end }}
        {{ range .removedShows }}
        <div class="removed-show">
            <div>
                {{ template "tvshow_miniature" . }}
            </div>
        </div>
        {{ end }}
    </div>
    {{ end }}
</div>
{{ end }}
