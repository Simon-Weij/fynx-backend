FROM golang:1.26.0-alpine3.23 AS build

RUN apk add --no-cache just

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY justfile ./
COPY src ./src
RUN just build /app/wayland-recorder-backend

FROM gcr.io/distroless/static-debian13

USER nonroot

WORKDIR /app

COPY --from=build /app/wayland-recorder-backend /app/wayland-recorder-backend

CMD [ "/app/wayland-recorder-backend" ]