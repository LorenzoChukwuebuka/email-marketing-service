import { Helmet, HelmetProvider } from "react-helmet-async";
import { ResetPasswordTemplate } from "../../../templates";
import useMetadata from "../../../hooks/useMetaData";


const ResetPassword: React.FC = () => {
    const metatdata = useMetadata()("ResetPassword")
    return (
        <>
            <HelmetProvider>
                <Helmet {...metatdata} />
                <ResetPasswordTemplate />
            </HelmetProvider>
        </>
    );
}


export default ResetPassword