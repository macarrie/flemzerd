import React from "react";

import Module from "../../types/module";

type Props = {
    module :Module,
};

class ModuleStatus extends React.Component<Props, any> {
    render() {
        const module = this.props.module;
        return (
            <li>
                <div className="columns is-multiline is-gapless">
                    <div className="column">{module.Name}</div>
                    {(!module.Status.Alive) ? (
                        <div className="column is-narrow">
                            <div className="tag is-danger">ERROR</div>
                        </div>
                        ) : (
                        <div className="column is-narrow">
                            <div className="tag is-success">OK</div>
                        </div>
                    )}

                    {(!module.Status.Alive) && (
                        <div className={"column is-full"}>
                            <hr className="" />
                            <article className="message is-danger">
                                <div className={"message-body"}>
                                    {module.Status.Message}
                                </div>
                            </article>
                        </div>
                    )}
                </div>
            </li>
        );
    }
}

 export default ModuleStatus;
