# Goji Starter Overview

The Goji Starter demonstrates a simple, reusable Goji web application.

## Run the app locally

1. [Install Go][]
2. Clone repo
3. Setup PostgreSQL(e.g. with [Postgres.app](http://postgresapp.com))
3. cd into the app directory
4. Run `go run *.go`
5. Access the running app in a browser at http://localhost:8080

[Install Go]: https://golang.org/doc/install

## Run the app on Bluemix

1. Create a web app in Bluemix UI
2. Clone repo
3. cd into the app direcotory
4. Connect to Bluemix `bluemix api https://api.ng.bluemix.net`
5. Login to Bluemix `bluemix login ...`
6. Create PostgreSQL service `cf create-service postgresql 100 postgresql01`
7. Bind service to your app `cf bind-service YOUR_APP_NAME postgresql01`
8. Push the app `cf push YOUR_APP_NAME`

