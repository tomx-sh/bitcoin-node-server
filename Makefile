# Default build command for your local development machine (Apple Silicon)
build-local:
	@mkdir -p dist/local
	GOARCH=arm64 GOOS=darwin go build -o dist/local/bitcoin-node-server ./app

# Build command for Intel-based macOS (target machine)
build-target:
	@mkdir -p dist/target
	GOARCH=amd64 GOOS=darwin go build -o dist/target/bitcoin-node-server ./app

# Clean all build artifacts
clean:
	rm -rf dist/*