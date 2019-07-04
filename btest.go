package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/pprof"
	"time"

	"github.com/dgraph-io/badger"
)

func open() *badger.DB {
	// Open the Badger database located in the /tmp/badger directory.
	// It will be created if it doesn't exist.
	opts := badger.DefaultOptions("/tmp/badger")
	db, err := badger.OpenManaged(opts)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func main() {

	db := open()

	t := uint64(time.Now().UnixNano())
	t2 := t - 10000

	// Start a writable transaction.

	log.Println("Putting 42")

	txn := db.NewTransactionAt(t, true)
	defer txn.Discard()
	txn.Delete([]byte("answer"))

	txn.Set([]byte("answer"), []byte("42"))

	// Commit the transaction and check for error.
	err := txn.CommitAt(t2, nil)
	if err != nil {
		log.Printf("Error committing %v\n", err)
		return
	}

	log.Println("Putting 43")
	txn = db.NewTransactionAt(t, true)
	defer txn.Discard()

	// Use the transaction...
	txn.Set([]byte("answer"), []byte("43"))
	// Commit the transaction and check for error.
	if err := txn.CommitAt(t, nil); err != nil {
		log.Printf("Error %v\n", err)
		return
	}

	log.Println("Reading 42")

	txn = db.NewTransactionAt(t2+500, false)
	defer txn.Discard()

	// try to get the 42 value
	item, err := txn.Get([]byte("answer"))
	if err != nil {
		log.Printf("Error %v\n", err)
		return
	}

	valCopy, err := item.ValueCopy(nil)
	if err != nil {
		log.Printf("Error %v\n", err)
		return
	}
	fmt.Printf("The answer which should be 42 is: %s\n", valCopy)

	txn.Discard()

	log.Println("Reading 43")

	txn = db.NewTransactionAt(t+500, false)

	// try to get the 43 value
	item, err = txn.Get([]byte("answer"))
	if err != nil {
		log.Printf("Error %v\n", err)
		return
	}

	valCopy, err = item.ValueCopy(nil)
	if err != nil {
		log.Printf("Error %v\n", err)
		return
	}
	fmt.Printf("The answer which should be 43 is: %s\n", valCopy)

	txn.Discard()
	txn = db.NewTransactionAt(t+500, false)

	fmt.Println("Logging the history of 'answer'")

	itr := txn.NewKeyIterator([]byte("answer"), badger.IteratorOptions{
		PrefetchValues: true,
		PrefetchSize:   100,
		Reverse:        false,
		AllVersions:    true,
	})
	for itr.Rewind(); itr.Valid(); itr.Next() {
		item := itr.Item()
		valCopy, err = item.ValueCopy(nil)
		log.Printf("val: %s\n", valCopy)
	}
	itr.Close()
    txn.Discard()

	f, err := os.Create("./btest.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	// stop the server without closing the db
	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		// db.Close()
		w.WriteHeader(200)
		pprof.StopCPUProfile()
		os.Exit(0)
	})

	err = http.ListenAndServe(":9999", nil)
	log.Fatal(err)

}
