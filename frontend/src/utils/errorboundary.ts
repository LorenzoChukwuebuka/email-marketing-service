import React, { ErrorInfo, ReactNode } from 'react';

class ErrorBoundary extends React.Component<{ children: ReactNode }, { hasError: boolean }> {
    constructor(props: { children: ReactNode }) {
        super(props);
        this.state = { hasError: false };
    }

    static getDerivedStateFromError(_: Error) {
        return { hasError: true };
    }

    componentDidCatch(error: Error, errorInfo: ErrorInfo) {
        console.log('Uncaught error:', error, errorInfo);
    }

    render() {
        if (this.state.hasError) {
            return "<h1>Something went wrong.</h1>";
        }

        return this.props.children;
    }
}

export default ErrorBoundary;