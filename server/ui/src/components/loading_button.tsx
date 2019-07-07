import React from "react";

type Props = {
    action(),
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
                <button className="uk-button uk-button-text">
                    <div className="uk-flex uk-flex-middle">
                        <div data-uk-spinner="ratio: 0.5"></div>
                        <div>
                            <i>{this.props.loading_text}</i>
                        </div>
                    </div>
                </button>
            );
        }

        return (
            <button className="uk-button uk-button-text"
                onClick={this.performAction}>
                <div className="uk-flex uk-flex-middle">
                    <div className="uk-float-left">
                        <span className="uk-icon" data-uk-icon="icon: refresh; ratio: 0.75"></span>
                    </div>
                    <div className="uk-text-nowrap">
                        {this.props.text}
                    </div>
                </div>
            </button>
        );
    }
}

export default LoadingButton;
