import { Controller } from "stimulus"

export default class extends Controller {
    static get targets() {
        return ["counter", "popup"];
    }

    initialize() {
        this.refreshInterval = this.data.get("refreshInterval");
        if (this.refreshInterval == 0) {
            this.refreshInterval = 10;
        }
    }

    connect() {
        this.countUnread();
        this.loadPopup();
        this.startRefreshing();
    }

    disconnect() {
        this.stopRefreshing();
    }

    startRefreshing() {
        this.refreshTimer = setInterval(() => {
            this.countUnread();
            this.loadPopup();
        }, this.data.get("refreshInterval") * 1000)
    }

    stopRefreshing() {
        if (this.refreshTimer) {
            clearInterval(this.refreshTimer)
        }
    }

    countUnread() {
        fetch('/api/v1/notifications/unread/', {
            method: "GET",
        }).then(response => {
            return response.json()
        }).then(json => {
            this.counterTarget.innerHTML = json.length;
            if (json.length > 0) {
                this.counterTarget.classList.add("alerted");
            }
        }).catch(response => { 
            console.log("Request error: ", response) ;
        });
    }

    loadPopup() {
        fetch('/api/v1/notifications/render/list')
            .then(response => response.text())
            .then(html => {
                this.popupTarget.innerHTML = html;
            }).catch(response => { 
                console.log("Request error: ", response) ;
            });
    }
}
