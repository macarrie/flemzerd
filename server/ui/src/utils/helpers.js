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

    static getMovieTitle = (movie) => {
        if (movie.CustomTitle !== "") {
            return movie.CustomTitle;
        }

        if (movie.UseDefaultTitle) {
            return movie.Title;
        }

        return movie.OriginalTitle;
    };

};
