{{ define "downloading_item" }}
<div class="uk-background-default uk-border-rounded uk-padding-small">
    <div class="col-12 my-3">
        <span class="uk-h4">Download status</span>
    </div>
    {{ if .Downloaded }}
    <div class="col-12"
         *ngIf="downloadingitem.Downloaded">
        <div class="mb-3">
            <span class="h6 text-success">
                <i class="material-icons md-18">check</i> Download successful
            </span>
        </div>
    </div>
    {{ end }}
    {{ if and .Downloading (not .Downloaded) }}
    <div class="col-12 my-3"
         *ngIf="downloadingitem.Downloading && !downloadingitem.Downloaded">
        <div class="mb-3">
            <span class="h6">
                <i class="material-icons md-18">swap_vert</i> Download in progress
            </span>
        </div>
        <table class="table table-sm">
            <tr *ngIf="!downloadingitem.Downloaded">
                <td><b>Downloading</b></td>
                <td>
                    <span class="text-danger"
                          *ngIf="!downloadingitem.Downloading"><i class="material-icons md-18 float-left mr-1">check_circle_outline</i></span>
                    <span class="text-success"
                          *ngIf="downloadingitem.Downloading"><i class="material-icons md-18 float-left mr-1">check_circle_outline</i>Started at utils.formatDate(downloadingitem.CurrentTorrent.CreatedAt, 'HH:mm DD/MM/YYYY')</span>
                </td>
            </tr>
            <tr>
                <td *ngIf="!downloadingitem.Downloaded"><b>Current download</b></td>
                {{ if .CurrentTorrent.Name }}
                <td class="font-italic text-muted"
                    *ngIf="downloadingitem.CurrentTorrent.Name">
                    {{ .CurrentTorrent.Name }}
                </td>
                {{ else }}
                <td class="font-italic text-muted"
                    *ngIf="!downloadingitem.CurrentTorrent.Name">
                    Unknown
                </td>
                {{ end }}
            </tr>
            <tr>
                <td><b>Download directory</b></td>
                {{ if .CurrentTorrent.DownloadDir }}
                <td class="text-muted font-italic"
                    *ngIf="downloadingitem.CurrentTorrent.DownloadDir">
                    <code>
                        {{ .CurrentTorrent.DownloadDir }}
                    </code>
                </td>
                {{ else }}
                <td class="text-muted font-italic"
                    *ngIf="!downloadingitem.CurrentTorrent.DownloadDir">
                    Unknown
                </td>
                {{ end }}
            </tr>
            {{ $failedClass := "text-danger" }}
            {{ if eq (len .DownloadingItem.FailedTorrents) 0 }}
                {{ $failedClass = "text-success" }}
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
        <div class="col-12 my-3"
             *ngIf="!downloadingitem.Downloading && !downloadingitem.Downloaded && !downloadingitem.Pending && downloadingitem.TorrentsNotFound">
            <div class="mb-3 text-center">
                <span class="h5 text-muted">
                    <i class="material-icons text-warning md-18">warning</i>
                    <span class="font-weight-bold">
                        No torrents found
                    </span>
                    <hr />
                    No torrents have been found during last download. <br />
                    Try again later (in case of recent releases) or after adding new indexers.
                </span>
            </div>
        </div>
        {{ else }}
        <div class="col-12 my-3"
             *ngIf="!downloadingitem.Downloading && !downloadingitem.Downloaded && !downloadingitem.Pending && !downloadingitem.TorrentsNotFound">
            <div class="mb-3 text-center">
                <span class="h5 text-muted">
                    Not downloading
                </span>
            </div>
        </div>
        {{ end }}
    {{ end }}
    {{ if .Pending }}
    <div class="col-12 my-3"
         *ngIf="downloadingitem.Pending">
        <div class="mb-3 text-center">
            <span class="h5 text-muted">
                Looking for torrents
            </span>
        </div>
    </div>
    {{ end }}
</div>
{{ end }}
