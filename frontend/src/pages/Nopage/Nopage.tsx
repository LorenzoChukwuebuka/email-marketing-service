import "./Nopage.css";
import { useNavigate } from "react-router-dom";
import { Button, Result } from "antd";
import Cookies from "js-cookie";

const Nopage = () => {
    const navigate = useNavigate();

    const handleRedirect = () => {
        const officeCookies = Cookies.get('OfficeCookies');
        if (officeCookies) {
            navigate('/app');
        } else {
            navigate('/');
        }
    };

    return (
        <div>
            <Result
                status="404"
                title="404"
                subTitle="Sorry, See you are on a broken page"
                extra={
                    <Button type="primary" onClick={handleRedirect}>
                        Go to Home
                    </Button>
                }
            />
        </div>
    );
};

export default Nopage;