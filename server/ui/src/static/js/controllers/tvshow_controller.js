import {Controller} from "stimulus"

export default class extends Controller {
    static get targets() {
        return [];
    }

    initialize() {
        this.showId = this.data.get("id");
    }

    delete() {
        fetch('/api/v1/tvshows/details/' + this.showId, {
            method: "DELETE",
        }).then(response => {
            Turbolinks.visit(window.location, { action: "replace" });
        }).catch(response => { 
            UIkit.notification({
                message: 'TV show deletion failed',
                pos: 'top-right',
                timeout: 5000
            });
            console.log("Request error: ", response) ;
        });
    }

    restore() {
        fetch('/api/v1/tvshows/restore/' + this.showId, {
            method: "POST",
        }).then(response => {
            Turbolinks.visit(window.location, { action: "replace" });
        }).catch(response => { 
            UIkit.notification({
                message: 'TV show restoration failed',
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
        fetch('/api/v1/tvshows/details/' + this.showId + '/use_default_title', {
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

    treatAsRegularShow() {
        this.changeAnimeState(false);
    }

    treatAsAnime() {
        this.changeAnimeState(true);
    }

    changeAnimeState(state) {
        fetch('/api/v1/tvshows/details/' + this.showId + '/change_anime_state', {
            method: "PUT",
            body: JSON.stringify({
                IsAnime: state,
            })
        }).then(response => {
            Turbolinks.visit(window.location, {action: "replace"});
        }).catch(response => {
            UIkit.notification({
                message: 'Could not change anime status for show (error during request)',
                status: 'danger',
                pos: 'top-right',
                timeout: 5000
            });
            console.log("Request error: ", response);
        });
    }
}
