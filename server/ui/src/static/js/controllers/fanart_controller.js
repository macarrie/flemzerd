import { Controller } from "stimulus"

export default class extends Controller {
    static get targets() {
        return ["container"];
    }

    initialize() {
        this.mediaId = this.data.get("media-id");
        this.mediaType = this.data.get("type");

        this.FANART_TV_KEY = "648ff4214a1eea4416ad51417fc8a4e4";
        this.FANART_BASE_URL = "http://webservice.fanart.tv/v3/";
    }

    connect() {
        this.load();
    }

    load() {
        switch (this.mediaType) {
            case "movie":
                this.getMovieFanart(this.mediaId);
                break;
            case "tvshow":
                this.getShowFanart(this.mediaId);
                break;
            default:
                console.log("No fanart found: media type unknown (" + this.mediaType + ")");
        }
    }

    getMovieFanart(id) {
        fetch(this.FANART_BASE_URL + "movies/" + this.mediaId + "?api_key=" + this.FANART_TV_KEY, {
            method: "GET",
        }).then(response => {
            if (response.ok) {
                return response.json();
            }
            throw new Error('Network response was not ok.');
        }).then(json => {
            if (json["moviebackground"] && json["moviebackground"].length > 0) {
                this.setFanart(json["moviebackground"][0]["url"]);
            }
        }).catch(response => { 
            console.log("Could not find fanart -> Request error: ", response) ;
        });
    }

    getShowFanart() {
        fetch(this.FANART_BASE_URL + "tv/" + this.mediaId + "?api_key=" + this.FANART_TV_KEY, {
            method: "GET",
        }).then(response => {
            if (response.ok) {
                return response.json();
            }
            throw new Error('Network response was not ok.');
        }).then(json => {
            if (json["showbackground"] && json["showbackground"].length > 0) {
                this.setFanart(json["showbackground"][0]["url"]);
            }
        }).catch(response => { 
            console.log("Could not find fanart -> Request error: ", response) ;
        });
    }

    setFanart(url) {
        this.containerTarget.style.backgroundImage = "url('" + url + "')";
    }
}
