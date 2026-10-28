//go:build linux && (amd64 || arm64)

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

package sw

import "testing"

func Test_getParams(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		params  []string
		regex   string
		manager string
		short   bool
	}{
		{"default", []string{}, "", "all", false},
		{"allKeyword", []string{"all"}, "", "all", false},
		{"allKeywordAndManager", []string{"all", "dpkg"}, "", "dpkg", false},
		{"allManager", []string{"^zabbix", "all", "short"}, "^zabbix", "all", true},
		{"regex", []string{"^(all|tall)$"}, "^(all|tall)$", "all", false},
		{"allSubstring", []string{"allure"}, "allure", "all", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			regex, manager, short, err := getParams(tt.params, 3)
			if err != nil {
				t.Fatalf("getParams() error = %v", err)
			}

			if regex != tt.regex || manager != tt.manager || short != tt.short {
				t.Errorf("getParams() = %q, %q, %v, want %q, %q, %v",
					regex, manager, short, tt.regex, tt.manager, tt.short)
			}
		})
	}
}
