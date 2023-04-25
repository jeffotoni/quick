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
		msgidValue string
	}{
		{
			name: "success",
			args: args{
				Config: Config{
					Name:  KeyMsgID,
					Start: 0,
					End:   0,
					Algo:  nil,
				},
			},
			testMsgID:  testMsgIDSuccess,
			msgidValue: "",
		},
	}

	for _, ti := range tests {
		t.Run(ti.name, func(*testing.T) {
			t.Logf("==== TEST %s ====", ti.name)
			h := New(ti.args.Config)
			a := h(ti.testMsgID.HandlerFunc)
			rec := httptest.NewRecorder()
			ti.testMsgID.Request.Header.Set(KeyMsgID, ti.msgidValue)
			a.ServeHTTP(rec, ti.testMsgID.Request)
			resp := rec.Result()

			if resp.Header.Get(KeyMsgID) == "" && len(ti.msgidValue) == 0 {
				t.Errorf("was expected a uuid and nothing came")
			}
		})
	}
}

// go test -bench=. -benchtime=1s -benchmem
func BenchmarkNew(b *testing.B) {
	for n := 0; n < b.N; n++ {
		New()
	}
	// b.StopTimer()
}

// go test -bench=. -benchtime=1s -benchmem
func BenchmarkAlgoDefault(b *testing.B) {
	for n := 0; n < b.N; n++ {
		AlgoDefault(9000000, 10000000)
	}
	// b.StopTimer()
}
