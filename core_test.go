package core

import (
	"fmt"
	"os"
	"testing"
)

const (
	testDir = "./test_data"
)

var c *Core

func TestNew(t *testing.T) {
	var (
		c   *Core
		err error
	)

	if c, err = testInit(); err != nil {
		t.Fatal(err)
	}

	if err = testTeardown(c); err != nil {
		t.Fatal(err)
	}
}

func TestCore_New(t *testing.T) {
	var (
		c   *Core
		err error
	)

	if c, err = testInit(); err != nil {
		t.Fatal(err)
	}
	defer testTeardown(c)

	foobar := newTestStruct("user_1", "contact_1", "FOO FOO", "bunny bar bar")

	var entryID string
	if entryID, err = c.New(&foobar); err != nil {
		t.Fatal(err)
	}

	if len(entryID) == 0 {
		t.Fatal("invalid entry id, expected non-empty value")
	}
}

func TestCore_Get(t *testing.T) {
	var (
		c   *Core
		err error
	)

	if c, err = testInit(); err != nil {
		t.Fatal(err)
	}
	defer testTeardown(c)

	foobar := newTestStruct("user_1", "contact_1", "FOO FOO", "bunny bar bar")

	var entryID string
	if entryID, err = c.New(&foobar); err != nil {
		t.Fatal(err)
	}

	var fb testStruct
	if err = c.Get(entryID, &fb); err != nil {
		t.Fatal(err)
	}

	if err = testCheck(&foobar, &fb); err != nil {
		t.Fatal(err)
	}
}

func TestCore_GetByRelationship_users(t *testing.T) {
	var (
		c   *Core
		err error
	)

	if c, err = testInit(); err != nil {
		t.Fatal(err)
	}
	defer testTeardown(c)

	foobar := newTestStruct("user_1", "contact_1", "FOO FOO", "bunny bar bar")

	if _, err = c.New(&foobar); err != nil {
		t.Fatal(err)
	}

	var foobars []*testStruct
	if err = c.GetByRelationship("users", "user_1", &foobars); err != nil {
		t.Fatal(err)
	}

	for _, fb := range foobars {
		if err = testCheck(&foobar, fb); err != nil {
			t.Fatal(err)
		}
	}
}

func TestCore_GetByRelationship_contacts(t *testing.T) {
	var (
		c   *Core
		err error
	)

	if c, err = testInit(); err != nil {
		t.Fatal(err)
	}
	defer testTeardown(c)

	foobar := newTestStruct("user_1", "contact_1", "FOO FOO", "bunny bar bar")

	if _, err = c.New(&foobar); err != nil {
		t.Fatal(err)
	}

	var foobars []*testStruct
	if err = c.GetByRelationship("contacts", "contact_1", &foobars); err != nil {
		t.Fatal(err)
	}

	for _, fb := range foobars {
		if err = testCheck(&foobar, fb); err != nil {
			t.Fatal(err)
		}
	}
}

func TestCore_GetByRelationship_invalid(t *testing.T) {
	var (
		c   *Core
		err error
	)

	if c, err = testInit(); err != nil {
		t.Fatal(err)
	}
	defer testTeardown(c)

	foobar := newTestStruct("user_1", "contact_1", "FOO FOO", "bunny bar bar")

	if _, err = c.New(&foobar); err != nil {
		t.Fatal(err)
	}

	var foobars []*testBadType
	if err = c.GetByRelationship("contacts", "contact_1", &foobars); err != ErrInvalidType {
		t.Fatalf("invalid error, expected %v and received %v", ErrInvalidType, err)
	}
}

