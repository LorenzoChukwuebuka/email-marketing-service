import { HelmetProvider, Helmet } from "react-helmet-async"
import ForgotPasswordComponent from "../components/forgotpasswordcomponent"
import useMetadata from "../../../hooks/useMetaData"

const ForgetPasswordTemplate: React.FC = () => {
    const meta = useMetadata("ForgotPassword")
    return (
        <HelmetProvider>
            <Helmet {...meta} />
            <ForgotPasswordComponent />
        </HelmetProvider>
    )
}

export default ForgetPasswordTemplate