// Copyright 2013, 2014 Peter Vasil, Tomo Krajina. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

package gpx

import (
	"encoding/xml"
	"strconv"
	"strings"
)

type NullableInt struct {
	data int
	null bool
}

func (n *NullableInt) Null() bool {
	return n.null
}

func (n *NullableInt) NotNull() bool {
	return !n.null
}

func (n *NullableInt) Value() int {
	return n.data
}

func (n *NullableInt) SetValue(data int) {
	n.data = data
}

func (n *NullableInt) SetNull() {
	var defaultValue int
	n.data = defaultValue
	n.null = true
}

func NewNullableInt(data int) *NullableInt {
	result := new(NullableInt)
	result.data = data
	result.null = false
	return result
}

func (n *NullableInt) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	t, err := d.Token()
	if err != nil {
		n.SetNull()
		return nil
	}
	if charData, ok := t.(xml.CharData); ok {
		strData := strings.Trim(string(charData), " ")
		value, err := strconv.ParseFloat(strData, 64)
		if err != nil {
			n.SetNull()
			return nil
		}
		n.SetValue(int(value))
	}
	return nil
}

func (n *NullableInt) UnmarshalXMLAttr(attr xml.Attr) error {
	strData := strings.Trim(string(attr.Value), " ")
	value, err := strconv.ParseFloat(strData, 64)
	if err != nil {
		n.SetNull()
		return nil
	}
	n.SetValue(int(value))
	return nil
}
