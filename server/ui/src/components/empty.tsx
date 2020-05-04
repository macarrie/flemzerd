import React from "react";

class Empty extends React.Component<any, any> {
    render() {
        return (
            <div className="container">
                <div className="columns is-centered empty-section-placeholder is-centered">
                    <div className={"column is-narrow"}>
                        <i>
                            <div className="title is-3 has-text-grey-light">
                                {this.props.label}
                            </div>
                        </i>
                    </div>
                </div>
            </div>
        );
    }
}

export default Empty;
