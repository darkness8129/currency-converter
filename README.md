## Commands

- `docker-compose up` - to start the service using docker-compose on 8080 port

## Possible Improvements (imho)

1. Currently, the email sending functionality works in such a way that the timer resets every time the server restarts. This can be problematic because if the service is deployed, the logic for sending emails once every 24 hours will break when new functionality is deployed and the server stops. This issue can be resolved by using a separate endpoint and a job scheduler (cron) that will call this endpoint at specific intervals.

2. If the server stops working at the moment emails are being sent, it will result in some subscribers not receiving their emails. This issue can be resolved by queuing all emails when the cron job makes a request to the endpoint. Processing emails from the queue will ensure that no email is missed.

3. The services for fetching currency rates and managing subscriptions are currently on the same HTTP server. It would be better to have a separate `currency-rate` microservice (with both HTTP and gRPC controllers for fetching rates) and a separate `subscriptions` microservice (with its own database, HTTP controller for managing subscriptions, and logic for sending emails every 24 hours). HTTP controllers would be used for external communication with clients, while for internal communication, the `subscriptions` service would call the `currency-rate` service via gRPC to fetch the currency rates and send them to users.
