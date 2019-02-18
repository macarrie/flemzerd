import { Controller } from "stimulus"

export default class extends Controller {
    static get targets() {
        return [];
    }

    initialize() {
        this.showId = this.data.get("showid");
        this.seasonNumber = this.data.get("number");
    }

    download() {
        UIkit.notification({
            message: 'Sending season download request',
            status: 'primary',
            pos: 'top-right',
            timeout: 5000
        });

        fetch('/api/v1/tvshows/details/' + this.showId + '/seasons/' + this.seasonNumber + '/download', {
            method: "POST",
        }).then(response => {
            Turbolinks.visit(window.location, { action: "replace" });
        }).catch(response => { 
            UIkit.notification({
                message: 'Season download request failed',
                status: 'danger',
                pos: 'top-right',
                timeout: 5000
            });
            console.log("Request error: ", response) ;
        });
    }

    markDownloaded() {
        this.changeDownloadedState(true);
    }

    markNotDownloaded() {
        this.changeDownloadedState(false);
    }

    changeDownloadedState(state) {
        fetch('/api/v1/tvshows/details/' + this.showId + '/seasons/' + this.seasonNumber + '/download_state', {
            method: "PUT",
            body: JSON.stringify({
                Downloaded: state,
            })
        }).then(response => {
            Turbolinks.visit(window.location, { action: "replace" });
        }).catch(response => { 
            UIkit.notification({
                message: 'Changing season downloaded state',
                status: 'danger',
                pos: 'top-right',
                timeout: 5000
            });
            console.log("Request error: ", response) ;
        });
    }
}
