package browsing

import "errors"

// History maintains a double-linked list of pages visited.
type History struct {
	currentNode *historyNode
}
type historyNode struct {
	Url      string
	Next     *historyNode
	Previous *historyNode
}

// NewHistory creates a new History.
func NewHistory() History {
	return History{
		currentNode: nil,
	}
}

// GetCurrentUrl gets the url of the current page.
func (h *History) GetCurrentUrl() (string, error) {
	if h.currentNode == nil {
		return "", errors.New("history is empty")
	}
	return h.currentNode.Url, nil
}

// CanGoBack checks if there is a previous page.
func (h *History) CanGoBack() bool {
	return h.currentNode != nil && h.currentNode.Previous != nil
}

// CanGoForward checks if there is a next page.
func (h *History) CanGoForward() bool {
	return h.currentNode != nil && h.currentNode.Next != nil
}

// GoBack moves the current page to the previous page.
func (h *History) GoBack() error {
	if !h.CanGoBack() {
		return errors.New("no previous page")
	}
	h.currentNode = h.currentNode.Previous
	return nil
}

// GoForward moves the current page to the next page.
func (h *History) GoForward() error {
	if !h.CanGoForward() {
		return errors.New("no next page")
	}
	h.currentNode = h.currentNode.Next
	return nil
}

// Push pushes and moves the current page to the provided url.
// This overwrites any history in front of the current page.
func (h *History) Push(url string) {
	// TODO: clean up list in front (might not be necessary)
	newNode := &historyNode{
		Url:      url,
		Next:     nil,
		Previous: h.currentNode,
	}
	if h.currentNode != nil {
		h.currentNode.Next = newNode
		h.currentNode = h.currentNode.Next
	} else {
		h.currentNode = newNode
	}
}
