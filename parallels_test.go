package main

import (
	"github.com/razshare/frizzante/parallels"
	"testing"
	"time"
)

func TestCompose(test *testing.T) {
	list := parallels.Compose[string](
		func(complete parallels.Complete[string]) {
			time.Sleep(2 * time.Millisecond)
			complete <- "bye!"
		},
		func(complete parallels.Complete[string]) {
			time.Sleep(time.Millisecond)
			complete <- "hello"
		},
	)

	count := len(list)

	if count != 2 {
		test.Fatalf("list was expected to contain %d items, received %d instead", 2, count)
	}

	if list[0] != "hello" {
		test.Fatalf("first item in list was expected to be `%s`, received `%s` instead", "hello", list[0])
	}

	if list[1] != "bye!" {
		test.Fatalf("first item in list was expected to be `%s`, received `%s` instead", "bye!", list[1])
	}
}

func TestFirst(test *testing.T) {
	text := parallels.First[string](
		func(complete parallels.Complete[string], completed parallels.Completed[string]) {
			time.Sleep(2 * time.Millisecond)
			select {
			case <-completed:
				return
			default:
				complete <- "bye!"
			}
		},
		func(complete parallels.Complete[string], completed parallels.Completed[string]) {
			time.Sleep(time.Millisecond)
			select {
			case <-completed:
				return
			default:
				complete <- "hello"
			}
		},
	)

	if text != "hello" {
		test.Fatalf("text was expected to be `%s`, received `%s` instead", "hello", text)
	}
}
