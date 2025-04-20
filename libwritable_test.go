package frizzante

import "testing"

func TestWritableCreate(t *testing.T) {
	s := WritableCreate("default")

	if "default" != s.value {
		t.Fatal("findSubscriber value is not default")
	}

	s = WritableCreate("default")
	WritableWrite(s, "updated")

	if "updated" != s.value {
		t.Fatal("findSubscriber value is not updated")
	}
}

func TestWritableSubscribe(t *testing.T) {
	s := WritableCreate("default")

	unsubscribe := WritableSubscribe(s, func(value string) {
		if "default" != value {
			t.Fatal("store value is not default")
		}
	})

	if 1 != len(s.subscribers) {
		t.Fatal("store does not have 1 findSubscriber")
	}

	unsubscribe()

	if 0 != len(s.subscribers) {
		t.Fatal("store does not have 0 subscribers")
	}

	s = WritableCreate("default")
	WritableWrite(s, "updated")

	unsubscribe = WritableSubscribe(s, func(value string) {
		if "updated" != value {
			t.Fatal("store value is not updated")
		}
	})

	unsubscribe()
}

func TestWritableRead(t *testing.T) {
	s := WritableCreate("default")
	value := WritableRead(s)

	if "default" != value {
		t.Fatal("store value is not default")
	}
}
