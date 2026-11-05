package exporttask

import "testing"

func TestValidateSQL(t *testing.T) {
	cases := []struct {
		name    string
		sql     string
		wantErr bool
	}{
		{name: "select", sql: "SELECT id, name FROM t WHERE id > 1", wantErr: false},
		{name: "with", sql: "WITH x AS (SELECT id FROM t) SELECT id FROM x", wantErr: false},
		{name: "empty", sql: "  ", wantErr: true},
		{name: "update", sql: "UPDATE t SET name = 'x'", wantErr: true},
		{name: "multiple", sql: "SELECT 1; SELECT 2", wantErr: true},
		{name: "limit", sql: "SELECT * FROM t LIMIT 10", wantErr: true},
		{name: "limit_offset", sql: "SELECT * FROM t LIMIT 10 OFFSET 20", wantErr: true},
		{name: "outfile", sql: "SELECT * FROM t INTO OUTFILE '/tmp/x'", wantErr: true},
		{name: "for_update", sql: "SELECT * FROM t FOR UPDATE", wantErr: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateSQL(tc.sql)
			if (err != nil) != tc.wantErr {
				t.Fatalf("ValidateSQL(%q) err = %v, wantErr = %v", tc.sql, err, tc.wantErr)
			}
		})
	}
}
