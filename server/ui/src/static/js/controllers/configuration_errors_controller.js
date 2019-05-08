import { Controller } from "stimulus"

export default class extends Controller {
    static get targets() {
        return ["counter"];
    }

    initialize() {
        this.refreshInterval = this.data.get("refreshInterval");
        if (this.refreshInterval == 0) {
            this.refreshInterval = 60;
        }
    }

    connect() {
        this.countErrors();
        this.startRefreshing();
    }

    disconnect() {
        this.stopRefreshing();
    }

    startRefreshing() {
        this.refreshTimer = setInterval(() => {
            this.countErrors();
        }, this.data.get("refreshInterval") * 1000)
    }

    stopRefreshing() {
        if (this.refreshTimer) {
            clearInterval(this.refreshTimer)
        }
    }

    countErrors() {
        fetch('/api/v1/config/check', {
            method: "GET",
        }).then(response => {
            return response.json()
        }).then(json => {
            this.counterTarget.innerHTML = json.length;
            if (json.length > 0) {
                this.counterTarget.classList.remove("uk-hidden");
                this.counterTarget.classList.add("alerted");
            } else {
                this.counterTarget.classList.add("uk-hidden")
            }
        }).catch(response => { 
            console.log("Request error: ", response) ;
        });
    }
}
