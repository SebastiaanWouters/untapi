# Untappd User Data Fetcher

This Go application fetches user data from the Untappd API and stores it in a database. The data is exposed via an API endpoint.

## Table of Contents

- [Installation](#installation)
  - [Cloning the Repository](#cloning-the-repository)
  - [Setting Up Dependencies](#setting-up-dependencies)
- [Configuration](#configuration)
  - [Environment Variables](#environment-variables)
- [How to Run](#how-to-run)
  - [Initialization](#initialization)
  - [Serving the App](#serving-the-app)
- [Hosting with Caddy](#hosting-with-caddy)

## Installation

### Cloning the Repository

To clone the repository, open your terminal and run:

    git clone https://github.com/Sebastiaan-Wouters/untapi.git

This will clone the project into a directory named `untapi`.

### Setting Up Dependencies

Navigate to the project directory:

    cd untapi

Install the Go dependencies:

    go mod download

Or you can build the project, and Go will automatically download the dependencies:

    go build

## Configuration

### Environment Variables

The application requires a `.env` file for configuration. Create a `.env` file in the root directory of the project and add the following variables:

- `ACCESS_TOKEN`: Your Untappd API access token.
- `USERS`: A comma-separated list of usernames to fetch from the Untappd API.

Example `.env` file:

    ACCESS_TOKEN=your_access_token_here
    USERS=username1,username2,username3

## How to Run

### Initialization

To initialize the application, run:

    go run main.go init

This command will create the required database collections and populate them with initial user data.

### Serving the App

After initialization, run:

    go run main.go serve

This starts the application. User data will be fetched and updated every 5 minutes (as per the cron job settings).

## Hosting with Caddy

To host the application using Caddy with HTTPS, you'll need to install Caddy and create a `Caddyfile` for your configuration.

Sample `Caddyfile`:

    your-domain.com {
      reverse_proxy localhost:PORT
    }

Replace `your-domain.com` with your actual domain and `your-email@example.com` with your email. Update `PORT` to the port where your application is running.

To run Caddy, execute:

    caddy run

This will start the Caddy server and your application will be available at `https://your-domain.com`. Make sure to keep the go program running in the background while executing the caddy server (using go run main.go serve & or other background services). 