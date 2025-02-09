# iot-telemetry
Real-time analytics platform for IoT devices.

## How it works
1. Create an account with the auth service.
2. Log into your Account to generate a cookie.
3. When logged in use auth API to issue access token.
4. Use access token to authenticate admin API requests.
5. Register a device with the admin API.
6. Send device telemetry to the telemetry service through the device with the API key.
7. Consume from the device topic to collect real-time analytics.
