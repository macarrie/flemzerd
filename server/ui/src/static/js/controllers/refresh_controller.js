import { Controller } from "stimulus"

export default class extends Controller {
    static get targets() {
        return [];
    }

    initialize() {
        //this.movieId = this.data.get("id");
    }

    connect() {
        if (this.data.has("interval")) {
            this.refreshInterval = this.data.get("interval");
        } else {
            this.refreshInterval = 30
        }
        this.refreshTimer = setInterval(() => {
            Turbolinks.visit(window.location, { action: "replace" });
        }, this.refreshInterval * 1000);
    }

    disconnect() {
        clearInterval(this.refreshTimer);
    }
}
