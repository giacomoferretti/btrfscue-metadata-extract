package main

// btrfs_timespec
// https://github.com/torvalds/linux/blob/master/include/uapi/linux/btrfs_tree.h#L817
type Timespec struct {
	Sec  uint64
	Nsec uint32
}

// btrfs_disk_key
// https://github.com/torvalds/linux/blob/master/include/uapi/linux/btrfs_tree.h#L441
type DiskKey struct {
	ObjectID uint64
	Type     uint8
	Offset   uint64
}

// btrfs_inode_item
// https://github.com/torvalds/linux/blob/master/include/uapi/linux/btrfs_tree.h#L822
type InodeItem struct {
	Generation uint64
	Transid    uint64
	Size       uint64
	Nbytes     uint64
	BlockGroup uint64
	Nlink      uint32
	Uid        uint32
	Gid        uint32
	Mode       uint32
	Rdev       uint64
	Flags      uint64
	Sequence   uint64
	Reserved   [4]uint64
	Atime      Timespec
	Ctime      Timespec
	Mtime      Timespec
	Otime      Timespec
}

// btrfs_inode_ref
// https://github.com/torvalds/linux/blob/master/include/uapi/linux/btrfs_tree.h#L803
type InodeRef struct {
	Index      uint64
	NameLength uint16
	Name       string
}

// btrfs_dir_item
// https://github.com/torvalds/linux/blob/master/include/uapi/linux/btrfs_tree.h#L855
type DirItem struct {
	Location   DiskKey
	Transid    uint64
	DataLength uint16
	NameLength uint16
	Type       uint8
	Name       string
}

// btrfs_file_extent_item
// https://github.com/torvalds/linux/blob/master/include/uapi/linux/btrfs_tree.h#L1019
type FileExtentItem struct {
	Generation    uint64
	RamBytes      uint64
	Compression   uint8
	Encryption    uint8
	OtherEncoding uint16
	Type          uint8
}
type FileExtentItemDisk struct {
	DiskBytenr   uint64
	DiskNumBytes uint64
	Offset       uint64
	NumBytes     uint64
}
