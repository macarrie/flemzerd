import React from "react";
import MediaIds from "../types/media_ids";

type Props = {
    ids: MediaIds,
    type: string,
};

type State = {
    ids: MediaIds,
};

class MediaIdsComponent extends React.Component<Props, State> {
    state: State = {
        ids: {} as MediaIds,
    };

    constructor(props: Props) {
        super(props);

        this.state.ids = this.props.ids;
    }

    componentWillReceiveProps(nextProps: Readonly<Props>): void {
        this.setState({ids: nextProps.ids});
    }

    getTraktUrl(type: string) :string {
        switch (type) {
            case "movie":
                return `https://trakt.tv/search/trakt/${this.state.ids.Trakt}?id_type=movie`;
            case "tvshow":
                return `https://trakt.tv/search/trakt/${this.state.ids.Trakt}?id_type=show`;
            case "episode":
                return `https://trakt.tv/search/trakt/${this.state.ids.Trakt}?id_type=episode`;
            default:
                return "";
        }
    }

    render() {
        let ids = this.state.ids;
        let type = this.props.type;

        if (typeof ids === "undefined" || ids.ID === 0) {
            return (
                <div>Loading</div>
            );
        }

        return (
            <div className="uk-grid uk-grid-small external-links-banner" data-uk-grid>
                <div>
                    {ids.Imdb !== "" ? (
                        <a href={`https://www.imdb.com/title/${ids.Imdb}`}
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
                        <a href={this.getTraktUrl(type)}
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
                        <a href={`https://themoviedb.org/tv/${ids.Tmdb}`}
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

export default MediaIdsComponent;
