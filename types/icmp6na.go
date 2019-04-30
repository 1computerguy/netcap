/*
 * NETCAP - Traffic Analysis Framework
 * Copyright (c) 2017 Philipp Mieden <dreadl0ck [at] protonmail [dot] ch>
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

package types

import "strings"

func (i ICMPv6NeighborAdvertisement) CSVHeader() []string {
	return filter([]string{
		"Timestamp",
		"Flags",         // int32
		"TargetAddress", // string
		"Options",       // []*ICMPv6Option
	})
}

func (i ICMPv6NeighborAdvertisement) CSVRecord() []string {
	var opts []string
	for _, o := range i.Options {
		opts = append(opts, o.ToString())
	}
	return filter([]string{
		formatTimestamp(i.Timestamp),
		formatInt32(i.Flags),
		i.TargetAddress,
		strings.Join(opts, ""),
	})
}

func (i ICMPv6NeighborAdvertisement) NetcapTimestamp() string {
	return i.Timestamp
}

func (a ICMPv6NeighborAdvertisement) JSON() (string, error) {
	return jsonMarshaler.MarshalToString(&a)
}
