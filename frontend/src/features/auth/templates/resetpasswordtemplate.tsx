import { HelmetProvider, Helmet } from "react-helmet-async"
import ResetPasswordComponent from "../components/resetpasswordcomponent"

const ResetPasswordTemplate: React.FC = () => {
    return (
        <HelmetProvider>
            <Helmet />

            <ResetPasswordComponent />

        </HelmetProvider>
    )
}

export default ResetPasswordTemplate