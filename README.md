# fb - File Bundler

A terminal-based application for bundling files into an XML format and copying them to the clipboard for use with LLMs.

## Features
- **File Selection**: Select files using spacebar or click in tree view.
- **Tree View**: Browse directory structure with a collapsible tree.
- **LLM-Friendly Export**: Generate bundles in XML.
- **Clipboard Support**: Copy generated bundle directly to clipboard.

## Installation
```bash
go install github.com/stalkaltan/fb@latest
```

If after installation, the application is still not found, make sure the `$GOPATH/bin` directory is in your `$PATH` environment variable.

## Usage

### Basic Commands
1. Launch the application:
   ```bash
   fb
   ```
2. Navigate through files using arrow keys or `h/j/k/l`.
3. Select files by pressing **Space**.
4. Bundle selected files to clipboard with **x**.

### Key Bindings
- **Up/Down**: Move selection (`k`/`j`)
- **Left/Right**: Collapse/Expand directory (`h`/`l`)
- **Enter**: Expand/Collapse directory
- **Space**: Select/Deselect file
- **Tab**: Switch between panels
- **x**: Bundle files to clipboard
- **?**: Toggle Help panel
- **q**: Quit application

## License
MIT License

## Contribution
Contributions are welcome! Fork the repository, create a feature branch, and submit a pull request.