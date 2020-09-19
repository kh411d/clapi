package db

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/spf13/cast"

	f "github.com/fauna/faunadb-go/v3/faunadb"
)

var errFaunaDocNotUniq = errors.New("Response error 400. Errors: [create](instance not unique): document is not unique.")
var errFaunaSetNotFound = errors.New("Response error 404. Errors: [](instance not found): Set not found.")

type faunaDB struct {
	client *f.FaunaClient
	ctx    context.Context

	//Collection used for kv documents
	coll string
	//Index for the key i.e. "url_idx"
	keyIndex string
}

//NewFaunaDB create faunadb client
func NewFaunaDB(secret, collection, idx string) (kv KV, err error) {
	fdb := f.NewFaunaClient(secret)

	kv = &faunaDB{
		client:   fdb,
		ctx:      context.TODO(),
		coll:     collection,
		keyIndex: idx,
	}
	return
}

func (m *faunaDB) WithContext(ctx context.Context) KV {
	m.ctx = ctx
	return m
}

// Get the item with the provided key.
func (m *faunaDB) Get(key string) (val []byte, err error) {
	// Retrieve profile by its ID
	res, err := m.client.Query(
		f.Get(
			f.MatchTerm(
				f.Index(m.keyIndex),
				key,
			),
		),
	)

	if err != nil {
		return nil, err
	}

	var clap Clap

	if err := res.At(f.ObjKey("data")).Get(&clap); err != nil {
		return nil, err
	}

	val = []byte(cast.ToString(clap.Count))
	return
}

// Add writes the given item, if no value already exists for its key.
func (m *faunaDB) Add(key string, val []byte, expiration time.Duration) (err error) {

	clap := Clap{
		URL:   key,
		Count: cast.ToInt64(val),
	}

	_, err = m.client.Query(
		f.Create(
			f.Collection("claps"),
			f.Obj{"data": clap},
		),
	)

	return
}

// Set writes the given item, unconditionally.
func (m *faunaDB) Set(key string, val []byte, expiration time.Duration) (err error) {
	return
}

// Incr increment key
// TODO: need to find efficient way to do increment
func (m *faunaDB) Incr(key string) (err error) {

	var refID f.RefV
	var clap Clap
	var res f.Value

	clap = Clap{
		URL:   key,
		Count: 1,
	}

	_, err = m.client.Query(
		f.Create(
			f.Collection("claps"),
			f.Obj{"data": clap},
		),
	)

	//Try to update
	if err != nil {
		log.Printf("Error Fauna Incr: %v", err)

		// Retrieve profile by its ID
		res, err = m.client.Query(
			f.Get(
				f.MatchTerm(
					f.Index(m.keyIndex),
					key,
				),
			),
		)

		if err != nil {
			return
		}

		if err := res.At(f.ObjKey("ref")).Get(&refID); err != nil {
			return err
		}

		if err := res.At(f.ObjKey("data")).Get(&clap); err != nil {
			return err
		}

		// Update existing profile entry
		_, err = m.client.Query(
			f.Update(
				refID,
				f.Obj{"data": f.Obj{
					"count": clap.Count + 1,
				}},
			),
		)
	}

	return

}

// IncrBy increment key
// TODO: need to find efficient way to do increment
func (m *faunaDB) IncrBy(key string, val int64) (err error) {

	var refID f.RefV
	var clap Clap
	var res f.Value

	if val == 0 {
		return nil
	}

	clap = Clap{
		URL:   key,
		Count: val,
	}

	_, err = m.client.Query(
		f.Create(
			f.Collection("claps"),
			f.Obj{"data": clap},
		),
	)

	//Try to update
	if err != nil {
		log.Printf("Error Fauna IncrBy: %v", err)

		// Retrieve profile by its ID
		res, err = m.client.Query(
			f.Get(
				f.MatchTerm(
					f.Index(m.keyIndex),
					key,
				),
			),
		)

		if err != nil {
			return
		}

		if err := res.At(f.ObjKey("ref")).Get(&refID); err != nil {
			return err
		}

		if err := res.At(f.ObjKey("data")).Get(&clap); err != nil {
			return err
		}

		// Update existing profile entry
		_, err = m.client.Query(
			f.Update(
				refID,
				f.Obj{"data": f.Obj{
					"count": clap.Count + val,
				}},
			),
		)
	}

	return

}

// Delete deletes the item with the provided key.
func (m *faunaDB) Delete(key string) (err error) {
	return
}
