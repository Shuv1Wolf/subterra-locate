# Start with the scratch image
FROM scratch

# Set a working directory
# RUN mkdir /app
WORKDIR /app

# Copy the entire project
COPY ./exe ./bin
COPY ./config ./config

# Set environment variables
ENV GRPC_PORT=8080

# Expose the port
EXPOSE "8080:8080"

# Specify the command to run
CMD ["./bin/app"]
