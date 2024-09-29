package main

import "strings"

type stringSlice []string

func (s *stringSlice) String() string {
	return strings.Join(*s, ", ")
}

func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

// type sliceFlag[T any] struct {
// 	slice []T
// }
// func (s *sliceFlag[T]) String() string {
// 	return fmt.Sprintf("%v", s.slice)
// }
// func (s *sliceFlag[T]) Set(value string) error {
// 	s.slice = append(s.slice, *(*T)(unsafe.Pointer(&value)))
// 	return nil
// }
