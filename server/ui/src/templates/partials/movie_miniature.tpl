{{ define "movie_miniature" }}
<div class="uk-inline media-miniature">
    <a href="/movies/{{ .ID }}">
        {{ $overlay := "" }}
        {{ $imgClass := "" }}
        {{ if isInFuture .Date }}
            {{ $overlay = "overlay-yellow" }}
            {{ $imgClass = "black_and_white" }}
        {{ end }}
        {{ if .DeletedAt }}
            {{ $overlay = "overlay-red" }}
            {{ $imgClass = "black_and_white" }}
        {{ end }}
        {{ if .DownloadingItem.Downloaded }}
            {{ $overlay = "overlay-green" }}
            {{ $imgClass = "black_and_white" }}
        {{ end }}
        {{ if or .DownloadingItem.Pending .DownloadingItem.Downloading }}
            {{ $overlay = "overlay-blue" }}
            {{ $imgClass = "black_and_white" }}
        {{ end }}
        <img src="{{ .Poster }}" alt="{{ .Title }}" class="uk-border-rounded {{ $imgClass }}" uk-img>

        <div class="uk-overlay uk-border-rounded uk-position-cover {{ $overlay }}">
            {{ if and (isInFuture .Date) (not .DeletedAt) }}
            <span class="uk-position-center">
                Future release
            </span>
            {{ end }}

            {{ if and (not .DeletedAt) (or .DownloadingItem.Downloading .DownloadingItem.Pending) }}
                {{ if .DownloadingItem.Downloading }}
                <span class="uk-position-center">
                    Downloading
                </span>
                {{ end }}

                {{ if .DownloadingItem.Pending }}
                <span class="uk-position-center">
                    Pending
                </span>
                {{ end }}
            {{ end }}

            {{ if .DeletedAt }}
            <span class="uk-position-center">
                Removed
            </span>
            {{ end }}

            {{ if and .DownloadingItem.Downloaded (not .DeletedAt) }}
            <span class="uk-position-center">
                Downloaded
            </span>
            {{ end }}
        </div>
    </a>

    <div class="uk-icon uk-position-top-right">
        {{ if and .DownloadingItem.Downloaded (not .DeletedAt) }}
        <button class="uk-icon"
                uk-tooltip="delay: 500; title: Unmark as downloaded"
                uk-icon="icon: push; ratio: 0.75"
                data-controller="movie"
                data-action="click->movie#markNotDownloaded"
                data-movie-id="{{ .ID }}">
        </button>
        {{ end }}

        {{ if and (and (and (not .DownloadingItem.Downloaded) (not .DownloadingItem.Downloading)) (not .DownloadingItem.Downloaded)) (not .DeletedAt) }}
        <button class="uk-icon"
                uk-tooltip="delay: 500; title: Mark as downloaded"
                uk-icon="icon: check; ratio: 0.75"
                (click)="changeDownloadedState(movie, true)"
                data-controller="movie"
                data-action="click->movie#markDownloaded"
                data-movie-id="{{ .ID }}">
        </button>
        {{ end }}

        <!--[>TODO: Handle configuration<]-->
        {{ if and (and (and (and (not .DownloadingItem.Downloaded) (not .DownloadingItem.Downloading)) (not .DownloadingItem.Pending)) (not .DeletedAt)) (not (isInFuture .Date)) }}
        <button class="uk-icon"
                uk-tooltip="delay: 500; title: Download"
                uk-icon="icon: download; ratio: 0.75"
                *ngIf="!movie.DownloadingItem.Downloaded && !movie.DownloadingItem.Downloading && !movie.DownloadingItem.Pending && movie.DeletedAt == null && !config?.System?.AutomaticMovieDownload && !utils.dateIsInFuture(movie.Date)"
                (click)="downloadMovie(movie.ID)"
                data-controller="movie"
                data-action="click->movie#download"
                data-movie-id="{{ .ID }}">
        </button>
        {{ end }}

        {{ if and (not .DeletedAt) (not .DownloadingItem.Downloaded) }}
        <button class="uk-icon"
                uk-tooltip="delay: 500; title: Remove"
                uk-icon="icon: close; ratio: 0.75"
                data-controller="movie"
                data-action="click->movie#delete"
                data-movie-id="{{ .ID }}">
        </button>
        {{ end }}

        {{ if and .DeletedAt (not .DownloadingItem.Downloaded) }}
        <button class="uk-icon"
                uk-tooltip="delay: 500; title: Restore"
                uk-icon="icon: reply; ratio: 0.75"
                data-controller="movie"
                data-action="click->movie#restore"
                data-movie-refresh-url="/movies"
                data-movie-id="{{ .ID }}">
        </button>
        {{ end }}
    </div>
</div>
{{ end }}
