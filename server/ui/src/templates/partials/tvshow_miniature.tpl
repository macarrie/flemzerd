{{ define "tvshow_miniature" }}
<div class="uk-inline media-miniature">
    <a href="/tvshows/{{ .ID }}">
        {{ $overlay := "" }}
        {{ $imgClass := "" }}
        {{ if or .DeletedAt (eq .Status 3) }}
            {{ $overlay = "overlay-gray" }}
            {{ $imgClass = "black_and_white" }}
        {{ end }}
        {{ if .DeletedAt }}
            {{ $overlay = "overlay-red" }}
        {{ end }}
        <img src="{{ .Poster }}" alt="{{ .Title }}" class="uk-border-rounded {{ $imgClass }}" uk-img>
        <div class="uk-overlay uk-border-rounded uk-position-cover {{ $overlay }}">
            {{ if .DeletedAt }}
            <span class="uk-position-center">
                Removed
            </span>
            {{ end }}
            {{ if and (eq .Status 3) (not .DeletedAt) }}
            <span class="uk-position-center">
                Ended
            </span>
            {{ end }}
        </div>
    </a>

    <div class="uk-icon uk-position-top-right">
        {{ if not .DeletedAt }}
        <button class="uk-icon"
                uk-tooltip="delay: 500; title: Remove"
                uk-icon="icon: close; ratio: 0.75"
                data-controller="tvshow"
                data-action="click->tvshow#delete"
                data-tvshow-id="{{ .ID }}">
        </button>
        {{ else }}
        <button class="uk-icon"
                uk-tooltip="delay: 500; title: Restore"
                uk-icon="icon: reply; ratio: 0.75"
                data-controller="tvshow"
                data-action="click->tvshow#restore"
                data-tvshow-id="{{ .ID }}">
        </button>
        {{ end }}
    </div>
</div>
{{ end }}
