flemzerd.directive('downloadingitem', function() {
    return {
        restrict: 'E',
        link: function($scope, $element, $attributes) {
            $scope.item = $scope.$eval($attributes.item);
            $scope.downloaded = $scope.$eval($attributes.downloaded);
        },
        templateUrl: "/static/templates/downloadingitem.directive.template.html"
    };
});
