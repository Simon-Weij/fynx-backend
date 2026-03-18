FROM golang:1.26.0-alpine3.23 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY src ./src
RUN CGO_ENABLED=0 go build -o /app/wayland-recorder-backend ./src

FROM gcr.io/distroless/static-debian13

USER nonroot

WORKDIR /app

COPY --from=build /app/wayland-recorder-backend /app/wayland-recorder-backend

CMD [ "/app/wayland-recorder-backend" ]