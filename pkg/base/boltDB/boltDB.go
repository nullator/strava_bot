package boltdb

import "github.com/boltdb/bolt"

type base struct {
	db *bolt.DB
}

func NewBase(db *bolt.DB) *base {
	return &base{db: db}
}

func (db *base) Save(key string, value string, bucket string) error {
	err := db.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		return b.Put([]byte(key), []byte(value))
	})
	return err
}

func (db *base) Get(key string, bucket string) (*string, error) {
	var value *string
	err := db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b != nil {
			data := b.Get([]byte(key))
			v := string(data)
			value = &v
		} else {
			value = nil
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return value, err
}
