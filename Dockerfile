FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/metrovanlibsearch .

FROM gcr.io/distroless/static-debian13:nonroot
WORKDIR /app
COPY --from=build /out/metrovanlibsearch /app/metrovanlibsearch
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app/metrovanlibsearch"]
CMD ["serve", "--addr", ":8080"]
