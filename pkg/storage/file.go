package storage

import (
	"bytes"
	"context"
	"fmt"

	"github.com/etcd-io/bbolt"
	"github.com/supergiant/control/pkg/sgerrors"
)

const bucketName = "supergiant.io"

type FileRepository struct {
	db *bbolt.DB
}

func NewFileRepository(fileName string) (*FileRepository, error) {
	db, err := bbolt.Open(fileName, 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))

		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &FileRepository{
		db: db,
	}, nil
}

func (i *FileRepository) Get(ctx context.Context, prefix string, key string) ([]byte, error) {
	var value []byte

	err := i.db.View(func(tx *bbolt.Tx) error {
		value = tx.Bucket([]byte(bucketName)).Get([]byte(prefix + key))

		if value == nil {
			return sgerrors.ErrNotFound
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return value, nil
}

func (i *FileRepository) Put(ctx context.Context, prefix string, key string, value []byte) error {
	err := i.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucketName))

		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		err = bucket.Put([]byte(prefix+key), value)
		return err
	})

	return err
}

func (i *FileRepository) Delete(ctx context.Context, prefix string, key string) error {
	err := i.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucketName))

		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		return bucket.Delete([]byte(prefix + key))
	})

	return err
}

func (i *FileRepository) GetAll(ctx context.Context, prefix string) ([][]byte, error) {
	values := make([][]byte, 0)

	err := i.db.View(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucketName))

		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		cursor := bucket.Cursor()
		prefixBytes := []byte(prefix)

		for k, v := cursor.Seek(prefixBytes); k != nil && bytes.HasPrefix(k, prefixBytes); k, v = cursor.Next() {
			values = append(values, v)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return values, nil
}
