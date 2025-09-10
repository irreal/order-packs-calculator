package models

// In a real world app, a pack would likely be a struct with more fields

// type Pack struct {
// 	Size int
// 	...other fields
// }

// However, here we will declare a pack to be an int, as we only care about it's size,
// and it makes it less awkward than a struct with a single field for primary key purposes such storing in maps
type Pack int
