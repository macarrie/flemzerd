import React from "react";

import API from "../../utils/api";
import Helpers from "../../utils/helpers";

import MediaIds from "../media_ids";
import Editable from "../editable";
import MediaActionBar from "../media_action_bar";

class TvShowDetails extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            movie: null,
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
            this.setState({show: show_result});
            this.getFanart();
        }).catch(error => {
            console.log("Get show details error: ", error);
        });
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
            </>
        );
    }
}

export default TvShowDetails;
