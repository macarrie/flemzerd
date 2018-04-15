flemzerd.component("status", {
    templateUrl: "/static/templates/status.template.html",
    controller: function StatusCtrl($rootScope, $scope, $http, $interval, modules) {
        var self = this;
        $scope.modules = {};

        self.refresh = function() {
            $scope.modules = $rootScope.modules;
        };

        self.refresh();
        $interval(self.refresh, 1000);

        return;
    }
});
