import { Helmet, HelmetProvider } from "react-helmet-async";
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import renderApiName from '../../../../utils/name';
import useSenderStore from "../../../../store/userstore/senderStore";


const VerifySenderComponent: React.FC = () => {
    const navigate = useNavigate();
    const [isLoading, setIsLoading] = useState<boolean>(false)
    const { setVerifySender, verifySender } = useSenderStore()

    const handleVerify = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();

        // Extract token from URL
        const token = new URLSearchParams(window.location.search).get("token");
        const userId = new URLSearchParams(window.location.search).get("userId")
        const email = new URLSearchParams(window.location.search).get("email")

        setVerifySender({
            email: email as string,
            user_id: userId as string,
            token: token as string
        })

        await verifySender()
    };


    return (
        <HelmetProvider>
            <Helmet title="Verify Sender" />

            <div className="container mx-auto mt-[10em] px-4">
                <div className="max-w-lg mx-auto mt-5">

                    <h3 className="text-2xl font-bold text-center mb-4">
                        <a href="/"> {renderApiName()} </a>
                    </h3>
                    <div className="bg-white shadow-md rounded-lg p-8">

                        <h3 className="text-2xl font-bold text-center mb-4">
                            Verify Sender Email
                        </h3>

                        <form onSubmit={handleVerify}>
                            <div className="text-center">
                                {!isLoading ? (
                                    <button
                                        className="w-full bg-gray-800 text-white py-2 px-4 rounded-md hover:bg-gray-700"
                                        type="submit"
                                    >
                                        Verify Sender
                                    </button>
                                ) : (
                                    <button className="w-full bg-gray-800 text-white py-2 px-4 rounded-md hover:bg-gray-700" disabled>
                                        Please wait{" "}
                                        <span className="loading loading-dots loading-sm"></span>
                                    </button>
                                )}
                            </div>
                        </form>


                    </div>
                </div>
            </div>
        </HelmetProvider>
    );
};

export default VerifySenderComponent;
