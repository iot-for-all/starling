import { Link } from 'react-router-dom';
import {
    Icon,
} from "tabler-react";
import "tabler-react/dist/Tabler.css";

import "./NoDataFoundCard.css";

const NoDataFoundCard = (props) => {
    const noDataImage = (props.noDataImage) ? props.noDataImage : "/images/nodata.svg";
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
                <Link
                    to={props.actionUrl}
                    className="btn btn-sm btn-primary"
                >
                    <Icon prefix="fe" name={props.actionIcon} />
                    {props.actionName}
                </Link>
            </div>
        </div>
    );
};

export default NoDataFoundCard;