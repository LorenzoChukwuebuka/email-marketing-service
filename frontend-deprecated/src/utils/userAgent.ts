interface UserAgentInfo {
    name: string;
    version: string;
}

const parseUserAgent = (userAgent: string): UserAgentInfo => {
    const info: UserAgentInfo = {
        name: 'Unknown Browser',
        version: 'Unknown',
    };

    if (userAgent.includes('Chrome')) {
        info.name = 'Google Chrome';
    } else if (userAgent.includes('Firefox')) {
        info.name = 'Mozilla Firefox';
    } else if (userAgent.includes('Safari') && !userAgent.includes('Chrome')) {
        info.name = 'Apple Safari';
    } else if (userAgent.includes('Edge')) {
        info.name = 'Microsoft Edge';
    } else if (userAgent.includes('MSIE') || userAgent.includes('Trident/')) {
        info.name = 'Internet Explorer';
    }

    // Extract version
    const versionMatch = userAgent.match(
        /(Chrome|Firefox|Safari|Edge|MSIE|rv:|Opera|OPR)\/([\d.]+)/
    );
    if (versionMatch && versionMatch.length >= 3) {
        info.version = versionMatch[2];
    }

    return info;
};

export default parseUserAgent;
