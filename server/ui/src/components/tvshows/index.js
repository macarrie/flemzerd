import React from "react";
import { Route } from "react-router-dom";

import API from "../../utils/api";

import MediaMiniature from "../media_miniature";
import TvShowDetails from "./tvshow_details";

class TvShows extends React.Component {
    state = {
        tracked: [],
        removed: [],
    };

    constructor(props) {
        super(props);

        this.renderShowList = this.renderShowList.bind(this);
    }

    componentDidMount() {
        this.getShows();
    }

    getShows() {
        API.Shows.tracked().then(response => {
            this.setState({tracked: response.data});
        }).catch(error => {
            console.log("Get tracked shows error: ", error);
        });

        API.Shows.removed().then(response => {
            this.setState({removed: response.data});
        }).catch(error => {
            console.log("Get removed shows error: ", error);
        });
    }

    renderShowList() {
        return (
            <div className="uk-container" data-uk-filter="target: .item-filter">
                <div className="uk-grid" data-uk-grid>
                    <div className="uk-width-expand">
                        <span className="uk-h3">TV Shows</span>
                    </div>
                    <ul className="uk-subnav uk-subnav-pill">
                        <li className="uk-visible@s uk-active" data-uk-filter-control="">
                            <button className="uk-button uk-button-text">
                                All
                            </button>
                        </li>
                        <li className="uk-visible@s" data-uk-filter-control="filter: .tracked-show">
                            <button className="uk-button uk-button-text">
                                Tracked
                            </button>
                        </li>
                        <li className="uk-visible@s" data-uk-filter-control="filter: .removed-show">
                            <button className="uk-button uk-button-text">
                                Removed
                            </button>
                        </li>
                        <li>
                            <button className="uk-button uk-button-text">
                                <div className="uk-flex uk-flex-middle">
                                    <div className="uk-float-left">
                                        <span className="uk-icon" data-uk-icon="icon: refresh; ratio: 0.75"></span>
                                    </div>
                                    <div className="uk-text-nowrap">
                                        Check now
                                    </div>
                                </div>
                                <div>
                                    <div data-uk-spinner="ratio: 0.5"></div>
                                    <div>
                                        <i>Checking</i>
                                    </div>
                                </div>
                            </button>
                        </li>
                        <li>
                            <button className="uk-button uk-button-text">
                                <div>
                                    <span className="uk-icon" data-uk-icon="icon: refresh; ratio: 0.75"></span>
                                    Refresh
                                </div>
                                <div>
                                    <div data-uk-spinner="ratio: 0.5"></div>
                                    <i>Refreshing</i>
                                </div>
                            </button>
                        </li>
                    </ul>
                </div>
                <hr />

                <div className="uk-grid uk-grid-small uk-child-width-1-2 uk-child-width-1-6@l uk-child-width-1-4@m uk-child-width-1-3@s item-filter" data-uk-grid>
                    {this.state.tracked ? (
                        this.state.tracked.map(show => (
                            <MediaMiniature key={show.ID} item={show} type="tvshow" />
                        ))
                    ) : ""
                    }
                    {this.state.removed ? (
                        this.state.removed.map(show => (
                            <MediaMiniature key={show.ID} item={show} type="tvshow" />
                        ))
                    ) : ""
                    }
                </div>
            </div>
        )
    }

    render() {
        const match = this.props.match;

        return (
            <>
                <Route path={`${match.path}/:id`} component={TvShowDetails} />
                <Route exact
                       path={match.path}
                       render={this.renderShowList} />
            </>
        );
    }
}

export default TvShows;
