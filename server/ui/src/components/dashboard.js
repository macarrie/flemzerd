import React from "react";
import { Link } from "react-router-dom";

import API from "../utils/api";

class Dashboard extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            config: null,
            stats: null,
        };
    }

    componentDidMount() {
        this.getConfig();
        this.getStats();
    }

    getConfig() {
        API.Config.get().then(response => {
            console.log("Config: ", response.data);
            this.setState({config: response.data});
        }).catch(error => {
            console.log("Get config error: ", error);
        });
    }

    getStats() {
        API.Stats.get().then(response => {
            console.log("stats: ", response.data);
            this.setState({stats: response.data.stats});
        }).catch(error => {
            console.log("Get stats error: ", error);
        });
    }

    render() {
        if (this.state.config == null || this.state.stats == null) {
            return <div>Loading</div>;
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
