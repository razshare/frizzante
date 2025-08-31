package text

import "testing"

func TestSha1(t *testing.T) {
	hash, err := Sha1("hello")
	if err != nil {
		t.Fatal(err)
	}

	if hash != "qvTGHdzF6KLavt4PO0gs2a6pQ00" {
		t.Fatal("hash should be qvTGHdzF6KLavt4PO0gs2a6pQ00")
	}
}
