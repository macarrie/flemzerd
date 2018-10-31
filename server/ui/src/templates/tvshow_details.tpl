{{ define "content" }}

<div id="full_background"
    data-controller="fanart"
    data-target="fanart.container"
    data-fanart-type="tvshow"
    data-fanart-media-id="{{ .item.MediaIds.Tvdb }}"></div>
{{ template "media_details_action_bar" . }}

<div class="uk-container uk-light mediadetails">
    <div class="uk-grid" uk-grid>
        <div class="uk-width-1-3">
            <img width="100%" src="{{ .item.Poster }}" alt="{{ .item.Title }}" class="uk-border-rounded" uk-img>
        </div>
        <div class="uk-width-expand">
            <div class="uk-grid uk-grid-collapse uk-flex uk-flex-middle" uk-grid>
                <div class="inlineedit"
                     data-controller="inlineedit"
                     data-inlineedit-type="tvshow"
                     data-inlineedit-id="{{ .item.ID }}">
                    <span class="uk-h2 d-inline"
                          data-target="inlineedit.title"
                          data-action="click->inlineedit#displaycontrols">{{ getShowTitle .item }}</span>
                    <div class="uk-grid uk-grid-small uk-flex uk-flex-middle" uk-grid
                         data-target="inlineedit.controls">
                        <div>
                            <input class="uk-border-rounded uk-h2 uk-margin-remove-top uk-text-light" 
                                   data-target="inlineedit.field"
                                   type="text" 
                                   placeholder="Enter custom title"
                                   value="{{ getShowTitle .item }}">
                        </div>
                        <div>
                            <button data-action="click->inlineedit#submit"
                               data-target="inlineedit.submit"
                               class="uk-button uk-border-rounded submit">
                               <span uk-icon="check"></span>
                            </button>
                        </div>
                        <div>
                            <button data-action="click->inlineedit#closeControls"
                               data-target="inlineedit.cancel"
                               class="uk-button uk-border-rounded cancel">
                               <span uk-icon="close"></span>
                            </button>
                        </div>
                    </div>
                </div>
                <div>
                    <span class="uk-h3 uk-text-muted ml-2 uk-margin-small-left media_title_details">({{ .item.FirstAired.Format "2006" }})</span>
                </div>
            </div>
            <div class="uk-grid uk-grid-medium" uk-grid>
                <div> See on </div>
                <div> {{ template "media_ids" . }} </div>
            </div>

            <div class="container">
                <div class="row">
                    <h4>Overview</h4>
                    <div class="col-12">
                        {{ .item.Overview}}
                    </div>
                </div>
                <div class="row mt-3">
                    <h4>Seasons</h4>
                    <div class="col-12">
                        {{ .item.NumberOfSeasons}} seasons, {{ .item.NumberOfEpisodes }} episodes
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>

