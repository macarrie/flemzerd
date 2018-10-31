import { Controller } from "stimulus"

export default class extends Controller {
    static get targets() {
        return ["title", "controls", "field", "submit", "cancel"];
    }

    initialize() {
        this.id = this.data.get("id");
        this.type = this.data.get("type");

        if (this.type == "movie") {
            this.url = '/api/v1/movies/details/' + this.id + '/custom_title';
        } else {
            this.url = '/api/v1/tvshows/details/' + this.id + '/custom_title';
        }

        this.titleTarget.classList.remove("uk-hidden");
        this.submitTarget.classList.remove("uk-hidden");
        this.cancelTarget.classList.remove("uk-hidden");
        this.controlsTarget.classList.add("uk-hidden");
    }

    displaycontrols() {
        this.titleTarget.classList.add("uk-hidden");
        this.controlsTarget.classList.remove("uk-hidden");
        this.fieldTarget.focus();
    }

    closeControls() {
        this.titleTarget.classList.remove("uk-hidden");
        this.controlsTarget.classList.add("uk-hidden");
    }

    submit() {
        this.submitTarget.classList.add("uk-hidden");
        this.cancelTarget.classList.add("uk-hidden");

        fetch(this.url, {
            method: "PUT",
            body: JSON.stringify({
                CustomTitle: this.fieldTarget.value,
            })
        }).then(response => {
            Turbolinks.visit(window.location, { action: "replace" });
        }).catch(response => { 
            UIkit.notification({
                message: 'Could not change custom title',
                status: 'danger',
                pos: 'top-right',
                timeout: 5000
            });
            console.log("Request error: ", response) ;
        });
    }
}
