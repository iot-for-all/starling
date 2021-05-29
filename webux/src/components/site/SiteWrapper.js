import {
    Page,
} from "tabler-react";

import Navbar from './Navbar';
import PageFooter from './PageFooter';
import "./SiteWrapper.css";

const SiteWrapper = (props) => {
    return (
        <Page>
            <Page.Main className="pageMain">
                <Navbar />
                {props.children}
            </Page.Main>
            <PageFooter />
        </Page>
    );
};

export default SiteWrapper;