FROM golang:1.14.9-alpine
USER root
RUN mkdir /app
COPY . /app
WORKDIR /app/
CMD ["go", "run", "backend/AuthService/service.go"]
