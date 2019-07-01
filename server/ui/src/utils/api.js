import axios from "axios";

axios.defaults.baseURL = "/api/v1/";

export default class API {
    static Config = {
        get: function() {
            return axios.get('/config');
        }
    };

    static Stats = {
        get: function() {
            return axios.get('/stats');
        }
    };

    static Movies = {
        tracked: function() {
            return axios.get('/movies/tracked');
        },
        downloading: function() {
            return axios.get('/movies/downloading');
        },
        downloaded: function() {
            return axios.get('/movies/downloaded');
        },
        removed: function() {
            return axios.get('/movies/removed');
        },

        get: function(id) {
            return axios.get('/movies/details/' + id);
        },
        delete: function(id) {
            return axios.delete('/movies/details/' + id);
        },
        restore: function(id) {
            return axios.post('/movies/restore/' + id);
        },
        download: function(id) {
            return axios.post('/movies/details/' + id + '/download');
        },
        abortDownload: function(id) {
            return axios.delete('/movies/details/' + id + '/download');
        },
        changeDownloadedState: function(id, downloaded_state) {
            return axios.put('/movies/details/' + id + '/download_state', {
                Downloaded: downloaded_state,
            });
        },
        changeCustomTitle: function(id, customTitle) {
            return axios.put('/movies/details/' + id + '/custom_title', {
                CustomTitle: customTitle,
            });
        },
        useDefaultTitle: function(id, useDefaultTitleBool) {
            return axios.put('/movies/details/' + id + '/use_default_title', {
                UseDefaultTitle: useDefaultTitleBool,
            });
        }
    };

    static Episodes = {
        downloading: function(id) {
            return axios.get('/tvshows/downloading');
        },
        get: function(id) {
            return axios.get('/tvshows/episodes/' + id);
        },
        download: function(id) {
            return axios.post('/tvshows/episodes/' + id + '/download');
        },
        abortDownload: function(id) {
            return axios.delete('/tvshows/episodes/' + id + '/download');
        },
        changeDownloadedState: function(id, downloaded_state) {
            return axios.put('/tvshows/episodes/' + id + '/download_state', {
                Downloaded: downloaded_state,
            });
        },
    };

    static Shows = {
        tracked: function() {
            return axios.get('/tvshows/tracked');
        },
        removed: function() {
            return axios.get('/tvshows/removed');
        },

        get: function(id) {
            return axios.get('/tvshows/details/' + id);
        },
        delete: function(id) {
            return axios.delete('/tvshows/details/' + id);
        },
        restore: function(id) {
            return axios.post('/tvshows/restore/' + id);
        },
        download: function(id) {
            return axios.post('/tvshows/details/' + id + '/download');
        },
        abortDownload: function(id) {
            return axios.delete('/tvshows/details/' + id + '/download');
        },
        changeDownloadedState: function(id, downloaded_state) {
            return axios.put('/tvshows/details/' + id + '/download_state', {
                Downloaded: downloaded_state,
            });
        },
        changeCustomTitle: function(id, customTitle) {
            return axios.put('/tvshows/details/' + id + '/custom_title', {
                CustomTitle: customTitle,
            });
        },
        useDefaultTitle: function(id, useDefaultTitleBool) {
            return axios.put('/tvshows/details/' + id + '/use_default_title', {
                UseDefaultTitle: useDefaultTitleBool,
            });
        },
        changeAnimeState: function(id, isAnime) {
            return axios.put('/tvshows/details/' + id + '/change_anime_state', {
                IsAnime: isAnime,
            });
        },
        getSeason: function(id, seasonNb) {
            return axios.get('/tvshows/details/' + id + '/seasons/' + seasonNb);
        }
    };

    static Notifications = {
        all: function() {
            return axios.get('/notifications/all');
        },
        read: function() {
            return axios.get('/notifications/read');
        },
        unread: function() {
            return axios.get('/notifications/unread');
        },
        markAllRead: function() {
            return axios.post('/notifications/read');
        },
        markRead: function(id) {
            return axios.post('/notifications/read/' +id, {
                Read: true,
            });
        },
        deleteAll: function() {
            return axios.delete('/notifications/all');
        },
    };

    static Fanart = {
        FANART_BASE_URL: "http://webservice.fanart.tv/v3/",
        FANART_TV_KEY: "648ff4214a1eea4416ad51417fc8a4e4",
        fanart_params: {
            "params": {
                "api_key": this.FANART_TV_KEY
            }
        },

        getTvShowFanart: function(show) {
            return axios.get(this.FANART_BASE_URL + "tv/" + show.MediaIds.Tvdb, {
                params: {
                    "api_key": this.FANART_TV_KEY
                }
            });
        },
        getMovieFanart: function(movie) {
            return axios.get(this.FANART_BASE_URL + "movies/" + movie.MediaIds.Tmdb, {
                params: {
                    "api_key": this.FANART_TV_KEY
                }
            });
        }
    };

    static Modules = {
        Watchlists: {
            status: function() {
                return axios.get('/modules/watchlists/status');
            },
            Trakt: {
                get_token: function () {
                    return axios.get('/modules/watchlists/trakt/token');
                },
                get_device_code: function () {
                    return axios.get('/modules/watchlists/trakt/devicecode');
                },
                auth: function () {
                    return axios.get('/modules/watchlists/trakt/auth');
                },
                get_auth_errors: function () {
                    return axios.get('/modules/watchlists/trakt/auth_errors');
                }
            }
        },
        Providers: {
            status: function() {
                return axios.get('/modules/providers/status');
            }
        },
        Notifiers: {
            status: function() {
                return axios.get('/modules/notifiers/status');
            }
        },
        Indexers: {
            status: function() {
                return axios.get('/modules/indexers/status');
            }
        },
        Downloaders: {
            status: function() {
                return axios.get('/modules/downloaders/status');
            }
        },
        Mediacenters: {
            status: function() {
                return axios.get('/modules/mediacenters/status');
            }
        }
    }
}
