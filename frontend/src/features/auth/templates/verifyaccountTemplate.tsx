import { Helmet, HelmetProvider } from "react-helmet-async"
import useMetadata from "../../../hooks/useMetaData"
import VerifyAccountComponent from "../components/verifyusercomponent"

const VerifyAccountTemplate = () => {
    const meta = useMetadata("AccountVerification")
    return (
        <HelmetProvider>
            <Helmet {...meta} />
            <VerifyAccountComponent />
        </HelmetProvider>
    )
}

export default VerifyAccountTemplate