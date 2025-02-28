/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package resid

import (
	"strings"

	"sigs.k8s.io/kustomize/pkg/gvk"
)

// ResId is an immutable identifier of a k8s resource object.
type ResId struct {
	ItemId
	// namePrefix of the resource.
	// An untransformed resource has no prefix.
	// A fully transformed resource has an arbitrary
	// number of prefixes concatenated together.
	Prefix string `json:"prefix,omitempty"`

	// nameSuffix of the resource.
	// An untransformed resource has no suffix.
	// A fully transformed resource has an arbitrary
	// number of suffixes concatenated together.
	Suffix string `json:"suffix,omitempty"`
}

// NewResIdWithPrefixSuffixNamespace creates new resource identifier with a prefix, suffix and a namespace
func NewResIdWithPrefixSuffixNamespace(k gvk.Gvk, n, p, s, ns string) ResId {
	return ResId{ItemId: ItemId{Gvk: k, Name: n, Namespace: ns}, Prefix: p, Suffix: s}
}

// NewResIdWithPrefixNamespace creates new resource identifier with a prefix and a namespace
func NewResIdWithPrefixNamespace(k gvk.Gvk, n, p, ns string) ResId {
	return ResId{ItemId: ItemId{Gvk: k, Name: n, Namespace: ns}, Prefix: p}
}

// NewResIdWithSuffixNamespace creates new resource identifier with a suffix and a namespace
func NewResIdWithSuffixNamespace(k gvk.Gvk, n, s, ns string) ResId {
	return ResId{ItemId: ItemId{Gvk: k, Name: n, Namespace: ns}, Suffix: s}
}

// NewResIdWithPrefixSuffix creates new resource identifier with a prefix and suffix
func NewResIdWithPrefixSuffix(k gvk.Gvk, n, p, s string) ResId {
	return ResId{ItemId: ItemId{Gvk: k, Name: n}, Prefix: p, Suffix: s}
}

// NewResId creates new resource identifier
func NewResId(k gvk.Gvk, n string) ResId {
	return ResId{ItemId: ItemId{Gvk: k, Name: n}}
}

// NewResIdKindOnly creates new resource identifier
func NewResIdKindOnly(k string, n string) ResId {
	return ResId{ItemId: ItemId{Gvk: gvk.FromKind(k), Name: n}}
}

const (
	noNamespace = "~X"
	noPrefix    = "~P"
	noName      = "~N"
	noSuffix    = "~S"
	separator   = "|"
)

// String of ResId based on GVK, name and prefix
func (n ResId) String() string {
	ns := n.ItemId.Namespace
	if ns == "" {
		ns = noNamespace
	}
	p := n.Prefix
	if p == "" {
		p = noPrefix
	}
	nm := n.ItemId.Name
	if nm == "" {
		nm = noName
	}
	s := n.Suffix
	if s == "" {
		s = noSuffix
	}

	return strings.Join(
		[]string{n.ItemId.Gvk.String(), ns, p, nm, s}, separator)
}

// GvknString of ResId based on GVK and name
func (n ResId) GvknString() string {
	return n.ItemId.Gvk.String() + separator + n.ItemId.Name
}

// GvknEquals returns true if the other id matches
// Group/Version/Kind/name.
func (n ResId) GvknEquals(id ResId) bool {
	return n.ItemId.Name == id.ItemId.Name && n.ItemId.Gvk.Equals(id.ItemId.Gvk)
}

// NsGvknEquals returns true if the other id matches
// namespace/Group/Version/Kind/name.
func (n ResId) NsGvknEquals(id ResId) bool {
	// TODO: same a n.ItemId.Equals(id.ItemId)
	return n.ItemId.Namespace == id.ItemId.Namespace && n.GvknEquals(id)
}

// Gvk returns Group/Version/Kind of the resource.
func (n ResId) Gvk() gvk.Gvk {
	return n.ItemId.Gvk
}

// Name returns resource name.
func (n ResId) Name() string {
	return n.ItemId.Name
}

// Namespace returns resource namespace.
func (n ResId) Namespace() string {
	return n.ItemId.Namespace
}

// CopyWithNewPrefixSuffix make a new copy from current ResId
// and append a new prefix and suffix
func (n ResId) CopyWithNewPrefixSuffix(p, s string) ResId {
	result := n
	if p != "" {
		result.Prefix = n.concatPrefix(p)
	}
	if s != "" {
		result.Suffix = n.concatSuffix(s)
	}
	return result
}

// CopyWithNewNamespace make a new copy from current ResId and set a new namespace
func (n ResId) CopyWithNewNamespace(ns string) ResId {
	result := n
	result.ItemId.Namespace = ns
	return result
}

// HasSameLeftmostPrefix check if two ResIds have the same
// left most prefix.
func (n ResId) HasSameLeftmostPrefix(id ResId) bool {
	prefixes1 := n.prefixList()
	prefixes2 := id.prefixList()
	return prefixes1[0] == prefixes2[0]
}

// HasSameRightmostSuffix check if two ResIds have the same
// right most suffix.
func (n ResId) HasSameRightmostSuffix(id ResId) bool {
	suffixes1 := n.suffixList()
	suffixes2 := id.suffixList()
	return suffixes1[len(suffixes1)-1] == suffixes2[len(suffixes2)-1]
}

func (n ResId) concatPrefix(p string) string {
	if p == "" {
		return n.Prefix
	}
	if n.Prefix == "" {
		return p
	}
	return p + ":" + n.Prefix
}

func (n ResId) concatSuffix(s string) string {
	if s == "" {
		return n.Suffix
	}
	if n.Suffix == "" {
		return s
	}
	return n.Suffix + ":" + s
}

func (n ResId) prefixList() []string {
	return strings.Split(n.Prefix, ":")
}

func (n ResId) suffixList() []string {
	return strings.Split(n.Suffix, ":")
}
