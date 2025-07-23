package notify

import (
	"os"
	"testing"
)

func TestAdd(t *testing.T) {
	n, err := NewNotify()
	if err != nil {
		t.Fatalf("Failed to create Notify instance: %v", err)
	}

	testPath := "../../test/notify"
	err = os.MkdirAll(testPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	err = n.Add(testPath)
	if err != nil {
		t.Errorf("Failed to add path %s: %v", testPath, err)
		return
	}

	go func() {
		for {
			select {
			case event, ok := <-n.Events:
				if !ok {
					return
				}
				t.Logf("Event: %s for path: %s", event.Op, event.Name)
			case err, ok := <-n.Errors:
				if !ok {
					return
				}
				t.Errorf("Error: %v", err)
			}
		}
	}()
	<-make(chan struct{})
}
