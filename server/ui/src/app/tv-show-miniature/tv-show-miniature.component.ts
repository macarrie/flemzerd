import { Component, OnInit, Input } from '@angular/core';

import { TvShow } from '../tvshow';
import { TvshowsService } from '../tvshows.service';

@Component({
    selector: 'tv-show-miniature',
    templateUrl: './tv-show-miniature.component.html',
    styleUrls: ['./tv-show-miniature.component.scss']
})
export class TvShowMiniatureComponent implements OnInit {
    @Input() tvshow :TvShow;

    constructor(private tvshowsService :TvshowsService) {}

    ngOnInit() {}

    refreshShow() :void {
        this.tvshowsService.getShow(this.tvshow.ID).subscribe(tvshow => this.tvshow = tvshow);
    }

    deleteShow(id :number) {
        this.tvshowsService.deleteShow(id).subscribe(response => this.refreshShow());
    }

    restoreShow(id :number) {
        this.tvshowsService.restoreShow(id).subscribe(response => this.refreshShow());
    }
}
