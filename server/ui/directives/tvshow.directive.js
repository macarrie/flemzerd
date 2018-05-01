flemzerd.directive('tvshow', function () {
    return {
        restrict: 'E',
        templateUrl: function(elem, attr) {
            if (attr.detailed != undefined) {
                return "/static/templates/tvshow.detailed.directive.template.html";
            }
            if (attr.downloading != undefined) {
                return "/static/templates/tvshow.downloading.directive.template.html";
            }
            return "/static/templates/tvshow.directive.template.html";
        }
    }
});
