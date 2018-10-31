{{ define "notifications_list" }}
<ul class="uk-list uk-list-striped no-stripes item-filter">
    {{ range . }}
        {{ template "notification" . }}
    {{ end }}
</ul>
{{ end }}
