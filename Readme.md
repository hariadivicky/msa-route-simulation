# Code Me - Route Simulator

![route result](/resource/result.gif "Route result animation")

This project simulates classic problem which to calculate meeting point between two riders on given routes.
Take point A and point B as route start and finish, one rider start from point A while other rider start from point B.
Find location of meeting point (coordinate) and time (second) required them to met.

original task stored in [task.md](task.md) file.

## Build Source

1. Clone this repository

    ```sh
    git clone https://github.com/hariadivicky/msa-route-simulation
    ```

2. Build frontend

    ```sh
    # go to ui folder
    cd ui
    # install NPM dependencies and build
    npm install && npm run build
    ```

3. Build backend

    ```sh
    # go to server folder
    cd server
    # build and run
    go build
    ```

4. Run

    ```sh
    # we assume you are in server/ directory
    ./server

    # server listening on localhost:8000
    ```

## Development Mode

Only for frontend

```sh
# go to ui directory
cd ui

# run watcher & hot reload
npm run serve
```