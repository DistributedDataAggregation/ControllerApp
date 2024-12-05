FROM ubuntu:latest

ENV GO_VERSION=1.23.2
ENV GOROOT=/usr/local/go
ENV GOPATH=/go
ENV PATH=$GOROOT/bin:$GOPATH/bin:$PATH

# Install dependencies
RUN apt-get update && apt-get install -y \
    curl \
    git \
    build-essential \
    && rm -rf /var/lib/apt/lists/*

# Download and install Go
RUN curl -OL https://go.dev/dl/go$GO_VERSION.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go$GO_VERSION.linux-amd64.tar.gz && \
    rm go$GO_VERSION.linux-amd64.tar.gz

# Create the application directory
WORKDIR /app

# Copy application files to the container
COPY . .

EXPOSE 3000

# Build the Go application
RUN go build -o app .

# Set the entry point to the application
ENTRYPOINT ["./app"]
