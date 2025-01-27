package main

import (
	"bytes"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"os"
)

type FileBundlerApp struct {
	app         *tview.Application
	ui          *UIComponents
	config      *Config
	selection   *SelectionManager
	currentPath string
}

type UIComponents struct {
	fileTree      *tview.TreeView
	statusBar     *tview.TextView
	searchBox     *tview.InputField
	layout        *tview.Flex
	selectionList *tview.List // New component for selected files
	mainPane      *tview.Flex // New container for main content,
	hintBar       *tview.TextView
}

func NewFileBundlerApp() *FileBundlerApp {
	app := &FileBundlerApp{
		app:       tview.NewApplication(),
		config:    NewDefaultConfig(),
		selection: NewSelectionManager(),
	}

	app.initUI()
	app.setupKeyBindings()
	return app
}

func (a *FileBundlerApp) initUI() {
	// Initialize UI components
	a.ui = &UIComponents{
		fileTree:      tview.NewTreeView(),
		statusBar:     tview.NewTextView(),
		searchBox:     tview.NewInputField(),
		layout:        tview.NewFlex(),
		selectionList: tview.NewList(),
		mainPane:      tview.NewFlex(),
		hintBar:       tview.NewTextView(),
	}

	// Configure file tree
	a.ui.fileTree.SetBorder(true).SetTitle("File Tree")
	a.ui.fileTree.SetGraphics(true)
	a.ui.fileTree.SetTitleColor(tcell.ColorBlue.TrueColor())

	// Configure selection list
	a.ui.selectionList.SetBorder(true).SetTitle("Selected Files")
	a.ui.selectionList.ShowSecondaryText(false)
	a.ui.selectionList.SetTitleColor(tcell.ColorGreen.TrueColor())

	// Configure status bar
	a.ui.statusBar.SetBorder(true).SetTitle("Status")
	a.ui.statusBar.SetText("Use arrow keys or h/j/k/l to navigate. Space to select. Enter to expand/collapse.")
	a.ui.statusBar.SetTextColor(tcell.ColorGray.TrueColor())

	// Create main pane with file tree and selection list
	a.ui.mainPane.SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().
			AddItem(a.ui.fileTree, 0, 2, true).       // File tree takes 2/3 of width
			AddItem(a.ui.selectionList, 0, 1, false), // Selection list takes 1/3 of width
			0, 1, true)

	a.ui.hintBar.SetTextAlign(tview.AlignCenter)
	a.ui.hintBar.SetText("Press ? for help")
	a.ui.hintBar.SetTextColor(tcell.ColorYellow)

	a.ui.layout.SetDirection(tview.FlexRow).
		AddItem(a.ui.mainPane, 0, 1, true).
		AddItem(a.ui.statusBar, 1, 0, false).
		AddItem(a.ui.hintBar, 1, 0, false)

	a.app.SetRoot(a.ui.layout, true)
	a.app.SetFocus(a.ui.fileTree)
}

// Add method to update selection list
func (a *FileBundlerApp) updateSelectionList() {
	a.ui.selectionList.Clear()
	for _, path := range a.selection.GetSelectedPaths() {
		a.ui.selectionList.AddItem(path, "", 0, nil)
	}
}

