import { Link } from 'react-router-dom';
import {
    Icon,
} from "tabler-react";
import "tabler-react/dist/Tabler.css";

import "./NoDataFoundCard.css";

const NoDataFoundCard = (props) => {
    const noDataImage = (props.noDataImage) ? props.noDataImage : "/images/nodata.svg";
    const actions = props.actions.map((action) => {
        return (<Link
            key={action.actionName}
            to={action.actionUrl}
            className="btn btn-sm btn-primary mr-2"
        >
            <Icon prefix="fe" name={action.actionIcon} />
            {action.actionName}
        </Link>)
    });
    return (
        <div className="empty">
            <div className="empty-image">
                <img src={noDataImage} alt={props.message} />
            </div>
            <h5>{props.message}</h5>
            <p className="empty-subtitle ">
                {props.description}
            </p>
            {
                props.description2 &&
                <p className="empty-subtitle ">
                    {props.description2}
                </p>
            }
            <div className="empty-action">
                {actions}
            </div>
        </div>
    );
};

export default NoDataFoundCard;