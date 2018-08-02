import { Injectable } from '@angular/core';

import * as moment from 'moment';

@Injectable({
    providedIn: 'root'
})
export class UtilsService {
    constructor() {}

    getYear(datestring) :number {
        var d = new Date(datestring);
        return d.getFullYear();
    }

    formatDate(date, format) :string {
        return moment(date).format(format);
    }

    formatNumber(nb) :string {
        if (nb < 10) {
            return '0' + nb;
        }
        return nb;
    }

    dateIsInFuture(dateString) :boolean {
        var d = new Date(dateString);
        return d.getTime() > Date.now();
    }
}
