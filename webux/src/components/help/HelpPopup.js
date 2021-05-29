import React from 'react';
import Tippy from '@tippyjs/react';
import "./HelpButtonPopup.css";
import 'tippy.js/dist/tippy.css';
import 'tippy.js/themes/light-border.css';

const HelpPopup = (props) => {
    return (
        <Tippy content={<div className="helpTooltip">{props.content}</div>} interactive={true} interactiveBorder={20} delay={100}
            placement={props.placement ? props.placement : "top"}
            theme="light-border">
            <div>
                <span className="form-help" data-bs-toggle="popover" data-bs-placement="top" data-bs-html="true" data-bs-original-title="" title="">?</span>
                {(props.label) ? props.label : ""}
            </div>
        </Tippy>
    );
};

export default HelpPopup;