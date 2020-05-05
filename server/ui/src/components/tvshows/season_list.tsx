import React from "react";

import API from "../../utils/api";
import Helpers from "../../utils/helpers";
import {SeasonDetails, TvSeason} from "../../types/tvshow";

import EpisodeTable from "./episode_table";

import {RiArrowDropUpLine} from "react-icons/ri";
import {RiArrowDropDownLine} from "react-icons/ri";
import {RiCheckDoubleLine} from "react-icons/ri";
import {RiDownloadLine} from "react-icons/ri";
import {RiEraserLine} from "react-icons/ri";

type Props = {
    refreshSeason(n: number): void,
    seasons :SeasonDetails[] | null,
};

type State = {
    seasons :SeasonDetails[] | null,
};

class SeasonList extends React.Component<Props, State> {
    state :State = {
        seasons: null,
    };

    constructor(props: Props) {
        super(props);

        this.state.seasons = this.props.seasons;
    }

    componentWillReceiveProps(nextProps :Props) {
        if (nextProps.seasons === null) {
            return;
        }
        console.log("Season list:", this.state.seasons);
        let seasons = nextProps.seasons;
        for (let i = 0; i < seasons.length; i++) {
            console.log("ITEM", seasons[i]);
            console.log(typeof seasons[i]);
            if (seasons[i] === null || typeof seasons[i] == "undefined") {
                console.log("EMPTY");
                 seasons[i] = {
                     Info: {SeasonNumber: i + 1} as TvSeason,
                     LoadPending: true,
                } as SeasonDetails;
            }
        }

        console.log("Handled seasons: ", seasons);
        this.setState({ seasons: seasons });
    }

    render() {
        if (this.state.seasons == null) {
            return (
                <div>No seasons</div>
            )
        }

        return (
            <div className={"columns is-mobile seasonlist-container"}>
                <div className={"column is-full"}>
                {this.state.seasons.map(season => (
                    <SeasonElt key={season.Info.SeasonNumber}
                               refreshSeason={this.props.refreshSeason}
                               season={season} />
                ))}
                </div>
            </div>
        )
    }
}

type SeasonProps = {
    refreshSeason(n: number): void,
    season :SeasonDetails | null,
};

type SeasonState = {
    season :SeasonDetails | null,
    show_episode_list :Boolean,
};

class SeasonElt extends React.Component<SeasonProps, SeasonState> {
    state :SeasonState = {
        season: null,
        show_episode_list: false,
    };

    constructor(props: SeasonProps) {
        super(props);

        this.state.season = this.props.season;

        this.markSeasonDownloaded = this.markSeasonDownloaded.bind(this);
        this.markSeasonNotDownloaded = this.markSeasonNotDownloaded.bind(this);
        this.changeSeasonDownloadedState = this.changeSeasonDownloadedState.bind(this);
        this.downloadSeason = this.downloadSeason.bind(this);
        this.toggleEpisodeList = this.toggleEpisodeList.bind(this);

        this.renderSeason = this.renderSeason.bind(this);
    }

    componentWillReceiveProps(nextProps :SeasonProps) {
        this.setState({ season: nextProps.season });
    }

    markSeasonDownloaded() {
        this.changeSeasonDownloadedState(true);
    }

    markSeasonNotDownloaded() {
        this.changeSeasonDownloadedState(false);
    }

    changeSeasonDownloadedState(downloaded_state :boolean) {
        if (this.state.season == null) {
            return;
        }

        let requests :any = [];
        let season = this.state.season;
        for (var index in this.state.season.EpisodeList) {
            let episode = this.state.season.EpisodeList[index];
            if (episode.DownloadingItem.Downloaded === downloaded_state || Helpers.dateIsInFuture(episode.Date)) {
                continue;
            }

            season.EpisodeList[index].ActionPending = true;
            requests.push(API.Episodes.changeDownloadedState(episode.ID, downloaded_state));
        }
        this.setState({season: season});

        Promise.all(requests).then(values => {
            this.props.refreshSeason(season.Info.SeasonNumber);
        });
    }

