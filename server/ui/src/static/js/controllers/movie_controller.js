import { Controller } from "stimulus"

export default class extends Controller {
    static get targets() {
        return [];
    }

    initialize() {
        this.movieId = this.data.get("id");
    }

    delete() {
        fetch('/api/v1/movies/details/' + this.movieId, {
            method: "DELETE",
        }).then(response => {
            Turbolinks.visit(window.location, { action: "replace" });
        }).catch(response => { 
            UIkit.notification({
                message: 'Movie deletion failed',
                pos: 'top-right',
                timeout: 5000
            });
            console.log("Request error: ", response) ;
        });
    }

    restore() {
        fetch('/api/v1/movies/restore/' + this.movieId, {
            method: "POST",
        }).then(response => {
            Turbolinks.visit(window.location, { action: "replace" });
        }).catch(response => { 
            UIkit.notification({
                message: 'Movie restore failed',
                status: 'danger',
                pos: 'top-right',
                timeout: 5000
            });
            console.log("Request error: ", response) ;
        });
    }

    download() {
        UIkit.notification({
            message: 'Sending download request',
            status: 'primary',
            pos: 'top-right',
            timeout: 5000
        });

        fetch('/api/v1/movies/details/' + this.movieId + '/download', {
            method: "POST",
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

    markDownloaded() {
        this.changeDownloadedState(true);
    }

    markNotDownloaded() {
        this.changeDownloadedState(false);
    }

    changeDownloadedState(state) {
        fetch('/api/v1/movies/details/' + this.movieId + '/download_state', {
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

    useOriginalTitle() {
        this.useDefaultTitleChange(false);
    }

    useDefaultTitle() {
        this.useDefaultTitleChange(true);
    }

    useDefaultTitleChange(state) {
        fetch('/api/v1/movies/details/' + this.movieId + '/use_default_title', {
            method: "PUT",
            body: JSON.stringify({
                UseDefaultTitle: state,
            })
        }).then(response => {
            Turbolinks.visit(window.location, { action: "replace" });
        }).catch(response => { 
            UIkit.notification({
                message: 'Could not change title used for movie (error during request)',
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

        fetch('/api/v1/movies/details/' + this.movieId + '/download', {
            method: "DELETE",
        }).then(response => {
            Turbolinks.visit(window.location, { action: "replace" });
        }).catch(response => { 
            UIkit.notification({
                message: 'Movie download abort failed',
                status: 'danger',
                pos: 'top-right',
                timeout: 5000
            });
            console.log("Request error: ", response) ;
        });
    }
}
