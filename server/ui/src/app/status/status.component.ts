import { Component, OnInit, OnDestroy } from '@angular/core';

import { ModulesService } from '../modules.service';

@Component({
    selector: 'app-status',
    templateUrl: './status.component.html',
    styleUrls: ['./status.component.scss']
})
export class StatusComponent implements OnInit {
    providers :any;
    indexers :any;
    downloaders :any;
    notifiers :any;
    watchlists :any;
    mediacenters :any;

    modulesRefresh :any;

    constructor(
        private modulesService :ModulesService
    ) {
        modulesService.providers.subscribe(providers => {
            this.providers = providers;
        });
        modulesService.indexers.subscribe(indexers => {
            this.indexers = indexers;
        });
        modulesService.downloaders.subscribe(downloaders => {
            this.downloaders = downloaders;
        });
        modulesService.notifiers.subscribe(notifiers => {
            this.notifiers = notifiers;
        });
        modulesService.watchlists.subscribe(watchlists => {
            this.watchlists = watchlists;
        });
        modulesService.mediacenters.subscribe(mediacenters => {
            this.mediacenters = mediacenters;
        });
    }

    ngOnInit() {
        this.modulesService.getProviders();
        this.modulesService.getIndexers();
        this.modulesService.getDownloaders();
        this.modulesService.getNotifiers();
        this.modulesService.getWatchlists();
        this.modulesService.getMediaCenters();

        this.modulesRefresh = setInterval(() => {
        this.modulesService.getProviders();
        this.modulesService.getIndexers();
        this.modulesService.getDownloaders();
        this.modulesService.getNotifiers();
        this.modulesService.getWatchlists();
        this.modulesService.getMediaCenters();
        }, 60000)
    }

    ngOnDestroy() {
         clearInterval(this.modulesRefresh);
    }
}
