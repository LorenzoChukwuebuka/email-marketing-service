// src/hooks/useTokenExpiration.ts
import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { jwtDecode } from 'jwt-decode';
import Cookies from 'js-cookie';
interface DecodedToken {
    exp: number;
}

const useTokenExpiration = () => {
    const navigate = useNavigate();

    useEffect(() => {
        const checkTokenExpiration = () => {
            const token = Cookies.get('Cookies');
            if (token) {
                const decodedToken = jwtDecode<DecodedToken>(token);
                const currentTime = Date.now() / 1000;

                if (decodedToken.exp < currentTime) {
                    Cookies.remove('Cookies')
                    navigate('/auth/login');
                }
            }
        };

        const interval = setInterval(checkTokenExpiration, 60000); // Check every minute

        return () => clearInterval(interval);
    }, [navigate]);
};

export default useTokenExpiration;