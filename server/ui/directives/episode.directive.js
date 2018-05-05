flemzerd.directive('episode', function () {
    return {
        restrict: 'E',
        templateUrl: function(elem, attr) {
            if (attr.detailed != undefined) {
                return "/static/templates/episode.detailed.directive.template.html";
            }
            if (attr.downloading != undefined) {
                return "/static/templates/tvshow.downloading.directive.template.html";
            }
            return "/static/templates/tvshow.directive.template.html";
        }
    }
});
