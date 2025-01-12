import { Helmet, HelmetProvider } from "react-helmet-async";
import useMetadata from "../../../hooks/useMetaData";
import { SignUpTemplate } from "../../../templates";

const SignUpPage: React.FC = () => {
    const metadata = useMetadata()("Signup")
    return (
        <>
            <HelmetProvider>
                <Helmet {...metadata} />
                <SignUpTemplate />
            </HelmetProvider>

        </>
    );
}

export default SignUpPage
