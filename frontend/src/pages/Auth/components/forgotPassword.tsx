
import { ForgotPasswordTemplate } from "../../../templates";
import useMetadata from '../../../hooks/useMetaData';
import { HelmetProvider, Helmet } from 'react-helmet-async';


const ForgotPassword: React.FC = () => {
    const metadata = useMetadata()("ForgotPassword")
    return (
        <>
            <HelmetProvider> <Helmet {...metadata} />
                <ForgotPasswordTemplate />
            </HelmetProvider>
        </>

    );
}

export default ForgotPassword;