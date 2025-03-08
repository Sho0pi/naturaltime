# Build the JavaScript file
build:
	@echo "Installing dependencies..."
	@npm ci  # Ensures consistent installs based on package-lock.json
	@mkdir -p dist
	@echo "Building JavaScript file..."
	@npx browserify naturaltime.js --standalone naturaltime > dist/naturaltime.out.js

# Package the generated JavaScript file for releases
package: build
	@tar -czvf naturaltime-js.tar.gz dist/naturaltime.out.js

# Clean build artifacts
clean:
	@rm -rf dist naturaltime-js.tar.gz node_modules
