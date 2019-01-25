package storage

import (
	"context"
	"fmt"

	"github.com/boltdb/bolt"
)

const bucketName = "supergiant.io"

type FileRepository struct {
	db *bolt.DB
}

func NewFileRepository(fileName string) (*FileRepository, error) {
	db, err := bolt.Open(fileName, 0600, nil)
	if err != nil {
		return nil, err
	}

	return &FileRepository{
		db: db,
	}, nil
}

func (i *FileRepository) Get(ctx context.Context, prefix string, key string) ([]byte, error) {
	err := i.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))

		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	return err
}

func (i *FileRepository) Put(ctx context.Context, prefix string, key string, value []byte) error {
	i.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))

		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
}

func (i *FileRepository) Delete(ctx context.Context, prefix string, key string) error {
	i.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))

		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
}

func (i *FileRepository) GetAll(ctx context.Context, prefix string) ([][]byte, error) {i.db.Update(func(tx *bolt.Tx) error {
	_, err := tx.CreateBucketIfNotExists([]byte(bucketName))

	if err != nil {
		return fmt.Errorf("create bucket: %s", err)
	}
	return nil
	})
}