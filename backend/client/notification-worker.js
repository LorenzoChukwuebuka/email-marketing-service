const NOTIFICATION_INTERVAL = 60000; // Check every minute
const API_BASE_URL = 'http://localhost:9000/api/v1'; // Replace with your actual API base URL
let jwtToken = null;

self.addEventListener('install', event => {
  console.log('Service Worker installed', event);
  self.skipWaiting();
});

self.addEventListener('activate', event => {
  console.log('Service Worker activated');
  event.waitUntil(self.clients.claim());
});

self.addEventListener('message', event => {
  if (event.data && event.data.type === 'START_NOTIFICATION_CHECK') {
    console.log('Received START_NOTIFICATION_CHECK message');
    jwtToken = event.data.token;
    startNotificationCheck();
  }
});

function startNotificationCheck() {
  console.log('Starting notification check');
  setInterval(() => {
    if (!jwtToken) {
      console.error('JWT token not available');
      return;
    }

    console.log('Fetching notifications');
    fetch(`${API_BASE_URL}/user-notifications`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${jwtToken}`,
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      mode: 'cors' // Explicitly set CORS mode
    })
      .then(response => {
        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
        return response.json();
      })
      .then(data => {
        if (data.status === true) {
        //  console.log('Notifications received:', data.payload);
          self.clients.matchAll().then(clients => {
            clients.forEach(client => {
              client.postMessage({
                type: 'NOTIFICATION_UPDATE',
                payload: data.payload
              });
            });
          });
        }
      })
      .catch(error => console.error('Error fetching notifications:', error));
  }, NOTIFICATION_INTERVAL);
}