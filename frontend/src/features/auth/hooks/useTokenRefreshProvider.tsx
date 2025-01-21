import { FC, ReactNode } from 'react';
import { useTokenRefresh } from './useTokenRefreshQuery';

interface TokenRefreshProviderProps {
    children: ReactNode;
    onRefreshError?: (error: Error) => void;
}

export const TokenRefreshProvider: FC<TokenRefreshProviderProps> = ({ children, onRefreshError }) => {
    const { error } = useTokenRefresh();

    // Handle refresh errors
    if (error && onRefreshError) {
        onRefreshError(error as Error);
    }

    return <div>{ children } </div>;
};