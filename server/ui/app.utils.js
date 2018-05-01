flemzerd.run(function($rootScope){
    utils = {};

    utils.getYear = function(movie) {
        var d = new Date(movie.Date);
        return d.getFullYear();
    };

    utils.displayDate = function(dateString) {
        var d = new Date(dateString);
        return d.toUTCString();
    };

    utils.objLength = function(obj) {
        if (typeof obj == undefined) {
            return 0;
        }
        return Object.keys(obj).length;
    };


    $rootScope.utils = utils;
});
