type EventCallback<T = unknown> = (payload: T) => void;

interface Events {
    [key: string]: EventCallback[];
}

const events: Events = {};

const eventBus = {
    emit<T = unknown>(event: string, payload: T) {
        if (events[event]) {
            events[event].forEach(callback => callback(payload));
        }
    },

    on<T = unknown>(event: string, callback: EventCallback<T>) {
        if (!events[event]) {
            events[event] = [];
        }
        events[event].push(callback as EventCallback);
    },

    off<T = unknown>(event: string, callback: EventCallback<T>) {
        if (!events[event]) return;
        events[event] = events[event].filter(cb => cb !== callback);
    }
};

export default eventBus;
