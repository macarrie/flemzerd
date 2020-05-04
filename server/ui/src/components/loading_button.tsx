import React from "react";

import {RiRefreshLine} from "react-icons/ri";

type Props = {
    action(): Promise<any>,
    text :string,
    loading_text :string,
};

type State = {
    loading :boolean,
};

class LoadingButton extends React.Component<Props, State> {
    state :State = {
        loading: false,
    };

    constructor(props :Props) {
        super(props);

        this.performAction = this.performAction.bind(this);
    }

    performAction() {
        this.setState({loading: true});
        this.props.action().then(() => {
            this.setState({loading: false});
        })
    }

    render() {
        if (this.state.loading) {
            return (
                <button className="button is-loading">
                    <i>{this.props.loading_text}</i>
                </button>
            );
        }

        return (
            <button className="button is-naked"
                onClick={this.performAction}>
                <span className="icon is-small has-text-grey-light">
                    <RiRefreshLine />
                </span>
                <span>
                    {this.props.text}
                </span>
            </button>
        );
    }
}

export default LoadingButton;
