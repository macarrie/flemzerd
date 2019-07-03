import React from "react";

class MediaIds extends React.Component {
    getTraktUrl(type, ids) {
        switch (type) {
            case "movie":
                return "https://trakt.tv/search/trakt/{{ .item.MediaIds.Trakt }}?id_type=movie";
            case "tvshow":
                return "https://trakt.tv/search/trakt/{{ .item.MediaIds.Trakt }}?id_type=show";
            case "episode":
                return "https://trakt.tv/search/trakt/{{ .item.MediaIds.Trakt }}?id_type=episode";
            default:
                return "";
        }
    }

    render() {
        let ids = this.props.ids;
        let type = this.props.type;

        return (
            <div className="uk-grid uk-grid-small external-links-banner" data-uk-grid>
                <div>
                    {ids.Imdb !== "" ? (
                        <a  href={`https://www.imdb.com/title/${ids.Imdb}`}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="imdb">
                            <span className="uk-icon uk-icon-link" data-uk-icon="link"></span>
                            IMDB
                        </a>
                    ) : ""}
                </div>
                <div>
                    {ids.Trakt !== 0 ? (
                        <a  href={this.getTraktUrl(type, ids)}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="trakt">
                            <span className="uk-icon uk-icon-link" data-uk-icon="link"></span>
                            Trakt
                        </a>
                    ) : ""}
                </div>
                <div>
                    {ids.Tmdb !== 0 ? (
                        <a  href={`https://themoviedb.org/tv/${ids.Tmdb}`}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="tmdb">
                            <span className="uk-icon uk-icon-link" data-uk-icon="link"></span>
                            TMDB
                        </a>
                    ) : ""}
                </div>
            </div>
        );
    }
}

export default MediaIds;
