# P8ion_backend
This is the backend of P8ion. It is written in Golang and uses Gin framework for HTTP Serving.

## Prerequisites
  - Go

## Setup
  1. Clone Repository
  2. Run **go mod download** to download dependencies
  3. Copy the **.env.example** file to **.env**
  4. Fill **APP_ENV** with **DEV**, **DOCKER**, or **PROD** depending on your environment. If you are running the backend locally, use **DEV**. If you are running the backend in a docker container, use **DOCKER** and fill in other values as well. If you are running the backend in production, use **PROD**. 
  5. Copy the **config.example.json** file to **config.example.json** and fill in values according to your environment.
  6. Have a MySql db running locally and create a database with the name specified in .env

## Run
  1. Run **make run** command to start the backend in dev mode.

## Rules 
  1. Commit messages accordting to the standard as specified here. http://karma-runner.github.io/6.4/dev/git-commit-msg.html
  2. Format code before commiting (use Prettier).
