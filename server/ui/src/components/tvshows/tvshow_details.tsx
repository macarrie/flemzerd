import React from "react";
import {Route} from "react-router-dom";

import API from "../../utils/api";
import Helpers from "../../utils/helpers";
import Const from "../../const";
import TvShow, {SeasonDetails} from "../../types/tvshow";

import Empty from "../empty";
import MediaIdsComponent from "../media_ids";
import Editable from "../editable";
import MediaActionBar from "../media_action_bar";
import SeasonList from "./season_list";
import EpisodeDetails from './episode_details';

type State = {
    show :TvShow | null,
    seasons :SeasonDetails[] | null,
    fanartURL :string,
};

class TvShowDetails extends React.Component<any, State> {
    tvshow_refresh_interval :number;
    state :State = {
        show: null,
        seasons: null,
        fanartURL: "",
    };

    constructor(props: any) {
        super(props);

        this.tvshow_refresh_interval = 0;

        this.getSeason = this.getSeason.bind(this);
        this.getSeasonList = this.getSeasonList.bind(this);

        this.exitTitleEdit = this.exitTitleEdit.bind(this);

        this.useOriginalTitle = this.useOriginalTitle.bind(this);
        this.useDefaultTitle = this.useDefaultTitle.bind(this);
        this.changeDefaultTitle = this.changeDefaultTitle.bind(this);

        this.treatAsRegularShow = this.treatAsRegularShow.bind(this);
        this.treatAsAnime = this.treatAsAnime.bind(this);
        this.changeAnimeState = this.changeAnimeState.bind(this);

        this.restoreShow = this.restoreShow.bind(this);
        this.deleteShow = this.deleteShow.bind(this);

        this.renderMediaDetails = this.renderMediaDetails.bind(this);
    }

    componentDidMount() {
        this.getShow();

        this.tvshow_refresh_interval = window.setInterval(this.getShow.bind(this), Const.DATA_REFRESH);
    }

    componentWillUnmount() {
        clearInterval(this.tvshow_refresh_interval);
    }

    getShow() {
        API.Shows.get(this.props.match.params.id).then(response => {
            let show_result :TvShow = response.data;
            show_result.DisplayTitle = Helpers.getMediaTitle(show_result)
            this.setState({
                show: show_result,
                seasons: this.state.seasons == null ? new Array<SeasonDetails>(show_result.NumberOfSeasons) : this.state.seasons,
            });
            this.getFanart();
            this.getSeasonList();
        }).catch(error => {
            console.log("Get show details error: ", error);
        });
    }

    getSeason(n: number) {
        API.Shows.getSeason(this.props.match.params.id, n).then(response => {
            if (this.state.seasons == null) {
                return;
            }

            let seasonlist :SeasonDetails[] = this.state.seasons;
            seasonlist[n - 1] = response.data;

            this.setState({seasons: seasonlist});
        }).catch(error => {
            console.log("Get season details error: ", error);
            if (this.state.seasons == null) {
                return;
            }

            let seasonlist :SeasonDetails[] = this.state.seasons;
            seasonlist[n - 1] = {LoadError: true} as SeasonDetails;

            this.setState({seasons: seasonlist});
        });
    }

    getSeasonList() {
        if (this.state.show == null) {
            return;
        }

        for (let season in this.state.show.Seasons) {
            this.getSeason(this.state.show.Seasons[season].SeasonNumber);
        }
    }

    getFanart() {
        if (this.state.show == null) {
            return;
        }

        API.Fanart.getTvShowFanart(this.state.show).then(response => {
            let fanartObj = response.data;
            if (fanartObj["showbackground"] && fanartObj["showbackground"].length > 0) {
                this.setState({fanartURL: fanartObj["showbackground"][0]["url"]});
            }
        }).catch(error => {
            console.log("Get show fanart error: ", error);
        });
    }

