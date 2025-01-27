package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type FileBundle struct {
	XMLName  xml.Name    `xml:"filebundle"`
	Version  string      `xml:"version,attr"`
	Metadata Metadata    `xml:"metadata"`
	Files    []FileEntry `xml:"files>file"`
}

type Metadata struct {
	Created    time.Time `xml:"created"`
	SourcePath string    `xml:"source_path"`
}

type FileEntry struct {
	Path     string    `xml:"path"`
	Size     int64     `xml:"size"`
	Modified time.Time `xml:"modified"`
	RelPath  string    `xml:"relative_path"` // Path relative to source directory
	Content  string    `xml:"content"`       // Base64 encoded content
}

func (a *FileBundlerApp) GenerateXMLToWriter(writer io.Writer) error {
	// Create a buffer for our custom formatting
	var buffer bytes.Buffer

	// Write XML header
	buffer.WriteString("<FileBundle>\n")

	// Write metadata
	buffer.WriteString("    <metadata>\n")
	buffer.WriteString(fmt.Sprintf("        <root_directory>%s</root_directory>\n", a.currentPath))
	buffer.WriteString("    </metadata>\n")

	// Write files
	buffer.WriteString("    <Files>\n")

	// Add selected files
	for _, path := range a.selection.GetSelectedPaths() {
		info, err := os.Stat(path)
		if err != nil {
			continue
		}

		// Read file content
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Calculate relative path
		relPath, err := filepath.Rel(a.currentPath, path)
		if err != nil {
			relPath = path
		}

		// Write file entry
		buffer.WriteString("        <File>\n")
		buffer.WriteString(fmt.Sprintf("            <Path>%s</Path>\n", path))
		buffer.WriteString(fmt.Sprintf("            <Size>%d</Size>\n", info.Size()))
		buffer.WriteString(fmt.Sprintf("            <Modified>%s</Modified>\n", info.ModTime().Format(time.RFC3339)))
		buffer.WriteString(fmt.Sprintf("            <Relative_Path>%s</Relative_Path>\n", relPath))
		buffer.WriteString("            <Content>\n")
		buffer.Write(content)
		buffer.WriteString("\n            </Content>\n")
		buffer.WriteString("        </File>\n")
	}

	buffer.WriteString("    </Files>\n")
	buffer.WriteString("</FileBundle>")

	// Write the final buffer to the writer
	_, err := writer.Write(buffer.Bytes())
	return err
}
