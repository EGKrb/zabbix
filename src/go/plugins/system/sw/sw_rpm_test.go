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

import (
	"encoding/json"
	"strings"
	"testing"
)

func Test_rpmDetails(t *testing.T) {
	t.Parallel()

	// LONGSIZE output, %{SIZE} prints "(none)" for packages larger than 4 GiB
	input := []string{
		"gpg-pubkey,b5333005-66308c87,(none),0,1714457735,1751463203",
		"fme-engine-2024,2024.2.2.0-24820.el9,x86_64,7995596553,1736897647,1747232166",
	}

	expected := []struct {
		name string
		size uint64
		arch string
	}{
		{"gpg-pubkey", 0, "(none)"},
		{"fme-engine-2024", 7995596553, "x86_64"},
	}

	out, err := rpmDetails("rpm", input, "")
	if err != nil {
		t.Fatalf("rpmDetails() error = %v", err)
	}

	var pd []PackageDetails

	err = json.Unmarshal([]byte(out), &pd)
	if err != nil {
		t.Fatalf("cannot unmarshal %s: %v", out, err)
	}

	if len(pd) != len(expected) {
		t.Fatalf("expected %d packages, got %d: %s", len(expected), len(pd), out)
	}

	for i, e := range expected {
		if pd[i].Name != e.name || pd[i].Size != e.size || pd[i].Arch != e.arch {
			t.Errorf("package %d: expected %s size %d arch %s, got %s size %d arch %s",
				i, e.name, e.size, e.arch, pd[i].Name, pd[i].Size, pd[i].Arch)
		}
	}
}

func Test_rpmQueryLongSize(t *testing.T) {
	t.Parallel()

	for _, m := range getManagers() {
		if m.name != "rpm" {
			continue
		}

		if !strings.Contains(m.detailsCmd, "%{LONGSIZE}") || strings.Contains(m.detailsCmd, "%{SIZE}") {
			t.Errorf("rpm query must use %%{LONGSIZE}, got %s", m.detailsCmd)
		}

		return
	}

	t.Fatal("rpm package manager not found")
}
