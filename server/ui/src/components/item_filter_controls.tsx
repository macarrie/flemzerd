import React from "react";
import {MediaMiniatureFilter} from "../types/media_miniature_filter";

type Props = {
    type :string,
    filterValue :MediaMiniatureFilter,
    updateFilter(value: MediaMiniatureFilter): void,
};

type State = {
    filter :MediaMiniatureFilter,
};

class ItemFilterControls extends React.Component<Props, State> {
    constructor(props: Props) {
        super(props);
        this.state = {
            filter: props.filterValue,
        };
    }

    componentWillReceiveProps(nextProps: Props) {
        this.setState({
            filter: nextProps.filterValue,
        });
    }

    getActiveFilterClass(targetFilter :MediaMiniatureFilter) :String {
        if (targetFilter === this.state.filter) {
            return "active";
        }

        return "";
    }

    render() {
        if (this.props.type === "movie") {
            return (
                <div className="column is-narrow buttons has-addons item-filter-controls">
                    <button className={`button is-naked is-underlined ${this.getActiveFilterClass(MediaMiniatureFilter.NONE)}`}
                        onClick={() => this.props.updateFilter(MediaMiniatureFilter.NONE)}>
                        All
                    </button>
                    <button className={`button is-naked is-underlined ${this.getActiveFilterClass(MediaMiniatureFilter.TRACKED)}`}
                        onClick={() => this.props.updateFilter(MediaMiniatureFilter.TRACKED)}>
                        Tracked
                    </button>
                    <button className={`button is-naked is-underlined ${this.getActiveFilterClass(MediaMiniatureFilter.FUTURE)}`}
                        onClick={() => this.props.updateFilter(MediaMiniatureFilter.FUTURE)}>
                        Future
                    </button>
                    <button className={`button is-naked is-underlined ${this.getActiveFilterClass(MediaMiniatureFilter.DOWNLOADED)}`}
                            onClick={() => this.props.updateFilter(MediaMiniatureFilter.DOWNLOADED)}>
                        Downloaded
                    </button>
                    <button className={`button is-naked is-underlined ${this.getActiveFilterClass(MediaMiniatureFilter.REMOVED)}`}
                        onClick={() => this.props.updateFilter(MediaMiniatureFilter.REMOVED)}>
                        Removed
                    </button>
                </div>
            );
        } else if (this.props.type === "tvshow") {
            return (
                <div className="column is-narrow buttons has-addons">
                    <button className={`button is-naked is-underlined ${this.getActiveFilterClass(MediaMiniatureFilter.NONE)}`}
                            onClick={() => this.props.updateFilter(MediaMiniatureFilter.NONE)}>
                        All
                    </button>
                    <button className={`button is-naked is-underlined ${this.getActiveFilterClass(MediaMiniatureFilter.TRACKED)}`}
                            onClick={() => this.props.updateFilter(MediaMiniatureFilter.TRACKED)}>
                        Tracked
                    </button>
                    <button className={`button is-naked is-underlined ${this.getActiveFilterClass(MediaMiniatureFilter.REMOVED)}`}
                            onClick={() => this.props.updateFilter(MediaMiniatureFilter.REMOVED)}>
                        Removed
                    </button>
                </div>
            )
        }

        return null;
    }
}

export default ItemFilterControls;
