import React from "react";

class Loading extends React.Component {
    state = {
        load_error: null,
    };

    constructor(props) {
        super(props);

        this.state.load_error = this.props.load_error;
    }

    componentWillReceiveProps(nextProps, nextContext) {
        this.setState({load_error: nextProps.load_error});
    }

    render() {
        if (this.state.load_error != null) {
            console.log("Load error: ", this.state.load_error);
            return (
                <div className="uk-container uk-height-1-1">
                    <div className="uk-flex uk-flex-middle uk-flex-center uk-margin-large-top">
                        <div className="uk-icon uk-text-muted" data-uk-icon="icon: close; ratio: 3"></div>
                        <div>
                            <span className="uk-h1 uk-text-muted">
                                <i>
                                    Loading error
                                </i>
                            </span>
                        </div>
                    </div>
                </div>
            );
        }

        return (
            <div className="uk-container uk-height-1-1">
                <div className="uk-flex uk-flex-middle uk-flex-center uk-margin-large-top">
                    <span className="uk-margin-right" data-uk-spinner="ratio: 2"></span>
                    <div>
                            <span className="uk-h1 uk-text-muted">
                                <i>
                                    Loading
                                </i>
                            </span>
                    </div>
                </div>
            </div>
        );
    }
}

export default Loading;
