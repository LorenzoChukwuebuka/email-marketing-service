import { Helmet, HelmetProvider } from "react-helmet-async";
import useMetadata from "../../../hooks/useMetaData";
import { OTPTemplate } from "../../../templates";

const OTPPage: React.FC = () => {
    const metaData = useMetadata()("AccountVerification")
    return (

        <HelmetProvider>
            <Helmet {...metaData} />
            <OTPTemplate />
        </HelmetProvider>

    )
}

export default OTPPage
