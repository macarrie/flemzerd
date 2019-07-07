import React from "react";

class Loading extends React.Component<any, any> {
    render() {
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
