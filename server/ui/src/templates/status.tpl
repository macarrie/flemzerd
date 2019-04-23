{{ define "content" }}
<div data-controller="refresh"
     data-refresh-interval="60">
</div>

<div class="uk-container">
    <div class="uk-grid uk-child-width-1-1 uk-child-width-1-2@s" uk-grid="masonry: true">
        <div>
            <h3>Watchlists</h3>
            <ul class="uk-list uk-list-striped no-stripes">
                {{ range .watchlists }}
                {{ template "module_status" . }}
                {{ end }}
            </ul>
        </div>
        <div>
            <h3>Providers</h3>
            <ul class="uk-list uk-list-striped no-stripes">
                {{ range .providers }}
                {{ template "module_status" . }}
                {{ end }}
            </ul>
        </div>
        <div>
            <h3>Notifiers</h3>
            <ul class="uk-list uk-list-striped no-stripes">
                {{ range .notifiers }}
                {{ template "module_status" . }}
                {{ end }}
            </ul>
        </div>
        <div>
            <h3>Indexers</h3>
            <ul class="uk-list uk-list-striped no-stripes">
                {{ range .indexers }}
                {{ template "module_status" . }}
                {{ end }}
            </ul>
        </div>
    
        <div>
            <h3>Downloaders</h3>
            <ul class="uk-list uk-list-striped no-stripes">
                {{ range .downloaders }}
                {{ template "module_status" . }}
                {{ end }}
            </ul>
        </div>
        <div>
            <h3>Media centers</h3>
            <ul class="uk-list uk-list-striped no-stripes">
                {{ range .mediacenters }}
                {{ template "module_status" . }}
                {{ end }}
            </ul>
        </div>
    </div>
</div>
{{ end }}
