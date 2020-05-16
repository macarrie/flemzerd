import React from "react";

import ModuleSettings from "../../types/module_settings";

import {RiArrowDropDownLine} from "react-icons/ri";
import {RiArrowDropUpLine} from "react-icons/ri";

interface Props {
    indexer :ModuleSettings;
};

type State = {
    indexer: ModuleSettings,
    display_detailed_contents :Boolean,
};

class TorznabIndexerComponent extends React.Component<Props, State> {
    state: State = {
        indexer: {} as ModuleSettings,
        display_detailed_contents: false,
    };

    constructor(props: Props) {
        super(props);

        this.state.indexer = this.props.indexer;

        this.toggleContent = this.toggleContent.bind(this);
    }

    componentWillReceiveProps(nextProps: Props) {
        this.setState({ indexer: nextProps.indexer });
    }

    toggleContent() {
        this.setState({ display_detailed_contents: !this.state.display_detailed_contents });
    }

    getContentToggleIcon() {
        return (
            <button className={"button is-naked is-small"}
                    onClick={this.toggleContent}>
                <span className={"icon"}>
                {this.state.display_detailed_contents ? (
                    <RiArrowDropUpLine />
                ) : (
                    <RiArrowDropDownLine />
                )}
                </span>
            </button>
        );
    }

    render() {
        if (this.state.indexer == null) {
            return (
                <li>Loading</li>
            );
        }

        return (
            <div className={"columns is-mobile is-gapless is-multiline is-vcentered"}>
                <div className="column"
                     onClick={this.toggleContent}>
                    {this.state.indexer.Config.name}
                </div>
                <div className="column is-narrow">
                    {this.getContentToggleIcon()}
                </div>
                {this.state.display_detailed_contents && (
                    <div className={"column is-full is-clipped notification-content"}>
                        <div>
                            <hr />
                            <code>{this.state.indexer.Config.url}</code>
                        </div>
                    </div>
                )}
            </div>
        );
    }
}

export default TorznabIndexerComponent;
