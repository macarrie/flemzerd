{{ define "module_status" }}
<li>
    <div class="uk-flex uk-flex-middle uk-width-expand">
        <div class="uk-width-expand">{{ .Name }}</div>
        {{ if not .Status.Alive }}
        <div class="uk-label uk-label-danger">ERROR</div>
        {{ else }}
        <div class="uk-label uk-label-success">OK</div>
        {{ end }}
    </div>

    {{ if not .Status.Alive }}
    <div>
        <hr class="uk-margin-small-top" />
        <div class="uk-alert uk-alert-danger">
            {{ .Status.Message }}
        </div>
    </div>
    {{ end }}
</li>
{{ end }}
