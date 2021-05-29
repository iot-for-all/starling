import {
    Text
} from "tabler-react";
import "tabler-react/dist/Tabler.css";
import "./Toolbar.css";

const Toolbar = (props) => {

    return (
        <div className="tableToolbar">
            <div className="tableToolbarButtons">
                {props.children}
            </div>
            <div className="tableToolbarStatus">
                <Text size="sm" >{props.countMessage}</Text>
            </div>
        </div>
    );
};

export default Toolbar;