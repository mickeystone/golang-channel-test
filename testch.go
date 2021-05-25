package main

import (
	"encoding/hex"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"sync"
)

type DBdata struct {
	key []byte
	value []byte
	obfuscateKey []byte
}

var chainstate = "./chainstate"
var dbchan = make(chan DBdata, 10)


func generate(wg *sync.WaitGroup){
	opts := &opt.Options{
		Compression: opt.NoCompression,
	}

	db, _ := leveldb.OpenFile(chainstate, opts) // You have got to dereference the pointer to get the actual value
	defer db.Close()

	iter := db.NewIterator(nil, nil)
	var obfuscateKey []byte
	n := 1
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()

		prefix := key[0]

		if prefix == 14 { // 14 = obfuscateKey
			obfuscateKey = value
		}
		if prefix == 67 {
			temp := DBdata{key, value, obfuscateKey}
			dbchan <- temp

			fmt.Println(n, hex.EncodeToString(key))
			n += 1
		}

	}
	close(dbchan)
	iter.Release()

	wg.Done()
}

func printchan(wg *sync.WaitGroup) {
	n := 1
	for x:= range dbchan {
		//fmt.Println(n, hex.EncodeToString(x.key))
		if len(x.key)>1000 {

		}
		n += 1
	}
	wg.Done()
}


func main() {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go generate(&wg)

	wg.Add(1)
	go printchan(&wg)

	wg.Wait()
}
