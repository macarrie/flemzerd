{{ define "content" }}
<div data-controller="refresh"
     data-refresh-interval="60">
</div>

<div class="uk-container" uk-filter="target: .item-filter">
    <div class="uk-grid" uk-grid>
        <div class="uk-width-expand">
            <span class="uk-h3">Notifications</span>
        </div>
        <div>
            <ul class="uk-subnav uk-subnav-pill">
                <li uk-filter-control class="uk-active">
                    <a href="">
                        All
                    </a>
                </li>
                <li uk-filter-control="filter: .read">
                    <a href="">
                        Read
                    </a>
                </li>
                <li>
                    <a data-controller="notifications"
                       data-action="click->notifications#markAllAsRead">
                        <span class="uk-icon"
                              uk-tooltip="delay: 500; title: Mark all notifications as read"
                              uk-icon="icon: check; ratio: 0.75"></span>
                        Mark all as read
                    </a>
                </li>
                <li>
                    <a data-controller="notifications"
                       data-action="click->notifications#deleteAll">
                       <span class="uk-icon"
                             uk-tooltip="delay: 500; title: Delete all notifications"
                             uk-icon="icon: trash; ratio: 0.75"></span>
                       Delete all
                    </a>
                </li>
            </ul>
        </div>
    </div>
    <hr />

    <div class="uk-grid uk-child-width-1-1">
        {{ if gt (len .notifications) 0 }}
            {{ template "notifications_list" .notifications }}
        {{ else }}
            <div class="uk-text-center uk-h2 uk-text-muted">
                <i>
                    No notifications
                </i>
            </div>
        {{ end }}
    </div>

</div>
{{ end }}
