import axios from "axios";

axios.defaults.baseURL = "/api/v1/";

export default class API {
    static Actions = {
        poll: function () {
            return axios.post('/actions/poll');
        }
    };

    static Config = {
        get: function() {
            return axios.get('/config');
        },
        check: function () {
            return axios.get('/config/check');
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

        get: function (id: number) {
            return axios.get('/movies/details/' + id);
        },
        delete: function (id: number) {
            return axios.delete('/movies/details/' + id);
        },
        restore: function (id: number) {
            return axios.post('/movies/restore/' + id);
        },
        download: function (id: number) {
            return axios.post('/movies/details/' + id + '/download');
        },
        abortDownload: function (id: number) {
            return axios.delete('/movies/details/' + id + '/download');
        },
        changeDownloadedState: function (id: number, downloaded_state: boolean) {
            return axios.put('/movies/details/' + id + '/download_state', {
                Downloaded: downloaded_state,
            });
        },
        changeCustomTitle: function (id: number, customTitle: string) {
            return axios.put('/movies/details/' + id + '/custom_title', {
                CustomTitle: customTitle,
            });
        },
        useDefaultTitle: function (id: number, useDefaultTitle: boolean) {
            return axios.put('/movies/details/' + id + '/use_default_title', {
                UseDefaultTitle: useDefaultTitle,
            });
        }
    };

    static Episodes = {
        downloading: function (id: number) {
            return axios.get('/tvshows/downloading');
        },
        get: function (id: number) {
            return axios.get('/tvshows/episodes/' + id);
        },
        download: function (id: number) {
            return axios.post('/tvshows/episodes/' + id + '/download');
        },
        abortDownload: function (id: number) {
            return axios.delete('/tvshows/episodes/' + id + '/download');
        },
        changeDownloadedState: function (id: number, downloaded_state: boolean) {
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

        get: function (id: number) {
            return axios.get('/tvshows/details/' + id);
        },
        delete: function (id: number) {
            return axios.delete('/tvshows/details/' + id);
        },
        restore: function (id: number) {
            return axios.post('/tvshows/restore/' + id);
        },
        download: function (id: number) {
            return axios.post('/tvshows/details/' + id + '/download');
        },
        abortDownload: function (id: number) {
            return axios.delete('/tvshows/details/' + id + '/download');
        },
        changeCustomTitle: function (id: number, customTitle: string) {
            return axios.put('/tvshows/details/' + id + '/custom_title', {
                CustomTitle: customTitle,
            });
        },
        useDefaultTitle: function (id: number, useDefaultTitle: boolean) {
            return axios.put('/tvshows/details/' + id + '/use_default_title', {
                UseDefaultTitle: useDefaultTitle,
            });
        },
        changeAnimeState: function (id: number, isAnime: boolean) {
            return axios.put('/tvshows/details/' + id + '/change_anime_state', {
                IsAnime: isAnime,
            });
        },
        getSeason: function (id: number, seasonNb: number) {
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
        markRead: function (id: number) {
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

        // TODO; Fix type
        getTvShowFanart: function (show: any) {
            return axios.get(this.FANART_BASE_URL + "tv/" + show.MediaIds.Tvdb, {
                params: {
                    "api_key": this.FANART_TV_KEY
                }
            });
        },
        // TODO; Fix type
        getMovieFanart: function (movie: any) {
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
            refresh: function() {
                return axios.post('/modules/watchlists/refresh');
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
            },
            Telegram: {
                auth: function () {
                    return axios.get('/modules/notifiers/telegram/auth');
                },
                get_auth_code: function () {
                    return axios.get('/modules/notifiers/telegram/auth_code');
                },
                get_chat_id: function () {
                    return axios.get('/modules/notifiers/telegram/chatid');
                },
            },
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
