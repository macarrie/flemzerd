import React from "react";

class Loading extends React.Component {
    constructor(props) {
        super(props);
    }

    render() {
        return (
            <div className="uk-container">
                <div className="uk-grid uk-child-width-1-1" data-uk-grid>
                    <div className="uk-text-center uk-flex uk-flex-middle">
                        <div className="uk-margin-right"
                            data-uk-spinner="ratio: 2"></div>
                        <div>
                            <h1 className="uk-text-muted">
                                    <i>
                                        Loading
                                    </i>
                            </h1>
                        </div>
                    </div>
                </div>
            </div>
        )
    }
}

export default Loading;
