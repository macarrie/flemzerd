import { Component, OnInit } from '@angular/core';

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

    constructor(
        private modulesService :ModulesService
    ) {}

    ngOnInit() {
        this.getProviders();
        this.getIndexers();
        this.getDownloaders();
        this.getNotifiers();
        this.getWatchlists();
        this.getMediaCenters();
    }

    getProviders() {
        this.modulesService.getProviders().subscribe(providers => {
            this.providers = providers;
        })
    }

    getIndexers() {
        this.modulesService.getIndexers().subscribe(indexers => {
            this.indexers = indexers;
        })
    }

    getDownloaders() {
        this.modulesService.getDownloaders().subscribe(downloaders => {
            this.downloaders = downloaders;
        })
    }

    getNotifiers() {
        this.modulesService.getNotifiers().subscribe(notifiers => {
            this.notifiers = notifiers;
        })
    }

    getWatchlists() {
        this.modulesService.getWatchlists().subscribe(watchlists => {
            this.watchlists = watchlists;
        })
    }

    getMediaCenters() {
        this.modulesService.getMediaCenters().subscribe(mediacenters => {
            this.mediacenters = mediacenters;
        })
    }
}
