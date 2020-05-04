import React from "react";
import MediaIds from "../types/media_ids";

import {FaExternalLinkAlt} from "react-icons/fa";

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
                <div className={"has-text-grey"}>---</div>
            );
        }

        return (
            <div className="columns external-links-banner">
                <div className={"column is-narrow"}>
                    {ids.Imdb !== "" ? (
                        <a href={`https://www.imdb.com/title/${ids.Imdb}`}
                           target="_blank"
                           rel="noopener noreferrer"
                           className="imdb">
                            <span className="icon is-small">
                                    <FaExternalLinkAlt/>
                            </span>
                            IMDB
                        </a>
                    ) : ""}
                </div>
                <div className={"column is-narrow"}>
                    {ids.Trakt !== 0 ? (
                        <a href={this.getTraktUrl(type)}
                           target="_blank"
                           rel="noopener noreferrer"
                           className="trakt">
                            <span className="icon is-small">
                                    <FaExternalLinkAlt/>
                            </span>
                            Trakt
                        </a>
                    ) : ""}
                </div>
                <div className={"column is-narrow"}>
                    {ids.Tmdb !== 0 ? (
                        <a href={`https://themoviedb.org/tv/${ids.Tmdb}`}
                           target="_blank"
                           rel="noopener noreferrer"
                           className="tmdb">
                            <span className="icon is-small">
                                    <FaExternalLinkAlt/>
                            </span>
                            TMDB
                        </a>
                    ) : ""}
                </div>
            </div>
        );
    }
}

export default MediaIdsComponent;
