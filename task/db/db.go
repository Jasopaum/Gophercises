package db

import (
    "fmt"
    "time"
    "encoding/binary"

    "github.com/boltdb/bolt"
)

var db *bolt.DB
var taskBucket = []byte("Tasks")

type Task struct {
    key int
    value string
}

func Init() error {
    var err error
    db, err = bolt.Open("tasks.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
    if err != nil {
        return err
    }
    return db.Update(func(tx *bolt.Tx) error {
        _, err := tx.CreateBucketIfNotExists(taskBucket)
        return err
    })
}

func CreateTask(task string) error {
    return db.Update(func(tx *bolt.Tx) error {
        b := tx.Bucket(taskBucket)
        id, _ := b.NextSequence()

        return b.Put(itob(int(id)), []byte(task))
    })
}

func ListTasks() error {
    return db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket(taskBucket)
        c := b.Cursor()

        for k, v := c.First(); k != nil; k, v = c.Next() {
            fmt.Printf("key=%d, value=%s\n", btoi(k), string(v))
        }

        return nil
    })
}

func itob(v int) []byte {
    b := make([]byte, 8)
    binary.BigEndian.PutUint64(b, uint64(v))
    return b
}

func btoi(b []byte) int {
    return int(binary.BigEndian.Uint64(b))
}
