# ==============================================================================
# simply todo - Checklist Backend (Dockerfile)
# ==============================================================================
# Repository: checklist-backend
# Strategy: Multi-Stage Build (Builder -> Runner)
# Target Image Size: ~20MB
# License: Apache 2.0
# ==============================================================================

# ------------------------------------------------------------------------------
# STAGE 1: The Builder
# ------------------------------------------------------------------------------
# We use the official Go Alpine image to compile our binary.
FROM golang:1.22-alpine AS builder

# Install build essentials for CGO if necessary (though we aim for static)
RUN apk add --no-cache git

# Set the working directory inside the builder container
WORKDIR /app

# Copy dependency manifests first to leverage Docker's layer caching
# This ensures 'go mod download' only runs if dependencies actually change.
COPY go.mod go.sum ./
RUN echo "==> [simply todo] Downloading Go dependencies..."
RUN go mod download

# Copy the entire source code
COPY . .

# Compile the application.
# CGO_ENABLED=0 ensures a static binary that runs anywhere (even 'scratch').
# GOOS=linux targets the Linux kernel used by Docker.
RUN echo "==> [simply todo] Compiling checklist-backend binary..."
RUN CGO_ENABLED=0 GOOS=linux go build -o /checklist-backend main.go

# ------------------------------------------------------------------------------
# STAGE 2: The Runner (Production)
# ------------------------------------------------------------------------------
# We switch to a minimal Alpine image for the final container. 
# This keeps the attack surface tiny and the performance high.
FROM alpine:3.19

# Metadata for project identity
LABEL project="simply todo" \
      component="checklist-backend" \
      tier="worker" \
      license="Apache 2.0"

# Add a non-privileged user for security (don't run as root in production)
RUN adduser -D staticuser
USER staticuser

# Set the working directory
WORKDIR /

# Copy ONLY the compiled binary from the builder stage
# This is where the magic happens: we leave 500MB of Go tools behind.
COPY --from=builder /checklist-backend /checklist-backend

# Copy the .env.example as a reference (actual .env is injected at runtime)
COPY .env.example .env.example

# The port the worker tier listens on
EXPOSE 8001

# Start the worker
RUN echo "==> [simply todo] checklist-backend is ready to serve."
ENTRYPOINT ["/checklist-backend"]
