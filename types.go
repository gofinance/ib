package trade

import (
	"bytes"
)

type Request interface {
	writable
	code() OutgoingMessageId
	version() int64
}

type Reply interface {
	readable
	code() IncomingMessageId
}

type MatchedRequest interface {
	Request
	SetId(id int64)
	Id() int64
}

type MatchedReply interface {
	Reply
	Id() int64
}

type clientHandshake struct {
	version int64
	id      int64
}

func (c *clientHandshake) write(b *bytes.Buffer) (err error) {
	if err = writeInt(b, c.version); err != nil {
		return
	}
	if err = writeInt(b, mStartAPI); err != nil {
		return
	}
	if err = writeInt(b, 1); err != nil {
		return
	}
	return writeInt(b, c.id)
}
