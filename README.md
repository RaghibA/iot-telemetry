# IoT Telemetry

IoT Telemetry is a tech demo for a real-time analytics platform. The purpose of this tool is to enable IoT devices to report device telemetry for real time access, but the concepts and tools used could be extended to any situation where real time analytics could be useful.

## Project Dependencies

This project and its dependencies run in Docker containers, orchestrated by Docker Compose. To use this tool, ensure Docker is installed on your system.

## Getting Started


Clone or download this repository and navigate to its root directory.  From there, run the following command to build and start the containers:

    docker compose up --build -d

Docker Compose will automatically handle volumes & networks as defined in the configuration.

Once the build process is complete, you should see the following **eight containers**:

 - auth-service-1
 - iot-admin-service-1
 - iot-data-service-1
 - consumer-service-1
 - iot-telem-db
 - kafka-1
 - zookeeper-1
 - nginx

The following sections will provide detailed instruction on how to use the tool. 

**If you prefer, you can use the provided postman collections to try it out in a more convenient way**

## Auth Service

The Auth-Service API manages user accounts & ALCs. It also handles the issuance of API keys, access tokens, & cookies. The following steps will create an account and generate an API key, a cookie, and an admin access token:

Create an account. Save the API key. This will be used to send telemetry from the IoT Device.

    POST: localhost/auth/register
    REQUEST BODY:
    {
      "email": "email@test.com",
      "username": "someusername",
      "password": "myPassword123"
    }

Log into the account. This will generate a cookie. You will need to be signed into your account to request an access token.

    POST: localhost/auth/login
    REQUEST BODY:
    {
      "username": "someusername",
      "password": "myPassword123"
    }

Generate an access token. This will be used for authorizing Admin API requests. 

**Make sure your cookie is being set with the value: refresh_token=<token_string>**

    POST: localhost/auth/access-token

If you forget your API Key you can get one with the following request. This **WILL** overwrite your old API key, so any devices you have set up previously will need to be updated to use the new API key.

    GET: localhost/auth/api-key

## Admin Service

To send telemetry from your IoT devices, you will need to first register them in the admin service.

**All admin requests must include an access token in the 'Authorization' header, formatted as a Bearer token.**

You can register a device with the following request:

    POST: localhost/admin/device
    REQUEST BODY: { "deviceName": "my-device-name"}

Creating a device will link you API key with the device & create a topic in the kafka server for your device. 

**All data from this device will ONLY be sent to its respective topic.**

If you need to send data from multiple devices to a single source, you must register 1 device & use that deviceID amongst all devices when sending telemetry. You can embed proprietary device identifiers in your data if needed.


You will need the deviceID to send telemetry. This can be found in the response body after creating a device, or by using the GET device request:

    GET: localhost/admin/device

## Data Service

The data service is where the device will report its telemetry. The data will then be forwarded by the service into the appropriate kafka topic.

To send telemetry, use the following request:

    POST: localhost/telemetry/send
    REQUEST BODY:
    {
      "deviceId": "your-device-id",
      "data": {
        "test1": [1,2,3,4,...],
        "data123": { "data": "123"},
        "some-id": "my-device-xyz"
      } // Any valid JSON object
    }

## Consumer Service

You will not be able to consume data directly from the kafka topics. In order to get real time data from your device, you will need to use the consumer service to establish a connection via websocket. 

To use the consumer service, provide 'Authorization' & 'x-device-id' as headers. This will create an instance of a kafka consumer subscribed to your device's topic.

    'Authorization': Bearer accessTokenString
    'x-device-id': device-id-for-real-time-analytics
    
    Connect: ws://localhost/consumer/telemetry/consume 



