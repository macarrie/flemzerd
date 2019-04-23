{{ define "content" }}
<div data-controller="refresh"
     data-refresh-interval="60">
</div>

<div class="uk-container">
    <div class="uk-grid uk-text-center" uk-grid>
        <!--if track shows-->
        <div class="uk-width-1-1@s uk-width-1-2@m">
            <div class="uk-heading-hero title-gradient">
                <a href="/tvshows">
                    {{ .trackedShowsNb }}
                </a>
            </div>
            <span class="uk-h4 uk-text-muted">
                Tracked shows
            </span>
        </div>
        <div class="uk-width-1-1@s uk-width-1-2@m">
            <div class="uk-heading-hero title-gradient">
                <a href="/tvshows">
                    {{ .downloadingEpisodesNb }}
                </a>
            </div>
            <span class="uk-h4 uk-text-muted">
                Downloading episodes
            </span>
        </div>

        <!--if track movies-->
        <div class="uk-width-1-1@s uk-width-1-3@m">
            <div class="uk-heading-hero title-gradient">
                <a href="/movies">
                    {{ .trackedMoviesNb }}
                </a>
            </div>
            <span class="uk-h4 uk-text-muted">
                New movies
            </span>
        </div>
        <div class="uk-width-1-1@s uk-width-1-3@m">
            <div class="uk-heading-hero title-gradient">
                <a href="/movies">
                    {{ .downloadingMoviesNb }}
                </a>
            </div>
            <span class="uk-h4 uk-text-muted">
                Downloading movies
            </span>
        </div>
        <div class="uk-width-1-1@s uk-width-1-3@m">
            <div class="uk-heading-hero title-gradient">
                <a href="/movies">
                    {{ .downloadedMoviesNb }}
                </a>
            </div>
            <span class="uk-h4 uk-text-muted">
                Downloaded movies
            </span>
        </div>
    </div>
</div>
{{ end }}