    downloadSeason() {
        if (this.state.season == null) {
            return;
        }

        let requests :any = [];
        let season = this.state.season;
        for (var index in this.state.season.EpisodeList) {
            let episode = this.state.season.EpisodeList[index];
            if (episode.DownloadingItem.Pending || episode.DownloadingItem.Downloading || episode.DownloadingItem.Downloaded || Helpers.dateIsInFuture(episode.Date)) {
                continue;
            }

            season.EpisodeList[index].ActionPending = true;
            requests.push(API.Episodes.download(episode.ID));
        }
        this.setState({season: season});

        Promise.all(requests).then(values => {
            this.props.refreshSeason(season.Info.SeasonNumber);
        });
    }

    toggleEpisodeList() {
        this.setState({show_episode_list: !this.state.show_episode_list});
    }

    renderSeason() {
        let season = this.state.season;
        if (season === null) {
            season = {
                LoadError: true,
            } as SeasonDetails;
        }

        if (season.LoadError) {
            return (
                <div className={"season-container"} key={season.Info.SeasonNumber}>
                    <div>
                        <div className="columns is-gapless is-vcentered is-mobile">
                            <div className="column">
                                <span className="title is-5">
                                    <i className={"has-text-grey"}>
                                        Could not load season details
                                    </i>
                                </span>
                            </div>
                        </div>
                    </div>
                </div>
            );
        }

        if (season.LoadPending) {
            return (
                <div className={"season-container"} key={season.Info.SeasonNumber}>
                    <div>
                        <div className="columns is-gapless is-vcentered is-mobile">
                            <div className="column is-narrow">
                                <button className={"button is-naked is-small is-loading"}>
                                </button>
                            </div>
                            <div className="column">
                                <span className="title is-5">
                                    <i className={"has-text-grey"}>
                                        Loading season details
                                    </i>
                                </span>
                            </div>
                        </div>
                    </div>
                </div>
            );
        }

        if (season.EpisodeList == null) {
            season.EpisodeList = [];
        }

        if (season.Info.SeasonNumber === 0) {
            return;
        }

        return (
            <div className={"season-container"} key={season.Info.SeasonNumber}>
                <div>
                    <div className="columns is-gapless is-vcentered is-mobile">
                        <div className="column is-narrow">
                            <button className={"button is-naked is-small season-deploy-button"}
                                    onClick={this.toggleEpisodeList}>
                                <span className={"icon"}>
                                {this.state.show_episode_list ? (
                                    <RiArrowDropUpLine />
                                ) : (
                                    <RiArrowDropDownLine />
                                )}
                                </span>
                            </button>
                        </div>
                        <div className="column"
                            onClick={this.toggleEpisodeList}>
                            <span className="title is-3">Season {season.Info.SeasonNumber} </span><span className="title is-6 has-text-grey">( {season.EpisodeList.length} episodes )</span>
                        </div>
                        <div className="column is-narrow">
                                <button data-tooltip="Mark whole season as not downloaded"
                                    onClick={this.markSeasonNotDownloaded}
                                    className="button is-small is-naked">
                                    <span className={"icon is-small"}>
                                        <RiEraserLine />
                                    </span>
                                </button>
                                <button data-tooltip="Mark whole season as downloaded"
                                    onClick={this.markSeasonDownloaded}
                                    className="button is-small is-naked">
                                    <span className={"icon is-small"}>
                                        <RiCheckDoubleLine />
                                    </span>
                                </button>
                                <button data-tooltip="Download the whole season"
                                    onClick={this.downloadSeason}
                                    className="button is-small is-naked">
                                    <span className={"icon is-small"}>
                                        <RiDownloadLine />
                                    </span>
                                </button>
                        </div>
                    </div>
                </div>
                {this.state.show_episode_list && (
                <div>
                    <EpisodeTable list={season.EpisodeList} />
                </div>
                )}
            </div>
        );
    }

    render() {
        return (
            <>
                {this.renderSeason()}
            </>
        )
    }
}

export default SeasonList;
