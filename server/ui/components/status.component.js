flemzerd.component("status", {
    templateUrl: "/static/templates/status.template.html",
    controller: function StatusCtrl($http) {
        var self = this;

        $http.get('/api/v1/modules/providers/status').then(function(response) {
            if (response.status == 200) {
                self.providers = response.data;
                return;
            }
        });
        $http.get('/api/v1/modules/notifiers/status').then(function(response) {
            if (response.status == 200) {
                self.notifiers = response.data;
                return;
            }
        });
        $http.get('/api/v1/modules/indexers/status').then(function(response) {
            if (response.status == 200) {
                self.indexers = response.data;
                return;
            }
        });
        $http.get('/api/v1/modules/downloaders/status').then(function(response) {
            if (response.status == 200) {
                self.downloaders = response.data;
                return;
            }
        });
        $http.get('/api/v1/modules/watchlists/status').then(function(response) {
            if (response.status == 200) {
                self.watchlists = response.data;
                return;
            }
        });
    }
});
