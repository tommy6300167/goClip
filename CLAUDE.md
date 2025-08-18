# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is "超吉貼" (SuperClip), a Go-based clipboard manager application built with Fyne GUI framework. The application monitors clipboard changes in real-time and stores both text and image clipboard items for later retrieval. It features a dual-pane interface with a list view for history and a detail view for content preview/editing.

## Development Commands

### Building the Application
```bash
# Build for ARM64 (Apple Silicon)
GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 fyne package --os darwin --icon assets/clip.png --name "超吉貼" --id com.yourdomain.superclip --release

# Build for AMD64 (Intel)
GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 fyne package --os darwin --icon assets/clip.png --name "超吉貼-intel" --id com.yourdomain.superclip --release
```

### Running the Application
```bash
go run main.go
```

### Standard Go Commands
```bash
go build          # Build the application
go test ./...     # Run tests (if any exist)
go mod tidy       # Clean up dependencies
go fmt ./...      # Format code
go vet ./...      # Lint code
```

### Remove Quarantine (macOS)
```bash
xattr -d com.apple.quarantine ~/Downloads/超吉貼.app
```

## Architecture

### MVC Pattern
The application follows a Model-View-Controller pattern:

- **Models** (`models/`): Core data structures and configuration
  - `clipboard.go` - ClipboardItem definitions and types
  - `config.go` - Application configuration with defaults
  - `history.go` - History management data structures

- **Views** (`views/`): UI components using Fyne
  - `main_view.go` - Main application window and layout
  - `list_view.go` - History items list display
  - `detail_view.go` - Content preview and editing pane
  - `toolbar.go` - Application toolbar with actions
  - `theme.go` - Custom theming with NotoSansTC font

- **Controllers** (`controllers/`): Business logic coordination
  - `clipboard_controller.go` - Main clipboard monitoring and management

### Service Layer
- **Services** (`services/`): Core functionality implementation
  - `clipboard_service.go` - System clipboard interaction
  - `history_service.go` - History persistence and management
  - `file_service.go` - File operations for image storage

### Utilities
- **Utils** (`utils/`): Helper functions
  - `text_utils.go` - Text processing utilities
  - `time_utils.go` - Time formatting and timezone handling

### Key Features
- Real-time clipboard monitoring (800ms polling interval)
- Supports both text and image clipboard content
- Persistent history storage in `~/Library/Application Support/ClipMini/`
- Image files saved to dedicated images directory
- Maximum 30 items in history (configurable)
- Export functionality for history data
- Edit capability for text items
- Dual-language support (English/Chinese)

### Configuration
- Default polling interval: 800ms
- Max history items: 30
- Max display length: 50 characters
- Data storage: `~/Library/Application Support/ClipMini/`
- Image storage: `~/Library/Application Support/ClipMini/images/`

### Dependencies
- Fyne v2.6.2 - GUI framework
- Go 1.24.3 - Runtime
- NotoSansTC font - Embedded for Chinese character support