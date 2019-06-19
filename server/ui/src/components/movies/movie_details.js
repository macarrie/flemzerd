import React from "react";

import API from "../../utils/api";
import Helpers from "../../utils/helpers";

import MediaIds from "../media_ids";
import Editable from "../editable";
import MediaActionBar from "../media_action_bar";
import DownloadingItem from "../downloading_item";

class MovieDetails extends React.Component {
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

        this.markNotDownloaded = this.markNotDownloaded.bind(this);
        this.markDownloaded = this.markDownloaded.bind(this);
        this.changeDownloadedState = this.changeDownloadedState.bind(this);

        this.downloadMovie = this.downloadMovie.bind(this);
        this.abortDownload = this.abortDownload.bind(this);

        this.restoreMovie = this.restoreMovie.bind(this);
        this.deleteMovie = this.deleteMovie.bind(this);
    }

    componentDidMount() {
        this.getMovie();
    }

    getMovie() {
        API.Movies.get(this.props.match.params.id).then(response => {
            let movie_result = response.data;
            movie_result.DisplayTitle = Helpers.getMovieTitle(movie_result)
            this.setState({movie: movie_result});
            this.getFanart();
        }).catch(error => {
            console.log("Get movie details error: ", error);
        });
    }

    getFanart() {
        API.Fanart.getMovieFanart(this.state.movie).then(response => {
            let fanartObj = response.data;
            if (fanartObj["moviebackground"] && fanartObj["moviebackground"].length > 0) {
                this.setState({fanartURL: fanartObj["moviebackground"][0]["url"]});
            }
        }).catch(error => {
            console.log("Get movie fanart error: ", error);
        });
    }

    exitTitleEdit(value) {
        API.Movies.changeCustomTitle(this.state.movie.ID, value).then(response => {
            this.getMovie();
        }).catch(error => {
            console.log("Update custom title error: ", error);
        });
    }

    downloadMovie() {
        API.Movies.download(this.state.movie.ID).then(response => {
            this.getMovie();
        }).catch(error => {
            console.log("Download movie error: ", error);
        });
    }

    useOriginalTitle() {
        this.changeDefaultTitle(false);
    }

    useDefaultTitle() {
        this.changeDefaultTitle(true);
    }

    changeDefaultTitle(state) {
        API.Movies.useDefaultTitle(this.state.movie.ID, state).then(response => {
            this.getMovie();
        }).catch(error => {
            console.log("Change use default title bool error: ", error);
        });
    }

    markNotDownloaded() {
        this.changeDownloadedState(false);
    }

    markDownloaded() {
        this.changeDownloadedState(true);
    }

    changeDownloadedState(downloaded_state) {
        API.Movies.changeDownloadedState(this.state.movie.ID, downloaded_state).then(response => {
            this.getMovie();
        }).catch(error => {
            console.log("Change movie downloaded state error: ", error);
        });
    }

    abortDownload() {
        API.Movies.abortDownload(this.state.movie.ID).then(response => {
            this.getMovie();
        }).catch(error => {
            console.log("Abort movie download error: ", error);
        });
    }

    restoreMovie() {
        API.Movies.restore(this.state.movie.ID).then(response => {
            this.getMovie();
        }).catch(error => {
            console.log("Restore movie error: ", error);
        });
    }

    deleteMovie() {
        API.Movies.delete(this.state.movie.ID).then(response => {
            this.getMovie();
        }).catch(error => {
            console.log("Delete movie error: ", error);
        });
    }

    render() {
        if (this.state.movie == null) {
            return <div>Loading</div>;
        }

        return (
            <>
            <div id="full_background"
                style={{backgroundImage: `url(${this.state.fanartURL})`}}></div>
            <MediaActionBar item={this.state.movie}
                            downloadItem={this.downloadMovie}
                            useOriginalTitle={this.useOriginalTitle}
                            useDefaultTitle={this.useDefaultTitle}
                            markNotDownloaded={this.markNotDownloaded}
                            markDownloaded={this.markDownloaded}
                            abortDownload={this.abortDownload}
                            restoreItem={this.restoreMovie}
                            deleteItem={this.deleteMovie}
                            type="movie"/>

            <div className="uk-container uk-light mediadetails">
                <div className="uk-grid" data-uk-grid>
                    <div className="uk-width-1-3">
                        <img width="100%" src={this.state.movie.Poster} alt="{this.state.movie.Title}" className="uk-border-rounded" data-uk-img />
                    </div>
                    <div className="uk-width-expand">
                        <div className="uk-grid uk-grid-collapse uk-flex uk-flex-middle" data-uk-grid>
                            <div className="inlineedit">
                                <span className="uk-h2">
                                    <Editable
                                        value={this.state.movie.DisplayTitle}
                                        onFocus={this.enterTitleEdit}
                                        onFocusOut={this.exitTitleEdit}
                                        editingClassName="uk-border-rounded uk-h2 uk-margin-remove-top uk-text-light"
                                    />
                                </span>
                            </div>
                            <div>
                                <span className="uk-h3 uk-text-muted uk-margin-small-left media_title_details">({Helpers.getYear(this.state.movie.Date)})</span>
                            </div>
                        </div>
                        <div className="uk-grid uk-grid-medium" data-uk-grid>
                            <div> See on </div>
                            <div><MediaIds ids={this.state.movie.MediaIds} type="movie" /></div>
                        </div>

                        <div className="container">
                            <div className="row">
                                <h5>Overview</h5>
                                <div className="col-12">
                                    {this.state.movie.Overview}
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <div className="uk-container uk-margin-medium-top uk-margin-medium-bottom">
                <DownloadingItem item={this.state.movie.DownloadingItem} />
            </div>
            </>
        );
    }
}

export default MovieDetails;
