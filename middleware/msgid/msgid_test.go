package msgid

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// go test -v -failfast -count=1 -run ^TestNew$
func TestNew(t *testing.T) {

	type args struct {
		Config Config
	}

	tests := []struct {
		name        string
		args        args
		wantHandler func(http.Handler) http.Handler
	}{
		{
			name: "success",
			args: args{
				Config: Config{
					UUID:       false,
					Name:       "",
					Start:      0,
					End:        0,
					Algo:       nil,
					ConfigUUID: NewUUID(),
				},
			},
			wantHandler: returnSuccessHandler,
		},
	}

	for _, ti := range tests {
		t.Run(ti.name, func(*testing.T) {
			h := New(ti.args.Config)

			come := httptest.NewServer(h)

			desired := httptest.NewServer(ti.wantHandler)

			if come != desired {
				t.Errorf("was expected %#v and %#v come", desired, come)
			}
		})
	}
}
