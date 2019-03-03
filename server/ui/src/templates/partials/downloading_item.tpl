{{ define "downloading_item" }}
<div class="uk-background-default uk-border-rounded uk-padding-small">
    <div class="col-12 my-3">
        <span class="uk-h4">Download status</span>
    </div>
    {{ if .Downloaded }}
    <div>
        <div class="uk-text-center uk-margin-small-top uk-margin-small-bottom">
            <span class="uk-h4 uk-text-success">
                <i class="material-icons md-18">check</i> Download successful
            </span>
        </div>
    </div>
    {{ end }}
    {{ if and .Downloading (not .Downloaded) }}
    <div>
        <div class="uk-text-center uk-margin-small-top uk-margin-small-bottom">
            <span class="uk-h4">
                <i class="material-icons md-18">swap_vert</i> Download in progress
            </span>
        </div>
        <table class="uk-table uk-table-small uk-table-divider">
            <tr>
                <td><b>Downloading</b></td>
                <td>
                    {{ if .Downloading }}
                    <span class="uk-text-success"><i class="material-icons md-18 float-left mr-1">check_circle_outline</i>Started at {{ .CurrentTorrent.CreatedAt.Format "15:04 02/01/2006" }}</span>
                    {{ else }}
                    <span class="uk-text-danger"><i class="material-icons md-18 float-left mr-1">check_circle_outline</i></span>
                    {{ end }}
                </td>
            </tr>
            <tr>
                <td><b>Current download</b></td>
                {{ if .CurrentTorrent.Name }}
                <td class="uk-font-italic uk-text-muted">
                    {{ .CurrentTorrent.Name }}
                </td>
                {{ else }}
                <td class="uk-font-italic uk-text-muted">
                    Unknown
                </td>
                {{ end }}
            </tr>
            <tr>
                <td><b>Download directory</b></td>
                {{ if .CurrentTorrent.DownloadDir }}
                <td class="uk-text-muted uk-font-italic">
                    <code>
                        {{ .CurrentTorrent.DownloadDir }}
                    </code>
                </td>
                {{ else }}
                <td class="uk-text-muted uk-font-italic">
                    Unknown
                </td>
                {{ end }}
            </tr>
            {{ $failedClass := "uk-text-danger" }}
            {{ if eq (len .FailedTorrents) 0 }}
                {{ $failedClass = "uk-text-success" }}
            {{ end }}
            <tr class="{{ $failedClass }}">
                <td><b>Failed torrents</b></td>
                <td>
                    <b>{{ len .FailedTorrents }}</b>
                </td>
            </tr>
        </table>
    </div>
    {{ end }}

    {{ if and (and (not .Downloading) (not .Downloaded) (not .Pending)) }}
        {{ if .TorrentsNotFound }}
        <div>
            <div class="uk-text-center uk-margin-small-top uk-margin-small-bottom">
                <span class="uk-text-muted">
                    <div class="uk-h4 uk-text-bold">
                        <i class="material-icons uk-text-warning md-18">warning</i>
                        No torrents found
                    </div>
                    No torrents have been found during last download. <br />
                    Try again later (in case of recent releases) or after adding new indexers.
                </span>
            </div>
        </div>
        {{ else }}
        <div>
            <div class="uk-text-center uk-margin-small-top uk-margin-small-bottom">
                <span class="uk-h4 uk-text-muted">
                    Not downloading
                </span>
            </div>
        </div>
        {{ end }}
    {{ end }}
    {{ if .Pending }}
    <div>
        <div class="uk-text-center uk-margin-small-top uk-margin-small-bottom">
            <span class="uk-h4 uk-text-muted">
                Looking for torrents
            </span>
        </div>
    </div>
    {{ end }}
</div>
{{ end }}