func TestCore_GetByRelationship_update(t *testing.T) {
	var (
		c   *Core
		err error
	)

	if c, err = testInit(); err != nil {
		t.Fatal(err)
	}
	defer testTeardown(c)

	foobar := newTestStruct("user_1", "contact_1", "FOO FOO", "bunny bar bar")

	var entryID string
	if entryID, err = c.New(&foobar); err != nil {
		t.Fatal(err)
	}

	foobar.UserID = "user_3"

	if err = c.Edit(entryID, &foobar); err != nil {
		t.Fatal(err)
	}

	var foobars []*testStruct
	if err = c.GetByRelationship("users", "user_1", &foobars); err != nil {
		t.Fatal(err)
	}

	if len(foobars) != 0 {
		t.Fatalf("invalid number of entries, expected %d and received %d", 0, len(foobars))
	}

	if err = c.GetByRelationship("users", "user_3", &foobars); err != nil {
		t.Fatal(err)
	}

	if len(foobars) != 1 {
		t.Fatalf("invalid number of entries, expected %d and received %d", 1, len(foobars))
	}

	for _, fb := range foobars {
		if err = testCheck(&foobar, fb); err != nil {
			t.Fatal(err)
		}
	}
}

func TestCore_Edit(t *testing.T) {
	var (
		c   *Core
		err error
	)

	if c, err = testInit(); err != nil {
		t.Fatal(err)
	}
	defer testTeardown(c)

	foobar := newTestStruct("user_1", "contact_1", "FOO FOO", "bunny bar bar")

	var entryID string
	if entryID, err = c.New(&foobar); err != nil {
		t.Fatal(err)
	}

	foobar.Foo = "HELLO"

	if err = c.Edit(entryID, &foobar); err != nil {
		t.Fatal(err)
	}

	var fb testStruct
	if err = c.Get(entryID, &fb); err != nil {
		t.Fatal(err)
	}

	if err = testCheck(&foobar, &fb); err != nil {
		t.Fatal(err)
	}
}

func TestCore_ForEach(t *testing.T) {
	var (
		c   *Core
		err error
	)

	if c, err = testInit(); err != nil {
		t.Fatal(err)
	}
	defer testTeardown(c)

	foobar := newTestStruct("user_1", "contact_1", "FOO FOO", "bunny bar bar")

	if _, err = c.New(&foobar); err != nil {
		t.Fatal(err)
	}

	if _, err = c.New(&foobar); err != nil {
		t.Fatal(err)
	}

	var cnt int
	if err = c.ForEach(func(key string, v Value) (err error) {
		fb := v.(*testStruct)
		// We are not checking ID correctness in this test
		foobar.ID = fb.ID

		if err = testCheck(&foobar, fb); err != nil {
			t.Fatal(err)
		}

		cnt++
		return
	}); err != nil {
		t.Fatal(err)
	}

	if cnt != 2 {
		t.Fatalf("invalid number of entries, expected %d and received %d", 2, cnt)
	}

	return
}

func TestCore_ForEachRelationship(t *testing.T) {
	var (
		c   *Core
		err error
	)

	if c, err = testInit(); err != nil {
		t.Fatal(err)
	}
	defer testTeardown(c)

	foobar := newTestStruct("user_1", "contact_1", "FOO FOO", "bunny bar bar")

	if _, err = c.New(&foobar); err != nil {
		t.Fatal(err)
	}

	foobar.UserID = "user_2"
	foobar.ContactID = "contact_3"

	if _, err = c.New(&foobar); err != nil {
		t.Fatal(err)
	}

	var cnt int
	if err = c.ForEachRelationship("contacts", foobar.ContactID, func(key string, v Value) (err error) {
		fb := v.(*testStruct)
		// We are not checking ID correctness in this test
		foobar.ID = fb.ID

		if err = testCheck(&foobar, fb); err != nil {
			t.Fatal(err)
		}

		cnt++
		return
	}); err != nil {
		t.Fatal(err)
	}

	if cnt != 1 {
		t.Fatalf("invalid number of entries, expected %d and received %d", 1, cnt)
	}

	return
}

