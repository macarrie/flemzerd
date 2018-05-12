flemzerd.factory('modules', ['$rootScope', '$interval', '$http', function($rootScope, $interval, $http) {
    var loadModules = function() {
        $http.get('/api/v1/modules/providers/status').then(function(response) {
            if (response.status == 200) {
                $rootScope.modules.providers = response.data;
                return;
            }
        });
        $http.get('/api/v1/modules/notifiers/status').then(function(response) {
            if (response.status == 200) {
                $rootScope.modules.notifiers = response.data;
                return;
            }
        });
        $http.get('/api/v1/modules/indexers/status').then(function(response) {
            if (response.status == 200) {
                $rootScope.modules.indexers = response.data;
                return;
            }
        });
        $http.get('/api/v1/modules/downloaders/status').then(function(response) {
            if (response.status == 200) {
                $rootScope.modules.downloaders = response.data;
                return;
            }
        });
        $http.get('/api/v1/modules/watchlists/status').then(function(response) {
            if (response.status == 200) {
                $rootScope.modules.watchlists = response.data;
                return;
            }
        });
        $http.get('/api/v1/modules/mediacenters/status').then(function(response) {
            if (response.status == 200) {
                $rootScope.modules.mediacenters = response.data;
                return;
            }
        });
    };

    $rootScope.modules = {};
    loadModules();
    $interval(loadModules, 30000);

    return {
        "modules": $rootScope.modules
    }
}]);
