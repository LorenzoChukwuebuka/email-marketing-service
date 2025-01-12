import { HelmetProvider, Helmet } from "react-helmet-async";
import UserDashLayout from "../layouts/userDashLayout";
import useMetadata from "../hooks/useMetaData";

const DashBoardTemplate: React.FC = () => {
    const meta = useMetadata('Dashboard')
    return (
        <HelmetProvider>
            <Helmet {...meta} />
            <UserDashLayout />
        </HelmetProvider>

    )
}

export default DashBoardTemplate;