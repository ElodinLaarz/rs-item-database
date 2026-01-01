package db

import (
	"strings"

	"github.com/nutsdb/nutsdb"
	"google.golang.org/protobuf/proto"
	"rs-item-database/pb"
)

const (
	BucketName = "items"
)

type Store struct {
	db *nutsdb.DB
}

func NewStore(path string) (*Store, error) {
	opt := nutsdb.DefaultOptions
	opt.Dir = path
	db, err := nutsdb.Open(opt)
	if err != nil {
		return nil, err
	}

	// Ensure bucket exists
	err = db.Update(func(tx *nutsdb.Tx) error {
		if ok := tx.ExistBucket(nutsdb.DataStructureBTree, BucketName); !ok {
			return tx.NewBucket(nutsdb.DataStructureBTree, BucketName)
		}
		return nil
	})
	if err != nil {
		db.Close()
		return nil, err
	}

	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) SaveItem(item *pb.Item) error {
	data, err := proto.Marshal(item)
	if err != nil {
		return err
	}

	return s.db.Update(func(tx *nutsdb.Tx) error {
		// Key is lowercased name for searchability
		key := []byte(strings.ToLower(item.Name))
		return tx.Put(BucketName, key, data, 0)
	})
}

func (s *Store) GetItem(name string) (*pb.Item, error) {
	var item pb.Item
	err := s.db.View(func(tx *nutsdb.Tx) error {
		key := []byte(strings.ToLower(name))
		value, err := tx.Get(BucketName, key)
		if err != nil {
			return err
		}
		return proto.Unmarshal(value, &item)
	})
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *Store) SearchItems(prefix string, limit int) ([]*pb.Item, error) {
	var items []*pb.Item
	err := s.db.View(func(tx *nutsdb.Tx) error {
		prefixBytes := []byte(strings.ToLower(prefix))
		entries, err := tx.PrefixScan(BucketName, prefixBytes, 0, limit)
		if err != nil {
			// Treat bucket not found or similar as empty result
			if strings.Contains(err.Error(), "bucket not found") || err == nutsdb.ErrBucketNotFound {
				return nil
			}
			return err
		}

		for _, entry := range entries {
			var item pb.Item
			if err := proto.Unmarshal(entry, &item); err == nil {
				items = append(items, &item)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return items, nil
}
