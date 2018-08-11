import { Component, OnInit, OnDestroy } from '@angular/core';
import { Observable, of, Subject } from 'rxjs';
import { takeUntil  } from 'rxjs/operators'; // for rxjs ^5.5.0 lettable operators

import { ModulesService } from '../modules.service';

@Component({
    selector: 'app-status',
    templateUrl: './status.component.html',
    styleUrls: ['./status.component.scss']
})
export class StatusComponent implements OnInit, OnDestroy {
    private subs :Subject<any> = new Subject();

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
        modulesService.providers.pipe(takeUntil(this.subs)).subscribe(providers => {
            this.providers = providers;
        });
        modulesService.indexers.pipe(takeUntil(this.subs)).subscribe(indexers => {
            this.indexers = indexers;
        });
        modulesService.downloaders.pipe(takeUntil(this.subs)).subscribe(downloaders => {
            this.downloaders = downloaders;
        });
        modulesService.notifiers.pipe(takeUntil(this.subs)).subscribe(notifiers => {
            this.notifiers = notifiers;
        });
        modulesService.watchlists.pipe(takeUntil(this.subs)).subscribe(watchlists => {
            this.watchlists = watchlists;
        });
        modulesService.mediacenters.pipe(takeUntil(this.subs)).subscribe(mediacenters => {
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
        this.subs.next();
        this.subs.complete();
    }
}
