import React from "react";

import Const from "../../const";

class ConfigCheck extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            item: null,
        };

        this.state.item = this.props.item;
    }

    componentWillReceiveProps(nextProps) {
        this.setState({item: nextProps.item});
    }

    getClassName() {
        switch (this.state.item.Status) {
            case Const.WARNING:
                return "uk-label-warning";
            case Const.CRITICAL:
                return "uk-label-danger";
            default:
                return "";
        }
    }

    getStatusStr() {
        switch (this.state.item.Status) {
            case Const.WARNING:
                return "WARNING";
            case Const.CRITICAL:
                return "ERROR";
            default:
                return "";
        }
    }

    render() {
        return (
            <li className="uk-flex uk-flex-middle">
                <span className={`configuration-error-status uk-label uk-margin-small-right ${this.getClassName()}`}>
                    {this.getStatusStr()}
                </span>
                {this.state.item.Message}

                {this.state.item.Key !== "" && (
                    <code>{this.state.item.Key}: {this.state.item.Value}</code>
                )}
            </li>
        );
    }
}

export default ConfigCheck;
