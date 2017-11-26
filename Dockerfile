FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
