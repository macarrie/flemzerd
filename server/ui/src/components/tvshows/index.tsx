import React from "react";
import {Route} from "react-router-dom";

import API from "../../utils/api";
import Helpers from "../../utils/helpers";
import Const from "../../const";
import TvShow from "../../types/tvshow";
import Episode from "../../types/episode";
import {MediaMiniatureFilter} from "../../types/media_miniature_filter";

import Empty from "../empty";
import MediaMiniature from "../media_miniature";
import TvShowDetails from "./tvshow_details";
import DownloadingItemTable from "../downloading_items_table";
import LoadingButton from "../loading_button";
import ItemFilterControls from "../item_filter_controls";

import {RiLayoutGridLine, RiLayoutRowLine} from "react-icons/ri";
import ItemSearchControls from "../item_search_controls";

type State = {
    tracked :TvShow[] | null,
    removed :TvShow[] | null,
    downloading_episodes :Episode[] | null,
    display_mode :number,
    item_filter :MediaMiniatureFilter,
    search_filter :string,
};

class TvShows extends React.Component<any, State> {
    tvshow_refresh_interval :number;
    state :State = {
        tracked: null,
        removed: null,
        downloading_episodes: null,
        display_mode: Const.DISPLAY_MINIATURES,
        item_filter: MediaMiniatureFilter.TRACKED,
        search_filter: "",
    };

    constructor(props: any) {
        super(props);

        this.tvshow_refresh_interval = 0;

        this.poll = this.poll.bind(this);
        this.refreshWatchlists = this.refreshWatchlists.bind(this);

        this.getActiveLayoutClass = this.getActiveLayoutClass.bind(this);
        this.setFilter = this.setFilter.bind(this);
        this.setSearch = this.setSearch.bind(this);
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

    setFilter(val :MediaMiniatureFilter) {
        this.setState({item_filter: val});
    }

    setSearch(val :string) {
        this.setState({search_filter: val});
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

    getActiveLayoutClass(targetLayout :number) :String {
        if (targetLayout === this.state.display_mode) {
            return "is-info is-outlined";
        }

        return "";
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
                    <div className="columns">
                        <div className="column">
                            <span className="title is-3">Downloading episodes</span>
                        </div>
                    </div>
                    <hr />
                    <DownloadingItemTable
                        list={this.state.downloading_episodes}
                        type="episode"
                        skipTorrent={this.skipTorrentDownload}
                        abortDownload={this.abortEpisodeDownload}/>
                    <hr />
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
                            filter={this.state.item_filter}
                            search={this.state.search_filter}
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
                            filter={this.state.item_filter}
                            search={this.state.search_filter}
                            display_mode={this.state.display_mode}
                        />
                    ))
                ) : ""
                }
            </>
        );
    }

    renderShowListLayout() {
        let total_count = Helpers.count(this.state.tracked) + Helpers.count(this.state.removed);
        if (this.state.item_filter === MediaMiniatureFilter.NONE && total_count === 0) {
            return <Empty label={"No shows"} />;
        }
        if (this.state.item_filter === MediaMiniatureFilter.TRACKED && Helpers.count(this.state.tracked) === 0) {
            return <Empty label={"No tracked shows"} />;
        }
        if (this.state.item_filter === MediaMiniatureFilter.REMOVED && Helpers.count(this.state.removed) === 0) {
            return <Empty label={"No removed shows"} />;
        }

        if (this.state.display_mode === Const.DISPLAY_MINIATURES) {
            return (
                <div className={`columns is-multiline is-mobile`}>
                    {this.renderShowListContent()}
                </div>
            );
        } else {
            return (
                <table className="table is-fullwidth media-miniature-table">
                    <thead>
                        <tr>
                            <th>Title</th>
                            <th className={"has-text-centered"}>State</th>
                            <th className={"has-text-centered is-narrow"}>Actions</th>
                        </tr>
                    </thead>
                    <tbody>
                        {this.renderShowListContent()}
                    </tbody>
                </table>
            );
        }
    }

    renderShowList() {
        if (this.state.downloading_episodes == null && this.state.removed == null && this.state.tracked == null) {
            return (
                <div className={"container"}>
                    <Empty label={"Loading"}/>
                </div>
            );
        }

        return (
            <div className="container">
                {this.renderDownloadingList()}

                <div className="columns is-vcentered is-multiline is-mobile">
                    <div className="column is-narrow">
                        <span className="title is-3">TV Shows</span>
                    </div>
                    <ItemSearchControls
                        updateSearch={this.setSearch}
                    />
                    <ItemFilterControls
                        type={"tvshow"}
                        filterValue={this.state.item_filter}
                        updateFilter={this.setFilter}
                    />
                    <div className="column is-narrow field is-marginless is-grouped">
                            <LoadingButton
                                text="Check now"
                                loading_text="Checking"
                                action={this.poll}
                            />
                            <LoadingButton
                                text="Refresh"
                                loading_text="Refreshing"
                                action={this.refreshWatchlists}
                            />
                    </div>
                    <div className="column is-narrow is-hidden-mobile">
                        <div className="buttons has-addons">
                            <button className={`button ${this.getActiveLayoutClass(Const.DISPLAY_MINIATURES)}`} onClick={() => this.changeDisplayMode(Const.DISPLAY_MINIATURES)}>
                                <span className="icon is-small">
                                    <RiLayoutGridLine />
                                </span>
                            </button>
                            <button className={`button ${this.getActiveLayoutClass(Const.DISPLAY_LIST)}`} onClick={() => this.changeDisplayMode(Const.DISPLAY_LIST)}>
                                <span className="icon is-small">
                                    <RiLayoutRowLine />
                                </span>
                            </button>
                        </div>
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
