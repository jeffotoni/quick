package msguuid

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
	}{
		{
			name: "success_with_uuid_defined_header_value",
			args: args{
				Config: Config{},
			},
			testMsgID: testMsgIDSuccess,
			uuidValue: "quick is awesome!",
		},
		{
			name: "success_with_uuid_defined_header_value_custom_setup_default_version",
			args: args{
				Config: Config{
					Name:      KeyMsgUUID,
					KeyString: "00000000-0000-0000-0000-000000000000",
				},
			},
			testMsgID: testMsgIDSuccess,
			uuidValue: "",
		},
		{
			name: "success_with_uuid_defined_header_value_custom_setup_1",
			args: args{
				Config: Config{
					Version:   1,
					Name:      KeyMsgUUID,
					KeyString: "00000000-0000-0000-0000-000000000000",
				},
			},
			testMsgID: testMsgIDSuccess,
			uuidValue: "",
		},
		{
			name: "success_with_uuid_defined_header_value_custom_setup_2",
			args: args{
				Config: Config{
					Version:   2,
					Name:      KeyMsgUUID,
					KeyString: "00000000-0000-0000-0000-000000000000",
				},
			},
			testMsgID: testMsgIDSuccess,
			uuidValue: "",
		},
		{
			name: "success_with_uuid_defined_header_value_custom_setup_3",
			args: args{
				Config: Config{
					Version:   3,
					Name:      KeyMsgUUID,
					KeyString: "00000000-0000-0000-0000-000000000000",
				},
			},
			testMsgID: testMsgIDSuccess,
			uuidValue: "",
		},
		{
			name: "success_with_uuid_defined_header_value_custom_setup_default_version_no_key",
			args: args{
				Config: Config{
					Name: KeyMsgUUID,
				},
			},
			testMsgID: testMsgIDSuccess,
			uuidValue: "",
		},
		{
			name: "success_with_uuid_defined_header_value_custom_setup_version_1_no_key",
			args: args{
				Config: Config{
					Version: 1,
					Name:    KeyMsgUUID,
				},
			},
			testMsgID: testMsgIDSuccess,
			uuidValue: "",
		},
		{
			name: "success_with_uuid_defined_header_value_custom_setup_version_2_no_key",
			args: args{
				Config: Config{
					Version: 2,
					Name:    KeyMsgUUID,
				},
			},
			testMsgID: testMsgIDSuccess,
			uuidValue: "",
		},
		{
			name: "success_with_uuid_defined_header_value_custom_setup_version_3_no_key",
			args: args{
				Config: Config{
					Version: 3,
					Name:    KeyMsgUUID,
				},
			},
			testMsgID: testMsgIDSuccess,
			uuidValue: "",
		},
	}

	for _, ti := range tests {
		t.Run(ti.name, func(tt *testing.T) {
			t.Logf("==== TEST %s ====", ti.name)
			h := New(ti.args.Config)
			a := h(ti.testMsgID.HandlerFunc)
			rec := httptest.NewRecorder()
			ti.testMsgID.Request.Header.Set(KeyMsgUUID, ti.uuidValue)
			a.ServeHTTP(rec, ti.testMsgID.Request)
			resp := rec.Result()

			if resp.Header.Get(KeyMsgUUID) == "" && len(ti.uuidValue) == 0 {
				tt.Errorf("was expected a uuid and nothing came")
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
func BenchmarkGenerateDefaultUUID(b *testing.B) {
	for n := 0; n < b.N; n++ {
	generateDefaultUUID(DefaultConfig)
	}
	// b.StopTimer()
}
