import React from "react";
import { Link } from "react-router-dom";

import API from "../../utils/api";
import Helpers from "../../utils/helpers";

import MediaIds from "../media_ids";
import Editable from "../editable";
import MediaActionBar from "../media_action_bar";

class EpisodeTableLine extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            item: null,
        };

        this.state.item = this.props.item;
    }

    componentWillReceiveProps(nextProps) {
        this.setState({ item: nextProps.item });
    }

    render() {
        return (
            <tr className={Helpers.dateIsInFuture(this.state.item.Date) ? "episode-list-future" : ""}>
                <td>
                    <span className="uk-text-small uk-border-rounded episode-label">{Helpers.formatNumber(this.state.item.Season)}x{Helpers.formatNumber(this.state.item.Number)}</span>
                    <Link to={`/episodes/${this.state.item.ID}`}>
                        {this.state.item.Title}
                    </Link>
                </td>
                <td className="uk-table-shrink uk-text-nowrap uk-text-center uk-hidden@s">{Helpers.formatDate(this.state.item.Date, 'DD MMM YYYY')}</td>
                <td className="uk-table-shrink uk-text-nowrap uk-text-center uk-visible@s uk-hidden@m">{Helpers.formatDate(this.state.item.Date, 'ddd DD MMM YYYY')}</td>
                <td className="uk-table-shrink uk-text-nowrap uk-text-center uk-visible@m">{Helpers.formatDate(this.state.item.Date, 'dddd DD MMM YYYY')}</td>
                <td className="uk-text-center">

                </td>
                <td className="uk-text-center">

                </td>
            </tr>
        );
    }
}

class EpisodeTable extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            list: null,
        };

        this.state.list = this.props.list;
    }

    componentWillReceiveProps(nextProps) {
        this.setState({ list: nextProps.list });
    }

    render() {
        return (
            <div>
                <table className="uk-table uk-table-divider uk-table-small">
                    <thead>
                        <tr>
                            <th>Title</th>
                            <th className="uk-hidden@m">Release</th>
                            <th className="uk-visible@m">Release date</th>

                            <th className="uk-table-shrink uk-text-nowrap uk-text-center uk-hidden@m">Status</th>
                            <th className="uk-table-shrink uk-text-nowrap uk-text-center uk-visible@m">Download status</th>

                            <th className="uk-table-shrink uk-text-nowrap uk-text-center">Actions</th>
                        </tr>
                    </thead>
                    <tbody>
                        {(this.state.list || []).map(episode => (
                            <EpisodeTableLine
                                key={episode.ID}
                                item={episode} />
                        ))}
                    </tbody>
                </table>
            </div>
        );
    }
}

class SeasonList extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            seasons: null,
        }

        this.state.seasons = this.props.seasons;
    }

    componentWillReceiveProps(nextProps) {
        console.log("state update seasons: ", nextProps);
        this.setState({ seasons: nextProps.seasons });
    }

    render() {
        console.log("Seasons: ", this.state.seasons);
        return (
            <div>
                {(this.state.seasons || []).filter((elt, index) => {
                    return index !== 0;
                }).map(season => (
                    <React.Fragment key={season.Info.SeasonNumber}>
                        <div>
                            <div className="uk-grid uk-grid-small uk-child-width-1-2 uk-flex uk-flex-middle" data-uk-grid>
                                <div className="uk-width-expand">
                                    <span className="uk-h3">Season {season.Info.SeasonNumber} </span><span className="uk-text-muted uk-h6">( {season.EpisodeList.length} episodes )</span>
                                </div>
                                <div className="uk-width-auto">
                                    <ul className="uk-iconnav">
                                        <li data-uk-tooltip="delay: 500; title: Mark whole season as not downloaded"
                                            className="uk-icon"
                                            data-uk-icon="icon: push; ratio: 0.75"></li>
                                        <li data-uk-tooltip="delay: 500; title: Mark whole season as downloaded"
                                            className="uk-icon"
                                            data-uk-icon="icon: check; ratio: 0.75"></li>
                                        <li data-uk-tooltip="delay: 500; title: Download the whole season"
                                            className="uk-icon"
                                            data-uk-icon="icon: download; ratio: 0.75"></li>
                                    </ul>
                                </div>
                            </div>
                        </div>
                        <EpisodeTable list={season.EpisodeList} />
                    </React.Fragment>
                ))}
            </div>
        )
    }
}

