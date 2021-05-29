import "tabler-react/dist/Tabler.css";
import Dashboard from '../components/dashboard/Dashboard';
import SiteWrapper from '../components/site/SiteWrapper';

const HomePage = () => {
    return (
        <SiteWrapper>
            <Dashboard />
        </SiteWrapper>
    );
};

export default HomePage;