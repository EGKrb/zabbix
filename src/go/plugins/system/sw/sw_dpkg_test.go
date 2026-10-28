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
	"reflect"
	"testing"
)

func Test_dpkgList(t *testing.T) {
	t.Parallel()

	// dpkg --get-selections
	input := []string{
		"adduser\t\t\t\t\tinstall",
		"nano\t\t\t\t\t\thold",
		"docutils-common\t\t\t\tdeinstall",
		"libfoo:amd64\t\t\t\t\tpurge",
	}

	out, err := dpkgList(input, "")
	if err != nil {
		t.Fatalf("dpkgList() error = %v", err)
	}

	if expected := []string{"adduser", "nano"}; !reflect.DeepEqual(out, expected) {
		t.Errorf("dpkgList() = %v, expected %v", out, expected)
	}
}

func Test_dpkgDetails(t *testing.T) {
	t.Parallel()

	input := []string{
		"install ok installed,adduser,3.134,all,849",
		"hold ok installed,nano,7.2-1+deb12u1,amd64,2758",
		"deinstall ok config-files,docutils-common,0.19+dfsg-6,all,670",
		"hold ok half-installed,libbar,1.0-1,amd64,10",
		"install ok not-installed,libfoo,,,",
	}

	out, err := dpkgDetails("dpkg", input, "")
	if err != nil {
		t.Fatalf("dpkgDetails() error = %v", err)
	}

	var pd []PackageDetails

	err = json.Unmarshal([]byte(out), &pd)
	if err != nil {
		t.Fatalf("cannot unmarshal %s: %v", out, err)
	}

	names := []string{}
	for _, p := range pd {
		names = append(names, p.Name)
	}

	if expected := []string{"adduser", "nano"}; !reflect.DeepEqual(names, expected) {
		t.Errorf("dpkgDetails() packages = %v, expected %v", names, expected)
	}
}
