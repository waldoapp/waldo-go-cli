package data

type CreateKind int

const (
	CreateKindNever CreateKind = iota
	CreateKindAlways
	CreateKindIfNeeded
)
