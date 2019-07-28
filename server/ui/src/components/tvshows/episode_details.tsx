import React from "react";

import API from "../../utils/api";
import Helpers from "../../utils/helpers";
import Const from "../../const";
import Episode from "../../types/episode";

import Loading from "../loading";
import MediaIdsComponent from "../media_ids";
import MediaActionBar from "../media_action_bar";
import DownloadingItemComponent from "../downloading_item";

type State = {
    episode :Episode | null,
    fanartURL :string,
};

class EpisodeDetails extends React.Component<any, State> {
    episode_refresh_interval :number;
    state :State = {
        episode: null,
        fanartURL: "",
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
            this.getFanart();
        }).catch(error => {
            console.log("Get episode details error: ", error);
        });
    }

    getFanart() {
        if (this.state.episode == null) {
            return;
        }

        API.Fanart.getTvShowFanart(this.state.episode.TvShow).then(response => {
            let fanartObj = response.data;
            if (fanartObj["showbackground"] && fanartObj["showbackground"].length > 0) {
                this.setState({fanartURL: fanartObj["showbackground"][0]["url"]});
            }
        }).catch(error => {
            console.log("Get show fanart error: ", error);
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
                <Loading/>
            );
        }

        return (
            <>
            <div id="full_background"
                style={{backgroundImage: `url(${this.state.fanartURL})`}}></div>
            <MediaActionBar item={this.state.episode}
                            downloadItem={this.downloadEpisode}
                            markNotDownloaded={this.markNotDownloaded}
                            markDownloaded={this.markDownloaded}
                            skipTorrent={this.skipTorrent}
                            abortDownload={this.abortDownload}
                            type="episode"/>

            <div className="uk-container uk-light mediadetails">
                <div className="uk-grid" data-uk-grid>
                    <div className="uk-width-1-3">
                        <img width="100%" src={this.state.episode.TvShow.Poster} alt="{this.state.epîsode.Title}" className="uk-border-rounded" data-uk-img />
                    </div>
                    <div className="uk-width-expand">
                        <div className="uk-grid uk-grid-collapse uk-flex uk-flex-middle" data-uk-grid>
                            <div className="inlineedit">
                                <span className="uk-h2">
                                    {Helpers.getMediaTitle(this.state.episode.TvShow)} - {this.state.episode.Title}
                                </span>
                            </div>
                            <div>
                                <span
                                    className="uk-h3 uk-text-muted uk-margin-small-left media_title_details">({this.getEpisodeNumber()})</span>
                            </div>
                        </div>
                        <div className="uk-grid uk-grid-medium" data-uk-grid>
                            <div> See on </div>
                            <div><MediaIdsComponent ids={this.state.episode.MediaIds} type="episode"/></div>
                        </div>

                        <div className="container">
                            <div className="row">
                                <h5>Overview</h5>
                                <div className="col-12">
                                    {this.state.episode.Overview}
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <div className="uk-container uk-margin-medium-top uk-margin-medium-bottom">
                <DownloadingItemComponent item={this.state.episode.DownloadingItem}/>
            </div>
            </>
        );
    }
}

export default EpisodeDetails;
