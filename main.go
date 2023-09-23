package main

import (
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"

	bolt "go.etcd.io/bbolt"
)

type Inode struct {
	Generation  uint64
	Inode       uint64
	ParentInode uint64
	Name        string
	IsFile      bool
	FullPath    string
	Size        uint64
}

var inodeData []Inode

func createDir(inode Inode) {
	f := path.Join(os.Args[1], inode.FullPath)
	if !inode.IsFile {
		os.MkdirAll(f, os.ModePerm)
	}
}

func createFile(inode Inode, data []byte) {
	f := path.Join(os.Args[1], inode.FullPath)
	log.Printf("[GEN:%v] Writing %v bytes to %v...", inode.Generation, inode.Size, f)
	if inode.IsFile {
		if err := os.WriteFile(f, data, 0644); err != nil {
			log.Fatal(err)
		}
	}
}

func readFileFromDisk(offset, length uint64) []byte {
	fmt.Fprintf(os.Stderr, "Reading offset %v and length %v\n", offset, length)
	f, err := os.Open(os.Args[4])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = f.Seek(int64(offset), io.SeekStart)
	if err != nil {
		log.Fatal(err)
	}

	ret := make([]byte, length)
	_, err = f.Read(ret)
	if err != nil {
		log.Fatal(err)
	}

	return ret
}

func findInode(inode uint64) (Inode, error) {
	for _, v := range inodeData {
		if v.Inode == inode {
			return v, nil
		}
	}

	return Inode{}, errors.New("inode not found")
}

func dbProcessEntry(key, value []byte) error {
	btrfsKeyV2 := BtrfscueParseDbKey(key)
	btrfsKey, data := BtrfscueParseDbValue(value)

	// Print data
	dataHex := hex.EncodeToString(data)
	dataSanitized := string(data)
	dataSanitized = strings.Replace(dataSanitized, "\n", "\\x0a", -1)
	dataSanitized = strings.Replace(dataSanitized, "\b", "\\x0b", -1)
	dataPrint := ""
	if btrfsKey.Size > 0 {
		dataPrint = fmt.Sprintf(" %s, %s", dataSanitized, dataHex)
	}

	if btrfsKeyV2.Type == 108 {
		fmt.Printf("[G:%v O:%v ID:%v OF1:%v T:%v OF2:%v S:%v]%s\n", btrfsKeyV2.Generation, btrfsKeyV2.Owner, btrfsKeyV2.ObjectID, btrfsKeyV2.Offset, btrfsKeyV2.Type, btrfsKey.Offset, btrfsKey.Size, dataPrint)

		// BTRFS_EXTENT_DATA_KEY
		inode := btrfsKeyV2.ObjectID
		offset := btrfsKeyV2.Offset
		file_extent_item := ParseFileExtentItem(data)

		targetInode, err := findInode(inode)
		if err != nil {
			log.Fatal(err)
		}

		// Inline data
		if file_extent_item.Type == 0 {
			// Check length
			if file_extent_item.RamBytes != uint64(len(data[21:])) {
				log.Fatal("WRONG LENGTH ON EXTENT_DATA")
			}

			if btrfsKeyV2.Generation == targetInode.Generation {
				createFile(targetInode, data[21:])
			}
		} else if file_extent_item.Type == 1 {
			file_extent_item_disk := ParseFileExtentItemDisk(data[21:])
			fmt.Printf(" └─ FILE_EXTENT_ITEM = %v\n", file_extent_item_disk)

			if btrfsKeyV2.Generation == targetInode.Generation {
				createFile(targetInode, readFileFromDisk(file_extent_item_disk.DiskBytenr, targetInode.Size))
			}
		}
		fmt.Printf(" └─ BTRFS_EXTENT_DATA_KEY - INODE:%v OFFSET:%v = %v\n", inode, offset, file_extent_item)
	}

	return nil
}

func main() {
	// Check arguments
	if len(os.Args) != 5 {
		fmt.Fprintf(os.Stderr, "Usage: %v <output_folder> <output.bin> <metadata.db> <disk.img>\n", os.Args[0])
		os.Exit(1)
	}

	// GOB DECODE
	dataFile, err := os.Open(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	dataDecoder := gob.NewDecoder(dataFile)
	err = dataDecoder.Decode(&inodeData)
	if err != nil {
		log.Fatal(err)
	}
	dataFile.Close()

	// Create folder
	for _, v := range inodeData {
		if !v.IsFile {
			createDir(v)
		}
	}

	// Open bbolt database
	db, err := bolt.Open(os.Args[3], 0600, &bolt.Options{ReadOnly: true})
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	// Process all entries on "index"
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("index"))

		b.ForEach(dbProcessEntry)
		return nil
	})
	if err != nil {
		log.Fatalln(err)
	}
}
