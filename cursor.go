package core

import (
	"encoding/json"

	"github.com/boltdb/bolt"
)

func newCursor(core *Core, txn *bolt.Tx, cur *bolt.Cursor, relationship bool) (c Cursor) {
	c.core = core
	c.txn = txn
	c.cur = cur
	c.relationship = relationship
	return
}

// Cursor is an iterating structure
type Cursor struct {
	core *Core
	txn  *bolt.Tx
	cur  *bolt.Cursor

	relationship bool
}

func (c *Cursor) get(key, bs []byte, val Value) (err error) {
	if !c.relationship {
		return json.Unmarshal(bs, val)
	}

	if err = c.core.get(c.txn, key, val); err != nil {
		val = nil
		return
	}

	return
}

func (c *Cursor) teardown() {
	c.core = nil
	c.txn = nil
	c.cur = nil
}

// Seek will seek the provided ID
func (c *Cursor) Seek(id string, val Value) (err error) {
	k, v := c.cur.Seek([]byte(id))
	if k == nil && v == nil {
		err = ErrEndOfEntries
		return
	}

	if err = c.get(k, v, val); err != nil {
		return
	}

	return
}

// First will return the first entry
func (c *Cursor) First(val Value) (err error) {
	k, v := c.cur.First()
	if k == nil && v == nil {
		err = ErrEndOfEntries
		return
	}

	if err = c.get(k, v, val); err != nil {
		return
	}

	return
}

// Last will return the last entry
func (c *Cursor) Last(val Value) (err error) {
	k, v := c.cur.Last()
	if k == nil && v == nil {
		err = ErrEndOfEntries
		return
	}

	if err = c.get(k, v, val); err != nil {
		return
	}

	return
}

// Next will return the next entry
func (c *Cursor) Next(val Value) (err error) {
	k, v := c.cur.Next()
	if k == nil && v == nil {
		err = ErrEndOfEntries
		return
	}

	if err = c.get(k, v, val); err != nil {
		return
	}

	return
}

// Prev will return the previous entry
func (c *Cursor) Prev(val Value) (err error) {
	k, v := c.cur.Prev()
	if k == nil && v == nil {
		err = ErrEndOfEntries
		return
	}

	if err = c.get(k, v, val); err != nil {
		return
	}

	return
}