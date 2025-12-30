# GoReleaser Dockerfile - uses pre-built binary
# Using base-debian12 which includes CA certificates needed for HTTPS
FROM gcr.io/distroless/base-debian12:nonroot

# Copy the pre-built binary from GoReleaser
COPY cw-agent /cw-agent

# Set user
USER nonroot:nonroot

# Set entrypoint
ENTRYPOINT ["/cw-agent"]

# Default command
CMD ["start", "-c", "/etc/certwatch/certwatch.yaml"]
