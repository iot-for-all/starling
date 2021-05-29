import { store } from 'react-notifications-component';
import 'react-notifications-component/dist/theme.css';

export function addNotification(type, title, message) {
    store.addNotification({
        title: title,
        message: message,
        type: type,
        insert: "top",
        container: "top-right",
        animationIn: ["animate__animated", "animate__fadeIn"],
        animationOut: ["animate__animated", "animate__fadeOut"],
        dismiss: {
            duration: 2500,
            onScreen: false,
            pauseOnHover: true,
            waitForAnimation: false,
            showIcon: true,
            click: true,
            touch: true
        },
        slidingEnter: {
            duration: 300,
            timingFunction: 'linear',
            delay: 0
          },
        
          slidingExit: {
            duration: 300,
            timingFunction: 'linear',
            delay: 0
          },
    });
};