    exitTitleEdit(value :string) {
        if (this.state.show == null) {
            return;
        }

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

    changeDefaultTitle(state :boolean) {
        if (this.state.show == null) {
            return;
        }

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

    changeAnimeState(state :boolean) {
        if (this.state.show == null) {
            return;
        }

        API.Shows.changeAnimeState(this.state.show.ID, state).then(response => {
            this.getShow();
        }).catch(error => {
            console.log("Change use default title bool error: ", error);
        });
    }

    restoreShow() {
        if (this.state.show == null) {
            return;
        }

        API.Shows.restore(this.state.show.ID).then(response => {
            this.getShow();
        }).catch(error => {
            console.log("Restore show error: ", error);
        });
    }

    deleteShow() {
        if (this.state.show == null) {
            return;
        }

        API.Shows.delete(this.state.show.ID).then(response => {
            this.getShow();
        }).catch(error => {
            console.log("Delete show error: ", error);
        });
    }

    renderMediaDetails() {
        if (this.state.show == null) {
            return (
                <div className={"container"}>
                    <Empty label={"Loading"}/>
                </div>
            );
        }

        return (
            <>
                <div id="full_background" className={"has-background-dark "}
                     style={{backgroundImage: `url(${this.state.fanartURL})`}}></div>
                <MediaActionBar item={this.state.show}
                                useOriginalTitle={this.useOriginalTitle}
                                useDefaultTitle={this.useDefaultTitle}
                                treatAsRegularShow={this.treatAsRegularShow}
                                treatAsAnime={this.treatAsAnime}
                                restoreItem={this.restoreShow}
                                deleteItem={this.deleteShow}
                                type="tvshow"/>

                <div className="container mediadetails">
                    <div className="columns">
                        <div className="column is-one-quarter">
                            <img width="100%" src={this.state.show.Poster} alt="{this.state.show.Title}" className="thumbnail" />
                        </div>
                        <div className="column">
                            <div className="columns">
                                <div className="column is-narrow inlineedit">
                                    <span className="title is-2 has-text-grey-light">
                                        <Editable
                                            value={this.state.show.DisplayTitle}
                                            onFocusOut={this.exitTitleEdit}
                                            editingClassName="title is-2 has-text-grey-light"
                                        />
                                    </span>
                                </div>
                                <div className={"column"}>
                                    <span className="title is-3 has-text-grey-light media_title_details">({Helpers.getYear(this.state.show.FirstAired)})</span>
                                </div>
                            </div>
                            <div className="columns">
                                <div className={"column is-narrow"}> See on</div>
                                <div className={"column"}><MediaIdsComponent ids={this.state.show.MediaIds} type="movie"/></div>
                            </div>

                            <div className="container">
                                <div className="row">
                                    <h5 className={"title is-5 has-text-grey-light"}>Overview</h5>
                                    <div className="col-12">
                                        {this.state.show.Overview}
                                    </div>
                                </div>
                                <div className="row">
                                    <h5 className={"title is-5 has-text-grey-light"}>Seasons</h5>
                                    <div className="col-12">
                                        {this.state.show.NumberOfSeasons} seasons, {this.state.show.NumberOfEpisodes} episodes
                                    </div>
                                </div>
                                <div className="row">
                                    <h5 className={"title is-5 has-text-grey-light"}>Additional info</h5>
                                    <div className="col-12">
                                        {this.state.show.IsAnime ? (
                                            <>
                                                This show is considered as an anime. Because anime episode are not
                                                usually classified into seasons, episodes absolute numbers will be used
                                                instead of the classic season/episode numbering to search for torrents.
                                            </>
                                        ) : (
                                            <>
                                                This show is a regular show (not an anime). Classic episode numbering
                                                (seasons number and episode number inside season) will be used to
                                                searched for torrents
                                            </>
                                        )}
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

                {!this.state.show.DeletedAt && (
                    <div className="container box">
                        <div className="columns">
                            <div className={"column"}>
                                <SeasonList refreshSeason={this.getSeason}
                                            seasons={this.state.seasons}/>
                            </div>
                        </div>
                    </div>
                )}

            </>
        );
    }

    render() {
        const match = this.props.match;

        return (
            <>
                <Route path={`${match.path}/episodes/:id`} component={EpisodeDetails}/>
                <Route exact
                       path={match.path}
                       render={this.renderMediaDetails}/>
            </>
        );
    }
}

export default TvShowDetails;
