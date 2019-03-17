{{ define "media_ids" }}
<div class="uk-grid uk-grid-small external-links-banner" uk-grid>
    <div>
        {{ if ne .item.MediaIds.Imdb "" }}
        <a  href="https://www.imdb.com/title/{{ .item.MediaIds.Imdb }}"
            target="_blank"
            class="imdb">
            <span class="uk-icon uk-icon-link" uk-icon="link"></span>
            IMDB
        </a>
        {{ end }}
    </div>
    <div>
        {{ if ne .item.MediaIds.Trakt 0 }}
            {{ $href := "" }}
            {{ if eq .type "movie" }}
                {{ $href = "https://trakt.tv/search/trakt/{{ .item.MediaIds.Trakt }}?id_type=movie" }}
            {{ else if eq .type "tvshow" }}
                {{ $href = "https://trakt.tv/search/trakt/{{ .item.MediaIds.Trakt }}?id_type=show" }}
            {{ else if eq .type "episode" }}
                {{ $href = "https://trakt.tv/search/trakt/{{ .item.MediaIds.Trakt }}?id_type=episode" }}
            {{ end }}
            <a  href="{{ $href }}"
            target="_blank"
            class="trakt">
            <span class="uk-icon uk-icon-link" uk-icon="link"></span>
            Trakt
            </a>
        {{ end }}
    </div>
    <div>
        {{ if ne .item.MediaIds.Tmdb 0 }}
        <a  href="https://themoviedb.org/tv/{{ .item.MediaIds.Tmdb }}"
            target="_blank"
            class="tmdb">
            <span class="uk-icon uk-icon-link" uk-icon="link"></span>
            TMDB
        </a>
        {{ end }}
    </div>
</div>
{{ end }}
