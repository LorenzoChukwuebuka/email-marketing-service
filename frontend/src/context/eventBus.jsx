import { createContext, useContext, useState } from "react";

const EventBusContext = createContext();

export const EventBusProvider = ({ children }) => {
  const [events, setEvents] = useState({});

  const emit = (event, payload) => {
    if (events[event]) {
      events[event].forEach((callback) => callback(payload));
    }
  };

  const on = (event, callback) => {
    setEvents((prevEvents) => {
      const newEvents = { ...prevEvents };
      if (!newEvents[event]) {
        newEvents[event] = [];
      }
      newEvents[event].push(callback);
      return newEvents;
    });
  };

  const off = (event, callback) => {
    setEvents((prevEvents) => {
      const newEvents = { ...prevEvents };
      if (!newEvents[event]) return newEvents;
      newEvents[event] = newEvents[event].filter((cb) => cb !== callback);
      return newEvents;
    });
  };

  const eventBus = { emit, on, off };

  return (
    <EventBusContext.Provider value={eventBus}>
      {children}
    </EventBusContext.Provider>
  );
};

export const useEventBus = () => useContext(EventBusContext);
