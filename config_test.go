package docbasecli

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLoadConfig(t *testing.T) {
	doc := []byte(`[default]
AccessToken = "access-token"
Domain      = "domain"
UserID      = "user-id"
Editor      = "vim"
[profile1]
AccessToken = "access-token1"
Domain      = "domain1"
UserID      = "user-id1"
`)
	want := &Config{
		AccessToken: "access-token",
		Domain:      "domain",
		UserID:      "user-id",
		Editor:      "vim",
	}
	got, err := LoadConfig(bytes.NewReader(doc))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	// cmp: configMap
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("configMap mismatch (-want, +got):%s\n", diff)
	}
}