class TvShowDetails extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            movie: null,
            seasons: [],
            fanartURL: "",
        };

        this.exitTitleEdit = this.exitTitleEdit.bind(this);

        this.useOriginalTitle = this.useOriginalTitle.bind(this);
        this.useDefaultTitle = this.useDefaultTitle.bind(this);
        this.changeDefaultTitle = this.changeDefaultTitle.bind(this);

        this.treatAsRegularShow = this.treatAsRegularShow.bind(this);
        this.treatAsAnime = this.treatAsAnime.bind(this);
        this.changeAnimeState = this.changeAnimeState.bind(this);

        this.restoreShow = this.restoreShow.bind(this);
        this.deleteShow = this.deleteShow.bind(this);
    }

    componentDidMount() {
        this.getShow();
    }

    getShow() {
        API.Shows.get(this.props.match.params.id).then(response => {
            let show_result = response.data;
            show_result.DisplayTitle = Helpers.getMediaTitle(show_result)
            this.setState({
                show: show_result,
                seasons: new Array(show_result.NumberOfSeasons),
            });
            this.getFanart();
            this.getSeasonList();
        }).catch(error => {
            console.log("Get show details error: ", error);
        });
    }

    getSeasonList() {
        for (var season in this.state.show.Seasons) {
            let n = this.state.show.Seasons[season].SeasonNumber;

            API.Shows.getSeason(this.props.match.params.id, n).then(response => {
                console.log("Season details (update state): ", response.data);

                let seasonlist = this.state.seasons;
                seasonlist[n] = response.data;
                console.log("Season list after modification: ", seasonlist);
                this.setState({seasons: seasonlist});
            }).catch(error => {
                console.log("Get show details error: ", error);
            });
        }
    }

    getFanart() {
        API.Fanart.getTvShowFanart(this.state.show).then(response => {
            let fanartObj = response.data;
            if (fanartObj["showbackground"] && fanartObj["showbackground"].length > 0) {
                this.setState({fanartURL: fanartObj["showbackground"][0]["url"]});
            }
        }).catch(error => {
            console.log("Get show fanart error: ", error);
        });
    }

    exitTitleEdit(value) {
        API.Shows.changeCustomTitle(this.state.show.ID, value).then(response => {
            this.getShow();
        }).catch(error => {
            console.log("Update custom title error: ", error);
        });
    }

    useOriginalTitle() {
        this.changeDefaultTitle(false);
    }

    useDefaultTitle() {
        this.changeDefaultTitle(true);
    }

    changeDefaultTitle(state) {
        API.Shows.useDefaultTitle(this.state.show.ID, state).then(response => {
            this.getShow();
        }).catch(error => {
            console.log("Change use default title bool error: ", error);
        });
    }

    treatAsAnime() {
        this.changeAnimeState(true);
    }

    treatAsRegularShow() {
        this.changeAnimeState(false);
    }

    changeAnimeState(state) {
        API.Shows.changeAnimeState(this.state.show.ID, state).then(response => {
            this.getShow();
        }).catch(error => {
            console.log("Change use default title bool error: ", error);
        });
    }

    restoreShow() {
        API.Shows.restore(this.state.show.ID).then(response => {
            this.getShow();
        }).catch(error => {
            console.log("Restore show error: ", error);
        });
    }

    deleteShow() {
        API.Shows.delete(this.state.show.ID).then(response => {
            this.getShow();
        }).catch(error => {
            console.log("Delete show error: ", error);
        });
    }

    render() {
        if (this.state.show == null) {
            return <div>Loading</div>;
        }

        return (
            <>
            <div id="full_background"
                style={{backgroundImage: `url(${this.state.fanartURL})`}}></div>
            <MediaActionBar item={this.state.show}
                useOriginalTitle={this.useOriginalTitle}
                useDefaultTitle={this.useDefaultTitle}
                treatAsRegularShow={this.treatAsRegularShow}
                treatAsAnime={this.treatAsAnime}
                restoreItem={this.restoreShow}
                deleteItem={this.deleteShow}
                type="tvshow"/>

            <div className="uk-container uk-light mediadetails">
                <div className="uk-grid" data-uk-grid>
                    <div className="uk-width-1-3">
                        <img width="100%" src={this.state.show.Poster} alt="{this.state.show.Title}" className="uk-border-rounded" data-uk-img />
                    </div>
                    <div className="uk-width-expand">
                        <div className="uk-grid uk-grid-collapse uk-flex uk-flex-middle" data-uk-grid>
                            <div className="inlineedit">
                                <span className="uk-h2">
                                    <Editable
                                        value={this.state.show.DisplayTitle}
                                        onFocus={this.enterTitleEdit}
                                        onFocusOut={this.exitTitleEdit}
                                        editingClassName="uk-border-rounded uk-h2 uk-margin-remove-top uk-text-light"
                                    />
                                </span>
                            </div>
                            <div>
                                <span className="uk-h3 uk-text-muted uk-margin-small-left media_title_details">({Helpers.getYear(this.state.show.FirstAired)})</span>
                            </div>
                        </div>
                        <div className="uk-grid uk-grid-medium" data-uk-grid>
                            <div> See on </div>
                            <div><MediaIds ids={this.state.show.MediaIds} type="movie" /></div>
                        </div>

                        <div className="container">
                            <div className="row">
                                <h5>Overview</h5>
                                <div className="col-12">
                                    {this.state.show.Overview}
                                </div>
                            </div>
                            <div className="row">
                                <h5>Seasons</h5>
                                <div className="col-12">
                                    {this.state.show.NumberOfSeasons} seasons, {this.state.show.NumberOfEpisodes} episodes
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <div className="uk-container uk-margin-medium-top uk-margin-medium-bottom">
                <div className="uk-background-default uk-border-rounded">
                    <div className="uk-grid uk-child-width-1-1 uk-padding-small" data-uk-grid>
                        <div>
                            <SeasonList seasons={this.state.seasons}/>
                        </div>
                    </div>
                </div>
            </div>

            </>
        );
    }
}

export default TvShowDetails;
