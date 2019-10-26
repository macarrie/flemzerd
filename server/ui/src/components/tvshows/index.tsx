import React from "react";
import {Route} from "react-router-dom";

import API from "../../utils/api";
import Const from "../../const";
import TvShow from "../../types/tvshow";
import Episode from "../../types/episode";

import Loading from "../loading";
import MediaMiniature from "../media_miniature";
import TvShowDetails from "./tvshow_details";
import DownloadingItemTable from "../downloading_items_table";
import LoadingButton from "../loading_button";

type State = {
    tracked :TvShow[] | null,
    removed :TvShow[] | null,
    downloading_episodes :Episode[] | null,
    display_mode :number,
};

class TvShows extends React.Component<any, State> {
    tvshow_refresh_interval :number;
    state :State = {
        tracked: null,
        removed: null,
        downloading_episodes: null,
        display_mode: Const.DISPLAY_MINIATURES,
    };

    constructor(props: any) {
        super(props);

        this.tvshow_refresh_interval = 0;

        this.poll = this.poll.bind(this);
        this.refreshWatchlists = this.refreshWatchlists.bind(this);

        this.renderShowList = this.renderShowList.bind(this);
        this.getDownloadingEpisodes = this.getDownloadingEpisodes.bind(this);
        this.abortEpisodeDownload = this.abortEpisodeDownload.bind(this);
        this.skipTorrentDownload = this.skipTorrentDownload.bind(this);
    }

    componentDidMount() {
        this.load();

        this.tvshow_refresh_interval = window.setInterval(this.load.bind(this), Const.DATA_REFRESH);
    }

    componentWillUnmount() {
        clearInterval(this.tvshow_refresh_interval);
    }

    load() {
        this.getShows();
        this.getDownloadingEpisodes();
    }

    getShows() {
        API.Shows.tracked().then(response => {
            this.setState({tracked: response.data});
        }).catch(error => {
            console.log("Get tracked shows error: ", error);
        });

        API.Shows.removed().then(response => {
            this.setState({removed: response.data});
        }).catch(error => {
            console.log("Get removed shows error: ", error);
        });
    }

    poll() {
        return API.Actions.poll();
    }

    refreshWatchlists() {
        return API.Modules.Watchlists.refresh();
    }

    getDownloadingEpisodes() {
        API.Episodes.downloading().then(response => {
            this.setState({downloading_episodes: response.data});
        }).catch(error => {
            console.log("Get downloading episodes error: ", error);
        });

    }

    skipTorrentDownload(id: number) {
        API.Episodes.skipTorrent(id).then(response => {
            this.getDownloadingEpisodes();
        }).catch(error => {
            console.log("Skip torrent download error: ", error);
        });
    }

    abortEpisodeDownload(id: number) {
        API.Episodes.abortDownload(id).then(response => {
            this.getDownloadingEpisodes();
        }).catch(error => {
            console.log("Abort episode download error: ", error);
        });
    }

    changeDisplayMode(display_mode :number) {
        this.setState({
            display_mode: display_mode,
        });
    }

    renderDownloadingList() {
        if (this.state.downloading_episodes) {
            return (
                <>
                <div className="uk-width-1-1">
                    <span className="uk-h3">Downloading episodes</span>
                    <hr />
                </div>
                <DownloadingItemTable
                    list={this.state.downloading_episodes}
                    type="episode"
                    skipTorrent={this.skipTorrentDownload}
                    abortDownload={this.abortEpisodeDownload}/>
                </>
            );
        }

        return "";
    }

    renderShowListContent() {
        if (this.state.tracked === null && this.state.removed === null) {
            return null;
        }

        return (
            <>
                {this.state.tracked ? (
                    this.state.tracked.map(show => (
                        <MediaMiniature key={show.ID}
                            item={show}
                            type="tvshow"
                            display_mode={this.state.display_mode}
                        />
                    ))
                ) : ""
                }
                {this.state.removed ? (
                    this.state.removed.map(show => (
                        <MediaMiniature key={show.ID}
                            item={show}
                            type="tvshow"
                            display_mode={this.state.display_mode}
                        />
                    ))
                ) : ""
                }
            </>
        );
    }

    renderShowListLayout() {
        let filterClass :string = "item-filter";
        if (this.state.tracked === null && this.state.removed === null) {
            filterClass = "";
        }

        if (this.state.display_mode === Const.DISPLAY_MINIATURES) {
            return (
                <div className={`uk-grid uk-grid-small uk-child-width-1-2 uk-child-width-1-6@l uk-child-width-1-4@m uk-child-width-1-3@s ${filterClass}`} data-uk-grid>
                    {this.renderShowListContent()}
                </div>
            );
        } else {
            return (
                <table className="uk-table uk-table-divider uk-table-small">
                    <thead>
                        <tr>
                            <th className="uk-width-expand">Title</th>
                            <th>State</th>
                            <th>Actions</th>
                        </tr>
                    </thead>
                    <tbody className={`${filterClass}`}>
                        {this.renderShowListContent()}
                    </tbody>
                </table>
            );
        }
    }

    renderShowList() {
        if (this.state.downloading_episodes == null && this.state.removed == null && this.state.tracked == null) {
            return (
                <Loading/>
            );
        }

        return (
            <div className="uk-container" data-uk-filter="target: .item-filter">
                {this.renderDownloadingList()}

                <div className="uk-grid" data-uk-grid>
                    <div className="uk-width-expand">
                        <span className="uk-h3">TV Shows</span>
                    </div>
                    <ul className="uk-subnav uk-subnav-pill filter-container">
                        <li className="uk-visible@s" data-uk-filter-control="">
                            <button className="uk-button uk-button-text">
                                All
                            </button>
                        </li>
                        <li className="uk-visible@s uk-active" data-uk-filter-control="filter: .tracked-show">
                            <button className="uk-button uk-button-text">
                                Tracked
                            </button>
                        </li>
                        <li className="uk-visible@s" data-uk-filter-control="filter: .removed-show">
                            <button className="uk-button uk-button-text">
                                Removed
                            </button>
                        </li>
                        <li>
                            <LoadingButton
                                text="Check now"
                                loading_text="Checking"
                                action={this.poll}
                            />
                        </li>
                        <li>
                            <LoadingButton
                                text="Refresh"
                                loading_text="Refreshing"
                                action={this.refreshWatchlists}
                            />
                        </li>
                    </ul>

                    <div className="uk-button-group">
                        <button className="uk-button uk-button-default uk-button-small" onClick={() => this.changeDisplayMode(Const.DISPLAY_MINIATURES)}>
                            <span uk-icon="icon: grid; ratio: 0.7"></span>
                        </button>
                        <button className="uk-button uk-button-default uk-button-small" onClick={() => this.changeDisplayMode(Const.DISPLAY_LIST)}>
                            <span uk-icon="icon: list; ratio: 0.7"></span>
                        </button>
                    </div>
                </div>
                <hr />

                {this.renderShowListLayout()}
            </div>
        )
    }

    render() {
        const match = this.props.match;

        return (
            <>
                <Route path={`${match.path}/:id`} component={TvShowDetails} />
                <Route exact
                       path={match.path}
                       render={this.renderShowList} />
            </>
        );
    }
}

export default TvShows;