func TestCore_Cursor(t *testing.T) {
	var (
		c   *Core
		err error
	)

	if c, err = testInit(); err != nil {
		t.Fatal(err)
	}
	defer testTeardown(c)

	foobar := newTestStruct("user_1", "contact_1", "FOO FOO", "bunny bar bar")

	if _, err = c.New(&foobar); err != nil {
		t.Fatal(err)
	}

	if _, err = c.New(&foobar); err != nil {
		t.Fatal(err)
	}

	var cnt int
	if err = c.Cursor(func(cursor *Cursor) (err error) {
		var v Value
		for _, v, err = cursor.Seek(""); err == nil; _, v, err = cursor.Next() {
			fb := v.(*testStruct)
			// We are not checking ID correctness in this test
			foobar.ID = fb.ID

			if err = testCheck(&foobar, fb); err != nil {
				break
			}

			cnt++
		}

		if err == ErrEndOfEntries {
			err = nil
		}

		return
	}); err != nil {
		t.Fatal(err)
	}

	if cnt != 2 {
		t.Fatalf("invalid number of entries, expected %d and received %d", 2, cnt)
	}

	return
}

func TestCore_Cursor_First(t *testing.T) {
	var (
		c   *Core
		err error
	)

	if c, err = testInit(); err != nil {
		t.Fatal(err)
	}
	defer testTeardown(c)

	foobar := newTestStruct("user_1", "contact_1", "FOO FOO", "bunny bar bar")

	if _, err = c.New(&foobar); err != nil {
		t.Fatal(err)
	}

	if _, err = c.New(&foobar); err != nil {
		t.Fatal(err)
	}

	if err = c.Cursor(func(cursor *Cursor) (err error) {
		var v Value
		if _, v, err = cursor.First(); err != nil {
			return
		}

		fb := v.(*testStruct)
		if fb.ID != "00000000" {
			return fmt.Errorf("invalid ID, expected \"%s\" and recieved \"%s\"", "00000000", fb.ID)
		}

		foobar.ID = fb.ID

		if err = testCheck(&foobar, fb); err != nil {
			t.Fatal(err)
		}

		return
	}); err != nil {
		t.Fatal(err)
	}

	return
}

func TestCore_Cursor_Last(t *testing.T) {
	var (
		c   *Core
		err error
	)

	if c, err = testInit(); err != nil {
		t.Fatal(err)
	}
	defer testTeardown(c)

	foobar := newTestStruct("user_1", "contact_1", "FOO FOO", "bunny bar bar")

	if _, err = c.New(&foobar); err != nil {
		t.Fatal(err)
	}

	if _, err = c.New(&foobar); err != nil {
		t.Fatal(err)
	}

	if err = c.Cursor(func(cursor *Cursor) (err error) {
		var v Value
		if _, v, err = cursor.Last(); err != nil {
			return
		}

		fb := v.(*testStruct)
		if fb.ID != "00000001" {
			return fmt.Errorf("invalid ID, expected \"%s\" and recieved \"%s\"", "00000001", fb.ID)
		}

		foobar.ID = fb.ID

		if err = testCheck(&foobar, fb); err != nil {
			t.Fatal(err)
		}

		return
	}); err != nil {
		t.Fatal(err)
	}

	return
}

func TestCore_Cursor_Seek(t *testing.T) {
	var (
		c   *Core
		err error
	)

	if c, err = testInit(); err != nil {
		t.Fatal(err)
	}
	defer testTeardown(c)

	foobar := newTestStruct("user_1", "contact_1", "FOO FOO", "bunny bar bar")

	if _, err = c.New(&foobar); err != nil {
		t.Fatal(err)
	}

	if _, err = c.New(&foobar); err != nil {
		t.Fatal(err)
	}

	if err = c.Cursor(func(cursor *Cursor) (err error) {
		var v Value
		if _, v, err = cursor.Seek("00000001"); err != nil {
			return
		}

		fb := v.(*testStruct)
		if fb.ID != "00000001" {
			return fmt.Errorf("invalid ID, expected \"%s\" and recieved \"%s\"", "00000001", fb.ID)
		}

		foobar.ID = fb.ID

		if err = testCheck(&foobar, fb); err != nil {
			t.Fatal(err)
		}

		return
	}); err != nil {
		t.Fatal(err)
	}

	return
}

