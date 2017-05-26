// Copyright 2016-2017 Authors of Cilium
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package endpoint

import (
	"os"

	"github.com/cilium/cilium/pkg/maps/cidrmap"
)

type L3MapType int

const (
	IPv6Ingress L3MapType = iota
	IPv4Ingress
	IPv6Egress
	IPv4Egress
	MapCount
)

type L3Maps [MapCount]*cidrmap.CIDRMap

func (l3 *L3Maps) DestroyBpfMap(mt L3MapType, path string) {
	if l3[mt] != nil {
		l3[mt].Close()
		l3[mt] = nil
	}
	os.RemoveAll(path)
}

func (l3 *L3Maps) CreateBpfMap(mt L3MapType, path string) error {
	var err error

	// LPM trie maps cannot be dumped, so we clear them before opening
	l3.DestroyBpfMap(mt, path)

	prefixlen := int(128)
	if mt == IPv4Ingress || mt == IPv4Egress {
		prefixlen = 32
	}
	l3[mt], _, err = cidrmap.OpenMap(path, prefixlen)
	return err
}

func (l3 *L3Maps) DeepCopy() L3Maps {
	var cpy L3Maps
	for i, v := range l3 {
		cpy[i] = v.DeepCopy()
	}
	return cpy
}

func (l3 *L3Maps) Close() {
	for _, v := range l3 {
		v.Close()
	}
}