import { HelmetProvider, Helmet } from "react-helmet-async"
import SignUpComponent from "../components/signupcomponent"
import useMetadata from "../../../hooks/useMetaData"

const SignupTemplate: React.FC = () => {
    const meta = useMetadata('Signup')
    return (
        <HelmetProvider>
            <Helmet {...meta} />
            <SignUpComponent />
        </HelmetProvider>
    )
}

export default SignupTemplate