func TestCore_CursorRelationship(t *testing.T) {
	var (
		c   *Core
		err error
	)

	if c, err = testInit(); err != nil {
		t.Fatal(err)
	}
	defer testTeardown(c)

	foobar := newTestStruct("user_1", "contact_1", "FOO FOO", "bunny bar bar")

	if _, err = c.New(&foobar); err != nil {
		t.Fatal(err)
	}

	foobar.UserID = "user_2"
	foobar.ContactID = "contact_3"

	if _, err = c.New(&foobar); err != nil {
		t.Fatal(err)
	}

	var cnt int
	if err = c.CursorRelationship("contacts", foobar.ContactID, func(cursor *Cursor) (err error) {
		var v Value
		for _, v, err = cursor.Seek(""); err == nil; _, v, err = cursor.Next() {
			fb := v.(*testStruct)
			// We are not checking ID correctness in this test
			foobar.ID = fb.ID

			if err = testCheck(&foobar, fb); err != nil {
				t.Fatal(err)
			}

			cnt++
		}

		if err == ErrEndOfEntries {
			err = nil
		}

		return
	}); err != nil {
		t.Fatal(err)
	}

	if cnt != 1 {
		t.Fatalf("invalid number of entries, expected %d and received %d", 1, cnt)
	}

	return
}

func TestCore_Lookups(t *testing.T) {
	var (
		c   *Core
		err error
	)

	if c, err = testInit(); err != nil {
		t.Fatal(err)
	}
	defer testTeardown(c)

	if err = c.SetLookup("test_lookup", "test_0", "foo"); err != nil {
		t.Fatal(err)
	}

	if err = c.SetLookup("test_lookup", "test_0", "bar"); err != nil {
		t.Fatal(err)
	}

	var keys []string
	if keys, err = c.GetLookup("test_lookup", "test_0"); err != nil {
		t.Fatal(err)
	}

	if len(keys) != 2 {
		t.Fatalf("invalid number of keys, expected %d and received %d (%+v)", 2, len(keys), keys)
	}

	for i, key := range keys {
		var expected string
		switch i {
		case 0:
			expected = "bar"
		case 1:
			expected = "foo"
		}

		if expected != key {
			t.Fatalf("invalid key, expected %s and recieved %s", expected, key)
		}
	}

	if err = c.RemoveLookup("test_lookup", "test_0", "foo"); err != nil {
		t.Fatal(err)
	}

	if keys, err = c.GetLookup("test_lookup", "test_0"); err != nil {
		t.Fatal(err)
	}

	if len(keys) != 1 {
		t.Fatalf("invalid number of keys, expected %d and received %d (%+v)", 1, len(keys), keys)
	}

	if keys[0] != "bar" {
		t.Fatalf("invalid key, expected %s and recieved %s", "bar", keys[0])
	}
}

func ExampleNew() {
	var (
		c   *Core
		err error
	)

	if c, err = New("example", "./data", &testStruct{}, "users", "contacts"); err != nil {
		return
	}

	fmt.Printf("Core! %v\n", c)
}

func ExampleCore_New() {
	var ts testStruct
	ts.Foo = "Foo foo"
	ts.Bar = "Bar bar"

	var (
		entryID string
		err     error
	)

	if entryID, err = c.New(&ts); err != nil {
		return
	}

	fmt.Printf("New entry! %s\n", entryID)
}

func ExampleCore_Get() {
	var (
		ts  testStruct
		err error
	)

	if err = c.Get("00000000", &ts); err != nil {
		return
	}

	fmt.Printf("Retrieved entry! %+v\n", ts)
}

