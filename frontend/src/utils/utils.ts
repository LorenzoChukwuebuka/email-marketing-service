export function convertToNormalTime(timestamp: string): string {
    if (!timestamp) return '';

    // Extract date and time parts
    const [datePart, timePart] = timestamp.split(' ');
    const [year, month, day] = datePart.split('-');
    const [time] = timePart.split('.');
    const [hour, minute, second] = time.split(':');

    // Create a new Date object
    const date = new Date(Number(year), Number(month) - 1, Number(day), Number(hour), Number(minute), Number(second));

    // Check if the date is valid
    if (isNaN(date.getTime())) return 'Invalid Date';

    const options: Intl.DateTimeFormatOptions = {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit',
        timeZone: 'Africa/Lagos' // WAT is West Africa Time
    };

    return date.toLocaleString('en-US', options);
}

export function maskAPIKey(apiKey: string): string {
    if (apiKey) {
        // Ensure the API key is at least 6 characters long
        if (apiKey.length <= 6) {
            return apiKey;
        }

        // Get the first 3 and last 3 characters
        const start = apiKey.substring(0, 3);
        const end = apiKey.substring(apiKey.length - 3);

        // Replace the middle part with asterisks
        const masked = apiKey.replace(
            apiKey.substring(3, apiKey.length - 3),
            '******'
        );

        return `${start}${masked}${end}`;
    }

    return '';
}

function fallbackCopyToClipboard(text: string): void {
    const textArea = document.createElement('textarea');
    textArea.value = text;

    // Make the textarea out of viewport
    textArea.style.position = 'fixed';
    textArea.style.left = '-999999px';
    textArea.style.top = '-999999px';
    document.body.appendChild(textArea);
    textArea.focus();
    textArea.select();

    try {
        const successful = document.execCommand('copy');
        const msg = successful ? 'successful' : 'unsuccessful';
        console.log('Fallback: Copying text command was ' + msg);
    } catch (err) {
        console.error('Fallback: Oops, unable to copy', err);
    }

    document.body.removeChild(textArea);
}

export function copyToClipboard(text: string): void {
    if (!navigator.clipboard) {
        fallbackCopyToClipboard(text);
        return;
    }
    navigator.clipboard
        .writeText(text)
        .then(() => {
            console.log('Text copied to clipboard');
        })
        .catch(err => {
            console.error('Failed to copy text: ', err);
        });
}
