################ Build & Dev ################
# Build stage will be used:
# - for building the application for production
FROM golang:1.20.2-alpine3.17 as build

# Create project directory (workdir)
WORKDIR /app

# Copy source code files to WORKDIR
COPY . .

# Build application
RUN go get ./... && go mod tidy && go build -ldflags '-s -w' -o main .

# Container start command for development
CMD ["go", "run", "main.go"]


################ Production ################
# Creates a minimal image for production using distroless base image
# More info here: https://github.com/GoogleContainerTools/distroless
FROM gcr.io/distroless/base-debian11:nonroot as production

# Copy application binary from build/dev stage to the distroless container
COPY --from=build /app/main /


# Application port (optional)
EXPOSE 8000




# Container start command for production
CMD ["/main"]
