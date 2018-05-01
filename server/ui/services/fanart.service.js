flemzerd.factory('fanart', ['$http', function($http) {
    var FANART_TV_KEY = "648ff4214a1eea4416ad51417fc8a4e4"
    var FANART_BASE_URL = "http://webservice.fanart.tv/v3/"
    //var FANART_BASE_URL = http://webservice.fanart.tv/v3/movies/1726?api_key=

    var GetMovieFanart = function(movie) {
        params = {
            "params": {
                "api_key": FANART_TV_KEY
            }
        };

        return $http.get(FANART_BASE_URL + "movies/" + movie.MediaIds.Tmdb, params).then(function(response) {
            return response.data;
        });
    };

    return {
        GetMovieFanart: GetMovieFanart
    }
}]);
