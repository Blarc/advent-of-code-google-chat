FROM golang:1.19

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/github.com/Blarc/aoc-bot

# Copy everything from the current directory to the PWD (Present Working Directory) inside the container
COPY . .

# Build the package
RUN go build -o aoc .

# Run the executable
CMD ["./aoc"]     