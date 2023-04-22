package msgid

import (
	"net/http/httptest"
	"testing"
)

// go test -v -failfast -count=1 -run ^TestNew$
// go test -v -count=1 -failfast -cover -coverprofile=coverage.out -run ^TestNew$; go tool cover -html=coverage.out
func TestNew(t *testing.T) {

	type args struct {
		Config Config
	}

	tests := []struct {
		name       string
		args       args
		testMsgID  testMsgID
		uuidValue  string
		msgidValue string
	}{}

	for _, ti := range tests {
		t.Run(ti.name, func(*testing.T) {
			t.Logf("==== TEST %s ====", ti.name)
			h := New(ti.args.Config)
			a := h(ti.testMsgID.HandlerFunc)
			rec := httptest.NewRecorder()
			ti.testMsgID.Request.Header.Set(KeyMsgUUID, ti.uuidValue)
			a.ServeHTTP(rec, ti.testMsgID.Request)
			resp := rec.Result()

			if resp.Header.Get(KeyMsgID) == "" && len(ti.msgidValue) == 0 {
				t.Errorf("was expected a uuid and nothing came")
			}
		})
	}
}
