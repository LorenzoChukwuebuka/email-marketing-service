import { Helmet, HelmetProvider } from "react-helmet-async"

const AdminUserSpecificCampaigns: React.FC = () => {
    return (
        <HelmetProvider>

            <Helmet title="user campaign" />

        </HelmetProvider>
    )
}

export default AdminUserSpecificCampaigns