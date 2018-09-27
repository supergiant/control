package node

import (
	"strings"
	"testing"
)

func TestNode_String(t *testing.T) {
	name := "node-name"
	id := "id"
	size := "node-size"

	n := &Node{
		Id:   id,
		Name: name,
		Size: size,
	}

	s := n.String()

	if !strings.Contains(s, name) {
		t.Errorf("name %s not found in %s", name, n.String())
	}

	if !strings.Contains(s, size) {
		t.Errorf("size %s not found in %s", size, n.String())
	}

	if !strings.Contains(s, id) {
		t.Errorf("id %s not found in %s", id, n.String())
	}
}
