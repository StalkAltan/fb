package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"os"
	"path/filepath"
)

func (a *FileBundlerApp) loadDirectory(path string) error {
	root := tview.NewTreeNode(filepath.Base(path))
	root.SetReference(path)
	root.SetSelectable(true)
	root.SetColor(tcell.ColorBlue.TrueColor())

	// Populate the root node immediately
	err := a.populateNode(root)
	if err != nil {
		return err
	}

	a.ui.fileTree.SetRoot(root)
	a.ui.fileTree.SetCurrentNode(root)
	a.currentPath = path
	return nil
}

func (a *FileBundlerApp) populateNode(node *tview.TreeNode) error {
	path := node.GetReference().(string)

	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if shouldSkip(entry.Name(), a.config) {
			continue
		}

		childPath := filepath.Join(path, entry.Name())
		child := tview.NewTreeNode(entry.Name())
		child.SetReference(childPath)
		child.SetSelectable(true)

		if entry.IsDir() {
			child.SetColor(tcell.ColorBlue.TrueColor())
			child.SetText("📁 " + entry.Name())
		} else {
			child.SetText("📄 " + entry.Name())
			child.SetColor(tcell.ColorWhite.TrueColor())
		}

		node.AddChild(child)
	}

	return nil
}

// Move shouldSkip to be a standalone function
func shouldSkip(name string, config *Config) bool {
	if !config.IncludeHidden && len(name) > 0 && name[0] == '.' {
		return true
	}

	for _, pattern := range config.ExcludePatterns {
		if pattern == name {
			return true
		}
	}

	return false
}
