package main

import "encoding/binary"

func ParseTimespec(data []byte) Timespec {
	return Timespec{
		Sec:  binary.LittleEndian.Uint64(data[:]),
		Nsec: binary.LittleEndian.Uint32(data[8:]),
	}
}

func ParseDiskKey(data []byte) DiskKey {
	return DiskKey{
		ObjectID: binary.LittleEndian.Uint64(data[:8]),
		Type:     uint8(data[8]),
		Offset:   binary.LittleEndian.Uint64(data[9:]),
	}
}

// BTRFS_INODE_ITEM_KEY = 1
// https://github.com/torvalds/linux/blob/master/include/uapi/linux/btrfs_tree.h#L137
func ParseInodeItem(data []byte) InodeItem {
	reserved := [4]uint64{}
	for i := 0; i < 4; i++ {
		reserved[i] = binary.LittleEndian.Uint64(data[80+i*8:])
	}

	return InodeItem{
		Generation: binary.LittleEndian.Uint64(data[:]),
		Transid:    binary.LittleEndian.Uint64(data[8:]),
		Size:       binary.LittleEndian.Uint64(data[16:]),
		Nbytes:     binary.LittleEndian.Uint64(data[24:]),
		BlockGroup: binary.LittleEndian.Uint64(data[32:]),
		Nlink:      binary.LittleEndian.Uint32(data[40:]),
		Uid:        binary.LittleEndian.Uint32(data[44:]),
		Gid:        binary.LittleEndian.Uint32(data[48:]),
		Mode:       binary.LittleEndian.Uint32(data[52:]),
		Rdev:       binary.LittleEndian.Uint64(data[56:]),
		Flags:      binary.LittleEndian.Uint64(data[64:]),
		Sequence:   binary.LittleEndian.Uint64(data[72:]),
		Reserved:   reserved,
		Atime:      ParseTimespec(data[112:]),
		Ctime:      ParseTimespec(data[124:]),
		Mtime:      ParseTimespec(data[136:]),
		Otime:      ParseTimespec(data[148:]),
	}
}

// BTRFS_INODE_REF_KEY = 12
// https://github.com/torvalds/linux/blob/master/include/uapi/linux/btrfs_tree.h#L138
func ParseInodeRef(data []byte) InodeRef {
	nameLength := binary.LittleEndian.Uint16(data[8:])

	return InodeRef{
		Index:      binary.LittleEndian.Uint64(data[:]),
		NameLength: nameLength,
		Name:       string(data[10 : 10+nameLength]),
	}
}

// BTRFS_XATTR_ITEM_KEY = 24
// https://github.com/torvalds/linux/blob/master/include/uapi/linux/btrfs_tree.h#L140
func ParseDirItem(data []byte) DirItem {
	nameLength := binary.LittleEndian.Uint16(data[27:])
	return DirItem{
		Location:   ParseDiskKey(data),
		Transid:    binary.LittleEndian.Uint64(data[17:]),
		DataLength: binary.LittleEndian.Uint16(data[25:]),
		NameLength: nameLength,
		Type:       uint8(data[29]),
		Name:       string(data[30 : 30+nameLength]),
	}
}

// BTRFS_DIR_ITEM_KEY = 84
// https://github.com/torvalds/linux/blob/master/include/uapi/linux/btrfs_tree.h#L175

// BTRFS_DIR_INDEX_KEY = 96
// https://github.com/torvalds/linux/blob/master/include/uapi/linux/btrfs_tree.h#L176

// BTRFS_EXTENT_DATA_KEY = 108
// https://github.com/torvalds/linux/blob/master/include/uapi/linux/btrfs_tree.h#L180
func ParseFileExtentItem(data []byte) FileExtentItem {
	return FileExtentItem{
		Generation:    binary.LittleEndian.Uint64(data[:]),
		RamBytes:      binary.LittleEndian.Uint64(data[8:]),
		Compression:   data[16],
		Encryption:    data[17],
		OtherEncoding: binary.LittleEndian.Uint16(data[18:]),
		Type:          data[20],
	}
}
func ParseFileExtentItemDisk(data []byte) FileExtentItemDisk {
	return FileExtentItemDisk{
		DiskBytenr:   binary.LittleEndian.Uint64(data[:]),
		DiskNumBytes: binary.LittleEndian.Uint64(data[8:]),
		Offset:       binary.LittleEndian.Uint64(data[16:]),
		NumBytes:     binary.LittleEndian.Uint64(data[24:]),
	}
}
