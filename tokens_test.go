package css

import (
	"testing"

	"github.com/tinywasm/fmt"
)

func TestClassAsAttr(t *testing.T) {
	cls := Class("btn-primary")
	attr := cls.AsAttr()

	if attr.Key != "class" {
		t.Errorf("Expected key='class', got '%s'", attr.Key)
	}
	if attr.Value != "btn-primary" {
		t.Errorf("Expected value='btn-primary', got '%s'", attr.Value)
	}
}

func TestClassAsAttrMultiple(t *testing.T) {
	cls := Class("btn btn-primary")
	attr := cls.AsAttr()

	if attr.Key != "class" {
		t.Errorf("Expected key='class', got '%s'", attr.Key)
	}
	if attr.Value != "btn btn-primary" {
		t.Errorf("Expected value='btn btn-primary', got '%s'", attr.Value)
	}
}

func TestClassAsAttrType(t *testing.T) {
	cls := Class("my-class")
	attr := cls.AsAttr()

	// Verify it returns a KeyValue
	if _, ok := interface{}(attr).(fmt.KeyValue); !ok {
		t.Error("AsAttr() should return a fmt.KeyValue")
	}
}
