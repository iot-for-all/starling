import * as React from "react";
import { Link } from 'react-router-dom';
import {
    Card,
    Icon
} from "tabler-react";
import "tabler-react/dist/Tabler.css";
import "./IntroCard.css";

const IntroCard = (props) => {
    let statusElement = <span className="text-warning" title="You need to complete this step.">
        <Icon prefix="fe" name={"clock"} /> Todo
    </span>;
    if (props.statusIsComplete === true) {
        statusElement = <span className="text-success" title="You have completed this step.">
            <Icon prefix="fe" name={"check"} /> Done
    </span>;
    }
    return (
        <Card>
            <div className="introImageContainer">
                <img className="card-img-top" src={props.imgSrc} alt={props.imgAlt} />
                <div className="introNumber">
                    {props.introNumber}
                </div>
            </div>
            <Card.Body className="d-flex flex-column">
                <h4>
                    {props.title}
                </h4>
                <div className="introDescription">
                    {props.description}
                </div>
                <div className="d-flex align-items-center pt-5 mt-auto">
                    <div>
                        {
                            !props.statusIsComplete
                            && props.actionName !== ""
                            && <span title={props.actionTooltip}>
                                <Link
                                    to={props.actionUrl}
                                    className="btn btn-sm btn-primary ml-2"
                                >
                                    <Icon prefix="fe" name={props.actionIcon} />
                                    {props.actionName}
                                </Link>
                            </span>
                        }
                    </div>
                    <div className="ml-auto ">
                        <span className="icon d-none d-md-inline-block ml-3 introStatus">
                            {statusElement}
                        </span>
                    </div>
                </div>
            </Card.Body>
        </Card>
    );
};

export default IntroCard;
