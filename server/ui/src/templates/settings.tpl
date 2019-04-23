{{ define "content" }}
<div class="uk-container">
    <div class="uk-grid uk-grid-small uk-child-width-1-1 uk-child-width-1-2@m" uk-grid>
        <div class="uk-width-1-1">
            <h3 class="uk-heading-divider">Global settings</h3>
        </div>

        <div>
            <ul class="uk-list uk-list-striped no-stripes">
                <li class="uk-flex uk-flex-middle">
                    <div class="uk-width-expand">
                        Polling interval
                    </div>
                    <div>
                        <code>
                            {{ .config.System.CheckInterval }}mn
                        </code>
                    </div>
                </li>
                <li class="uk-flex uk-flex-middle">
                    <div class="uk-width-expand">
                        Torrent download attempts limit
                    </div>
                    <div>
                        <code>
                            {{ .config.System.TorrentDownloadAttemptsLimit }}
                        </code>
                    </div>
                </li>
            </ul>
        </div>

        <div>
            <ul class="uk-list uk-list-striped no-stripes">
                <li class="uk-flex uk-flex-middle">
                    <div class="uk-width-expand">
                        Interface enabled
                    </div>
                    {{ if .config.Interface.Enabled }}
                    <div class="uk-text-success"><i class="material-icons float-right md-18">check_circle_outline</i></div>
                    {{ else }}
                    <div class="uk-text-danger"><i class="material-icons float-right md-18">clear</i></div>
                    {{ end }}
                </li>
                <li class="uk-flex uk-flex-middle">
                    <div class="uk-width-expand">
                        Interface port
                    </div>
                    <div>
                        <code>
                            {{ .config.Interface.Port }}
                        </code>
                    </div>
                </li>
            </ul>
        </div>

        <div>
            <ul class="uk-list uk-list-striped no-stripes">
                <li class="uk-flex uk-flex-middle">
                    <div class="uk-width-expand">
                        TV Shows tracking enabled
                    </div>
                    {{ if .config.System.TrackShows }}
                    <div class="uk-text-success"><i class="material-icons float-right md-18">check_circle_outline</i></div>
                    {{ else }}
                    <div class="uk-text-danger"><i class="material-icons float-right md-18">clear</i></div>
                    {{ end }}
                </li>
                <li class="uk-flex uk-flex-middle">
                    <div class="uk-width-expand">
                        Movie tracking enabled
                    </div>
                    {{ if .config.System.TrackMovies }}
                    <div class="uk-text-success"><i class="material-icons float-right md-18">check_circle_outline</i></div>
                    {{ else }}
                    <div class="uk-text-danger"><i class="material-icons float-right md-18">clear</i></div>
                    {{ end }}
                </li>
                <li class="uk-flex uk-flex-middle">
                    <div class="uk-width-expand">
                        Automatic show download
                    </div>
                    {{ if .config.System.AutomaticShowDownload }}
                    <div class="uk-text-success"><i class="material-icons float-right md-18">check_circle_outline</i></div>
                    {{ else }}
                    <div class="uk-text-danger"><i class="material-icons float-right md-18">clear</i></div>
                    {{ end }}
                </li>
                <li class="uk-flex uk-flex-middle">
                    <div class="uk-width-expand">
                        Automatic movie download
                    </div>
                    {{ if .config.System.AutomaticMovieDownload }}
                    <div class="uk-text-success"><i class="material-icons float-right md-18">check_circle_outline</i></div>
                    {{ else }}
                    <div class="uk-text-danger"><i class="material-icons float-right md-18">clear</i></div>
                    {{ end }}
                </li>
                <li class="uk-flex uk-flex-middle">
                    <div class="uk-width-expand">
                        Preferred media quality
                    </div>
                    {{ if .config.System.PreferredMediaQuality }}
                    <div>
                        <code>
                            {{ .config.System.PreferredMediaQuality }}
                        </code>
                    </div>
                    {{ else }}
                    <div>
                        None
                    </div>
                    {{ end }}
                </li>
            </ul>
        </div>

        <div>
            <ul class="uk-list uk-list-striped no-stripes">
                <li class="uk-flex uk-flex-middle">
                    <div class="uk-width-expand">
                        TV Shows library path
                    </div>
                    <div class="uk-text-small">
                        <code>
                            {{ .config.Library.ShowPath }}
                        </code>
                    </div>
                </li>
                <li class="uk-flex uk-flex-middle">
                    <div class="uk-width-expand">
                        Movies library path
                    </div>
                    <div>
                        <code class="uk-text-small">
                            <span class="uk-text-small">
                                {{ .config.Library.MoviePath }}
                            </span>
                        </code>
                    </div>
                </li>
            </ul>
        </div>

        <div class="uk-width-1-1">
            <h3 class="uk-heading-divider">Notifications</h3>
        </div>
        <div>
            <ul class="uk-list uk-list-striped no-stripes">
                <li class="uk-flex uk-flex-middle">
                    <div class="uk-width-expand">
                        Notifications enabled 
                    </div>
                    {{ if .config.Notifications.Enabled }}
                    <div class="uk-text-success"><i class="material-icons float-right md-18">check_circle_outline</i></div>
                    {{ else }}
                    <div class="uk-text-danger"><i class="material-icons float-right md-18">clear</i></div>
                    {{ end }}
                </li>
                <li class="uk-flex uk-flex-middle">
                    <div class="uk-width-expand">
                        Notify on new episode 
                    </div>
                    {{ if .config.Notifications.NotifyNewEpisode }}
                    <div class="uk-text-success"><i class="material-icons float-right md-18">check_circle_outline</i></div>
                    {{ else }}
                    <div class="uk-text-danger"><i class="material-icons float-right md-18">clear</i></div>
                    {{ end }}
                </li>
                <li class="uk-flex uk-flex-middle">
                    <div class="uk-width-expand">
                        Notify on download complete 
                    </div>
                    {{ if .config.Notifications.NotifyDownloadComplete }}
                    <div class="uk-text-success"><i class="material-icons float-right md-18">check_circle_outline</i></div>
                    {{ else }}
                    <div class="uk-text-danger"><i class="material-icons float-right md-18">clear</i></div>
                    {{ end }}
                </li>
                <li class="uk-flex uk-flex-middle">
                    <div class="uk-width-expand">
                        Notify on download failure 
                    </div>
                    {{ if .config.Notifications.NotifyFailure }}
                    <div class="uk-text-success"><i class="material-icons float-right md-18">check_circle_outline</i></div>
                    {{ else }}
                    <div class="uk-text-danger"><i class="material-icons float-right md-18">clear</i></div>
                    {{ end }}
                </li>
            </ul>
        </div>

        <div class="uk-width-1-1">
            <h3 class="uk-heading-divider">Modules</h3>
        </div>
        <div>
            <h4>Watchlists</h4>
            <ul class="uk-list uk-list-striped no-stripes">
                {{ range $name, $watchlistConfig := .config.Watchlists }}
                    {{ if eq $name "trakt" }}
                    <li data-controller="trakt">
                    {{ else }}
                    <li>
                    {{ end }}

                    {{ $name }}
                    {{ if eq $name "trakt" }}
                        <div class="uk-float-right">
                            <div class="uk-flex uk-flex-middle"
                                 data-target="trakt.loading">
                                <div class="uk-margin-small-right"
                                     uk-spinner="ratio: 0.5"></div>
                                <div>
                                    <small><i>Loading</i></small>
                                </div>
                            </div>
                            <span class="uk-text-success"
                                data-target="trakt.done">
                                <small><i>configured</i></small>
                            </span>
                            <button
                                class="uk-button uk-button-default uk-text-danger uk-button-small uk-text-right"
                                data-target="trakt.start"
                                data-action="click->trakt#start_auth"
                                (click)="startTraktAuth()">
                                Connect to trakt.tv
                            </button>
                        </div>
                        <div data-target="trakt.validation">
                            <hr class="uk-margin-small-top"/>
                            Enter the following validation code into <a target="_blank" href="trakt_device_code.verification_url" data-target="trakt.verificationurl">trakt_device_code.verification_url</a>
                            <h5 data-target="trakt.usercode"
                                class="uk-h2 uk-text-center uk-text-primary">trakt_device_code.user_code</h5>
                        </div>
                    {{ end }}
                </li>
                {{ end }}
            </ul>
        </div>
        <div>
            <h4>Providers</h4>
            <ul class="uk-list uk-list-striped no-stripes">
                {{ range $name, $moduleConfig := .config.Providers }}
                <li>
                    {{ $name }}
                </li>
                {{ end }}
            </ul>
        </div>
        <div>
            <h4>Notifiers</h4>
            <ul class="uk-list uk-list-striped no-stripes">
                {{ range $name, $moduleConfig := .config.Notifiers }}
                {{ if eq $name "telegram" }}
                <li data-controller="telegram">
                {{ else }}
                <li>
                {{ end }}
                    {{ $name }}
                    {{ if eq $name "telegram" }}
                        <div class="uk-float-right">
                            <div class="uk-flex uk-flex-middle"
                                 data-target="telegram.loading">
                                <div class="uk-margin-small-right"
                                     uk-spinner="ratio: 0.5"></div>
                                <div>
                                    <small><i>Loading</i></small>
                                </div>
                            </div>
                            <span class="uk-text-success"
                                data-target="telegram.done">
                                <small><i>configured</i></small>
                            </span>
                            <button
                                class="uk-button uk-button-default uk-text-danger uk-button-small uk-text-right"
                                data-target="telegram.start"
                                data-action="click->telegram#start_auth"
                                (click)="startTraktAuth()">
                                Connect to Telegram
                            </button>
                        </div>
                        <div data-target="telegram.validation">
                            <hr class="uk-margin-small-top"/>
                            Send the following code to <code>@flemzerd_bot</code> Telegram bot to link flemzerd with Telegram
                            <h5 data-target="telegram.usercode"
                                class="uk-h2 uk-text-center uk-text-primary">telegram_user_code</h5>
                        </div>
                    {{ end }}
                </li>
                {{ end }}
            </ul>
        </div>
        <div>
            <h4>Indexers <span class="text-muted">(torznab)</span></h4>
            <ul class="uk-list uk-list-striped no-stripes">
                {{ range $name, $moduleConfig := .config.Indexers.torznab }}
                <li>
                    {{ $moduleConfig.name }}
                    <span class="float-right text-muted">
                        <small>
                            <code>{{ $moduleConfig.url }}</code>
                        </small>
                    </span>
                </li>
                {{ end }}
            </ul>
        </div>
        <div>
            <h4>Downloaders</h4>
            <ul class="uk-list uk-list-striped no-stripes">
                {{ range $name, $moduleConfig := .config.Downloaders }}
                <li>
                    {{ $name }}
                </li>
                {{ end }}
            </ul>
        </div>
        <div>
            <h4>Media centers</h4>
            <ul class="uk-list uk-list-striped no-stripes">
                {{ range $name, $moduleConfig := .config.MediaCenters }}
                <li>
                    {{ $name }}
                </li>
                {{ end }}
            </ul>
        </div>

        <div class="uk-width-1-1">
            <h4 class="uk-heading-divider">About flemzerd</h4>
            <ul class="uk-list uk-list-striped no-stripes">
                <li>Version: <code>{{ .config.Version }}</code></li>
                <li class="uk-text-center uk-padding-small">flemzerd is an Open Source software. For feature ideas, notifying bugs or contributing, do not hesitate to head to the <a target="blank" href="https://github.com/macarrie/flemzerd">GitHub repo</a></li>
            </ul>
        </div>
    </div>
</div>
{{ end }}
