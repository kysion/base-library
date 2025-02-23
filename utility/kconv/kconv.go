// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package kconv

import (
	"github.com/gogf/gf/v2/util/gconv"
)

// Struct maps the params key-value pairs to the corresponding struct object's attributes.
// The third parameter `mapping` is unnecessary, indicating the mapping rules between the
// custom key name and the attribute name(case-sensitive).
//
// Note:
//  1. The `params` can be any type of map/struct, usually a map.
//  2. The `pointer` should be type of *struct/**struct, which is a pointer to struct object
//     or struct pointer.
//  3. Only the public attributes of struct object can be mapped.
//  4. If `params` is a map, the key of the map `params` can be lowercase.
//     It will automatically convert the first letter of the key to uppercase
//     in mapping procedure to do the matching.
//     It ignores the map key, if it does not match.
func Struct[T any](params interface{}, pointer T, mapping ...map[string]string) (res T) {
	_ = gconv.Struct(params, pointer, mapping...)

	return pointer
}
func StructWithError[T any](params interface{}, pointer T, mapping ...map[string]string) (res T, err error) {
	err = gconv.Struct(params, pointer, mapping...)
	if err == nil {
		return pointer, err
	}
	return pointer, err
}

// StructTag acts as Struct but also with support for priority tag feature, which retrieves the
// specified tags for `params` key-value items to struct attribute names mapping.
// The parameter `priorityTag` supports multiple tags that can be joined with char ','.
func StructTag[T any](params interface{}, pointer T, priorityTag string) (res T) {
	_ = gconv.StructTag(params, pointer, priorityTag)
	return pointer
}
func StructTagWithError[T any](params interface{}, pointer T, priorityTag string) (res T, err error) {
	err = gconv.StructTag(params, pointer, priorityTag)
	if err == nil {
		return pointer, err
	}
	return pointer, nil
}
