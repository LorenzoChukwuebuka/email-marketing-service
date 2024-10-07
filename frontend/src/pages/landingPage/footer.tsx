import { Link } from "react-router-dom"

const Footer: React.FC = () => {
    return (
        <>

            <footer className="bg-gray-800 text-white py-12">
                <div className="container mx-auto px-4">
                    <div className="grid grid-cols-2 md:grid-cols-4 gap-8">
                        <div>
                            <h3 className="text-lg font-semibold mb-4">Product</h3>
                            <ul className="space-y-2">
                                <li><a href="#" className="hover:text-gray-300 transition duration-300">Features</a></li>
                                <li><a href="#" className="hover:text-gray-300 transition duration-300">Pricing</a></li>
                                <li><a href="#" className="hover:text-gray-300 transition duration-300">Integrations</a></li>
                            </ul>
                        </div>
                        <div>
                            <h3 className="text-lg font-semibold mb-4">Resources</h3>
                            <ul className="space-y-2">
                                <li><a href="#" className="hover:text-gray-300 transition duration-300">Blog</a></li>
                                <li><a href="#" className="hover:text-gray-300 transition duration-300">Help Center</a></li>
                                <li><a href="#" className="hover:text-gray-300 transition duration-300">Guides</a></li>
                            </ul>
                        </div>
                        <div>
                            <h3 className="text-lg font-semibold mb-4">Company</h3>
                            <ul className="space-y-2">
                                <li><a href="#" className="hover:text-gray-300 transition duration-300">About Us</a></li>
                                <li><a href="#" className="hover:text-gray-300 transition duration-300">Careers</a></li>
                                <li><a href="#" className="hover:text-gray-300 transition duration-300">Contact</a></li>
                            </ul>
                        </div>
                        <div>
                            <h3 className="text-lg font-semibold mb-4">Legal</h3>
                            <ul className="space-y-2">
                                <li><Link to="/privacy" className="hover:text-gray-300 transition duration-300">Privacy Policy</Link></li>
                                <li><Link to="/tos" className="hover:text-gray-300 transition duration-300">Terms of Service </Link> </li>
                                <li><a href="#" className="hover:text-gray-300 transition duration-300">GDPR</a></li>
                            </ul>
                        </div>
                    </div>
                    <div className="mt-12 pt-8 border-t border-gray-700 text-center">
                        <p>&copy; {new Date().getFullYear()} CrabMailer. All rights reserved.</p>
                    </div>
                </div>
            </footer>

        </>
    )
}

export default Footer