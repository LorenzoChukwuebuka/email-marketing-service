import { useMemo } from 'react';
import { HelmetProps } from 'react-helmet-async';
import { metadataConfig } from '../utils/meta.config';

// Create a custom hook to handle the metadata configuration
const useMetadata = <K extends keyof typeof metadataConfig>(key: K): HelmetProps => {
    // Memoize the result to avoid unnecessary re-renders
    return useMemo(() => {
        // Retrieve the metadata for the given key
        const metadata = metadataConfig[key];

        // Return the metadata in the format expected by the Helmet component
        return {
            title: metadata.title,
            meta: Object.entries(metadata).reduce((acc, [metaKey, value]) => {
                // Convert camelCase to colon-case for Open Graph and Twitter tags
                if (metaKey.startsWith('og')) {
                    const formattedKey = metaKey.replace(/([A-Z])/g, ':$1').toLowerCase(); // Convert ogImage -> og:image
                    acc.push({
                        property: formattedKey,
                        content: String(value),
                    });
                } else if (metaKey.startsWith('twitter')) {
                    acc.push({
                        name: metaKey,
                        content: String(value),
                    });
                } else if (metaKey === 'description' || metaKey === 'keywords') {
                    acc.push({
                        name: metaKey,
                        content: String(value),
                    });
                }
                return acc;
            }, [] as { name?: string; property?: string; content: string }[]),
        };
    }, [key]);
};

export default useMetadata;
