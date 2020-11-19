# Code Me - Route Simulator ###

## finish incomplete project
This project simulates classic problem which to calculate meeting point between two riders on given routes.
Take point A and point B as route start and finish, one rider start from point A while other rider start from point B.
Find location of meeting point (coordinate) and time (second) required them to met.

### task 
1. Implement route reader interface and must fulfill unit test.
2. Implement API endpoint for server route data.
3. Implement frontend codes to simulate route path and rider traveling path.
4. Calculate and display information of meeting point and time. 

## goal ilustration
![alt text](/resource/route-animation.gif "Route animation")

## how to build the project
### Frontend app
install `node`, `npm` and `yarn`

go to `ui` folder

run `yarn build`

### Backend / API 
go to `server` folder 

build using `go build` 

start project by using `./server` use available options if required

### More chalenge.
webpack custom configuration to ease development experience, eg: using `webpack --watch`

configurable "speed" which can represented by delay between route point

use combination of websocket and goroutine workers to publish updated location of rider to frontend application. 