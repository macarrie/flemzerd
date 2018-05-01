flemzerd.directive('movie', function () {
    return {
        restrict: 'E',
        templateUrl: function(elem, attr) {
            if (attr.detailed != undefined) {
                return "/static/templates/movie.detailed.directive.template.html";
            }
            if (attr.downloading != undefined) {
                return "/static/templates/movie.downloading.directive.template.html";
            }
            return "/static/templates/movie.directive.template.html";
        }
    }
});
