import { Controller } from "stimulus"

export default class extends Controller {
    static get targets() {
        return [];
    }

    initialize() {
        this.notificationId = this.data.get("id");
    }

    markAllAsRead() {
        fetch('/api/v1/notifications/read/', {
            method: "POST",
        }).then(response => {
            Turbolinks.visit(window.location, { action: "replace" });
        }).catch(response => { 
            UIkit.notification({
                message: 'Notifications could not be marked as read',
                status: 'danger',
                pos: 'top-right',
                timeout: 5000
            });
            console.log("Request error: ", response) ;
        });
    }

    deleteAll() {
        fetch('/api/v1/notifications/all', {
            method: "DELETE",
        }).then(response => {
            Turbolinks.visit(window.location, { action: "replace" });
        }).catch(response => { 
            console.log("Request error: ", response);
        });
    }

    markAsRead() {
        fetch('/api/v1/notifications/read/' + this.notificationId, {
            method: "POST",
            body: JSON.stringify({
                Read: true,
            })
        }).then(response => {
            Turbolinks.visit(window.location, { action: "replace" });
        }).catch(response => { 
            UIkit.notification({
                message: 'Notification could not be marked as read',
                status: 'danger',
                pos: 'top-right',
                timeout: 5000
            });
            console.log("Request error: ", response) ;
        });
    }
}