{{ if not .item.DeletedAt }}
<div class="uk-container uk-margin-medium-top uk-margin-medium-bottom">
    <div class="uk-background-default uk-border-rounded">
        <div class="uk-grid uk-child-width-1-1 uk-padding-small" uk-grid>
            <div>
                <ul uk-accordion="multiple: true">
                    {{ range .seasonsdetails }}
                    <li>
                        {{ if gt .Info.SeasonNumber 0 }}
                        <div class="uk-accordion-title" *ngIf="season && i > 0">
                            <div class="uk-grid uk-grid-small uk-child-width-1-2 uk-flex uk-flex-middle" uk-grid>
                                <div class="uk-width-expand" (click)="season.collapse = !season.collapse">
                                    <span class="uk-h3">Season {{ .Info.SeasonNumber }} </span><span class="uk-text-muted uk-h6">( {{ len .EpisodeList }} episodes )</span>
                                </div>
                                <div class="uk-width-auto">
                                    <ul class="uk-iconnav">
                                        <li><a uk-tooltip="delay: 500; title: Mark whole season as not downloaded" 
                                               class="uk-icon"
                                               uk-icon="icon: push; ratio: 0.75"
                                               href=""></a></li>
                                        <li><a uk-tooltip="delay: 500; title: Mark whole season as downloaded"
                                               class="uk-icon"
                                               uk-icon="icon: check; ratio: 0.75"
                                               href=""></a></li>
                                        <li><a uk-tooltip="delay: 500; title: Download the whole season"
                                               class="uk-icon"
                                               uk-icon="icon: download; ratio: 0.75"
                                               href=""></a></li>
                                    </ul>
                                </div>
                            </div>
                        </div>
                        <div class="uk-accordion-content" *ngIf="season && season.collapse && i > 0">
                            <table class="uk-table uk-table-divider uk-table-small">
                                <thead>
                                    <tr>
                                        <th>Title</th>
                                        <th>Release date</th>
                                        <th class="uk-table-shrink uk-text-nowrap uk-text-center">Download status</th>
                                        <th class="uk-table-shrink uk-text-nowrap uk-text-center">Actions</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {{ range .EpisodeList }}
                                    {{ $futureClass := "" }}
                                    {{ if isInFuture .Date }}
                                        {{ $futureClass = "episode-list-future" }}
                                    {{ end }}
                                    <tr class="{{ $futureClass }}">
                                        <td>
                                            <span class="uk-text-small uk-border-rounded episode-label">{{ printf "%02d" .Season }}x{{ printf "%02d" .Number }}</span>
                                            <a href="/episodes/{{ .ID }}">
                                                {{ .Title }}
                                            </a>
                                        </td>
                                        <td class="uk-table-shrink uk-text-nowrap uk-text-center">{{ .Date.Format "Monday 02 Jan 2006" }}</td>
                                        <td class="uk-text-center">
                                            {{ if .DownloadingItem.Downloaded }}
                                            <a class="btn btn-sm p-0 text-success"
                                               *ngIf="episode.DownloadingItem.Downloaded">
                                                <i class="material-icons">check</i>
                                            </a>
                                            {{ end }}
                                            {{ if not (or (or (or (isInFuture .Date) .DownloadingItem.Downloaded) .DownloadingItem.Downloading) .DownloadingItem.Pending) }}
                                            <div *ngIf="!utils.dateIsInFuture(episode.Date) && !episode.DownloadingItem.Downloaded && !episode.DownloadingItem.Downloading && !episode.DownloadingItem.Pending">
                                                <button class="uk-icon"
                                                        uk-tooltip="delay: 500; title: Download"
                                                        uk-icon="icon: download; ratio: 0.75"
                                                        *ngIf="!movie.DownloadingItem.Downloaded && !movie.DownloadingItem.Downloading && !movie.DownloadingItem.Pending && movie.DeletedAt == null && !config?.System?.AutomaticMovieDownload && !utils.dateIsInFuture(movie.Date)"
                                                        (click)="downloadMovie(movie.ID)"
                                                        data-controller="episode"
                                                        data-action="click->episode#download"
                                                        data-episode-id="{{ .ID }}">
                                                </button>
                                                <div class="btn btn-sm">
                                                    {{ if .DownloadingItem.TorrentsNotFound }}
                                                    <i class="material-icons uk-text-warning md-18"
                                                       *ngIf="episode.DownloadingItem.TorrentsNotFound"
                                                       title="No torrents found for episode on last download">warning</i>
                                                    {{ end }}
                                                </div>
                                            </div>
                                            {{ end }}
                                            {{ if or .DownloadingItem.Downloading .DownloadingItem.Pending }}
                                            <span class="btn btn-sm p-0" *ngIf="episode.DownloadingItem.Downloading || episode.DownloadingItem.Pending">
                                                <i class="material-icons uk-text-small">swap_vert</i>
                                            </span>
                                            {{ end }}
                                        </td>
                                        <td class="uk-text-center">
                                            {{ if not (isInFuture .Date) }}
                                                {{ if .DownloadingItem.Downloaded }}
                                                <button class="uk-icon"
                                                        uk-tooltip="delay: 500; title: Unmark as downloaded"
                                                        uk-icon="icon: push; ratio: 0.75"
                                                        data-controller="episode"
                                                        data-action="click->episode#markNotDownloaded"
                                                        data-episode-id="{{ .ID }}">
                                                </button>
                                                {{ end }}
                                                    {{ if or .DownloadingItem.Downloading .DownloadingItem.Pending }}
                                                    <span class="btn btn-sm p-0" *ngIf="episode.DownloadingItem.Downloading || episode.DownloadingItem.Pending"
                                                        data-controller="episode"
                                                        data-action="click->episode#abortDownload"
                                                        data-episode-id="{{ .ID }}">
                                                        <span class="uk-icon uk-icon-link uk-text-danger" uk-icon="icon: close; ratio: 0.75"></span>
                                                    </span>
                                                    {{ end }}
                                                {{ if and (and (not .DownloadingItem.Downloaded) (not .DownloadingItem.Pending)) (not .DownloadingItem.Downloading) }}
                                                <button class="uk-icon"
                                                        uk-tooltip="delay: 500; title: Mark as downloaded"
                                                        uk-icon="icon: check; ratio: 0.75"
                                                        (click)="changeDownloadedState(movie, true)"
                                                        data-controller="episode"
                                                        data-action="click->episode#markDownloaded"
                                                        data-episode-id="{{ .ID }}">
                                                </button>
                                                {{ end }}
                                            {{ end }}
                                        </td>
                                    </tr>
                                    {{ end }}
                                </tbody>
                            </table>
                        </div>
                        {{ end }}
                    </li>
                    {{ end }}
                </ul>
            </div>
        </div>
    </div>
</div>
{{ end }}

{{ end }}
