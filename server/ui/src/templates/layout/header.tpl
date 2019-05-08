<nav class="uk-navbar-container uk-navbar-transparent topbar uk-margin-small-bottom" uk-sticky>
    <div class="uk-container uk-navbar uk-flex uk-flex-middle">
        <div class="uk-navbar-left">
            <a class="uk-navbar-item uk-logo" href="/dashboard">
                <span class="navbar-brand" href="#">flemzer</span>
            </a>
            <ul class="uk-navbar-nav uk-visible@m">
                <li><a href="/dashboard">   <span class="uk-icon uk-margin-small-right" uk-icon="ratio: 0.8; icon: home"></span>Dashboard</a></li>
                <li><a href="/tvshows">     <span class="uk-icon uk-margin-small-right" uk-icon="ratio: 0.8; icon: laptop"></span>TVShows</a></li>
                <li><a href="/movies">      <span class="uk-icon uk-margin-small-right" uk-icon="ratio: 0.8; icon: video-camera"></span>Movies</a></li>
                <li><a href="/status">    <span class="uk-icon uk-margin-small-right" uk-icon="ratio: 0.8; icon: heart"></span>Status</a></li>
                <li class="notifications_counter configuration_errors_notification"
                    data-controller="configuration-errors"
                    data-configuration-errors-refresh-interval="120">
                    <a href="/settings">
                        <span class="uk-icon uk-margin-small-right" uk-icon="ratio: 0.8; icon: settings"></span>
                        Settings
                        <span class="uk-badge uk-border-pill uk-text-small"
                            data-target="configuration-errors.counter">
                            0
                        </span>
                    </a>
                </li>
            </ul>
        </div>
        <div class="uk-navbar-right">
            <ul class="uk-navbar-nav">
                <li data-controller="notifications-watcher"
                    data-notifications-watcher-refresh-interval="15">
                    <a type="button"
                       class="notifications_counter"
                       href="/notifications">
                        <span uk-icon="bell"></span>
                        <span class="uk-badge uk-border-pill uk-text-small"
                            data-target="notifications-watcher.counter">
                            0
                        </span>
                    </a>
                    <div class="uk-width-large" uk-drop>
                        <div class="uk-box-shadow-medium"
                            data-target="notifications-watcher.popup">Loading notifications</div>
                    </div>
                </li>
                <li>
                    <a class="uk-icon-link uk-hidden@m" uk-icon="menu" uk-toggle="target: #offcanvas-nav-primary"></a>
                </li>
            </ul>
        </div>
    </div>

    <div id="offcanvas-nav-primary" uk-offcanvas="mode: push; overlay: true; flip: true">
        <div class="uk-offcanvas-bar uk-flex uk-flex-column">

            <ul class="uk-nav uk-nav-default">
                <li><a href="/dashboard"><span class="uk-icon uk-margin-small-right" uk-icon="home"></span>Dashboard</a></li>
                <li><a href="/tvshows"><span class="uk-icon uk-margin-small-right" uk-icon="laptop"></span>TVShows</a></li>
                <li><a href="/movies"><span class="uk-icon uk-margin-small-right" uk-icon="video-camera"></span>Movies</a></li>
                <li><a href="/status"><span class="uk-icon uk-margin-small-right" uk-icon="heart"></span>Status</a></li>
                <li><a href="/settings"><span class="uk-icon uk-margin-small-right" uk-icon="settings"></span>Settings</a></li>
            </ul>

        </div>
    </div>
</nav>
