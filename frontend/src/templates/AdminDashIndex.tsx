import { Helmet, HelmetProvider } from "react-helmet-async";

const AdminDashIndexTemplate = () => {
    return (
        <HelmetProvider>
            <Helmet title="Admin Dashboard" />

            <div className="bg-white rounded-lg shadow-md p-6">
                <h2 className="text-2xl font-bold mb-4">Welcome Admin</h2>
            </div>
        </HelmetProvider>
    );
};

export default AdminDashIndexTemplate;
