// Package concurrency contain helpers to help with parallel, rate limited, and worker oriented concurrency.
//
// Most helper functions will take an input slice, and perform a transform function on each element of the slice.
// The typical signature will look like something this:
// `ConcurrentFunc(input []T, func(in T) (out O, ok bool)) output []O`
//
// If you don't care about the input or output, set the types to empty structs (struct{}) to reduce the memory footprint.
// On the output, have the transform function return false in the bool field on every invocation to discard the outputs
// and reduce the memory footprint even further.
//
// If you want the transform function to return an error for you to handle:
// make the output type an error value or a struct with an error field.
// In that way you can handle the errors when handling the result slice.
package concurrency
