{{ define "media_details_action_bar" }}
{{ $removedClass := "" }}
{{ if .item.DeletedAt }}
{{ $removedClass = "removed-element" }}
{{ end }}
<div class="uk-light action-bar {{ $removedClass }}">
    <div class="uk-container">
        <div class="uk-grid uk-grid-collapse" uk-grid>
            <div class="uk-width-expand">
                {{ $href := "" }}
                {{ if eq .type "movie" }}
                    {{ $href = "/movies" }}
                {{ else if eq .type "tvshow" }}
                    {{ $href = "/tvshows" }}
                {{ else if eq .type "episode" }}
                    {{ $href = printf "/tvshows/%d" .item.TvShow.ID }}
                {{ end }}
                <a href="{{ $href }}">
                    <span class="uk-icon uk-icon-link" uk-icon="icon: arrow-left; ratio: 0.75"></span>
                    Back
                </a>
            </div>
            <div class="uk-width-auto">
                {{ if or (eq .type "movie") (eq .type "episode") }}
                    {{ if and (and (and (and (not .item.DeletedAt) (not .item.DownloadingItem.Pending)) (not .item.DownloadingItem.Downloading)) (not .item.DownloadingItem.Downloaded)) (not (isInFuture .item.Date)) }}
                    <a data-controller="{{ .type }}"
                       data-action="click->{{ .type }}#download"
                       data-{{ .type }}-id="{{ .item.ID }}">
                        <span class="uk-icon uk-icon-link" uk-icon="icon: download; ratio: 0.75"></span>
                        Download
                    </a>
                    {{ end }}
                {{ end }}

                {{ if and (not (eq .type "episode")) (not .item.DeletedAt) }}
                    {{ if ne .item.Title .item.OriginalTitle }}
                        {{ if .item.UseDefaultTitle }}
                        <a data-controller="{{ .type }}"
                           data-action="click->{{ .type }}#useOriginalTitle"
                           data-{{ .type }}-id="{{ .item.ID }}">
                            <span class="uk-icon uk-icon-link" uk-icon="icon: move; ratio: 0.75"></span>
                            Use original title
                        </a>
                        {{ else }}
                        <a data-controller="{{ .type }}"
                           data-action="click->{{ .type }}#useDefaultTitle"
                           data-{{ .type }}-id="{{ .item.ID }}">
                            <span class="uk-icon uk-icon-link" uk-icon="icon: move; ratio: 0.75"></span>
                            Use default title
                        </a>
                        {{ end }}
                    {{ end }}
                {{ end }}

                {{ if .item.DeletedAt }}
                <a (click)="restoreMovie(movie)"
                   data-controller="{{ .type }}"
                   data-action="click->{{ .type }}#restore"
                   data-{{ .type }}-id="{{ .item.ID }}">
                    <span class="uk-icon uk-icon-link" uk-icon="icon: reply; ratio: 0.75"></span>
                    Restore
                </a>
                {{ else }}
                    {{ if or (eq .type "movie") (eq .type "episode") }}
                        {{ if .item.DownloadingItem.Downloaded }}
                        <a data-controller="{{ .type }}"
                            data-action="click->{{ .type }}#markNotDownloaded"
                            data-{{ .type }}-id="{{ .item.ID }}">
                            <span class="uk-icon uk-icon-link" uk-icon="icon: push; ratio: 0.75"></span>
                            Unmark as downloaded
                        </a>
                        {{ else }}
                            {{ if and (not .item.DownloadingItem.Pending) (not .item.DownloadingItem.Downloading) }}
                            <a data-controller="{{ .type }}"
                                data-action="click->{{ .type }}#markDownloaded"
                                data-{{ .type }}-id="{{ .item.ID }}">
                                <span class="uk-icon uk-icon-link" uk-icon="icon: check; ratio: 0.75"></span>
                                Mark as downloaded
                            </a>
                            {{ end }}
                        {{ end }}
                        {{ if or .item.DownloadingItem.Pending .item.DownloadingItem.Downloading }}
                        <a data-controller="{{ .type }}"
                           data-action="click->{{ .type }}#abortDownload"
                           data-{{ .type }}-id="{{ .item.ID }}">
                            <span class="uk-icon uk-icon-link" uk-icon="icon: close; ratio: 0.75"></span>
                            Abort download
                        </a>
                        {{ end }}
                    {{ end }}
                    {{ if or (eq .type "movie") (eq .type "tvshow") }}
                    <a data-controller="{{ .type }}"
                       data-action="click->{{ .type }}#delete"
                       data-{{ .type }}-id="{{ .item.ID }}">
                        <span class="uk-icon uk-icon-link" uk-icon="icon: close; ratio: 0.75"></span>
                        Remove
                    </a>
                    {{ end }}
                {{ end }}
            </div>
        </div>
    </div>
</div>
{{ end }}
