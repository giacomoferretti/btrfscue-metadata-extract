package main

import (
	"encoding/binary"
)

const (
	KeyLen     = 8 + 1 + 8
	itemKey    = 0
	itemOffset = itemKey + KeyLen
	itemSize   = itemOffset + 4
	ItemLen    = itemSize + 4

	CSumSize             = 32
	LabelSize            = 256
	SystemChunkArraySize = 2048

	UUIDSize = 16

	headerCSum = 0
	// The following three fields must match struct SuperBlock
	headerFSID   = headerCSum + CSumSize
	headerByteNr = headerFSID + UUIDSize
	headerFlags  = headerByteNr + 8 // Includes 1 byte backref rev.
	// Allowed to be different from SuperBlock from here on
	headerChunkTreeUUID = headerFlags + 8
	headerGeneration    = headerChunkTreeUUID + UUIDSize
	headerOwner         = headerGeneration + 8
	headerNrItems       = headerOwner + 8
	headerLevel         = headerNrItems + 4
	HeaderLen           = headerLevel + 1

	keyV2Owner      = 0
	keyV2Type       = keyV2Owner + 8
	keyV2ObjectID   = keyV2Type + 1
	keyV2Offset     = keyV2ObjectID + 8
	keyV2Generation = keyV2Offset + 8
	keyV2End        = keyV2Generation + 8
)

type KeyV2 struct {
	Owner      uint64
	Type       uint8
	ObjectID   uint64
	Offset     uint64
	Generation uint64
}

type Key struct {
	ObjectID uint64
	Type     uint8
	Offset   uint64
}

type Item struct {
	Key    Key
	Offset uint32
	Size   uint32
}

func BtrfscueParseDbKey(input []byte) KeyV2 {
	return KeyV2{
		Owner:      binary.BigEndian.Uint64(input[keyV2Owner:]),
		Type:       input[keyV2Type],
		ObjectID:   binary.BigEndian.Uint64(input[keyV2ObjectID:]),
		Offset:     binary.BigEndian.Uint64(input[keyV2Offset:]),
		Generation: binary.BigEndian.Uint64(input[keyV2Generation:]),
	}
}

func BtrfscueParseDbValue(input []byte) (Item, []byte) {
	k := Key{
		ObjectID: binary.LittleEndian.Uint64(input[:8]),
		Type:     uint8(input[8]),
		Offset:   binary.LittleEndian.Uint64(input[9:]),
	}

	i := Item{
		Key:    k,
		Offset: binary.LittleEndian.Uint32(input[itemOffset:]),
		Size:   binary.LittleEndian.Uint32(input[itemSize:]),
	}

	if i.Key.Type == 0 {
		return i, []byte{}
	} else {
		return i, input[ItemLen : ItemLen+binary.LittleEndian.Uint32(input[itemSize:])]
	}
}
