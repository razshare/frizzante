package frizzante

import "testing"

func TestReadableCreate(t *testing.T) {
	s := ReadableCreate("default", func(set func(value string)) func() {
		return func() {}
	})

	if "default" != s.value {
		t.Fatal("findSubscriber value is not default")
	}

	s = ReadableCreate("default", func(set func(value string)) func() {
		set("updated")
		return func() {}
	})

	if "updated" != s.value {
		t.Fatal("findSubscriber value is not updated")
	}

	s = ReadableCreate("default", func(set func(value string)) func() {
		return func() {
			set("destroyed")
		}
	})

	unsubscribe := ReadableSubscribe(s, func(value string) {})
	unsubscribe()

	if "destroyed" != s.value {
		t.Fatal("findSubscriber value is not destroyed")
	}
}

func TestReadableSubscribe(t *testing.T) {
	s := ReadableCreate("default", func(set func(value string)) func() {
		return func() {}
	})

	unsubscribe := ReadableSubscribe(s, func(value string) {
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

	s = ReadableCreate("default", func(set func(value string)) func() {
		set("updated")
		return func() {}
	})

	unsubscribe = ReadableSubscribe(s, func(value string) {
		if "updated" != value {
			t.Fatal("store value is not updated")
		}
	})

	unsubscribe()
}

func TestReadableRead(t *testing.T) {
	s := ReadableCreate("default", func(set func(value string)) func() {
		return func() {

		}
	})

	value := ReadableGet(s)

	if "default" != value {
		t.Fatal("store value is not default")
	}
}