func ExampleCore_GetByRelationship() {
	var (
		tss []*testStruct
		err error
	)

	if err = c.GetByRelationship("users", "user_1", &tss); err != nil {
		return
	}

	for i, ts := range tss {
		fmt.Printf("Retrieved entry #%d! %+v\n", i, ts)
	}
}

func ExampleCore_ForEach() {
	var err error
	if err = c.ForEach(func(key string, val Value) (err error) {
		fmt.Printf("Iterating entry (%s)! %+v\n", key, val)
		return
	}); err != nil {
		return
	}
}

func ExampleCore_ForEachRelationship() {
	var err error
	if err = c.ForEachRelationship("users", "user_1", func(key string, val Value) (err error) {
		fmt.Printf("Iterating entry (%s)! %+v\n", key, val)
		return
	}); err != nil {
		return
	}
}

func ExampleCore_Edit() {
	var (
		ts  *testStruct
		err error
	)

	// We will pretend the test struct is already populated

	// Let's update the Foo field to "New foo value"
	ts.Foo = "New foo value"

	if err = c.Edit("00000000", ts); err != nil {
		return
	}

	fmt.Printf("Edited entry %s!\n", "00000000")
}

func ExampleCore_Remove() {
	var err error
	if err = c.Remove("00000000"); err != nil {
		return
	}

	fmt.Printf("Removed entry %s!\n", "00000000")
}

func testInit() (c *Core, err error) {
	if err = os.MkdirAll(testDir, 0744); err != nil {
		return
	}

	return New("test", testDir, &testStruct{}, "users", "contacts")
}

func testTeardown(c *Core) (err error) {
	if err = c.Close(); err != nil {
		return
	}

	return os.RemoveAll(testDir)
}

func testCheck(a, b *testStruct) (err error) {
	if a.ID != b.ID {
		return fmt.Errorf("invalid id, expected %s and received %s", a.ID, b.ID)
	}

	if a.UserID != b.UserID {
		return fmt.Errorf("invalid user id, expected %s and received %s", a.UserID, b.UserID)
	}

	if a.ContactID != b.ContactID {
		return fmt.Errorf("invalid contact id, expected %s and received %s", a.ContactID, b.ContactID)
	}

	if a.Foo != b.Foo {
		return fmt.Errorf("invalid foo, expected %s and received %s", a.Foo, b.Foo)
	}

	if a.Bar != b.Bar {
		return fmt.Errorf("invalid bar, expected %s and received %s", a.Bar, b.Bar)
	}

	return
}

func newTestStruct(userID, contactID, foo, bar string) (t testStruct) {
	t.UserID = userID
	t.ContactID = contactID
	t.Foo = foo
	t.Bar = bar
	return
}

type testStruct struct {
	ID        string `json:"id"`
	UserID    string `json:"userID"`
	ContactID string `json:"contactID"`

	Foo string `json:"foo"`
	Bar string `json:"bar"`

	UpdatedAt int64 `json:"updatedAt"`
	CreatedAt int64 `json:"createdAt"`
}

func (t *testStruct) SetID(id string) {
	t.ID = id
}

func (t *testStruct) GetUpdatedAt() (updatedAt int64) {
	return t.UpdatedAt
}

func (t *testStruct) GetCreatedAt() (createdAt int64) {
	return t.CreatedAt
}

func (t *testStruct) GetID() (id string) {
	return t.ID
}

func (t *testStruct) GetRelationshipIDs() (ids []string) {
	ids = append(ids, t.UserID)
	ids = append(ids, t.ContactID)
	return
}

func (t *testStruct) SetUpdatedAt(updatedAt int64) {
	t.UpdatedAt = updatedAt
}

func (t *testStruct) SetCreatedAt(createdAt int64) {
	t.CreatedAt = createdAt
}

type testBadType struct {
	Foo string
	Bar string
}
