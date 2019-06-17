import React from "react";

class ModuleStatus extends React.Component {
    render() {
        const module = this.props.module;
        return (
            <li>
                <div className="uk-flex uk-flex-middle uk-width-expand">
                    <div className="uk-width-expand">{module.Name}</div>
                    {(!module.Status.Alive) ? (
                        <div className="uk-label uk-label-danger">ERROR</div>
                    ) : (
                        <div className="uk-label uk-label-success">OK</div>
                    )}
                </div>

                {(!module.Status.Alive) && (
                    <div>
                        <hr className="uk-margin-small-top" />
                        <div className="uk-alert uk-alert-danger">
                            {module.Status.Message}
                        </div>
                    </div>
                )}
            </li>
        );
    }
}

 export default ModuleStatus;
