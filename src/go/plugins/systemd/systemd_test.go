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

package systemd

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"golang.zabbix.com/sdk/log"
	"golang.zabbix.com/sdk/plugin"
)

func TestPlugin_setUnitStates(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input map[string]any
		want  map[string]any
	}{
		{
			"+knownStates",
			map[string]any{"LoadState": "loaded", "ActiveState": "inactive", "UnitFileState": "invalid"},
			map[string]any{
				"LoadState":     &state{1, "loaded"},
				"ActiveState":   &state{3, "inactive"},
				"UnitFileState": &state{9, "invalid"},
			},
		},
		{
			"+generatedUnit",
			map[string]any{"LoadState": "loaded", "ActiveState": "active", "UnitFileState": "generated"},
			map[string]any{
				"LoadState":     &state{1, "loaded"},
				"ActiveState":   &state{1, "active"},
				"UnitFileState": &state{12, "generated"},
			},
		},
		{
			"+lastStates",
			map[string]any{"LoadState": "merged", "ActiveState": "refreshing", "UnitFileState": "bad"},
			map[string]any{
				"LoadState":     &state{7, "merged"},
				"ActiveState":   &state{8, "refreshing"},
				"UnitFileState": &state{14, "bad"},
			},
		},
		{
			"-unknownState",
			map[string]any{"LoadState": "loaded", "ActiveState": "active", "UnitFileState": "something-new"},
			map[string]any{
				"LoadState":     &state{1, "loaded"},
				"ActiveState":   &state{1, "active"},
				"UnitFileState": &state{0, "something-new"},
			},
		},
		{
			"-notString",
			map[string]any{"LoadState": 1, "ActiveState": "active", "UnitFileState": "static"},
			map[string]any{
				"LoadState":     1,
				"ActiveState":   &state{1, "active"},
				"UnitFileState": &state{7, "static"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			p := &Plugin{Base: plugin.Base{Logger: log.New("test")}}
			p.setUnitStates(tt.input)

			if diff := cmp.Diff(tt.want, tt.input); diff != "" {
				t.Fatalf("setUnitStates() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_isDiscoverableUnitFile(t *testing.T) {
	t.Parallel()

	loaded := []unitJson{{Name: "loaded.service"}, {Name: "mnt-data.mount"}}

	tests := []struct {
		name string
		file unitFile
		ext  string
		want bool
	}{
		{"+disabled", unitFile{"/usr/lib/systemd/system/foo.service", "disabled"}, ".service", true},
		{"+generated", unitFile{"/run/systemd/generator/mnt-backup.mount", "generated"}, ".mount", true},
		{"+allTypes", unitFile{"/run/systemd/generator/mnt-backup.mount", "generated"}, "", true},
		{"-enabled", unitFile{"/usr/lib/systemd/system/foo.service", "enabled"}, ".service", false},
		{"-static", unitFile{"/usr/lib/systemd/system/foo.service", "static"}, ".service", false},
		{"-otherType", unitFile{"/run/systemd/generator/mnt-backup.mount", "generated"}, ".service", false},
		{"-template", unitFile{"/usr/lib/systemd/system/foo@.service", "disabled"}, ".service", false},
		{"-alreadyLoaded", unitFile{"/run/systemd/generator/mnt-data.mount", "generated"}, ".mount", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := isDiscoverableUnitFile(tt.file, tt.ext, loaded); got != tt.want {
				t.Fatalf("isDiscoverableUnitFile(%v, %q) = %v, want %v", tt.file, tt.ext, got, tt.want)
			}
		})
	}
}
