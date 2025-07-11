import { Helmet, HelmetProvider } from "react-helmet-async"
import useMetadata from "../../../hooks/useMetaData"
import VerifyLoginComponent from "../components/verifyLoginComponent"

const VerifyLoginTemplate = ()=>{
    const meta = useMetadata("VerifyLogin")
    return (
        <HelmetProvider>
            <Helmet {...meta} />
            <VerifyLoginComponent />
        </HelmetProvider>
    )
}

export default VerifyLoginTemplate