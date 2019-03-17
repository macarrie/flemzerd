{{ define "content" }}

<div id="full_background"
    data-controller="fanart"
    data-target="fanart.container"
    data-fanart-type="tvshow"
    data-fanart-media-id="{{ .item.TvShow.MediaIds.Tvdb }}"></div>
{{ template "media_details_action_bar" . }}

<div class="uk-container uk-light mediadetails">
    <div class="uk-grid" uk-grid>
        <div class="uk-width-1-3">
            <img width="100%" src="{{ .item.TvShow.Poster }}" alt="{{ .item.TvShow.Title }}" class="uk-border-rounded" uk-img>
        </div>
        <div class="uk-width-expand">
            <div class="uk-grid uk-grid-collapse uk-flex uk-flex-middle" uk-grid>
                <div>
                    <span class="uk-h2 d-inline">{{ .item.TvShow.GetTitle }} - {{ .item.Title }}</span>
                </div>
                <div>
                    <span class="uk-h3 uk-text-muted uk-margin-small-left media_title_details">(S{{ printf "%02d" .item.Season }}E{{ printf "%02d" .item.Number }})</span>
                </div>
            </div>
            <div class="uk-grid uk-grid-medium" uk-grid>
                <div> See on </div>
                <div> {{ template "media_ids" . }} </div>
            </div>

            <div class="container">
                <div class="row">
                    <h5>Overview</h5>
                    <div class="col-12">
                        {{ .item.Overview}}
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>

<div class="uk-container uk-margin-medium-top uk-margin-medium-bottom">
    {{ template "downloading_item" .item.DownloadingItem }}
</div>
{{ end }}
