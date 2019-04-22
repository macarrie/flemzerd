{{ define "content" }}
<div data-controller="refresh"
     data-refresh-interval="60">
</div>

{{ if gt (len .trackedMovies) 0 }}
<!--Add filter only if some movies are present to avoid uikit js errors-->
<div class="uk-container" uk-filter="target: .item-filter">
{{ else }}
<div class="uk-container">
{{ end }}

    {{ if eq (len .trackedMovies) 0 }}
    <div class="uk-text-center">
        <span>
            <div class="uk-h3 uk-text-muted">
                No tracked movies
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

    {{ if .downloadingMovies }}
    <div class="uk-width-1-1">
        <span class="uk-h3">Downloading movies</span>
        <hr />
    </div>
    <div class="uk-width-1-1">
        <table class="uk-table uk-table-small uk-table-divider">
            <thead>
                <th>Movie name</th>
                <th class="uk-text-nowrap">current download</th>
                <th class="uk-text-nowrap">Download status</th>
                <th class="uk-text-center uk-table-shrink uk-text-nowrap">Failed torrents</th>
                <th></th>
            </thead>

            <tbody>
                {{ range .downloadingMovies }}
                <tr>
                    <td>
                        <a href="/movies/{{ .ID }}">
                            {{ .GetTitle }}
                        </a>
                    </td>
                    <td class="uk-text-muted">
                        <i>
                            {{ .DownloadingItem.CurrentTorrent.Name }}
                        </i>
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
                        <a class="uk-icon uk-icon-link uk-text-danger abort-download-button" uk-icon="icon: close; ratio: 0.75" (click)="stopDownload(episode)"
                            data-controller="movie"
                            data-action="click->movie#abortDownload"
                            data-movie-id="{{ .ID }}">
                        </a>
                    </td>
                </tr>
                {{ end }}
            </tbody>
        </table>
    </div>
    {{ end }}

    {{ if gt (len .trackedMovies) 0 }}
    <div class="uk-grid" uk-grid>
        <div class="uk-width-expand">
            <span class="uk-h3">Movies</span>
        </div>
        <ul class="uk-subnav uk-subnav-pill">
            <li class="uk-visible@s uk-active" uk-filter-control>
                <a href="">
                    All
                </a>
            </li>
            <li class="uk-visible@s" uk-filter-control="filter: .tracked-movie">
                <a href="">
                    Tracked
                </a>
            </li>
            <li class="uk-visible@s" uk-filter-control="filter: .future-movie">
                <a href="">
                    Future
                </a>
            </li>
            <li class="uk-visible@s" uk-filter-control="filter: .downloaded-movie">
                <a href="">
                    Downloaded
                </a>
            </li>
            <li class="uk-visible@s" uk-filter-control="filter: .removed-movie">
                <a href="">
                    Removed
                </a>
            </li>
            <li data-controller="request"
                data-request-url="/api/v1/actions/poll"
                data-request-method="POST">
                <a data-action="click->request#send">
                    <div data-target="request.label"
                         class="uk-flex uk-flex-middle">
                            <div class="uk-float-left">
                                <span class="uk-icon" uk-icon="icon: refresh; ratio: 0.75"></span>
                            </div>
                            <div class="uk-text-nowrap">
                                Check now
                            </div>
                    </div>
                    <div data-target="request.wait">
                        <div uk-spinner="ratio: 0.5"></div>
                        <div>
                            <i>Checking</i>
                        </div>
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

    <div class="uk-grid uk-grid-small uk-child-width-1-2 uk-child-width-1-6@l uk-child-width-1-4@m uk-child-width-1-3@s item-filter" uk-grid>
        {{ range .downloadingMovies }}
        <div>
            {{ template "movie_miniature" . }}
        </div>
        {{ end }}

        {{ range .trackedMovies }}
        {{ if isInFuture .Date }}
        <div class="future-movie">
            {{ template "movie_miniature" . }}
        </div>
        {{ else }}
        <div class="tracked-movie">
            {{ template "movie_miniature" . }}
        </div>
        {{ end }}
        {{ end }}

        {{ range .downloadedMovies }}
        <div class="downloaded-movie">
            {{ template "movie_miniature" . }}
        </div>
        {{ end }}

        {{ range .removedMovies }}
        <div class="removed-movie">
            {{ template "movie_miniature" . }}
        </div>
        {{ end }}
    </div>
    {{ end }}
</div>
{{ end }}
