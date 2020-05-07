import React from "react";

import API from "../../utils/api";
import Helpers from "../../utils/helpers";
import Const from "../../const";
import Episode from "../../types/episode";

import Empty from "../empty";
import MediaIdsComponent from "../media_ids";
import MediaActionBar from "../media_action_bar";
import DownloadingItemComponent from "../downloading_item";

type State = {
    episode :Episode | null,
};

class EpisodeDetails extends React.Component<any, State> {
    episode_refresh_interval :number;
    state :State = {
        episode: null,
    };

    constructor(props: any) {
        super(props);

        this.episode_refresh_interval = 0;

        this.markNotDownloaded = this.markNotDownloaded.bind(this);
        this.markDownloaded = this.markDownloaded.bind(this);
        this.changeDownloadedState = this.changeDownloadedState.bind(this);

        this.downloadEpisode = this.downloadEpisode.bind(this);
        this.abortDownload = this.abortDownload.bind(this);
        this.skipTorrent = this.skipTorrent.bind(this);

        this.getEpisodeNumber = this.getEpisodeNumber.bind(this);
    }

    componentDidMount() {
        this.getEpisode();

        this.episode_refresh_interval = window.setInterval(this.getEpisode.bind(this), Const.DATA_REFRESH);
    }

    componentWillUnmount() {
        clearInterval(this.episode_refresh_interval);
    }

    getEpisode() {
        API.Episodes.get(this.props.match.params.id).then(response => {
            let episode_result = response.data;
            episode_result.DisplayTitle = Helpers.getMediaTitle(episode_result)
            console.log(episode_result);
            this.setState({episode: episode_result});
        }).catch(error => {
            console.log("Get episode details error: ", error);
        });
    }

    downloadEpisode() {
        if (this.state.episode == null) {
            return;
        }

        API.Episodes.download(this.state.episode.ID).then(response => {
            this.getEpisode();
        }).catch(error => {
            console.log("Download episode error: ", error);
        });
    }

    markNotDownloaded() {
        this.changeDownloadedState(false);
    }

    markDownloaded() {
        this.changeDownloadedState(true);
    }

    changeDownloadedState(downloaded_state: boolean) {
        if (this.state.episode == null) {
            return;
        }

        API.Episodes.changeDownloadedState(this.state.episode.ID, downloaded_state).then(response => {
            this.getEpisode();
        }).catch(error => {
            console.log("Change episode downloaded state error: ", error);
        });
    }

    skipTorrent() {
        if (this.state.episode == null) {
            return;
        }

        API.Episodes.skipTorrent(this.state.episode.ID).then(response => {
            this.getEpisode();
        }).catch(error => {
            console.log("Skip torrent download error: ", error);
        });
    }

    abortDownload() {
        if (this.state.episode == null) {
            return;
        }

        API.Episodes.abortDownload(this.state.episode.ID).then(response => {
            this.getEpisode();
        }).catch(error => {
            console.log("Abort episode download error: ", error);
        });
    }

    getEpisodeNumber() {
        if (this.state.episode == null) {
            return;
        }

        if (this.state.episode.TvShow.IsAnime) {
            return (
                <>
                    {Helpers.formatAbsoluteNumber(this.state.episode.AbsoluteNumber)}
                </>
            );
        }

        return (
            <>
                S{Helpers.formatNumber(this.state.episode.Season)}E{Helpers.formatNumber(this.state.episode.Number)}
            </>
        );
    }

    render() {
        if (this.state.episode == null) {
            return (
                <div className={"container"}>
                    <Empty label={"Loading"}/>
                </div>
            );
        }

        return (
            <>
            <div id="full_background" className={"has-background-dark"}
                style={{backgroundImage: `url(${this.state.episode?.TvShow.Background})`}}></div>
            <MediaActionBar item={this.state.episode}
                            downloadItem={this.downloadEpisode}
                            markNotDownloaded={this.markNotDownloaded}
                            markDownloaded={this.markDownloaded}
                            skipTorrent={this.skipTorrent}
                            abortDownload={this.abortDownload}
                            type="episode"/>

            <div className="container mediadetails">
                <div className="columns is-mobile">
                    <div className="column is-one-quarter-desktop is-one-third-tablet is-hidden-mobile">
                        <div className={"poster-container"}>
                            <span className={"poster-alt"}>{this.state.episode.Title}</span>
                            <img width="100%" src={this.state.episode.TvShow.Poster} alt={""} className="poster" />
                        </div>
                    </div>
                    <div className="column">
                        <div className="columns">
                            <div className="column is-narrow inlineedit">
                                <span className="title is-2 has-text-grey-light">
                                    {Helpers.getMediaTitle(this.state.episode.TvShow)} - {this.state.episode.Title}
                                </span>
                            </div>
                            <div className={"column"}>
                                <span className="title is-3 has-text-grey-light media_title_details">({this.getEpisodeNumber()})</span>
                            </div>
                        </div>
                        <div className="columns is-mobile">
                            <div className={"column is-narrow"}> See on </div>
                            <div className={"column"}><MediaIdsComponent ids={this.state.episode.MediaIds} type="episode"/></div>
                        </div>

                        <div className="container">
                            <div className="row">
                                <h5 className={"title is-5 has-text-grey-light"}>Overview</h5>
                                <div className="col-12">
                                    {this.state.episode.Overview}
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>


            <div className="container">
                <div className="columns is-mobile">
                    <div className={"column is-full"}>
                        <div className={"box"}>
                            <DownloadingItemComponent item={this.state.episode.DownloadingItem}/>
                        </div>
                    </div>
                </div>
            </div>
            </>
        );
    }
}

export default EpisodeDetails;
