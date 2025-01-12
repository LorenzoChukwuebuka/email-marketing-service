// eventBus.ts
type EventCallback = (payload: any) => void;

interface Events {
    [key: string]: EventCallback[];
}

let events: Events = {};

const eventBus = {
    emit(event: string, payload: any) {
        if (events[event]) {
            events[event].forEach(callback => callback(payload));
        }
    },

    on(event: string, callback: EventCallback) {
        if (!events[event]) {
            events[event] = [];
        }
        events[event].push(callback);
    },

    off(event: string, callback: EventCallback) {
        if (!events[event]) return;
        events[event] = events[event].filter(cb => cb !== callback);
    }
};

export default eventBus;
