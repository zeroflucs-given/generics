package bplustree

import collections "github.com/zeroflucs-given/generics/collections"

type record[V any] struct {
	RecordID collections.RecordID `json:"rid"`   // Record ID
	Value    V                    `json:"value"` // Data stored in the record
}
