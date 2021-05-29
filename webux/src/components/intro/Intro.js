import { useContext, useEffect } from 'react';
import {
    Card,
    Grid,
    Header,
    Text
} from "tabler-react";
import "tabler-react/dist/Tabler.css";
import GlobalContext from '../../context/globalContext';
import IntroCard from './IntroCard';



const Intro = (props) => {
    const globalContext = useContext(GlobalContext)

    // Called on mount to ensure reference data is loaded if coming from shortcut
    useEffect(() => {
        if (!globalContext.initialized) {
            globalContext.initializeData();
        }
    }, [globalContext])

    const stepStatus = [
        {stepCompleted: globalContext.models.length > 0},
        {stepCompleted: globalContext.apps.length > 0},
        {stepCompleted: globalContext.simulations.length > 0},
        {stepCompleted: globalContext.simulations.length > 0},
    ];
    let remainingSteps = stepStatus.reduce((currentNumber, obj) => {
        return currentNumber+(!obj.stepCompleted ? 1 : 0)}, 0);

    return (
        <Card>
            <Card.Body>
                <div>
                    <Grid.Row>
                        <Grid.Col>
                            <Header size={3}>Get started with Starling in 4 easy steps</Header>
                        </Grid.Col>
                        <Grid.Col>
                            <div className="float-right">
                                <Text size="sm"><strong>{remainingSteps} steps remaining</strong></Text>
                            </div>
                        </Grid.Col>
                    </Grid.Row>
                </div>
                <p>
                    Starling is a device simulator for IoT Central. Starling can simulate several types of devices
                    sending data at different rates. Several simulations can be executed simultaneously.
                    Complete the steps below to get started.
                </p>
                <Grid.Row cards deck>
                    <Grid.Col sm={6} xl={3}>
                        <IntroCard
                            imgSrc={"./images/intro1.jpg"}
                            imgAlt={"Add an IoT Central app"}
                            title={"Add an IoT Central app"}
                            description={
                                "Devices are created in an IoT Central application. Here you can configure your application credentials."
                            }
                            introNumber={1}
                            actionName={"Add Application"}
                            actionIcon={"plus"}
                            actionUrl={"/app/add?new&intro"}
                            actionTooltip={"Add an IoT Central application"}
                            statusIsComplete={stepStatus[1].stepCompleted}
                        />
                    </Grid.Col>
                    <Grid.Col sm={6} xl={3}>
                        <IntroCard
                            imgSrc={"./images/intro2.jpg"}
                            imgAlt={"Add a Device Model"}
                            title={"Add a Device Model"}
                            description={<div>
                                Device model describes a device using <a 
                                    href="https://github.com/Azure/opendigitaltwins-dtdl/blob/master/DTDL/v2/dtdlv2.md">DTDL</a>
                                . You only need to add device models once, as they can be shared across applications.
                            </div>}
                            introNumber={2}
                            actionName={"Add Device Model"}
                            actionIcon={"plus"}
                            actionUrl={"/model/add?new&intro"}
                            actionTooltip={"Add a device model"}
                            statusIsComplete={stepStatus[0].stepCompleted}
                        />
                    </Grid.Col>
                    <Grid.Col sm={6} xl={3}>
                        <IntroCard
                            imgSrc={"./images/intro3.jpg"}
                            imgAlt={"Add a Simulation"}
                            title={"Add a Simulation"}
                            description={
                                "Add a simulation targeting an IoT Central application. You can simulate several devices per device model, control rates of telemetry etc. per simulation"
                            }
                            introNumber={3}
                            actionName={"Add Simulation"}
                            actionIcon={"plus"}
                            actionUrl={"/sim/add?new&intro"}
                            actionTooltip={"Create a new simulation"}
                            statusIsComplete={stepStatus[2].stepCompleted}
                        />
                    </Grid.Col>
                    <Grid.Col sm={6} xl={3}>
                        <IntroCard
                            imgSrc={"./images/intro4.jpg"}
                            imgAlt={"Start Simulation"}
                            title={"Start Simulation"}
                            description={
                                "Once a simulation is created, you can start the simulation. Devices will be automatically provisioned when the simulation is started."
                            }
                            introNumber={4}
                            actionName={""}
                            actionIcon={""}
                            actionUrl={""}
                            actionTooltip={"Start Simulation"}
                            statusIsComplete={stepStatus[3].stepCompleted}
                        />
                    </Grid.Col>
                </Grid.Row>
            </Card.Body>
        </Card>
    );
};

export default Intro;