# API Gateway

This is an API Gateway service built in Node.js, designed to handle incoming requests and route them to the appropriate services. The service acts as a middleman between clients, other services, and the database, providing a unified entry point for accessing multiple services and managing data.

## Environment Variables

The following environment variables are required to configure the API Gateway service:

-   HOST: The host address on which the API Gateway service will run.
-   PORT: The port number on which the API Gateway service will listen for incoming requests.
-   REDIS_HOST: The host address of the Redis server, used for caching and session management.
-   REDIS_PORT: The port number of the Redis server.
-   REDIS_USERNAME: The username to authenticate with the Redis server.
-   REDIS_PASSWORD: The password to authenticate with the Redis server.
-   SECRET_ASSISTANT_URL: The URL to service for secret assistance.

## Scripts

The following scripts are available for managing and running the API Gateway service:

-   build: Builds the project, preparing it for deployment.
-   build:proto: Builds the protocol buffer files.
-   lint: Runs the linter to check for coding style and quality issues.
-   start: Starts the API Gateway service.
-   test: Runs the automated tests to ensure the functionality of the service.
