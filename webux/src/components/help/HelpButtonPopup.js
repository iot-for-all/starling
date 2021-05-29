import { React, useState } from 'react';
import Tippy from '@tippyjs/react';
import {
    Button,
} from "tabler-react"; import "./HelpPopup.css";
import 'tippy.js/dist/tippy.css';
import 'tippy.js/themes/light-border.css';

const HelpButtonPopup = (props) => {
    const [visible, setVisible] = useState(false);
    const show = () => setVisible(true);
    const hide = () => setVisible(false);

    return (
        <Tippy content={<div className="helpButtonTooltip">{props.content}</div>} interactive={true} interactiveBorder={20} delay={100}
            placement={props.placement ? props.placement : "top"}
            theme="light-border"
            visible={visible} onClickOutside={hide}
            >
            <div>
                <Button
                    type="button"
                    color="primary"
                    size="sm"
                    icon="info"
                    outline
                    className="ml-2"
                    onClick={visible ? hide : show}
                >
                    {(props.label) ? props.label : ""}
                </Button>
            </div>
        </Tippy>
    );
};

export default HelpButtonPopup;