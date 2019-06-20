import axios from "axios";

axios.defaults.baseURL = "/api/v1/";

export default class API {
    static Config = {
        get: function() {
            return axios.get('/config');
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
        }
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
