# Golang Application Setup Guide

## Prerequisites

Before starting, make sure you have the following installed:
- [Git](https://git-scm.com/downloads)
- [Go](https://golang.org/dl/) (version 1.16 or higher)
- [PostgreSQL](https://www.postgresql.org/download/) 

## Getting Started

Follow these steps to set up and run the application:

1. Clone the repository:
   ```
   git clone https://github.com/Arkavidia-Academy-SofEng/server.git
   cd your-repo-name
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

3. Set up the database:
   ```
   make migrate-up
   ```

4. Run the application:
   ```
   make run
   ```

## Available Make Commands

- `make migrate-up`: Applies all database migrations
- `make migrate-down`: Reverts the last database migration
- `make run`: Starts the application
- `make build`: Builds the application binary

## Configuration

1. Copy the example environment file:
   ```
   cp .env.example .env
   ```

2. Edit the `.env` file with your configuration settings

## API Documentation

The API documentation is available at `https://documenter.getpostman.com/view/32354585/2sAYkBu2n4#e621a063-f6e5-4186-ab2f-74b71aae4913` when the application is running.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