func (a *FileBundlerApp) showBundleConfirmation() {
	selectedPaths := a.selection.GetSelectedPaths()
	numFiles := len(selectedPaths)

	// Calculate total size
	var totalSize int64
	for _, path := range selectedPaths {
		info, err := os.Stat(path)
		if err != nil {
			continue
		}
		totalSize += info.Size()
	}

	// Format size for display
	var sizeStr string
	switch {
	case totalSize < 1024:
		sizeStr = fmt.Sprintf("%d B", totalSize)
	case totalSize < 1024*1024:
		sizeStr = fmt.Sprintf("%.1f KB", float64(totalSize)/1024)
	default:
		sizeStr = fmt.Sprintf("%.1f MB", float64(totalSize)/(1024*1024))
	}

	// Show warning for large files
	var message string
	if totalSize > 10*1024*1024 { // Warning for files over 10MB
		message = fmt.Sprintf("Warning: Large selection!\n\nBundle %d selected file(s)?\nTotal size: %s\n\nThis may take a moment and use significant memory.",
			numFiles, sizeStr)
	} else {
		message = fmt.Sprintf("Bundle %d selected file(s)?\nTotal size: %s",
			numFiles, sizeStr)
	}

	doneFunc := func(buttonIndex int, buttonLabel string) {
		if buttonIndex == 0 { // Bundle
			// Close confirmation modal
			a.app.SetRoot(a.ui.layout, true)

			if totalSize > 1024*1024 { // 1MB
				a.showProcessingMessage()
			}
			a.generateAndCopyToClipboard()
		} else {
			a.app.SetRoot(a.ui.layout, true)
			a.app.SetFocus(a.ui.fileTree)
		}
	}

	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"Bundle", "Cancel"}).
		SetDoneFunc(doneFunc)

	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			doneFunc(0, "Bundle")
			return nil
		}
		if event.Key() == tcell.KeyEscape {
			doneFunc(1, "Cancel")
			return nil
		}
		// Pass through all other keys
		return event
	})

	a.app.SetRoot(modal, true)
}

func (a *FileBundlerApp) showProcessingMessage() {
	modal := tview.NewModal().
		SetText("Processing files...\nThis may take a moment.").
		SetBackgroundColor(tcell.ColorBlue)

	a.app.SetRoot(modal, true)
	a.app.Draw()
}

func (a *FileBundlerApp) generateAndCopyToClipboard() {
	// Create a buffer to store XML
	var buffer bytes.Buffer

	if err := a.GenerateXMLToWriter(&buffer); err != nil {
		a.showMessage("Error generating XML: "+err.Error(), true)
		return
	}

	// Copy to clipboard
	if err := clipboard.WriteAll(buffer.String()); err != nil {
		a.showMessage("Error copying to clipboard: "+err.Error(), true)
		return
	}

	// Show success message
	a.showMessage(fmt.Sprintf("Successfully bundled %d file(s) to clipboard!",
		len(a.selection.GetSelectedPaths())), false)
}

func (a *FileBundlerApp) showMessage(message string, isError bool) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"})

	returnToMain := func(buttonIndex int, buttonLabel string) {
		a.app.SetRoot(a.ui.layout, true)
		a.app.SetFocus(a.ui.fileTree)
	}

	modal.SetDoneFunc(returnToMain)

	// Add input capture with logging
	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return event
	})

	if isError {
		modal.SetBackgroundColor(tcell.ColorRed)
	} else {
		modal.SetBackgroundColor(tcell.ColorGreen)
	}

	a.app.SetRoot(modal, true)
}

func (a *FileBundlerApp) toggleHelpPanel() {
	modal := tview.NewModal().
		SetText(`
Navigation Controls:
──────────────────
↑/k                         Move up
↓/j                       Move down
←/h                   Collapse node
→/l                     Expand node
Space          Select/Deselect file
Enter     Expand/Collapse directory
Tab           Switch between panels
x         Bundle files to clipboard
?                 Toggle Help panel
q                  Quit application

 Press Escape to close this window
`)
	// Add comprehensive input capture
	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			a.app.SetRoot(a.ui.layout, true)
			a.app.SetFocus(a.ui.fileTree)
			return nil
		case tcell.KeyEnter:
			a.app.SetRoot(a.ui.layout, true)
			a.app.SetFocus(a.ui.fileTree)
			return nil
		case tcell.KeyTab:
			// Prevent tab from changing focus
			return nil
		}

		switch event.Rune() {
		case '?':
			a.app.SetRoot(a.ui.layout, true)
			a.app.SetFocus(a.ui.fileTree)
			return nil
		case ' ':
			a.app.SetRoot(a.ui.layout, true)
			a.app.SetFocus(a.ui.fileTree)
			return nil
		}

		return event
	})

	// Show the modal and focus on it
	a.app.SetRoot(modal, true)
	a.app.SetFocus(modal)
}

func (a *FileBundlerApp) Run() error {
	// Load initial directory
	if err := a.loadDirectory(a.config.DefaultPath); err != nil {
		return err
	}

	return a.app.Run()
}
