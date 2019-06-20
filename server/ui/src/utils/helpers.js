import * as moment from "moment";

export default class Helpers {
    static getYear = (datestring) => {
        var d = new Date(datestring);
        return d.getFullYear();
    };

    static formatDate = (date, format) => {
        return moment(date).format(format);
    };

    static formatNumber = (nb) => {
        if (nb < 10) {
            return '0' + nb;
        }
        return nb;
    };

    static dateIsInFuture = (dateString) => {
        var d = new Date(dateString);
        return d.getTime() > Date.now();
    };

    static getMediaTitle = (media) => {
        if (media.CustomTitle !== "") {
            return media.CustomTitle;
        }

        if (media.UseDefaultTitle) {
            return media.Title;
        }

        return media.OriginalTitle;
    };

};
