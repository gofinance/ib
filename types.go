package ib

import (
	"bytes"
)

// Request .
type Request interface {
	writable
	code() OutgoingMessageID
	version() int64
}

// Reply .
type Reply interface {
	readable
	code() IncomingMessageID
}

// MatchedRequest .
type MatchedRequest interface {
	Request
	SetID(id int64)
	ID() int64
}

// MatchedReply .
type MatchedReply interface {
	Reply
	ID() int64
}

type clientHandshake struct {
	version int64
}

func (c *clientHandshake) write(b *bytes.Buffer) error {
	return writeInt(b, c.version)
}
