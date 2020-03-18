# Base image for GolangCI lint
FROM golangci/golangci-lint:v1.24.0

# Set working directory. Source code will be copied to this directory.
WORKDIR ${GOPATH}/github.com/rodsher/sqlstat

# Copy source code
COPY . .

# Run GolangCI lint
RUN ["golangci-lint", "run", "-v", "--timeout", "10m", "--skip-dirs-use-default", "--enable-all"]
