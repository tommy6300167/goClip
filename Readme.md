# SuperClip (超吉貼)

A powerful clipboard manager application built with Go and Fyne framework for macOS. SuperClip monitors your clipboard in real-time and stores both text and image content for easy retrieval and management.

## Features

### Core Functionality

- **Real-time Clipboard Monitoring**: Automatically detects and stores clipboard changes every 800ms(You can setting up your own config)
- **Text & Image Support**: Handles both text content and image files from clipboard
- **Persistent History**: Stores clipboard history with configurable maximum items (default: 30)
- **Dual-pane Interface**: List view for history browsing and detail view for content preview/editing
- **Search Functionality**: Quick search through clipboard history
- **Export Capability**: Export clipboard history data

### User Interface

- **Native macOS Design**: Built with Fyne framework for native look and feel
- **Chinese Font Support**: Embedded NotoSansTC font for proper Chinese character display
- **Responsive Layout**: Adaptive interface that works with different window sizes
- **Custom Theming**: Clean, modern interface optimized for productivity

## Usage Guide

### Keyboard Shortcuts

- **Cmd+S**: Save edited content
- **Cmd+F**: Quick search through clipboard history
- **Shift+Enter**: Insert line break when editing text
- **Cmd+E**: Quick edit mode (planned feature)

### Basic Operations

1. **Launch the application**: SuperClip will start monitoring your clipboard automatically
2. **View history**: All clipboard items appear in the left panel, sorted by recency
3. **Preview content**: Click any item to see full content in the right panel
4. **Edit text**: Double-click text items to edit them inline
5. **Search**: Use Cmd+F to quickly find specific clipboard content

### Data Storage

- History data: `~/Library/Application Support/ClipMini/`
- Image files: `~/Library/Application Support/ClipMini/images/`
- Configuration: Stored within application bundle

## Technical Documentation

### Architecture

SuperClip follows a clean MVC (Model-View-Controller) architecture:

#### Models (`models/`)

- **`clipboard.go`**: Core clipboard item data structures and types
- **`config.go`**: Application configuration with sensible defaults
- **`history.go`**: History management data structures

#### Views (`views/`)

- **`main_view.go`**: Main application window and layout management
- **`list_view.go`**: Clipboard history list display component
- **`detail_view.go`**: Content preview and editing panel
- **`toolbar.go`**: Application toolbar with action buttons
- **`theme.go`**: Custom theming with embedded NotoSansTC font

#### Controllers (`controllers/`)

- **`clipboard_controller.go`**: Central clipboard monitoring and management logic

#### Services (`services/`)

- **`clipboard_service.go`**: Low-level system clipboard interaction
- **`history_service.go`**: History persistence and management
- **`file_service.go`**: File operations for image storage

#### Utilities (`utils/`)

- **`text_utils.go`**: Text processing and formatting utilities
- **`time_utils.go`**: Time formatting and timezone handling

### Configuration

- **Polling Interval**: 800ms (configurable)
- **Max History Items**: 30 (configurable)
- **Max Display Length**: 50 characters for list preview
- **Window Size**: 800x600 pixels (resizable)

### Dependencies

- **Go**: 1.24.3+
- **Fyne**: v2.6.2 (Cross-platform GUI framework)
- **NotoSansTC**: Embedded font for Chinese character support

## Installation & Building

### Prerequisites

- Go 1.24.3 or later
- Fyne requirements for macOS development
- Xcode Command Line Tools

### Building for ARM64 (Apple Silicon)

```bash
GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 \
fyne package --os darwin \
--icon assets/clip.png \
--name "超吉貼" \
--id com.yourdomain.superclip \
--release
```

### Building for AMD64 (Intel Mac)

```bash
GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 \
fyne package --os darwin \
--icon assets/clip.png \
--name "超吉貼-intel" \
--id com.yourdomain.superclip \
--release
```

### Development Commands

```bash
# Run in development mode
go run main.go

# Build executable
go build

# Run tests
go test ./...

# Format code
go fmt ./...

# Lint code
go vet ./...

# Clean dependencies
go mod tidy
```

### Post-Installation (macOS Security)

If macOS prevents the app from running due to security restrictions:

```bash
xattr -d com.apple.quarantine ~/Downloads/超吉貼.app
```

## Project Structure

```
goClip/
├── main.go                 # Application entry point
├── assets/                 # Application resources
│   └── clip.png           # Application icon
├── controllers/            # Business logic controllers
├── models/                # Data models and structures
├── services/              # Core service implementations
├── utils/                 # Helper utilities
├── views/                 # UI components
│   ├── dialogs/          # Dialog windows
│   └── theme/            # Custom theming
└── CLAUDE.md             # Development guidelines
```

## Contributing

This project follows Go best practices and uses the Fyne framework for cross-platform GUI development. When contributing:

1. Follow the existing MVC architecture pattern
2. Maintain code formatting with `go fmt`
3. Ensure all tests pass with `go test ./...`
4. Use meaningful commit messages
5. Test on both ARM64 and AMD64 architectures when possible

## License

This project is developed as a personal clipboard management tool. Please refer to the license file for usage terms.
