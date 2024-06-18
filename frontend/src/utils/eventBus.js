// eventBus.js
let events = {};

const eventBus = {
  emit(event, payload) {
    if (events[event]) {
      events[event].forEach(callback => callback(payload));
    }
  },

  on(event, callback) {
    if (!events[event]) {
      events[event] = [];
    }
    events[event].push(callback);
  },

  off(event, callback) {
    if (!events[event]) return;
    events[event] = events[event].filter(cb => cb !== callback);
  }
};

export default eventBus;