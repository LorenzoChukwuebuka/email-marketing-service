import { Link } from "react-router-dom"

const IndexLandingPage: React.FC = () => {
    return (
        <>
            <header className="bg-white py-4 shadow-sm p-4">
                <div className="container mx-auto flex justify-between items-center">
                    <div className="text-2xl font-bold">{import.meta.env.VITE_API_NAME}</div>
                    <div>
                        <button className="bg-black hover:bg-gray-600 text-white px-4 py-2 rounded-md mr-2">
                            <Link to="/auth/login"> Login </Link>
                        </button>
                        <button className="bg-black hover:bg-gray-600 text-white px-4 py-2 rounded-md mr-2">
                            <Link to="/auth/sign-up"> Sign up </Link>
                        </button>
                    </div>
                </div>
            </header>
        </>
    )
}
export default IndexLandingPage