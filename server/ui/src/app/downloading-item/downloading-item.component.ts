import { Component, OnInit, Input } from '@angular/core';

import { DownloadingItem } from '../downloadingitem';

import { UtilsService } from '../utils.service';

@Component({
    selector: 'downloading-item',
    templateUrl: './downloading-item.component.html',
    styleUrls: ['./downloading-item.component.scss']
})
export class DownloadingItemComponent implements OnInit {
    @Input() downloadingitem :DownloadingItem;

    constructor(
        private utils :UtilsService
    ) {}

    ngOnInit() {}
}
