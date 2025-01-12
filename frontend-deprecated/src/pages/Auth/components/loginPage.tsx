import { Helmet, HelmetProvider } from "react-helmet-async";
import { LoginTemplate } from "../../../templates";
import useMetadata from "../../../hooks/useMetaData";


const LoginPage: React.FC = () => {
    const metadata = useMetadata()('Login');

    return (
        <>
            <HelmetProvider>
                <Helmet {...metadata} />
                <LoginTemplate />
            </HelmetProvider>
        </>

    )
}


export default LoginPage
