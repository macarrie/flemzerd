import { Controller } from "stimulus"

export default class extends Controller {
    static get targets() {
        return ["label", "wait"];
    }

    initialize() {
        this.url = this.data.get("url");
        this.method = this.data.get("method");

        this.labelTarget.style.display = "block";
        this.waitTarget.style.display = "none";
    }

    send() {
        this.labelTarget.style.display = "none";
        this.waitTarget.style.display = "block";

        if (!this.data.has("url")) {
            console.error("No URL defined for controller");
            return;
        }

        if (!this.data.has("method")) {
            this.method = "GET";
        }

        fetch(this.url, {
            method: this.method,
        }).then(response => {
            setTimeout(() => {
                this.labelTarget.style.display = "block";
                this.waitTarget.style.display = "none";
                Turbolinks.visit(window.location, { action: "replace" });
            }, 1000);
        }).catch(response => { 
            console.log("Request error: ", response) ;
        });
    }
}
