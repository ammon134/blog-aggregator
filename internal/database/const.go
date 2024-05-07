package database

import "github.com/lib/pq"

const (
	ErrConstraintUnique     = pq.ErrorCode("23505")
	ErrConstraintForeignKey = pq.ErrorCode("23503")
)
