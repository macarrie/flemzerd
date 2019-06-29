import React from "react";

import API from "../../utils/api";
import Helpers from "../../utils/helpers";

import EpisodeTable from "./episode_table";

class SeasonList extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            seasons: null,
        }

        this.state.seasons = this.props.seasons;

        this.markSeasonDownloaded = this.markSeasonDownloaded.bind(this);
        this.markSeasonNotDownloaded = this.markSeasonNotDownloaded.bind(this);
        this.changeSeasonDownloadedState = this.changeSeasonDownloadedState.bind(this);
        this.downloadSeason = this.downloadSeason.bind(this);
    }

    componentWillReceiveProps(nextProps) {
        this.setState({ seasons: nextProps.seasons });
    }

    markSeasonDownloaded(season_number) {
        this.changeSeasonDownloadedState(season_number, true);
    }

    markSeasonNotDownloaded(season_number) {
        this.changeSeasonDownloadedState(season_number, false);
    }

    changeSeasonDownloadedState(season_number, downloaded_state) {
        let requests = [];
        for (var index in this.state.seasons[season_number].EpisodeList) {
            let episode = this.state.seasons[season_number].EpisodeList[index];
            requests.push(API.Episodes.changeDownloadedState(episode.ID, downloaded_state));
        }

        Promise.all(requests).then(values => {
            this.props.refreshSeason(season_number);
        });
    }

    downloadSeason(season_number) {
        let requests = [];
        for (var index in this.state.seasons[season_number].EpisodeList) {
            let episode = this.state.seasons[season_number].EpisodeList[index];
            if (episode.DownloadingItem.Pending || episode.DownloadingItem.Downloading || episode.DownloadingItem.Downloaded || Helpers.dateIsInFuture(episode.Date)) {
                continue;
            }

            requests.push(API.Episodes.download(episode.ID));
        }
        Promise.all(requests).then(values => {
            this.props.refreshSeason(season_number);
        });
    }

    render() {
        return (
            <div data-uk-accordion
                data-toggle="> div > .uk-grid > .uk-accordion-title"
                data-multiple="true"
                data-collapsible="true">
                {(this.state.seasons || []).filter((elt, index) => {
                    return index !== 0;
                }).map(season => (
                    <div key={season.Info.SeasonNumber}>
                        <div>
                            <div className="uk-grid uk-grid-small uk-child-width-1-2 uk-flex uk-flex-middle" data-uk-grid>
                                <div className="uk-width-expand uk-accordion-title">
                                        <span className="uk-h3">Season {season.Info.SeasonNumber} </span><span className="uk-text-muted uk-h6">( {season.EpisodeList.length} episodes )</span>
                                </div>
                                <div className="uk-width-auto">
                                    <ul className="uk-iconnav">
                                        <li data-uk-tooltip="delay: 500; title: Mark whole season as not downloaded"
                                            onClick={() => this.markSeasonNotDownloaded(season.Info.SeasonNumber)}
                                            className="uk-icon"
                                            data-uk-icon="icon: push; ratio: 0.75"></li>
                                        <li data-uk-tooltip="delay: 500; title: Mark whole season as downloaded"
                                            onClick={() => this.markSeasonDownloaded(season.Info.SeasonNumber)}
                                            className="uk-icon"
                                            data-uk-icon="icon: check; ratio: 0.75"></li>
                                        <li data-uk-tooltip="delay: 500; title: Download the whole season"
                                            onClick={() => this.downloadSeason(season.Info.SeasonNumber)}
                                            className="uk-icon"
                                            data-uk-icon="icon: download; ratio: 0.75"></li>
                                    </ul>
                                </div>
                            </div>
                        </div>
                        <div className="uk-accordion-content">
                            <EpisodeTable list={season.EpisodeList} />
                        </div>
                    </div>
                ))}
            </div>
        )
    }
}

export default SeasonList;
