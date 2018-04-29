flemzerd.directive('movie', function () {
    return {
        restrict: 'E',
        scope:{
            movie:'=ngModel'
        },
        templateUrl: function(elem, attr) {
            console.log(attr);
            if (attr.downloading != undefined) {
                return "/static/templates/movie.downloading.directive.template.html";
            }
            return "/static/templates/movie.directive.template.html";
        }
    }
});
