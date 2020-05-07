import React from "react";

import {RiRefreshLine} from "react-icons/ri";

type Props = {
    action(): Promise<any>,
    text :string,
    loading_text :string,
    styleClass? :string,
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
        let classes = "button is-naked";
        if (this.props.styleClass) {
            classes = this.props.styleClass;
        }

        if (this.state.loading) {
            return (
                <button className={`${classes} is-loading`}>
                    <i>{this.props.loading_text}</i>
                </button>
            );
        }

        return (
            <button className={classes}
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
