import { HelmetProvider, Helmet } from "react-helmet-async"
import LoginComponent from "../components/logincomponent"
import useMetadata from "../../../hooks/useMetaData"

const LoginTemplate: React.FC = () => {
    const meta = useMetadata("Login")
    return (
        <HelmetProvider>
            <Helmet {...meta} />
            <LoginComponent />
        </HelmetProvider>
    )
}

export default LoginTemplate