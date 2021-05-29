import {
    Card,
    Container,
    List,
    Grid,
} from "tabler-react";
import "tabler-react/dist/Tabler.css";

import "./ListDetails.css";
const ListDetails = (props) => {
    return (
        <Card className="listCard">
            <Card.Body>
                <div className="my-3 my-md-5 listDetails">
                    <Container>
                        <Grid.Row>
                            <Grid.Col md={3}>
                                <h4 className="mb-5">{props.listTitle}</h4>
                                <div>
                                    <List.Group>
                                        {props.list}
                                    </List.Group>
                                </div>
                            </Grid.Col>
                            <Grid.Col md={9}>
                                {props.detailsForm}
                            </Grid.Col>
                        </Grid.Row>
                    </Container>
                </div>
            </Card.Body>
        </Card>
    );
};

export default ListDetails;