flemzerd.run(function($rootScope){
    utils = {};

    utils.getYear = function(datestring) {
        var d = new Date(datestring);
        return d.getFullYear();
    };

    utils.displayDate = function(dateString) {
        var d = new Date(dateString);
        return d.toUTCString();
    };

    utils.objLength = function(obj) {
        if (typeof obj == undefined || obj == null) {
            return 0;
        }
        return Object.keys(obj).length;
    };

    utils.formatNumber = function(nb) {
        if (nb < 10) {
            return '0' + nb;
        }
        return nb;
    };


    $rootScope.utils = utils;
});
