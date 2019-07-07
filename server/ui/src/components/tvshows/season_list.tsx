import React from "react";

import API from "../../utils/api";
import Helpers from "../../utils/helpers";
import {SeasonDetails} from "../../types/tvshow";

import EpisodeTable from "./episode_table";

type Props = {
    refreshSeason(n :number),
    seasons :SeasonDetails[] | null,
};

type State = {
    seasons :SeasonDetails[] | null,
};

class SeasonList extends React.Component<Props, State> {
    state :State = {
        seasons: null,
    };

    constructor(props) {
        super(props);

        this.state.seasons = this.props.seasons;

        this.markSeasonDownloaded = this.markSeasonDownloaded.bind(this);
        this.markSeasonNotDownloaded = this.markSeasonNotDownloaded.bind(this);
        this.changeSeasonDownloadedState = this.changeSeasonDownloadedState.bind(this);
        this.downloadSeason = this.downloadSeason.bind(this);

        this.renderSeason = this.renderSeason.bind(this);
    }

    componentWillReceiveProps(nextProps :Props) {
        this.setState({ seasons: nextProps.seasons });
    }

    markSeasonDownloaded(season_number :number) {
        this.changeSeasonDownloadedState(season_number, true);
    }

    markSeasonNotDownloaded(season_number :number) {
        this.changeSeasonDownloadedState(season_number, false);
    }

    changeSeasonDownloadedState(season_number :number, downloaded_state :boolean) {
        if (this.state.seasons == null || this.state.seasons[season_number] == null) {
            return;
        }

        let requests :any = [];
        for (var index in this.state.seasons[season_number].EpisodeList) {
            let episode = this.state.seasons[season_number].EpisodeList[index];
            requests.push(API.Episodes.changeDownloadedState(episode.ID, downloaded_state));
        }

        Promise.all(requests).then(values => {
            this.props.refreshSeason(season_number);
        });
    }

    downloadSeason(season_number :number) {
        if (this.state.seasons == null || this.state.seasons[season_number] == null) {
            return;
        }

        let requests :any = [];
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

    make_uuid() {
        let length = 8;
        let timestamp = +new Date();

        var _getRandomInt = function( min, max ) {
            return Math.floor( Math.random() * ( max - min + 1 ) ) + min;
        }

        let generate = function() {
            var ts = timestamp.toString();
            var parts = ts.split( "" ).reverse();
            var id = "";

            for( var i = 0; i < length; ++i ) {
                var index = _getRandomInt( 0, parts.length - 1 );
                id += parts[index];	 
            }

            return id;
        }

        return generate();
    }

    renderSeason(season :SeasonDetails) {
        if (season == null) {
            return (
                <div 
                    key={this.make_uuid()}
                    className="uk-text-muted uk-flex uk-flex-middle">
                <span className="uk-margin-right" data-uk-spinner="ratio: 0.5"></span>
                <i>
                        Loading season details
                </i>
            </div>
            );
        }

        if (season.LoadError) {
            return (
                <div 
                    key={this.make_uuid()}
                    className="uk-text-muted uk-flex uk-flex-middle">
                    <span className="uk-margin-right" data-uk-icon="close"></span>
                    <i>
                            Could not load season
                    </i>
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
        );
    }

    renderSeasons() {
        if (this.state.seasons == null) {
            return;
        }

        let seasonRenders :any = [];

        for (var i = 0; i < this.state.seasons.length; i++) {
            if (i === 0 && this.state.seasons[i] == null) {
                continue;
            }

            seasonRenders.push(this.renderSeason(this.state.seasons[i]));
        }

        return seasonRenders;
    }

    render() {
        return (
            <div data-uk-accordion
                data-toggle="> div > .uk-grid > .uk-accordion-title"
                data-multiple="true"
                data-collapsible="true">
                {this.renderSeasons()}
        </div>
        )
    }
}

export default SeasonList;
