/*
** Copyright (C) 2001-2026 Zabbix SIA
**
** This program is free software: you can redistribute it and/or modify it under the terms of
** the GNU Affero General Public License as published by the Free Software Foundation, version 3.
**
** This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY;
** without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
** See the GNU Affero General Public License for more details.
**
** You should have received a copy of the GNU Affero General Public License along with this program.
** If not, see <https://www.gnu.org/licenses/>.
**/

package handlers

import (
	"encoding/json"
	"testing"
)

func Test_blkioStats_setTotals(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		stats string
		want  blkioStats
	}{
		{
			"cgroupV2",
			`{"io_service_bytes_recursive":[{"major":8,"minor":0,"op":"read","value":4096},` +
				`{"major":8,"minor":0,"op":"write","value":8192},{"major":8,"minor":16,"op":"read","value":1000}],` +
				`"io_serviced_recursive":[{"major":8,"minor":0,"op":"read","value":3},` +
				`{"major":8,"minor":0,"op":"write","value":5},{"major":8,"minor":16,"op":"read","value":1}],` +
				`"io_queue_recursive":null}`,
			blkioStats{ReadBytes: 5096, WriteBytes: 8192},
		},
		{
			"cgroupV1",
			`{"io_service_bytes_recursive":[{"major":8,"minor":0,"op":"Read","value":100},` +
				`{"major":8,"minor":0,"op":"Write","value":200},{"major":8,"minor":0,"op":"Sync","value":250},` +
				`{"major":8,"minor":0,"op":"Async","value":50},{"major":8,"minor":0,"op":"Discard","value":0},` +
				`{"major":8,"minor":0,"op":"Total","value":300}],` +
				`"io_serviced_recursive":[{"major":8,"minor":0,"op":"Read","value":1},` +
				`{"major":8,"minor":0,"op":"Write","value":2},{"major":8,"minor":0,"op":"Total","value":3}]}`,
			blkioStats{ReadBytes: 100, WriteBytes: 200},
		},
		{
			"empty",
			`{"io_service_bytes_recursive":null}`,
			blkioStats{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var b blkioStats

			err := json.Unmarshal([]byte(tt.stats), &b)
			if err != nil {
				t.Fatalf("cannot unmarshal: %v", err)
			}

			b.setTotals()

			if b.ReadBytes != tt.want.ReadBytes || b.WriteBytes != tt.want.WriteBytes {
				t.Errorf("setTotals() = read %d, write %d, want read %d, write %d",
					b.ReadBytes, b.WriteBytes, tt.want.ReadBytes, tt.want.WriteBytes)
			}
		})
	}
}
