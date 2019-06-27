import React from "react";

class TraktAuth extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            config: null,
        };
    }

    componentDidMount() {
    }

    render() {
        return (
            <>
                <div className="uk-float-right">
                    <div className="uk-flex uk-flex-middle">
                        <div className="uk-margin-small-right"
                             data-uk-spinner="ratio: 0.5"></div>
                        <div>
                            <small><i>Loading</i></small>
                        </div>
                    </div>
                    <span className="uk-text-success">
                                <small><i>configured</i></small>
                            </span>
                    <button className="uk-button uk-button-default uk-text-danger uk-button-small uk-text-right">
                        Connect to trakt.tv
                    </button>
                </div>
                <div>
                    <hr className="uk-margin-small-top"/>
                    Enter the following validation code into <a target="_blank"
                                                                href="trakt_device_code.verification_url">trakt_device_code.verification_url</a>
                    <h5 className="uk-h2 uk-text-center uk-text-primary">trakt_device_code.user_code</h5>
                </div>
            </>
        );
    }
}

export default TraktAuth;
