import React from "react";

import API from "../../utils/api";
import Helpers from "../../utils/helpers";
import Const from "../../const";
import Movie from "../../types/movie";

import Empty from "../empty";
import MediaIdsComponent from "../media_ids";
import Editable from "../editable";
import MediaActionBar from "../media_action_bar";
import DownloadingItemComponent from "../downloading_item";

type State = {
    movie: Movie,
    fanartURL: string,
};

class MovieDetails extends React.Component<any, State> {
    movie_refresh_interval: number;
    state: State = {
        movie: {} as Movie,
        fanartURL: "",
    };

    constructor(props: any) {
        super(props);

        this.movie_refresh_interval = 0;

        this.exitTitleEdit = this.exitTitleEdit.bind(this);

        this.useOriginalTitle = this.useOriginalTitle.bind(this);
        this.useDefaultTitle = this.useDefaultTitle.bind(this);
        this.changeDefaultTitle = this.changeDefaultTitle.bind(this);

        this.markNotDownloaded = this.markNotDownloaded.bind(this);
        this.markDownloaded = this.markDownloaded.bind(this);
        this.changeDownloadedState = this.changeDownloadedState.bind(this);

        this.downloadMovie = this.downloadMovie.bind(this);
        this.skipTorrent = this.skipTorrent.bind(this);
        this.abortDownload = this.abortDownload.bind(this);

        this.restoreMovie = this.restoreMovie.bind(this);
        this.deleteMovie = this.deleteMovie.bind(this);
    }

    componentDidMount() {
        this.getMovie();

        this.movie_refresh_interval = window.setInterval(this.getMovie.bind(this), Const.DATA_REFRESH);
    }

    componentWillUnmount() {
        clearInterval(this.movie_refresh_interval);
    }

    getMovie() {
        API.Movies.get(this.props.match.params.id).then(response => {
            let movie_result: Movie = response.data;
            movie_result.DisplayTitle = Helpers.getMediaTitle(movie_result);
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

    exitTitleEdit(value: string) {
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

    changeDefaultTitle(state: boolean) {
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

    changeDownloadedState(downloaded_state: boolean) {
        API.Movies.changeDownloadedState(this.state.movie.ID, downloaded_state).then(response => {
            this.getMovie();
        }).catch(error => {
            console.log("Change movie downloaded state error: ", error);
        });
    }

    skipTorrent() {
        API.Movies.skipTorrent(this.state.movie.ID).then(response => {
            this.getMovie();
        }).catch(error => {
            console.log("Skip torrent download error: ", error);
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
        if (this.state.movie.ID === 0) {
            return (
                <Empty label={"Loading"}/>
            );
        }

        return (
            <>
            <div id="full_background" className={"has-background-dark"}
                style={{backgroundImage: `url(${this.state.fanartURL})`}}></div>
            <MediaActionBar item={this.state.movie}
                            downloadItem={this.downloadMovie}
                            useOriginalTitle={this.useOriginalTitle}
                            useDefaultTitle={this.useDefaultTitle}
                            markNotDownloaded={this.markNotDownloaded}
                            markDownloaded={this.markDownloaded}
                            skipTorrent={this.skipTorrent}
                            abortDownload={this.abortDownload}
                            restoreItem={this.restoreMovie}
                            deleteItem={this.deleteMovie}
                            type="movie"/>

            <div className="container mediadetails">
                <div className="columns is-mobile">
                    <div className="column is-one-quarter-desktop is-one-third-tablet is-hidden-mobile">
                        <img width="100%" src={this.state.movie.Poster} alt="{this.state.movie.Title}" className="thumbnail" />
                    </div>
                    <div className="column">
                        <div className="columns is-mobile">
                            <div className="column is-narrow inlineedit">
                                <span className="title is-2 has-text-grey-light">
                                    <Editable
                                        value={this.state.movie.DisplayTitle}
                                        onFocusOut={this.exitTitleEdit}
                                        editingClassName="title is-2 has-text-grey-light"
                                    />
                                </span>
                            </div>
                            <div className={"column"}>
                                <span className="title is-3 has-text-grey-light media_title_details">({Helpers.getYear(this.state.movie.Date)})</span>
                            </div>
                        </div>
                        <div className="columns is-mobile">
                            <div className={"column is-narrow"}> See on </div>
                            <div className={"column"}><MediaIdsComponent ids={this.state.movie.MediaIds} type="movie"/></div>
                        </div>

                        <div className="container">
                            <div className="row">
                                <h5 className={"title is-5 has-text-grey-light"}>Overview</h5>
                                <div className="col-12">
                                    {this.state.movie.Overview}
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <div className="container">
                <div className="columns is-mobile">
                    <div className={"column is-full"}>
                        <div className={"box"}>
                            <DownloadingItemComponent item={this.state.movie.DownloadingItem}/>
                        </div>
                    </div>
                </div>
            </div>
            </>
        );
    }
}

export default MovieDetails;
