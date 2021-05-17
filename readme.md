# Blog Application - Using Golang-GRPC, Progessive webapp and Mongo Atlas.

Here is an attempt to create a blog application using mentioned tech-stack.
For now only the authentication portion is complete.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

Docker is only app needs to available in your system [to run the application].

Inside docker, below primary dependencies will automatically get installed:

```
npm:
    parcel-bundler
golang/alpine/apk:
    protoc
    grpc
    github.com/improbable-eng/grpc-web/go/grpcweb
    golang.org/x/net/http2/h2c
    go.mongodb.org/mongo-driver/mongo
```

### Installing

To run:

1. Clone the repo:

```
git clone https://github.com/HiteshRepo/blog-application.git
```

2. cd to blog-application directory.
3. docker image for grpc backend:

```
docker build -t blog-grpc-app-backend.
```

4. docker container for grpc backend image:

```
docker run -d -p 9001:9001 blog-grpc-app-backend
```

5. cd to 'frontend' sub-directory.
6. docker image for grpc frontend:

```
docker build -t blog-grpc-app-frontend.
```

7. docker container for grpc backend image:

```
docker run -d -p 1234:1234 -p 7001:7001 blog-grpc-app-frontend.
```

8. After the containers have successfully started: go to http://localhost:1234
9. Try Signup, Login and Logout actions.

## Running the tests

17th May, 2021
I have added test cases only for service.go with coverage of 73.8%

## Developments in line

1. Configure prometheus, zap.logger, kibana.
2. Have view user-based posts.
3. Also comment on posts.

## Acknowledgments

- Hat tip to anyone whose code was used
- Inspiration
- etc
