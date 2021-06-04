import React, { useState, Fragment } from "react";
import { NavLink, withRouter } from "react-router-dom";

import {
    Button,
    Container,
    Grid,
    Nav,
    Site,
} from "tabler-react";
import "./Navbar.css";

const navBarItems = [
    {
        value: "Home",
        to: "/",
        icon: "home",
        LinkComponent: withRouter(NavLink),
        useExact: true,
    },
    {
        value: "Applications",
        to: "/app",
        LinkComponent: withRouter(NavLink),
        useExact: false,
    },
    {
        value: "Models",
        to: "/model",
        LinkComponent: withRouter(NavLink),
        useExact: false,
    },
    {
        value: "Simulations",
        to: "/sim",
        LinkComponent: withRouter(NavLink),
        useExact: false,
    },
    {
        value: "Metrics",
        to: "/metrics",
        LinkComponent: withRouter(NavLink),
        icon: "trending-up",
        useExact: false,
    },
    {
        value: "Settings",
        to: "/settings",
        LinkComponent: withRouter(NavLink),
        icon: "settings",
        useExact: false,
    },
];


const Navbar = (props) => {
    const [collapse, setcollapse] = useState(false);

    const handleCollapseMobileMenu = () => {
        setcollapse(!collapse);
    };
    const navbarClasses = collapse ? "header d-lg-flex p-0 collapse" : "header d-lg-flex p-0";

    return (
        <Fragment>
            <div className="header darkHeader">
                <Container className={""}>
                    <div className="d-flex">
                        <Site.Logo href={"/"} alt={"Starling"} src={"/starling-light.png"} />
                        <div className="d-flex order-lg-2 ml-auto">
                            <Nav.Item type="div" className="d-none d-md-flex">
                                <Button
                                    href="https://github.com/iot-for-all/starling"
                                    target="_blank"
                                    size="sm"
                                    RootComponent="a"
                                    color="light"
                                    icon="github"
                                >Source code</Button>
                            </Nav.Item>
                        </div>
                        <Button
                            className="header-toggler d-lg-none ml-3 ml-lg-0 hamburgerBtn"
                            type="button"
                            outline
                            onClick={handleCollapseMobileMenu}
                        >
                            <span className="header-toggler-icon hamburger"></span>
                        </Button>
                    </div>
                </Container>
            </div>
            <div className={navbarClasses}>
                <Container>
                    <Grid.Row className="align-items-center">
                        <Grid.Col className="col-lg order-lg-first">
                            <Nav
                                tabbed
                                className="border-0 flex-column flex-lg-row"
                                collapse={collapse}
                                itemsObjects={navBarItems}
                            />
                        </Grid.Col>
                    </Grid.Row>
                </Container>
            </div>
        </Fragment>
    );
};

export default Navbar;