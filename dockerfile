# Use a minimal base image for the final executable
FROM alpine:latest

# Set the working directory for the final image
WORKDIR /root/

# Copy the compiled binary from the local directory to the container
COPY cortexForge .

# Specify the command to run the executable, GOODLUCK
CMD ["./cortexForge"]
