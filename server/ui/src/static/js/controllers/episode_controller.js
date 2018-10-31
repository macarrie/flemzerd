import { Controller } from "stimulus"

export default class extends Controller {
    static get targets() {
        return [];
    }

    initialize() {
        this.episodeId = this.data.get("id");
    }

    download() {
        UIkit.notification({
            message: 'Sending download request',
            status: 'primary',
            pos: 'top-right',
            timeout: 5000
        });

        fetch('/api/v1/tvshows/episodes/' + this.episodeId + '/download', {
            method: "POST",
        }).then(response => {
            Turbolinks.visit(window.location, { action: "replace" });
        }).catch(response => { 
            UIkit.notification({
                message: 'Episode download request failed',
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
        fetch('/api/v1/tvshows/episodes/' + this.episodeId + '/download_state', {
            method: "PUT",
            body: JSON.stringify({
                Downloaded: state,
            })
        }).then(response => {
            Turbolinks.visit(window.location, { action: "replace" });
        }).catch(response => { 
            UIkit.notification({
                message: 'Movie download request failed',
                status: 'danger',
                pos: 'top-right',
                timeout: 5000
            });
            console.log("Request error: ", response) ;
        });
    }

    abortDownload() {
        UIkit.notification({
            message: 'Aborting download',
            status: 'primary',
            pos: 'top-right',
            timeout: 5000
        });

        fetch('/api/v1/tvshows/episodes/' + this.episodeId + '/download', {
            method: "DELETE",
        }).then(response => {
            Turbolinks.visit(window.location, { action: "replace" });
        }).catch(response => { 
            UIkit.notification({
                message: 'Episode download abort failed',
                status: 'danger',
                pos: 'top-right',
                timeout: 5000
            });
            console.log("Request error: ", response) ;
        });
    }
}
