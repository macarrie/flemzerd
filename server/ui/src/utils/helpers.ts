import * as moment from "moment";

import DownloadingItem from "../types/downloading_item";
import Torrent from "../types/torrent";

export default class Helpers {
    static getYear = (datestring: Date) => {
        let d = new Date(datestring);
        return d.getFullYear();
    };

    static formatDate = (date: Date, format: string) => {
        return moment.default(date).format(format);
    };

    static formatNumber = (nb: number) => {
        if (nb < 10) {
            return '0' + nb;
        }
        return nb;
    };

    static formatAbsoluteNumber = (nb: number) => {
        if (nb < 10) {
            return '00' + nb;
        }
        if (nb < 100) {
            return '0' + nb;
        }

        return nb;
    };

    static dateIsInFuture = (dateString: Date) => {
        let d = new Date(dateString);
        return d.getTime() > Date.now();
    };

    static getMediaTitle = (media: any): string => {
        if (media.ID === 0) {
            return "Unknown media"
        }

        if (media.hasOwnProperty("Season")) {
            return `${Helpers.getMediaTitle(media.TvShow)} S${Helpers.formatNumber(media.Season)}E${Helpers.formatNumber(media.Number)} - ${media.Title}`;
        }

        if (media.CustomTitle !== "") {
            return media.CustomTitle;
        }

        if (media.UseDefaultTitle) {
            return media.Title;
        }

        return media.OriginalTitle;
    };

    static getCurrentTorrent(d :DownloadingItem) :Torrent {
        for (let i in d.TorrentList) {
            console.log("Torrent in torrent list: ", d.TorrentList[i]);
            if (!d.TorrentList[i].Failed) {
                return d.TorrentList[i];
            }
        }

        return {} as Torrent;
    };

    static getFailedTorrents(d :DownloadingItem) :Torrent[] {
        let retList :Torrent[] = [] as Torrent[];
        for (let i in d.TorrentList) {
            console.log("Torrent in torrent list: ", d.TorrentList[i]);
            if (d.TorrentList[i].Failed) {
                retList.push(d.TorrentList[i]);
            }
        }

        return retList;
    };

    static createCookie(name :string, value :string, days :number) {
        let expires :string;
        if (days) {
            let date = new Date();
            date.setTime(date.getTime()+(days*24*60*60*1000));
            expires = "; expires="+date.toUTCString();
        } else {
            expires = "";
        }

        document.cookie = name+"="+value+expires+"; path=/";
    };

    static readCookie(name :string) {
        let nameEQ = name + "=";
        let ca = document.cookie.split(';');
        for(let i=0;i < ca.length;i++) {
            let c = ca[i];
            while (c.charAt(0) === ' ') {
                c = c.substring(1,c.length);
            }

            if (c.indexOf(nameEQ) === 0) {
                return c.substring(nameEQ.length, c.length);
            }
        }
        return null;
    };

    static eraseCookie(name :string) {
        Helpers.createCookie(name,"",-1);
    };
};
