package db

import (
    "time"
    "encoding/binary"

    "github.com/boltdb/bolt"
)

var db *bolt.DB
var taskBucket = []byte("Tasks")

type Task struct {
    Key int
    Value string
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

func ListTasks() ([]Task, error) {
    var tasks []Task
    err := db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket(taskBucket)
        c := b.Cursor()

        for k, v := c.First(); k != nil; k, v = c.Next() {
            tasks = append(tasks,
                           Task {Key: btoi(k), Value: string(v)})
        }
        return nil
    })
        return tasks, err
}

func DeleteTask(key int) error {
    return db.Update(func(tx *bolt.Tx) error {
        b := tx.Bucket(taskBucket)
        return b.Delete(itob(key))
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
