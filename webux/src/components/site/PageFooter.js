import {Fragment} from  "react";
import {
    Container,
    Grid,
    List
} from "tabler-react";

const PageFooter = (props) => {
    return (
        <Fragment>
            <footer className="footer">
                <Container>
                    <Grid.Row className="align-items-center flex-row-reverse">
                        <Grid.Col auto={true} className="ml-auto">
                            <Grid.Row className="align-items-center">
                                <Grid.Col auto={true}>
                                    <List className="list-inline list-inline-dots mb-0">
                                        <List.Item className="list-inline-item">
                                            <a
                                                href="https://github.com/iot-for-all/starling/readme.md"
                                                target="_blank"
                                                rel="noopener noreferrer"
                                            >Documentation</a>
                                        </List.Item>
                                        <List.Item className="list-inline-item">
                                            <a
                                                href="https://github.com/iot-for-all/starling/docs/faq.md"
                                                target="_blank"
                                                rel="noopener noreferrer"
                                            >FAQ</a>
                                        </List.Item>
                                        <List.Item className="list-inline-item">
                                            <a href="https://github.com/iot-for-all/starling/issues">Issues</a>
                                        </List.Item>
                                    </List>
                                </Grid.Col>
                            </Grid.Row>
                        </Grid.Col>
                        <Grid.Col width={12} lgAuto className="mt-3 mt-lg-0 text-center">
                            Open source Azure IoT Central Device Simulator
                    </Grid.Col>
                    </Grid.Row>
                </Container>
            </footer>
        </Fragment>
    );
};

export default PageFooter;