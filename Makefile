.PHONY: build install uninstall clean quotes-pratchett quotes-clean test fmt lint path-check source-path

# Application name
APP_NAME = personal-disorganiser

# Installation directory
INSTALL_DIR = /usr/local/bin

# Build target
build:
	go build -o $(APP_NAME) ./cmd

# Install to system PATH
install: build
	@echo "Installing $(APP_NAME) to $(INSTALL_DIR)..."
	@if [ -w $(INSTALL_DIR) ]; then \
		cp $(APP_NAME) $(INSTALL_DIR)/$(APP_NAME); \
	else \
		echo "Requires sudo permission to install to $(INSTALL_DIR)"; \
		sudo cp $(APP_NAME) $(INSTALL_DIR)/$(APP_NAME); \
	fi
	@echo "$(APP_NAME) installed to $(INSTALL_DIR)"
	@echo ""
	@echo "Checking PATH configuration..."
	@if echo $$PATH | grep -q "$(INSTALL_DIR)"; then \
		echo "✓ $(INSTALL_DIR) is already in your PATH"; \
	else \
		echo "⚠ $(INSTALL_DIR) is not in your PATH"; \
		echo "Adding $(INSTALL_DIR) to your shell configuration..."; \
		if [ "$$SHELL" = "/bin/zsh" ] || [ "$$SHELL" = "/usr/bin/zsh" ]; then \
			echo 'export PATH="$(INSTALL_DIR):$$PATH"' >> ~/.zshrc; \
			echo "Added to ~/.zshrc"; \
			echo "Sourcing ~/.zshrc to update current session..."; \
			export PATH="$(INSTALL_DIR):$$PATH"; \
		elif [ "$$SHELL" = "/bin/bash" ] || [ "$$SHELL" = "/usr/bin/bash" ]; then \
			echo 'export PATH="$(INSTALL_DIR):$$PATH"' >> ~/.bash_profile; \
			echo "Added to ~/.bash_profile"; \
			echo "Sourcing ~/.bash_profile to update current session..."; \
			export PATH="$(INSTALL_DIR):$$PATH"; \
		else \
			echo "Please add 'export PATH=\"$(INSTALL_DIR):$$PATH\"' to your shell's RC file"; \
			export PATH="$(INSTALL_DIR):$$PATH"; \
		fi; \
	fi
	@echo ""
	@echo "✓ Installation complete!"
	@echo "Testing command availability..."
	@if command -v $(APP_NAME) >/dev/null 2>&1; then \
		echo "✓ $(APP_NAME) is ready to use!"; \
		echo "Run: $(APP_NAME)"; \
	else \
		echo "⚠ $(APP_NAME) command not immediately available in this session"; \
		echo "Please restart your terminal or run:"; \
		if [ "$$SHELL" = "/bin/zsh" ] || [ "$$SHELL" = "/usr/bin/zsh" ]; then \
			echo "  source ~/.zshrc"; \
		elif [ "$$SHELL" = "/bin/bash" ] || [ "$$SHELL" = "/usr/bin/bash" ]; then \
			echo "  source ~/.bash_profile"; \
		else \
			echo "  source your shell's RC file"; \
		fi; \
		echo "Then run: $(APP_NAME)"; \
	fi

# Remove from system PATH
uninstall:
	@echo "Removing $(APP_NAME) from $(INSTALL_DIR)..."
	@if [ -w $(INSTALL_DIR) ]; then \
		rm -f $(INSTALL_DIR)/$(APP_NAME); \
	else \
		echo "Requires sudo permission to remove from $(INSTALL_DIR)"; \
		sudo rm -f $(INSTALL_DIR)/$(APP_NAME); \
	fi
	@echo "$(APP_NAME) removed from $(INSTALL_DIR)"

# Clean build artifacts
clean:
	rm -f $(APP_NAME)
	go clean

# Check and fix PATH configuration
path-check:
	@echo "Checking PATH configuration for $(INSTALL_DIR)..."
	@if echo $$PATH | grep -q "$(INSTALL_DIR)"; then \
		echo "✓ $(INSTALL_DIR) is already in your PATH"; \
	else \
		echo "⚠ $(INSTALL_DIR) is not in your PATH"; \
		echo "To add it, run one of these commands:"; \
		echo ""; \
		if [ "$$SHELL" = "/bin/zsh" ] || [ "$$SHELL" = "/usr/bin/zsh" ]; then \
			echo "  echo 'export PATH=\"$(INSTALL_DIR):\$$PATH\"' >> ~/.zshrc"; \
			echo "  source ~/.zshrc"; \
		elif [ "$$SHELL" = "/bin/bash" ] || [ "$$SHELL" = "/usr/bin/bash" ]; then \
			echo "  echo 'export PATH=\"$(INSTALL_DIR):\$$PATH\"' >> ~/.bash_profile"; \
			echo "  source ~/.bash_profile"; \
		else \
			echo "  Add 'export PATH=\"$(INSTALL_DIR):\$$PATH\"' to your shell's RC file"; \
		fi; \
		echo ""; \
		echo "Or run 'make install' to automatically configure it."; \
	fi

# Source shell configuration to update PATH in current session
source-path:
	@echo "To update PATH in your current terminal session, run:"
	@if [ "$$SHELL" = "/bin/zsh" ] || [ "$$SHELL" = "/usr/bin/zsh" ]; then \
		echo "  source ~/.zshrc"; \
	elif [ "$$SHELL" = "/bin/bash" ] || [ "$$SHELL" = "/usr/bin/bash" ]; then \
		echo "  source ~/.bash_profile"; \
	else \
		echo "  source your shell's RC file"; \
	fi
	@echo ""
	@echo "Or simply restart your terminal for the PATH to take effect."

# Download and parse Terry Pratchett quotes
quotes-pratchett:
	@echo "Creating quotes directory..."
	@mkdir -p ~/.config/personal-disorganizer/quotes
	@echo "Downloading Terry Pratchett quotes..."
	@curl -s https://www.lspace.org/ftp/words/pqf/pqf -o ~/.config/personal-disorganizer/quotes/pratchett.pqf
	@echo "Parsing quotes to JSON format..."
	@go run scripts/parse-pratchett.go
	@echo "Configuring Terry Pratchett quotes in user config..."
	@./scripts/configure-quotes.sh
	@echo "Terry Pratchett quotes installed and configured!"

# Clean quote files
quotes-clean:
	rm -rf ~/.config/personal-disorganizer/quotes
	@echo "Quote files removed"

# Run tests
test:
	go test ./...

# Format Go code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run