{{ define "content" }}

<div id="full_background"
    data-controller="fanart"
    data-target="fanart.container"
    data-fanart-type="movie"
    data-fanart-media-id="{{ .item.MediaIds.Tmdb }}"></div>
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
                     data-inlineedit-type="movie"
                     data-inlineedit-id="{{ .item.ID }}">
                    <span class="uk-h2 d-inline"
                          data-target="inlineedit.title"
                          data-action="click->inlineedit#displaycontrols">{{ getMovieTitle .item }}</span>
                    <div class="uk-grid uk-grid-small uk-flex uk-flex-middle" uk-grid
                         data-target="inlineedit.controls">
                        <div>
                            <input class="uk-border-rounded uk-h2 uk-margin-remove-top uk-text-light" 
                                   data-target="inlineedit.field"
                                   type="text" 
                                   placeholder="Enter custom title"
                                   value="{{ getMovieTitle .item }}">
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
                    <span class="uk-h3 uk-text-muted uk-margin-small-left media_title_details">({{ .item.Date.Format "2006" }})</span>
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
