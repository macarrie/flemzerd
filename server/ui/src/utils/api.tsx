import axios from "axios";

import Auth from "../auth";
import MediaInfo from "../types/media_info";

axios.defaults.baseURL = "/api/v1/";

export default class API {
    static auth = function() {
        return axios.create({
            baseURL: "/api/v1/",
            headers: {'Authorization': 'Bearer ' + Auth.getToken()}
        })
    };

    static Auth = {
        login: function (username :string, password :string) {
            return axios.post('/auth/login', {
                username: username,
                password: password,
            });
        },
        refresh_token: function () {
            return axios.get('/auth/refresh_token');
        },
    };

    static Actions = {
        poll: function () {
            return API.auth().post('/actions/poll');
        },
        refreshAllMetadata: function () {
            return API.auth().post('/actions/refresh_all_metadata');
        },
    };

    static Config = {
        get: function() {
            return API.auth().get('/config');
        },
        check: function () {
            return API.auth().get('/config/check');
        }
    };

    static Stats = {
        get: function() {
            return API.auth().get('/stats');
        }
    };

    static Movies = {
        tracked: function() {
            return API.auth().get('/movies/tracked');
        },
        downloading: function() {
            return API.auth().get('/movies/downloading');
        },
        downloaded: function() {
            return API.auth().get('/movies/downloaded');
        },
        removed: function() {
            return API.auth().get('/movies/removed');
        },

        get: function (id: number) {
            return API.auth().get('/movies/details/' + id);
        },
        delete: function (id: number) {
            return API.auth().delete('/movies/details/' + id);
        },
        restore: function (id: number) {
            return API.auth().post('/movies/restore/' + id);
        },
        download: function (id: number) {
            return API.auth().post('/movies/details/' + id + '/download');
        },
        skipTorrent: function (id: number) {
            return API.auth().post('/movies/details/' + id + '/download/skip_torrent');
        },
        abortDownload: function (id: number) {
            return API.auth().delete('/movies/details/' + id + '/download');
        },
        changeDownloadedState: function (id: number, downloaded_state: boolean) {
            return API.auth().put('/movies/details/' + id + '/download_state', {
                Downloaded: downloaded_state,
            });
        },
        changeCustomTitle: function (id: number, customTitle: string) {
            return API.auth().put('/movies/details/' + id + '/custom_title', {
                CustomTitle: customTitle,
            });
        },
        useDefaultTitle: function (id: number, useDefaultTitle: boolean) {
            return API.auth().put('/movies/details/' + id + '/use_default_title', {
                UseDefaultTitle: useDefaultTitle,
            });
        },
        refreshMetadata: function (id: number) {
            return API.auth().post('/movies/details/' + id + '/refresh_metadata');
        },
    };

    static Episodes = {
        downloading: function () {
            return API.auth().get('/tvshows/downloading');
        },
        get: function (id: number) {
            return API.auth().get('/tvshows/episodes/' + id);
        },
        download: function (id: number) {
            return API.auth().post('/tvshows/episodes/' + id + '/download');
        },
        skipTorrent: function (id: number) {
            return API.auth().post('/tvshows/episodes/' + id + '/download/skip_torrent');
        },
        abortDownload: function (id: number) {
            return API.auth().delete('/tvshows/episodes/' + id + '/download');
        },
        changeDownloadedState: function (id: number, downloaded_state: boolean) {
            return API.auth().put('/tvshows/episodes/' + id + '/download_state', {
                Downloaded: downloaded_state,
            });
        },
    };

    static Shows = {
        tracked: function() {
            return API.auth().get('/tvshows/tracked');
        },
        removed: function() {
            return API.auth().get('/tvshows/removed');
        },

        get: function (id: number) {
            return API.auth().get('/tvshows/details/' + id);
        },
        delete: function (id: number) {
            return API.auth().delete('/tvshows/details/' + id);
        },
        restore: function (id: number) {
            return API.auth().post('/tvshows/restore/' + id);
        },
        changeCustomTitle: function (id: number, customTitle: string) {
            return API.auth().put('/tvshows/details/' + id + '/custom_title', {
                CustomTitle: customTitle,
            });
        },
        useDefaultTitle: function (id: number, useDefaultTitle: boolean) {
            return API.auth().put('/tvshows/details/' + id + '/use_default_title', {
                UseDefaultTitle: useDefaultTitle,
            });
        },
        changeAnimeState: function (id: number, isAnime: boolean) {
            return API.auth().put('/tvshows/details/' + id + '/change_anime_state', {
                IsAnime: isAnime,
            });
        },
        getSeason: function (id: number, seasonNb: number) {
            return API.auth().get('/tvshows/details/' + id + '/seasons/' + seasonNb);
        },
        refreshMetadata: function (id: number) {
            return API.auth().post('/tvshows/details/' + id + '/refresh_metadata');
        },
    };

    static Notifications = {
        all: function() {
            return API.auth().get('/notifications/all');
        },
        read: function() {
            return API.auth().get('/notifications/read');
        },
        unread: function() {
            return API.auth().get('/notifications/unread');
        },
        markAllRead: function() {
            return API.auth().post('/notifications/read');
        },
        markRead: function (id: number) {
            return API.auth().post('/notifications/read/' +id, {
                Read: true,
            });
        },
        deleteAll: function() {
            return API.auth().delete('/notifications/all');
        },
    };

    static Modules = {
        Watchlists: {
            status: function() {
                return API.auth().get('/modules/watchlists/status');
            },
            refresh: function() {
                return API.auth().post('/modules/watchlists/refresh');
            },
            Trakt: {
                get_token: function () {
                    return API.auth().get('/modules/watchlists/trakt/token');
                },
                get_device_code: function () {
                    return API.auth().get('/modules/watchlists/trakt/devicecode');
                },
                auth: function () {
                    return API.auth().get('/modules/watchlists/trakt/auth');
                },
                get_auth_errors: function () {
                    return API.auth().get('/modules/watchlists/trakt/auth_errors');
                }
            }
        },
        Providers: {
            status: function() {
                return API.auth().get('/modules/providers/status');
            }
        },
        Notifiers: {
            status: function() {
                return API.auth().get('/modules/notifiers/status');
            },
            Telegram: {
                auth: function () {
                    return API.auth().get('/modules/notifiers/telegram/auth');
                },
                get_auth_code: function () {
                    return API.auth().get('/modules/notifiers/telegram/auth_code');
                },
                get_chat_id: function () {
                    return API.auth().get('/modules/notifiers/telegram/chatid');
                },
            },
        },
        Indexers: {
            status: function() {
                return API.auth().get('/modules/indexers/status');
            }
        },
        Downloaders: {
            status: function() {
                return API.auth().get('/modules/downloaders/status');
            }
        },
        Mediacenters: {
            status: function() {
                return API.auth().get('/modules/mediacenters/status');
            }
        },
        Scanner: {
            scan_movies: function() {
                return API.auth().get('/modules/scan/movies');
            },
            import_movie: function(movie :MediaInfo) {
                return API.auth().post('/modules/scan/movies/import', movie);
            },
            scan_shows: function() {
                return API.auth().get('/modules/scan/tvshows');
            },
            import_show: function(show :MediaInfo) {
                return API.auth().post('/modules/scan/tvshows/import', show);
            },
        },
    }
}
