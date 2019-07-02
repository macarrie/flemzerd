import React from "react";
import {Link} from "react-router-dom";

import API from "../utils/api";
import Const from "../const";
import Loading from "./loading";

class Dashboard extends React.Component {
    refresh_interval;

    constructor(props) {
        super(props);

        this.state = {
            config: null,
            stats: null,
        };
    }

    componentDidMount() {
        this.load();

        this.refresh_interval = setInterval(this.load.bind(this), Const.DATA_REFRESH);
    }

    componentWillUnmount() {
        clearInterval(this.refresh_interval);
    }

    load() {
        console.log("refresh");
        this.getConfig();
        this.getStats();
    }

    getConfig() {
        API.Config.get().then(response => {
            this.setState({config: response.data});
        }).catch(error => {
            console.log("Get config error: ", error);
        });
    }

    getStats() {
        API.Stats.get().then(response => {
            this.setState({stats: response.data.stats});
        }).catch(error => {
            console.log("Get stats error: ", error);
        });
    }

    render() {
        if (this.state.config == null || this.state.stats == null) {
            return <Loading />;
        }

        return (
            <div className="uk-container">
                <div className="uk-grid uk-text-center" data-uk-grid>
                    {this.state.config.System.TrackShows && (
                        <>
                            <div className="uk-width-1-1@s uk-width-1-2@m">
                                <div className="uk-heading-hero title-gradient">
                                    <Link to="/tvshows">
                                        {this.state.stats.Shows.Tracked}
                                    </Link>
                                </div>
                                <span className="uk-h4 uk-text-muted">
                                Tracked shows
                            </span>
                            </div>
                            <div className="uk-width-1-1@s uk-width-1-2@m">
                                <div className="uk-heading-hero title-gradient">
                                    <Link to="/tvshows">
                                        {this.state.stats.Episodes.Downloading}
                                    </Link>
                                </div>
                                <span className="uk-h4 uk-text-muted">
                                Downloading episodes
                            </span>
                            </div>
                        </>
                    )}

                    {this.state.config.System.TrackMovies && (
                        <>
                            <div className="uk-width-1-1@s uk-width-1-3@m">
                                <div className="uk-heading-hero title-gradient">
                                    <Link to="/movies">
                                        {this.state.stats.Movies.Tracked}
                                    </Link>
                                </div>
                                <span className="uk-h4 uk-text-muted">
                                New movies
                            </span>
                            </div>
                            <div className="uk-width-1-1@s uk-width-1-3@m">
                                <div className="uk-heading-hero title-gradient">
                                    <Link to="/movies">
                                        {this.state.stats.Movies.Downloading}
                                    </Link>
                                </div>
                                <span className="uk-h4 uk-text-muted">
                                Downloading movies
                            </span>
                            </div>
                            <div className="uk-width-1-1@s uk-width-1-3@m">
                                <div className="uk-heading-hero title-gradient">
                                    <Link to="/movies">
                                        {this.state.stats.Movies.Downloaded}
                                    </Link>
                                </div>
                                <span className="uk-h4 uk-text-muted">
                                Downloaded movies
                            </span>
                            </div>
                        </>
                    )}
                </div>
            </div>
        );
    }
}

export default Dashboard;
