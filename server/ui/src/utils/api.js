import axios from "axios";

axios.defaults.baseURL = "/api/v1/";

export default class API {
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
        changeDownloadedState: function(id, downloaded_state) {
            return axios.put('/movies/details/' + id + '/download_state', {
                Downloaded: downloaded_state,
            });
        },
    };
}
