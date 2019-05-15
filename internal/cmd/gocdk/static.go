// Code generated by vfsgen; DO NOT EDIT.

package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	pathpkg "path"
	"time"
)

// static statically implements the virtual filesystem provided to vfsgen.
var static = func() http.FileSystem {
	fs := vfsgen۰FS{
		"/": &vfsgen۰DirInfo{
			name:    "/",
			modTime: time.Time{},
		},
		"/demo": &vfsgen۰DirInfo{
			name:    "demo",
			modTime: time.Time{},
		},
		"/demo/blob": &vfsgen۰DirInfo{
			name:    "blob",
			modTime: time.Time{},
		},
		"/demo/blob/demo_blob.go": &vfsgen۰CompressedFileInfo{
			name:             "demo_blob.go",
			modTime:          time.Time{},
			uncompressedSize: 7171,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xd4\x58\xdb\x6e\xdb\x48\xd2\xbe\x16\x9f\xa2\x42\xe0\xcf\x90\xf3\x2b\xe4\x5c\xec\xc5\x42\x2b\x29\x88\x0f\x99\x31\xe2\x89\x83\xd8\x99\xc1\x62\x30\x48\x5a\x64\x51\xea\x88\xec\x66\xba\x9b\x92\x05\x41\xef\xbe\xa8\x6e\x92\xa2\x0e\x76\xe4\x78\x6f\x56\x17\x36\x0f\xdd\x5f\x55\xd7\xe1\xab\x2a\x96\x2c\x99\xb3\x29\x42\xc1\xb8\xf0\x3c\x5e\x94\x52\x19\x08\xbc\x9e\x9f\x48\x61\xf0\xde\xf8\x5e\xcf\x47\xa5\xa4\xd2\x74\x95\x15\xf6\xc1\xcc\x14\x79\x6c\xb0\x28\x73\x66\x90\x1e\x70\x49\x7f\x05\x9a\x78\x66\x4c\x49\xd7\x52\xfb\x9e\xd7\xf3\xa7\x32\xc9\x65\x95\x46\x29\x2e\xe2\x49\x2e\x27\xbe\xd7\xfb\x0c\x07\x4f\xe3\x8c\xe7\xf8\xc8\xeb\x02\x0b\xf7\x36\xf4\xbc\x38\x86\xbb\x9b\x8b\x9b\x40\x2d\x98\x98\xa2\x30\xe1\x00\xee\x66\x5c\x03\x41\x00\xd7\x50\x69\x54\xaf\x16\x5c\xf3\x49\x8e\x7d\x60\x69\x0a\x05\x13\x2b\x48\x64\x51\xa0\x30\x1a\xf0\xbe\xcc\x19\x17\x5c\x4c\x09\x6a\x26\x97\xc0\x0d\x2c\xa5\x9a\xeb\xc8\xf3\xb2\x4a\x24\xc0\x05\x37\x41\x08\x6b\xaf\x47\x87\x89\x7e\x63\x22\xcd\xf1\x6d\x25\x92\xc0\x8f\x53\x2c\xa4\xd3\xc9\xef\x03\xfd\x3f\x63\x1a\xdd\x0a\x15\x3e\xbe\x21\xe7\xda\xd4\x9b\xae\xb9\x36\xa7\x6d\x5a\x70\x5c\xd6\x9b\xfe\xe0\xb8\x3c\x6d\xd3\x52\x71\x83\xf5\xae\x3f\xe9\xba\xdd\xb6\xf1\xbc\x05\x53\x30\xa9\x92\x39\x9a\x4f\x1f\xaf\x41\x1b\x45\x96\xd8\x3e\x84\x9f\x69\x57\x74\x66\x6f\x3a\xcf\x2f\x95\x02\x1b\x07\x07\x46\x3a\xe2\x90\xb3\xeb\x9b\xb3\xcf\x67\x9f\xce\xdf\x5d\xde\x7d\xb6\x52\x66\xb2\xca\x53\x98\x20\x70\x01\xc8\x92\x19\x4c\xb8\x2c\xf0\x27\x0d\x89\x14\x19\x9f\x82\xae\x92\x99\x45\x32\x33\x66\xc8\x21\x09\x13\xb4\x5c\x21\x4b\x61\x86\x0a\x23\xaf\xb7\x55\x7a\x04\x52\x47\xbf\xa2\x41\xb1\x08\xfc\x3d\x59\x7e\xe8\xf5\x78\xd6\x39\xe1\x68\x04\xbe\x4f\x7a\xee\x00\xf8\x05\x16\x83\x38\xf6\xbd\xde\xa6\x41\xee\x77\x4e\x3a\xb2\xb6\x8b\x6e\x4a\x14\xce\x12\x41\x9d\x0e\xd1\x19\x4b\xe6\x53\x25\x2b\x91\x06\x61\x7f\x2b\xc6\x9a\xf6\xa1\xd0\xe4\x1a\x4a\x85\xc6\xac\x40\xb1\x25\xfc\x76\xf7\xfb\x75\x04\xb7\xce\x24\x4b\x84\x19\x5b\x20\x30\x1b\xa0\x52\x80\x36\xab\x1c\x41\xcf\x10\x4d\x8c\x26\x89\x20\x93\x0a\xc8\xb5\xfa\xb5\xe7\x99\x55\x89\x6d\xd4\x5d\x30\xc3\xc8\x7f\x55\x62\xe8\x78\x1d\x6f\xf6\xb6\xce\xda\x74\x36\x51\xd4\x1d\xd9\xd4\xfc\xba\x9b\x9b\x9f\x03\xe9\xd1\xce\x9b\xc9\x57\x4c\x8c\x86\xbf\xfe\x76\x11\xb2\x7d\xb6\x23\x84\xa2\xf4\xc7\x84\xbc\xc3\xd5\xc1\xb2\x27\x08\xb6\x81\xfe\x88\xe4\xe3\xb2\x8f\x4a\xdf\x2e\xb5\xa0\xe7\xe4\x7a\xa2\x8e\x9d\xa7\xb7\x55\x92\xa0\xd6\x00\x13\x29\x73\xd2\x24\x91\x42\x5b\xea\x24\x6d\xee\x6a\x72\xfc\xa0\x30\xe3\xf7\x30\x82\x2f\xde\xf0\xc5\xc5\xcd\xf9\xdd\xbf\x3f\x5c\x02\xd1\xe7\xd8\x1b\x36\xff\x90\xa5\x63\x0f\x60\x58\xa0\x61\x90\xcc\x98\xd2\x68\x46\xfe\xa7\xbb\xb7\xaf\xfe\xe9\xdb\x17\x86\x9b\x1c\xc7\xfb\xac\x68\xa3\x62\x18\xbb\x97\xde\x30\x76\x38\xc3\x89\x4c\x57\x76\x57\x49\x7f\xc1\x45\x60\x49\x14\x4f\xeb\x85\x36\x8a\x19\xd4\x60\x66\x48\x54\x09\x32\x83\x5f\x25\x9c\x5f\xbc\xfb\x49\xc3\x90\xc1\x4c\x61\x36\xf2\x89\x5c\xf4\x20\x8e\xa7\x32\x95\x49\x24\xd5\x34\x3e\xe0\xf1\x31\xfd\x1d\xc6\x6c\x0c\x75\x05\x89\x48\x66\x5c\x76\x45\x5f\x19\x0a\xfd\xa4\x52\x0a\x85\xc9\x57\x50\x69\x2e\xa6\xc0\xa0\xc3\x30\x30\x61\x1a\x53\x90\xc2\x2a\x44\xfe\xf2\xd7\x6b\x88\xe8\x62\xb3\xf1\xfb\xb0\x9c\xf1\x64\x66\xc1\x6a\x46\x70\x84\x51\x29\x4c\x61\xc1\x99\xdd\x85\x62\xc1\x95\x14\x44\xef\xb0\x60\x8a\xb3\x49\x8e\x70\xc0\x0b\x47\xf4\xbb\x45\x3c\x76\xe6\xed\x49\x13\x29\x12\x2c\x8d\x8e\x2b\x95\xeb\xd8\x1f\x13\x13\xd9\x33\x53\x52\x16\x52\xa1\x85\xe1\x22\x93\xaa\x60\x86\x4b\x01\x6c\x22\x2b\x43\xe7\xd0\x44\x75\xce\xb2\xf0\xe6\xc3\x95\xee\x8a\xaf\x72\x27\x7f\x98\xf3\x71\x2b\x3f\x72\x05\x62\x4c\xe1\x6d\x65\xd0\xd1\x92\x26\xf4\x64\x66\xef\x1d\xdf\x0c\xe3\x9c\x1f\x45\xb0\xd5\x62\x4c\x29\x78\x14\x81\x81\x2e\x31\xe1\x19\x4f\xac\x0b\x48\xc3\x13\x40\x5d\x35\x19\xdb\xa0\xb7\xb0\x0c\x04\x2e\x1b\x04\x23\x0f\x30\x86\x31\x1d\x70\xbd\xe6\x19\x44\x97\x4a\x6d\x36\xce\xe4\x43\x6d\x94\x14\xd3\xf1\x7a\xed\x9e\x0e\xe3\xfa\x81\x35\xcb\x7a\x8d\x22\xdd\x6c\xbe\x78\xbb\xf9\x73\x5b\x65\x6d\xfe\xc4\x2e\xb4\x87\xb1\x4d\x1c\x5a\x19\xc7\x70\x25\xca\xca\x0c\x5c\xd5\x6a\x58\x31\x72\x18\x74\xdb\xe0\xd4\x8c\xbe\x97\x96\xff\x0f\x87\xb2\x0e\x61\x1b\xde\xac\x61\xe9\xf6\x7b\xb0\x5f\x3c\x80\xf5\x5a\x51\x19\x80\x0e\x61\x69\x6b\x0b\x80\x61\xca\x17\xce\xd6\xb4\x8c\xec\x74\xa5\x2f\xb8\xaa\xdf\xda\x15\xbb\x61\xf1\xba\xb4\xc8\x23\xca\x0d\x62\xaa\xcd\xc6\x1f\x6f\xaf\xc9\x2b\x2d\x1a\xe6\x1a\x8f\x03\x51\x74\xbc\x9e\xe3\xea\x24\x14\x72\x86\xd3\x35\xae\x95\x6d\x3c\x74\xa2\xd5\x9a\x42\x50\x5b\x8d\x6e\x4f\xb3\x1a\x99\xe3\x88\xc9\x28\xc7\xea\xf8\x2c\xc7\xc3\x9c\x4d\x30\x6f\xd4\x3d\x9f\x49\xa9\xb1\x26\x16\x30\x12\xe8\xa4\x83\xfa\xe5\x50\x63\x8e\x89\x01\xc1\x0a\x1c\xf9\x73\x5c\xf9\xe3\xd6\x38\x8f\xb8\xa8\xde\x2c\x4b\x9b\xd6\x0b\x96\x57\x38\xf2\xd7\x6b\xb2\x95\x33\x9b\xbb\x1a\xc6\x6e\x45\x17\x73\x6b\x3b\xb2\x9e\x13\x5f\x6b\x1e\x3b\xbd\x6b\x22\x00\x18\x72\xb2\x18\x50\x11\x1b\xf9\xba\x9a\x14\xdc\xf8\xcd\xd2\xe6\xc0\x4f\xb4\x7b\x5b\x07\x6b\xc3\xdb\xfb\xd3\x2d\xdf\xad\x6d\xf5\x29\xfe\x54\xd2\x20\x70\xf3\xc2\xdb\x8b\xaf\xc7\x7c\x72\x46\xae\x98\xe3\x8a\xdc\x61\x19\x84\x2e\x02\xea\xc4\x4b\x85\x0b\x2e\x2b\xed\xbc\xb5\xe4\x79\x4e\xbc\x2e\x17\xa8\x68\x9d\x41\x11\xb6\xae\x9b\xa8\xb8\xc1\xdb\xb1\x94\x1d\x4c\x3a\x1e\xdd\xfa\x67\x1b\xd8\x0f\x18\x7c\x5f\x4f\xdb\x31\x34\x3a\x1e\x15\x4c\xc2\x98\x42\x56\xcb\x6b\x08\xd5\x07\x25\x97\x7a\xe4\xff\xc3\x87\x44\xe6\x74\xf1\x8b\xcb\xa6\xdd\x9e\x81\x42\xa4\x41\x78\x42\x10\x34\x27\xb2\x60\x70\x65\x5e\x3c\x31\x2c\x42\xd7\xe8\x07\x1d\x26\x2c\xca\x1c\x60\x04\xcd\xd8\x16\xfd\x5e\x69\x13\xb4\x77\xef\x71\x19\xf8\xd6\x23\x54\x91\xfd\x30\xfa\x40\x7d\x48\xb0\xcf\xa3\x61\xd8\x21\xc1\x53\x11\x6d\x61\xeb\x22\x76\x29\xb4\x41\xb4\x04\x71\x2a\xa2\x2d\x74\x5d\xc4\x2e\xbd\x34\x88\x2e\xf2\x09\xf2\x04\x44\x57\xe5\xba\x90\x3b\x89\x13\x86\x64\x52\x3b\xf9\xec\x4d\x7d\xc1\x12\xec\x2c\xf6\x11\x75\x29\x85\x46\xbb\x4d\xf5\x41\xe1\x37\xf8\xb9\x7e\xf3\xad\x42\x6d\xec\xb0\xe4\x5c\x3d\x18\xc1\xcb\x6e\xc1\x5a\x7f\xfa\x78\x3d\xd8\x0e\x13\x1b\x3b\xc2\xa0\x52\xb4\xb0\xeb\xbf\xe8\xf2\x1e\x93\xca\x60\xb0\xec\x83\x45\x0a\xff\x65\x97\xbd\x18\x81\xe0\xb9\x1d\x72\xac\xc4\x4b\xea\x69\x69\x11\x2a\x55\xdf\x84\x7d\xa7\xe6\xad\x61\xa6\xd2\x57\xc2\xa0\x12\x2c\xbf\x45\xb5\x40\x65\x57\x84\x34\x0a\xb9\x21\x66\x6f\x44\xb5\xfe\x73\x1d\x23\x37\x58\xd8\xce\x86\x41\x33\x34\x95\x52\xd3\xac\xbd\x82\x4a\xa4\xa8\x80\x81\xef\x6a\x95\x4f\x48\xdf\x2a\x54\x2b\x28\x99\x62\x05\x1a\x54\x11\x5c\xd2\xf0\x47\x78\x98\x42\xca\x15\x26\x46\xaa\x15\x75\x89\x0c\x72\x2e\xe6\x94\x89\xf4\xd6\xcd\x81\xed\x82\x3e\x61\x31\x91\xba\xd9\x51\x48\xf1\xea\x81\xbd\x14\x17\x6e\x6f\xc6\x73\x8c\xb6\x0e\xeb\x1c\xe7\x39\x0e\x6b\x5a\x81\x43\x87\xa5\x98\xa1\x02\x92\xe7\xa6\xe2\x3d\x0f\x36\xf9\x72\x82\x07\x9f\xe5\x42\x1a\x67\x37\x41\xe8\x75\x86\xe0\xcb\x5d\x78\x2b\x35\xaa\x87\xdc\x66\x81\xd7\xeb\x29\x34\x95\x12\x76\x1e\x96\xa5\xd1\xed\x91\x5d\x6d\xb4\x45\x4e\xd3\xfe\x0b\xcc\x79\x41\x16\x1b\x80\x1f\xfb\x7d\xaf\xd7\x73\x55\x64\x40\xf4\xa4\xf0\x5b\xf4\x56\xaa\xe2\x0f\x62\xaf\xa0\x89\x84\xb0\x6f\x61\x69\x93\xb5\x87\x95\x6a\x71\x03\x12\x15\x7a\x3d\xea\xa7\x09\x5c\x4e\xbe\xf6\x1b\xb3\xd1\xf2\xe8\x3d\xde\x9b\x80\x50\xcf\xdd\xfc\x1d\x50\x72\x37\xb6\x1d\x8d\x80\xcb\xe8\xf2\xe6\xad\xb3\xdb\x44\x21\x9b\x3b\x13\x34\x2b\xba\x66\xed\x1e\x3c\x2b\x8c\x33\x6a\x16\xf8\x6f\x19\xcf\x31\xa5\xe8\x21\x89\xcc\x55\x2a\x81\xf7\x66\x67\x4e\x99\xe3\x6a\x00\xff\xb7\xf0\xad\x7a\xa4\x43\x6b\x30\x27\xcf\x82\x77\xc7\xd5\x11\xb0\xb2\x44\x91\x06\x07\xaf\xfa\x20\x27\x5f\x43\x67\x92\x0c\x72\x14\x87\x4b\x42\x3a\xdb\x2f\x07\xfe\x72\x1f\xe2\x1c\x6b\x09\x69\xf5\xb3\xd9\xe8\x2c\xea\x37\x29\xdc\x46\x7d\xe7\x93\xd1\x73\xa2\xbe\x69\xe5\x48\x9d\xdd\xc0\x27\xf7\xbf\x23\xcb\xec\x39\x9e\x8a\x72\xed\x75\x3d\xe7\x65\xdb\x7e\x0c\x46\x90\xb1\x5c\xe3\xf1\x7c\x79\xb1\xb3\xd6\x39\x6d\x27\x8b\x9a\x1a\x71\x4a\x16\x3d\x2f\x8d\xac\x53\xff\x1b\xb9\xc4\x33\xa7\xa1\xed\x4b\xb6\xdf\xa2\xe2\x18\xde\x4b\xdb\x1f\xb9\xfe\x10\xd3\x08\x3e\x62\xcd\xa0\x54\xe2\x61\xc9\xcd\x0c\x18\xa4\x4a\x96\xa9\x5c\x0a\x0a\xca\xc4\xf5\xb9\x52\x60\x44\xc2\x8f\xa4\x93\xe0\x39\x69\xdf\xa4\xd3\xa9\xf9\xf4\x40\x42\xb5\x19\x65\xad\x71\x2c\xa7\x1e\x4c\xaa\xec\x07\x93\x6a\x9b\x55\xb5\xcc\x1f\xc8\xab\x86\x00\xbe\x97\x59\x4f\x4c\x2d\x1b\x0e\x40\xed\x6f\xe3\xc1\x37\xd6\x81\x4b\xa6\xa1\x54\x72\xc1\x53\x72\xe2\x85\x5c\x8a\x5c\xb2\xd4\x4d\xc5\xd4\x59\x90\x33\x6c\x49\x9a\xe3\x8a\xfc\xb6\x13\xe4\x23\x30\xaa\x42\x1b\x32\x2c\xa5\x7c\x6c\xa2\xdd\x39\xf5\x3d\x2e\x3f\xda\x17\xbb\x3e\xeb\x6f\x63\xaa\x0f\xb5\xcf\x9f\xc0\x78\x5b\xe7\x24\x0a\x49\x0b\x27\xe3\x31\x7e\x73\xe9\xea\x94\x8c\xce\x73\xa9\x31\xb0\x42\x65\x74\x2e\xcb\x15\x25\x98\x7b\x17\x3a\xcb\xec\x7f\x06\x3d\x97\x42\x73\x0a\x6e\x8d\xc6\x70\x31\x85\xba\x3f\x7e\x75\xb7\x2a\xb1\xdf\xde\x5d\xa3\x98\x9a\x19\xcc\x2c\x94\x8e\xf6\xe9\xac\xfb\x2d\xfb\x39\x7c\xd6\x8e\x48\x2d\xa1\x6d\x7f\x87\xd4\xb6\xfd\x3d\x40\x72\xbb\x5f\x09\x0f\xb8\xb0\x1d\x18\x6a\x42\xfc\x4e\xa7\xd0\x76\xad\xff\x13\xad\xc2\x11\x7a\x7b\xf9\xb2\x7e\xb4\x3b\x08\x6d\xb9\xef\xfb\xec\xf8\x50\x62\x22\x1d\x04\x98\xed\x00\xb1\x28\xcd\x6a\x7f\xc2\xb4\x79\x7a\x0c\xff\x21\x5d\x1e\x97\xa4\x65\xd1\x7e\x41\x6b\xc5\x1c\x91\xb1\x9b\xb4\x56\xd6\x9b\x3c\x7f\x24\x67\xff\xfa\x7b\xb2\x32\x18\x1c\xd1\x2d\x74\x09\x7d\xe0\xea\x87\x72\xd9\x9d\xdc\x65\xf4\x4e\x02\x77\xa8\xaa\x23\xa6\xf9\x6a\xdd\xf0\x0e\x25\xd8\x7f\x02\x00\x00\xff\xff\x31\x3d\x36\x53\x03\x1c\x00\x00"),
		},
		"/demo/runtimevar": &vfsgen۰DirInfo{
			name:    "runtimevar",
			modTime: time.Time{},
		},
		"/demo/runtimevar/demo_runtimevar.go": &vfsgen۰CompressedFileInfo{
			name:             "demo_runtimevar.go",
			modTime:          time.Time{},
			uncompressedSize: 2284,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x84\x55\x41\x6f\xdb\x38\x13\x3d\x5b\xbf\x62\xca\x43\x3f\xa9\x70\xc5\xef\xb0\x0b\x2c\xbc\x92\x17\x69\xe3\x76\x8d\x4d\x9b\xc2\x71\xb2\xd8\x53\x4b\x4b\x63\x9b\xa8\x44\xaa\xc3\x91\xdd\x20\xf0\x7f\x5f\x90\xb2\x62\xb9\x4d\xb0\x3a\xd8\x12\x39\xf3\xe6\xcd\x70\xe6\xb1\x51\xc5\x57\xb5\x41\xa8\x95\x36\x51\xa4\xeb\xc6\x12\x43\x1c\x8d\x44\x61\x0d\xe3\x77\x16\xd1\x48\x6c\xb9\xae\x24\x63\xdd\x54\x8a\xd1\x2f\x18\x64\xb9\x65\x6e\xfc\xbb\x75\x22\x8a\x46\x62\x63\x8b\xca\xb6\x65\x5a\xe2\x4e\x52\x6b\x58\xd7\xb8\x53\x24\xa2\xd1\x67\x78\x66\x4f\x16\xd6\x38\x56\x86\xff\xcb\x6e\xad\xab\x0e\x2b\x89\x22\x29\x61\x79\x7d\x79\x1d\xd3\x4e\x99\x0d\x1a\x4e\x26\xb0\xdc\x6a\x07\xde\x06\xb4\x83\xd6\x21\xbd\xde\x69\xa7\x57\x15\x8e\x41\x95\x25\xd4\xca\xdc\x43\x61\xeb\x1a\x0d\x3b\xc0\xef\x4d\xa5\xb4\xd1\x66\xe3\xa1\xb6\x76\x0f\x9a\x61\x6f\xe9\xab\x4b\xa3\x68\xdd\x9a\x02\xb4\xd1\x1c\x27\xf0\x10\x8d\x7c\x86\xe9\x9f\xca\x94\x15\xbe\x6b\x4d\x11\x0b\x59\x62\x6d\x87\xcc\xc4\x18\x4e\x5f\x9d\x25\x25\xd1\x21\x8a\x76\x8a\x60\xa7\x48\xab\x55\x85\xb7\x8b\x2b\x70\x4c\x3e\xe4\x70\x19\x5e\x9d\x5c\xd3\xbb\xe3\xe2\x99\xc5\x8c\x08\x90\xc8\xd2\x4f\xcc\x86\xd0\x39\x58\x97\xbe\x47\x46\xb3\x8b\xc5\xe2\xf6\xe3\x72\xfe\x61\x76\x77\xb1\xf8\x7c\x77\xb1\x98\x5f\xbc\xb9\x9a\x7d\xbe\x5d\x5c\x89\x24\x1a\xe9\xf5\x19\xa3\x3c\x07\x21\x3c\xd6\x0f\x60\xa2\x3f\x95\x89\x94\x7f\xec\x54\x95\xd7\xf7\xaf\x7b\x8b\x97\x25\x16\xb6\x44\xca\xbb\x74\x44\x34\x3a\x9c\xb8\x8c\xcf\x78\xe7\x83\xc2\xa4\xd7\x0d\x9a\x3e\xc3\xf8\xd8\x57\xe9\x1b\x55\x7c\xdd\x90\x6d\x4d\x19\x27\xe3\x21\xb5\x50\x40\xbe\x6f\x70\x00\x71\xa9\x58\xf9\x22\xb6\x05\x7b\xce\x9e\x6a\x78\x8e\x75\x1d\xf9\x90\xe1\xe9\xea\x35\xba\x31\xaa\x71\x5b\xcb\x67\x45\xee\x17\x3d\xbe\x94\x30\x37\x4d\xcb\x93\xa1\x85\x8f\x92\x46\xa1\x00\x83\xd8\xcb\x63\xe3\x43\x0e\x5f\xa2\xec\xc5\xe5\xf5\xdb\xe5\x3f\x9f\x66\xe0\x87\x62\x1a\x65\xfd\x1f\xaa\x72\x1a\x01\x64\x35\xb2\x82\x62\xab\xc8\x21\xe7\xe2\x76\xf9\xee\xf5\x6f\x22\x6c\xb0\xe6\x0a\xa7\x4f\xb7\x38\xf8\xc6\xca\x64\x67\x12\x65\xb2\x43\xcb\x56\xb6\xbc\x0f\xbe\x8d\xff\x85\xae\xd3\x1b\x3f\xad\xde\xde\x38\x26\xc5\xe8\x80\xb7\xe8\xfb\x1e\xec\x1a\xde\x5b\x78\x7b\xf9\xd7\xff\x1c\x64\x0a\xb6\x84\xeb\x5c\xf8\x26\x76\x13\x29\x37\xb6\xb4\x45\x6a\x69\x23\x9f\x99\xd4\xe9\xe9\x3d\x93\x6a\x0a\x47\x61\x48\x7d\x7c\xd9\x0c\x69\xcc\xd9\x4f\x5a\xd1\x12\xa1\xe1\xea\x1e\x5a\xa7\xcd\x06\x14\x3c\xd1\xcf\xb0\x52\x0e\x4b\xb0\x26\x90\xf4\xc7\x26\x1e\x1e\x20\xf5\x2f\x87\x83\x18\xc3\x7e\xab\x8b\x6d\x00\x2d\x94\x81\x15\x42\x61\xcd\x5a\x6f\x5a\xc2\x12\x76\x5a\x05\x2f\x34\x3b\x4d\xd6\xf8\xf9\x3d\x0d\xcf\xb3\x9d\xfe\x04\xdf\x1b\xc4\xa7\xea\x71\xaa\x42\x61\x4d\x81\x0d\x3b\xd9\x52\xe5\xa4\x98\x6e\x91\x30\xd4\x60\x6d\x09\x6a\x4b\x18\x60\xb4\x59\x5b\xaa\x15\x6b\x6b\x40\xad\x6c\xcb\x3e\x1f\x07\xda\x1c\xab\x0e\x17\x9f\xe6\x6e\x10\xfe\xe1\x41\xaf\x21\x9d\x11\x1d\x0e\xc1\x3f\x6b\xa6\x99\x63\xb2\x66\x33\xf5\x35\xf0\x2d\x7b\x38\x64\xf2\xb8\xf4\xe8\x84\xa6\x0c\x0e\x9d\x7b\xdf\xb3\x03\x8c\x4a\xad\xb0\xea\x32\xf3\x2d\x81\xfd\x49\xc0\x4e\x55\x6d\xe8\x02\x5f\xb5\xc7\x4a\x69\x37\x39\xda\x66\x2b\x92\xbd\x5f\xe6\x67\x50\x11\x2a\x20\xbb\x77\xb9\xf8\x55\x40\x61\x2b\x97\x8b\x5f\xfe\x2f\x80\x50\x95\xd6\x54\xf7\xb9\x60\x6a\x51\x04\xba\x3d\x91\xf4\x2e\x44\xf1\xcc\x7b\x88\x0e\x33\x93\x1d\xb1\x63\x22\x4f\x70\x9d\x33\xec\x95\x83\x4a\x39\x86\xda\x96\x7a\xad\xb1\x04\xc5\x13\x38\xc3\xbf\x6d\x4a\xc5\xb8\xd4\xb5\x0f\x92\x3e\x05\xdd\xd7\x28\x93\xdd\x80\x64\x32\x0c\xe1\x97\x4e\x76\x07\x83\x5b\x37\x15\xe4\xd0\x5f\x5c\xe9\x87\xd6\x71\xfc\xf8\xf5\x11\xf7\xb1\x18\xf4\x7f\x92\x7e\xf2\x43\x1b\xff\x3c\xf8\x49\x72\x14\xe0\x9f\xb4\x3e\xde\x43\xb8\x24\x16\xe8\x1a\x6b\x1c\xfe\x4d\x9a\x91\xc6\x40\xf8\x0d\x5e\x1d\x77\xbe\xb5\xe8\x38\xa8\xb6\xf6\x82\x03\x93\x1c\x5e\x9e\x4b\x8e\x17\xe1\xdb\xc5\xd5\x64\x28\x82\xe3\xa0\xad\x25\xae\x91\xc0\xc7\xee\x74\xdf\x8b\x38\x12\x79\x8c\xf3\x34\xd3\xd9\x77\x2c\x5a\xc6\x78\x3f\x86\x10\x26\xf9\x3d\x18\xbe\xc8\xc1\xe8\x2a\xb8\x76\xd7\xd9\xcc\xeb\xa3\xb7\x42\xa2\xe3\x47\x32\xee\x92\xb8\x61\xc5\xad\x9b\x1b\x46\x32\xaa\xba\x41\xda\x21\x05\x8b\x24\x1a\x79\x32\x87\x38\x89\xce\xae\x91\xd9\x79\x80\x10\x37\xed\xc4\x7f\x60\x12\x8d\x46\x84\xdc\x92\x09\x19\xb9\xe3\x31\x8f\xfb\x3c\x7a\xcb\xf4\xca\x6b\x19\xc7\x84\xdf\xd2\xb7\xdd\x15\x11\x27\xdd\xb5\x85\xcf\xc7\xc1\x1f\xf1\xbb\xcd\x47\xf9\xcf\xe1\xa5\x3b\xa9\xfe\xbf\x01\x00\x00\xff\xff\xd4\x98\xc0\xd4\xec\x08\x00\x00"),
		},
		"/demo/secrets": &vfsgen۰DirInfo{
			name:    "secrets",
			modTime: time.Time{},
		},
		"/demo/secrets/demo_secrets.go": &vfsgen۰CompressedFileInfo{
			name:             "demo_secrets.go",
			modTime:          time.Time{},
			uncompressedSize: 5401,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xd4\x57\x6d\x6f\xe3\xc6\x11\xfe\x4c\xfe\x8a\x39\x02\x4d\xc9\x44\x21\x13\xe0\x92\x36\xae\xa8\xe0\x62\xcb\x57\xe3\xae\xb1\x63\xf9\x5a\x14\x45\x71\xb7\x22\x87\xd2\xc6\xe4\x2e\xb3\xbb\x94\xec\x18\xfc\xef\xc5\xec\x2e\x69\xc9\x56\x5a\xfb\x90\xa0\xed\x7d\x38\x53\xe4\xec\xbc\x3c\x33\xcf\xec\x4c\xcb\x8a\x6b\xb6\x42\x68\x18\x17\x61\xc8\x9b\x56\x2a\x03\x71\x18\x44\x85\x14\x06\x6f\x4c\x14\x06\x11\x8a\x42\x96\x5c\xac\xb2\x25\xd3\xf8\xf5\x4b\x7a\x55\x35\xf6\xcb\xda\x34\x75\x66\xb0\x69\x6b\x66\x90\x5e\x08\x34\xd9\xda\x98\x96\x9e\xa5\x8e\xc2\x30\x88\x56\xb2\xa8\x65\x57\xa6\x25\x6e\x32\x8d\x85\x42\xa3\xa3\x30\x78\x0f\x87\x3e\x64\xb5\x2c\x58\x3d\x4a\x25\x61\x98\x65\x70\x75\x7e\x72\x1e\xab\x0d\x13\x2b\x14\x26\x39\x82\xab\x35\xd7\x50\xf1\x1a\x81\x6b\xe8\x34\xaa\xcf\x37\x5c\xf3\x65\x8d\x13\x60\x65\x09\x0d\x13\xb7\x50\xc8\xa6\x41\x61\x34\xe0\x4d\x5b\x33\x2e\xb8\x58\x91\xaa\xb5\xdc\x02\x37\xb0\x95\xea\x5a\xa7\x61\x58\x75\xa2\x00\x2e\xb8\x89\x13\xb8\x0b\x03\x72\x3c\xfd\x33\x13\x65\x8d\xa7\x9d\x28\xe2\x28\x2b\xb1\x91\xa3\x6f\xd1\x04\xfc\xe3\x5c\x14\xea\xb6\x35\x4e\x54\x25\xff\xf1\x24\x3a\xf9\x8f\x57\x50\xe2\x03\x05\x27\xb8\xaf\xa0\x0f\xc3\x0d\x53\x70\x8d\xd8\xa2\x7a\x77\xf9\x16\xb4\x51\x14\xf3\xfd\x4b\xf8\xd4\x1f\x4d\xdf\xd8\xdf\x3b\x9f\xe6\x4a\x01\x2a\x25\xd5\x23\x44\xee\xf5\xe5\x20\x75\xfa\x1a\x0d\x8a\x4d\x1c\x2d\xe6\xc7\x97\xf3\xab\xc5\xfb\x37\xf3\xf9\xc5\xfc\xf2\xfd\xbb\xcb\xb7\x51\x12\x06\xbc\xda\xb1\x9f\xe7\x10\x45\xa4\x22\x38\x90\xc1\x4b\x6c\xe4\x06\xa1\xc4\x8a\x75\xb5\x01\x2a\x1e\x95\x86\xc1\x9e\xb9\xc8\x15\xdb\x35\xde\x1e\x65\x99\x6e\x5e\x2f\x7f\x6c\xfe\xf0\xe5\xf7\x37\xe5\x97\x67\xab\xaf\x4e\x17\x5f\x6c\x7f\xfc\x66\x51\x2f\x7f\x7e\x75\xa6\x84\xac\x8f\x7f\xfe\x66\xf9\xc3\x0f\x5f\x77\xaf\xd6\xf5\xcb\x3c\x0a\x83\x7e\xf0\x7c\xb2\x13\x61\x3e\x60\x97\x9e\xb7\x28\x1c\x08\xb1\xaf\xf2\xf4\x3b\x56\x5c\xaf\x94\xec\x44\x19\x27\x93\xfb\x30\x2c\xb0\xe6\xb6\xc5\x11\x77\x66\x18\x61\xdb\x15\x86\x82\xdb\x41\x3a\xd8\x41\x31\x38\x13\x40\xff\xdc\x27\xc8\x32\x5b\xa6\xc0\x45\xdb\x99\x30\x38\xef\xcc\xfe\x47\xd9\x19\xfb\xe1\x3b\x1b\x31\x2c\xa5\xac\xc9\x6e\x21\x85\xb6\x5c\xf4\xb6\xaf\x3c\xcf\x2e\x14\x56\xfc\x06\x72\xf8\x10\x4e\x5f\x9c\x9c\x1f\x5f\xfd\xfd\x62\x0e\xc4\xc4\x59\x38\x1d\xfe\x20\x2b\x67\x21\xc0\xb4\x41\xc3\xa0\x58\x33\xa5\xd1\xe4\xd1\xbb\xab\xd3\xcf\xff\x18\xd9\x0f\x86\x9b\x1a\x67\x07\x28\x08\x54\x78\xd3\xcc\x7d\x0f\xa7\x99\x53\x35\x5d\xca\xf2\xd6\x1e\x6c\xe9\x7f\x70\x24\x6c\xa9\x73\x90\xbc\xd0\x46\x31\x83\x1a\xcc\x1a\x29\x56\x90\x15\xbc\x96\x70\x7c\xf2\xe6\xf7\x1a\xa6\x0c\xd6\x0a\xab\x3c\xa2\x22\xd7\x47\x59\xb6\x92\xa5\x2c\x52\xa9\x56\xd9\xa1\xde\x30\xf3\x0f\xd3\x8c\xcd\xc0\xb7\xa7\x94\x2c\x67\xed\xae\x03\x67\x86\xe8\x5f\x74\x4a\xa1\x30\xf5\x2d\x74\x9a\xd0\x64\xb0\x5f\xe5\x40\x55\x54\x82\x14\xd6\x33\x4a\x57\x74\x77\x07\x29\x3d\xf4\x7d\x34\x81\xed\x9a\x17\x6b\xab\xaf\x60\x02\x96\x08\x85\x14\x15\x5f\x75\x0a\x4b\xd8\x70\x66\x4f\xa1\xd8\x70\x25\x05\xf5\x13\xd8\x30\xc5\xd9\xb2\x46\x38\xc4\x81\x03\x5e\x2e\x10\x0f\xc5\x7f\x1f\x75\x21\x45\x81\xad\xd1\x59\xa7\x6a\x9d\x45\xb3\x35\x2a\xb4\x91\x57\x52\x41\x23\x15\x5a\x35\x5c\x54\x52\x35\xcc\x70\x29\x80\x2d\x65\x67\x28\x14\x0d\x5c\x78\x94\xe1\xd5\xc5\x99\xde\x35\xdf\xd5\xce\xfe\xb4\xe6\xb3\xd1\x7e\x3a\x36\xa2\x99\x6f\x40\xd6\x92\x43\x8e\x22\xdd\xc7\x6e\x9a\xd5\xfc\xa0\x96\xa1\x1b\xcd\x7c\x17\x7a\x92\x96\x69\xe6\x7c\xba\xbb\xe3\x15\xa4\x73\xa5\xfa\xde\xe9\x6e\x67\x53\x6d\x94\x14\xab\x19\x65\x86\x58\xd4\xf7\xd3\xcc\xbf\xf2\xf1\xdc\xdd\xa1\x28\xfb\xfe\x43\xf8\x88\x0c\x8b\xae\x1a\xc9\x90\xb9\x22\x9d\x66\x96\x05\x24\x9c\x65\x70\x46\xa4\x3b\x1a\xbb\x1f\x11\x38\x1d\xb5\x78\x18\x06\x65\xf7\x3d\xe2\x01\xd7\x3e\x83\x0f\x14\x02\x25\x61\x36\x7a\x5d\xb3\x25\x7a\x98\x01\xe6\xc2\xa0\x02\x7b\xd7\x50\x43\x81\x92\x1a\x85\x91\xe0\x11\x3f\xf2\x62\xd3\xa5\xca\x86\x23\x53\x12\x64\x0a\x19\x28\xb9\xd5\x79\xf4\x32\x82\x42\xd6\x3a\x8f\xbe\xfa\x22\x02\xc1\x1a\xcc\xa3\x51\x5f\x64\xc1\x39\x13\x16\x9b\xe1\x9c\x77\x25\x73\x8e\x78\xa8\x28\x5b\xee\x77\x3b\x1a\xb2\x8d\x07\xa8\x8f\xe5\x51\xb1\xc6\xe2\x7a\x29\x6f\x06\x13\xfe\x4e\x87\x0d\xab\x3b\xcc\x23\xa3\x3a\x8c\x7c\x92\x5c\x3f\xea\x7b\x7b\x04\x4b\x9f\x84\x41\xeb\xd5\x1a\x5d\x94\x6c\x49\xcd\x9c\x6b\x70\xaa\xc0\x8e\x0b\x58\xfe\x82\x73\xbb\xbe\xe8\x6e\xd9\x70\x33\xda\xf6\xd9\x78\xe1\x3a\x54\xe6\xd0\x1e\x4b\xe6\xbc\x33\x3b\x25\xf3\x00\x7c\x7b\x10\x4b\x50\xa8\xe9\x46\x89\xf7\x5d\x49\x9e\x0d\xbf\x42\x56\x4a\x51\xdf\x7a\x40\x2c\xfa\xd4\xb7\x9f\x02\x7f\xc9\x37\xa3\x8d\x47\xa4\xf9\xb6\xe0\xed\x1a\x15\xe9\xc8\xef\x95\x7e\xe2\xfc\xb5\x6f\xfc\x2d\xd0\xf7\x23\xbf\x80\x5b\x8a\x79\x73\x5e\xfd\x40\x09\xf8\xec\x61\xcd\x3a\x4a\x3c\xa1\xfa\xbd\xfa\x5f\xad\xfa\x5d\x10\x9f\x7b\xd0\x47\x0a\xf8\xc8\x3f\x92\x02\xf7\x78\x3d\x8b\x03\x0f\x3d\xfc\x6d\x38\xb0\xa0\x81\xd2\xd7\xdc\x7e\xf4\xcf\xaf\x7e\x9f\x8d\xe7\x57\xbf\x3f\x38\x56\xff\xff\x40\xb1\xfb\xae\xf7\xed\xd8\xbf\x0e\xd6\xfa\x08\xee\x78\x1f\x7d\x5c\xa5\x27\x6e\xf8\x8d\x1f\xb5\xf5\xa6\xad\x21\x87\x61\x3f\x49\xff\xd2\x69\x13\x8f\xbf\xbe\xc7\x6d\x1c\x0d\x43\xcf\x70\x31\x26\xe9\x05\xcd\x4a\xf1\xe1\xfb\x21\x49\x1e\x71\xe7\xe9\x26\x86\x5b\xf3\x81\x89\x07\x24\x4c\x12\xbf\xf5\x1c\xdc\x14\x80\xd5\xb5\xdc\x8e\x63\x96\x72\x37\x0c\x71\x4f\xcb\x06\x47\xc6\x2d\x71\x08\x08\xcb\xd4\xcd\xf4\x07\xd5\xc5\x5b\xb0\x8b\xc7\x25\xea\x56\x0a\x8d\x7f\x53\xdc\xd0\xd0\xac\xf0\x27\xf8\xd4\x7f\xf9\xa9\x43\x6d\xec\x32\xe0\x6a\xf7\x28\x87\x4f\x76\xba\x09\x4d\xf8\xef\x2e\xdf\x1e\x51\xc6\xc6\xb1\x79\x12\x06\xc1\x99\xb0\xef\x48\x57\x7a\x2a\x55\xf3\x57\xaa\xf2\x78\xe7\x46\x4b\x48\xca\x55\xc0\xd1\x43\x29\x4f\xc8\xc4\xae\x11\xb6\x28\x27\x76\xaa\x2f\xb1\x42\x05\x14\x90\xdb\x4f\x68\xe7\x40\xa5\xc8\xa9\xc7\xa9\x4f\xe7\x37\x58\x74\x06\xe3\xed\xc4\x8d\xde\xc9\x9f\xac\xf0\x8b\x1c\x04\xaf\xed\x71\xb7\x77\xcd\x69\x6a\x27\x29\x54\xca\xff\x48\x26\x0e\x99\x85\x61\xa6\xd3\x67\x04\xb1\x60\xf5\x02\xd5\xc6\x2e\x13\x92\x76\x36\x72\xa8\x8f\x93\x70\x67\xf3\x99\xef\xab\xb7\x56\x53\xb7\x7c\x8c\x02\x61\x10\x28\x34\x9d\x12\x36\x22\x5e\x39\xdf\xa8\xaf\xbd\x18\x77\x26\x2a\x67\x2e\xe0\x1f\xff\x5c\xde\x1a\x74\x61\x3a\x29\x7f\x39\x58\xdf\x49\x08\xc7\xa5\x23\x20\x73\x36\x04\xc8\x7d\x2f\x4a\x17\xa6\x9c\xfb\xdd\x3d\x3d\x41\xea\x4b\x0b\xbb\x74\xc4\x83\xc9\xc4\x1e\xab\x1e\xc1\xb2\xe7\x79\xd5\x18\x87\x4a\x15\x47\x17\xfb\x03\xce\x96\x69\x10\x92\xc6\xe2\x9a\x97\x30\x24\xf3\x77\x9b\xc8\x3a\x62\xb5\x8f\xc1\x5a\xbc\x82\x1e\xb0\xd6\xe8\xac\x70\x01\xb9\x8f\x71\xcf\x23\x92\x1b\xeb\x77\x32\x24\xd8\x01\x98\xfa\x04\xc7\x54\x31\xc7\x6e\x7b\xa3\x6c\x71\x7b\xf2\x40\x28\xbf\x14\xc9\x29\xe3\x35\x96\xbb\x23\xda\xbe\xdf\xa3\xdb\xfd\x98\x48\xea\x5f\x07\xb1\xb5\x0f\x78\x25\x3d\xba\xa3\xef\x09\xa5\xb8\xdf\xa5\xf3\xfe\xde\xfe\x74\x3a\x97\x78\x90\xce\xfb\xea\xfe\x0b\x74\xde\xb9\x9d\x7f\x1b\x3e\xef\xf4\xd9\xff\x63\x3e\x0f\xcc\x3c\x7a\x1e\x35\x9f\x51\xce\xc7\x63\x22\x9e\xcc\xcc\xbd\x0a\x1f\x0b\xec\x21\xdf\x7c\x02\x7e\x45\xbe\x0d\xf3\xe0\xbf\xe3\xdb\xa1\x7e\xf7\x3c\x12\x8e\x01\x25\x8f\xba\xce\xbd\x1a\x7d\x48\xd6\x91\xf6\x5f\x01\x00\x00\xff\xff\xf5\x0a\x3c\x9c\x19\x15\x00\x00"),
		},
		"/init": &vfsgen۰DirInfo{
			name:    "init",
			modTime: time.Time{},
		},
		"/init/.dockerignore": &vfsgen۰CompressedFileInfo{
			name:             ".dockerignore",
			modTime:          time.Time{},
			uncompressedSize: 159,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x44\x8c\x31\xaa\xc3\x30\x10\x44\xfb\x3d\xc5\x82\xbb\x2d\xd6\xc7\xf8\xe5\x6f\x72\x01\x39\x5e\x39\x06\x25\x1b\x56\x23\x13\x37\x39\x7b\x10\x22\xa4\x19\x66\x1e\xcc\x9b\xf8\xcf\x49\xd4\x5e\x36\xf2\x4d\xa2\x6b\x29\x24\x5a\x3b\x5f\xcf\xb2\x2f\x44\x13\x5f\xac\x82\x97\xfd\x91\xe2\x24\x51\x58\x45\xa7\xff\x0d\xcf\x06\xf6\xcc\xb8\x19\x6f\xce\x57\x3f\x2c\xd2\x66\x0c\xf7\x2e\xf1\x86\xf1\x8e\x48\xd9\xe3\x4e\x22\xb3\xe2\xbb\x66\xe9\xae\x5c\x91\x60\xbf\xa6\x83\x1e\x29\x2a\x7d\x02\x00\x00\xff\xff\x7c\x7f\x80\xd9\x9f\x00\x00\x00"),
		},
		"/init/.gitignore": &vfsgen۰CompressedFileInfo{
			name:             ".gitignore",
			modTime:          time.Time{},
			uncompressedSize: 159,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x44\x8c\x31\xaa\xc3\x30\x10\x44\xfb\x3d\xc5\x82\xbb\x2d\xd6\xc7\xf8\xe5\x6f\x72\x01\x39\x5e\x39\x06\x25\x1b\x56\x23\x13\x37\x39\x7b\x10\x22\xa4\x19\x66\x1e\xcc\x9b\xf8\xcf\x49\xd4\x5e\x36\xf2\x4d\xa2\x6b\x29\x24\x5a\x3b\x5f\xcf\xb2\x2f\x44\x13\x5f\xac\x82\x97\xfd\x91\xe2\x24\x51\x58\x45\xa7\xff\x0d\xcf\x06\xf6\xcc\xb8\x19\x6f\xce\x57\x3f\x2c\xd2\x66\x0c\xf7\x2e\xf1\x86\xf1\x8e\x48\xd9\xe3\x4e\x22\xb3\xe2\xbb\x66\xe9\xae\x5c\x91\x60\xbf\xa6\x83\x1e\x29\x2a\x7d\x02\x00\x00\xff\xff\x7c\x7f\x80\xd9\x9f\x00\x00\x00"),
		},
		"/init/Dockerfile": &vfsgen۰CompressedFileInfo{
			name:             "Dockerfile",
			modTime:          time.Time{},
			uncompressedSize: 517,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x64\x90\xb1\x6e\xdb\x30\x10\x86\x77\x3e\xc5\xdf\x18\x48\x27\x51\x51\xa6\x40\x40\x97\x5a\x6a\x10\xa4\x16\x05\xc5\x49\x63\x14\x1d\x68\xfa\x2c\xb1\x91\x78\x06\xc9\x34\x35\x82\xbc\x7b\x61\xc9\xe8\xa2\xf1\x8e\xc7\xef\xbe\xff\x16\x28\x14\x2a\xb5\x46\x53\xae\xd4\x53\xf9\x09\xeb\x8e\xb0\xe7\xbe\xe7\x37\xeb\x5a\xf4\xd6\x11\xac\xdb\xb3\x1f\x02\x62\x47\xb8\x65\x2c\x8b\x7b\xf0\x7e\xac\x0a\x36\x2f\xe4\x61\x07\xdd\xd2\xe7\x20\x16\x70\x7a\x20\xbc\xd9\xd8\xf1\x6b\x84\x76\x47\xfc\x21\x1f\x2c\xbb\x33\x43\x47\xcb\x4e\x8a\x05\x5a\x36\xbb\x97\x64\xfc\x97\xe3\xfd\x5d\xd6\x9e\x7f\x93\x89\x95\x1e\xe8\xe3\x43\x88\x05\x1e\x22\x1d\x90\xe5\xf8\xfa\x6a\xfb\xdd\x69\xeb\xd6\x3a\xed\x8f\x52\x7c\x6b\xd4\x0a\x2d\xf7\xda\xb5\x79\x26\xb3\x6b\xe8\x80\xed\x69\x48\x2c\x55\xbd\x81\x44\x1a\xbc\x11\x3f\x54\x73\x5f\xdc\x35\x53\x51\x56\x4f\xb8\x55\x59\x96\xad\x54\xf1\xf8\xbd\x04\xbb\x73\xab\x6e\xd4\xf3\x06\x5d\x8c\x87\x90\xa7\xe9\xc1\xf3\xdf\xa3\x9c\xd0\x92\x7d\x9b\x8a\xe6\xb1\x42\xcb\x13\x1e\x09\xcf\x44\x71\x79\x79\x7a\x37\x3d\x69\x87\x64\xe0\x9d\xd1\xa6\xa3\xff\xfa\xd7\x39\x96\x9e\x74\xa4\xe9\x40\xe3\x5d\x46\x58\x9c\x07\x32\x5e\x5a\x4e\x77\x36\x44\xcf\x3d\x85\x90\x6e\x75\xa0\x29\x52\x92\xec\x3d\x0f\x5f\x26\x8b\x53\xa0\x74\xe6\x31\xeb\x88\xf2\xb9\x56\x0f\x25\x6e\xae\x6e\xae\x44\x59\xad\x9b\x4d\xad\xee\xaa\x35\x7e\x5e\xcc\x46\x2f\x7e\x89\x7f\x01\x00\x00\xff\xff\x35\x65\x45\x4d\x05\x02\x00\x00"),
		},
		"/init/README.md": &vfsgen۰FileInfo{
			name:    "README.md",
			modTime: time.Time{},
			content: []byte("\x49\x27\x6d\x20\x61\x20\x72\x65\x61\x64\x6d\x65\x20\x61\x62\x6f\x75\x74\x20\x75\x73\x69\x6e\x67\x20\x74\x68\x65\x20\x63\x6c\x69\x0a"),
		},
		"/init/biomes": &vfsgen۰DirInfo{
			name:    "biomes",
			modTime: time.Time{},
		},
		"/init/biomes/README.md": &vfsgen۰FileInfo{
			name:    "README.md",
			modTime: time.Time{},
			content: []byte("\x54\x68\x69\x73\x20\x64\x69\x72\x65\x63\x74\x6f\x72\x79\x20\x63\x6f\x6e\x74\x61\x69\x6e\x73\x20\x74\x68\x65\x20\x70\x72\x6f\x6a\x65\x63\x74\x27\x73\x20\x62\x69\x6f\x6d\x65\x73\x2c\x20\x74\x68\x65\x20\x75\x6e\x69\x74\x20\x6f\x66\x0a\x63\x6f\x6e\x66\x69\x67\x75\x72\x61\x74\x69\x6f\x6e\x20\x61\x6e\x64\x20\x70\x72\x6f\x76\x69\x73\x69\x6f\x6e\x69\x6e\x67\x20\x69\x6e\x20\x74\x68\x65\x20\x47\x6f\x20\x43\x44\x4b\x2e\x0a"),
		},
		"/init/biomes/dev": &vfsgen۰DirInfo{
			name:    "dev",
			modTime: time.Time{},
		},
		"/init/biomes/dev/biome.json": &vfsgen۰FileInfo{
			name:    "biome.json",
			modTime: time.Time{},
			content: []byte("\x7b\x0a\x20\x20\x22\x73\x65\x72\x76\x65\x5f\x65\x6e\x61\x62\x6c\x65\x64\x22\x20\x3a\x20\x74\x72\x75\x65\x2c\x0a\x20\x20\x22\x6c\x61\x75\x6e\x63\x68\x65\x72\x22\x20\x3a\x20\x22\x6c\x6f\x63\x61\x6c\x22\x0a\x7d\x0a"),
		},
		"/init/biomes/dev/main.tf": &vfsgen۰CompressedFileInfo{
			name:             "main.tf",
			modTime:          time.Time{},
			uncompressedSize: 148,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x3c\xcb\x31\xae\xc2\x30\x0c\x06\xe0\xdd\xa7\xf8\xd5\xee\xc9\xeb\xfa\x24\xb8\x04\xdd\x51\x14\x1c\xc5\x52\x1b\x07\x3b\x65\x41\x70\x76\x26\x98\x3f\x7d\x33\xd6\x2a\x8e\x22\x1b\xc3\x3b\x67\x29\xc2\x8e\x95\xcd\x52\x51\xdb\x61\xec\x7a\x58\x66\x0f\x34\xe3\xc2\x8c\x3a\x46\xf7\xff\x18\x37\x4e\xd6\x42\x4d\x5e\x25\xab\xf5\x90\x75\x8f\xe3\xdb\x22\x8a\x1a\x52\x83\xb4\x61\x1a\x88\x7e\x82\x27\x01\xc6\xf7\x43\x8c\x6f\xd7\x07\x9b\x8b\x36\x9c\x30\xbd\xcf\xf8\x0b\xcb\x32\xd1\x8b\x3e\x01\x00\x00\xff\xff\xac\x85\x07\x29\x94\x00\x00\x00"),
		},
		"/init/biomes/dev/outputs.tf": &vfsgen۰CompressedFileInfo{
			name:             "outputs.tf",
			modTime:          time.Time{},
			uncompressedSize: 302,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x3c\x8e\x31\x6e\xc3\x30\x0c\x45\x67\xeb\x14\x1f\xc9\xd0\xa5\xb0\x33\x06\x05\xba\xf5\x06\xcd\x1e\xa8\x32\x5d\x09\x95\x45\x81\xa4\xea\xa1\xc8\xdd\x0b\xc7\x89\x37\x02\x7c\x0f\xff\x1d\x71\x89\x49\x31\xa5\x4c\xd0\x4a\x21\x4d\x89\x14\x17\x12\xf1\x13\xcb\x8c\x99\xc7\x96\x09\xdc\xac\x36\xd3\xde\x1d\xf1\x49\x84\x68\x56\xf5\x6d\x18\x96\x65\xe9\xed\xc9\xf6\x89\x87\x91\x83\x0e\x81\xcb\x94\xbe\x9b\x78\x4b\x5c\x86\xa7\x1a\x6d\xce\xce\xad\x7b\x84\xec\x5b\x09\x71\x1f\x14\x28\x99\x82\xeb\x2a\x28\x26\x16\x58\x24\x7c\x25\x9e\xe9\x45\x1f\x34\x49\xef\x1e\xb5\x49\xe1\x91\x39\xf8\xbc\xff\x5e\xb1\xc4\x14\x22\xa4\x15\x05\x97\xbb\xbf\x11\x1f\x1c\x7e\x48\x30\x7a\x9a\xb9\xf4\x6e\xcb\xc1\x61\x13\xaf\x7b\xc2\x01\x7f\xae\xfb\xf5\xb9\xd1\x7a\x74\x91\xd5\xae\x95\xc5\xf0\x8e\xf3\xe9\x7c\x72\xdd\xcd\xdd\xdc\x7f\x00\x00\x00\xff\xff\x1a\x2a\xa0\x2e\x2e\x01\x00\x00"),
		},
		"/init/biomes/dev/secrets.auto.tfvars": &vfsgen۰CompressedFileInfo{
			name:             "secrets.auto.tfvars",
			modTime:          time.Time{},
			uncompressedSize: 261,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x44\x8e\xb1\x6e\xc3\x40\x0c\x43\x77\x7f\x05\x11\xcf\xb6\xf7\x00\x1d\x3b\x76\x6a\x7f\x40\xf6\xe9\x62\x21\xe7\x93\x21\xc9\xb9\xf6\xef\x0b\x03\x4d\xb3\x91\x04\x1f\xc1\x1e\x5f\xab\x38\x9a\x1e\x25\x61\x66\xb4\x95\x8d\xf1\xa3\xc7\x5f\xb4\x1f\x01\xe7\xc5\x38\x1c\x45\xee\x8c\x44\x41\x33\x39\x63\x27\xf7\xa6\x96\x7c\xec\x7a\x7c\xd0\x9d\xe1\x87\x31\xe2\x9c\x5b\xc9\x41\x58\xd4\x8c\x7d\xd7\x9a\xa4\xde\xf0\x20\x13\x9a\x0b\x43\xea\xbf\xf6\x31\xf2\x89\x7f\x32\x63\x8d\xd8\xfd\x3a\x4d\xad\xb5\x31\xd8\x8c\xb2\xda\x36\x8a\x4e\x49\x17\x9f\x16\xad\x59\x6e\x87\x51\x88\xd6\xe9\xc5\xaf\xb1\x95\xfe\x69\x87\xc4\x59\xaa\x9c\x15\x1f\x22\x3f\xc8\x7c\xc8\x52\xd8\xbb\xae\xc7\xfb\x37\x6d\x7b\xe1\x6b\xd7\x23\xcd\xcf\xf3\x78\xc3\xe5\xd2\xfd\x06\x00\x00\xff\xff\x40\x5d\xc8\x55\x05\x01\x00\x00"),
		},
		"/init/biomes/dev/variables.tf": &vfsgen۰CompressedFileInfo{
			name:             "variables.tf",
			modTime:          time.Time{},
			uncompressedSize: 263,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x44\xcc\x41\x6a\xc4\x30\x0c\x85\xe1\xbd\x4f\xf1\x48\xb6\x25\xd9\x0f\x74\xd9\x13\x74\x2e\xa0\xd8\xf2\x44\xe0\xd8\x46\x52\x26\x2d\xa5\x77\x2f\x19\x3a\xad\x56\x82\xc7\xff\x8d\xb8\xae\x62\xc8\x52\x18\xd6\x39\x4a\x16\x36\x5c\x59\x95\x72\xd3\x0d\x5b\x4b\x7b\x61\x48\xed\xbb\xe3\x4e\x2a\xb4\x14\xb6\x29\x8c\x78\x67\xc6\xea\xde\xed\x32\xcf\xc7\x71\x4c\xfe\x6c\x26\x69\x73\x6a\xd1\xe6\xd8\x6a\x96\xdb\xae\xe4\xd2\xea\xfc\x1f\xaf\xbe\x95\x10\x46\xbc\x7d\xd0\xd6\x0b\x5f\xc2\xf8\x27\x63\x48\x4b\x27\xb3\xa3\x69\x1a\xf0\x15\x46\x00\xfe\xd9\x19\xbf\xf7\x8a\xc1\x5c\xa5\xde\x86\xc7\x94\xd8\xa2\x4a\x3f\xfd\x73\xba\xae\x8c\x44\x4e\x0b\x19\xe3\xc9\xbc\x3c\x3e\x4e\x90\x0a\xe3\xa8\xec\x36\xd1\xee\x6d\xf2\x7c\x27\xb5\x13\xfa\x0e\x3f\x01\x00\x00\xff\xff\xce\xb7\x2d\x0a\x07\x01\x00\x00"),
		},
		"/init/go.mod": &vfsgen۰FileInfo{
			name:    "go.mod",
			modTime: time.Time{},
			content: []byte("\x6d\x6f\x64\x75\x6c\x65\x20\x7b\x7b\x2e\x4d\x6f\x64\x75\x6c\x65\x50\x61\x74\x68\x7d\x7d\x0a\x0a\x67\x6f\x20\x31\x2e\x31\x32\x0a\x0a\x72\x65\x71\x75\x69\x72\x65\x20\x28\x0a\x09\x67\x6f\x63\x6c\x6f\x75\x64\x2e\x64\x65\x76\x20\x76\x30\x2e\x31\x33\x2e\x30\x0a\x29\x0a"),
		},
		"/init/main.go": &vfsgen۰CompressedFileInfo{
			name:             "main.go",
			modTime:          time.Time{},
			uncompressedSize: 533,

			compressedContent: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x5c\x90\xc1\x6a\xdc\x30\x10\x86\xcf\x9a\xa7\x98\xcc\xa1\x48\xad\xf1\xa6\xb7\xb0\xc5\x87\xd0\x36\xcd\x21\x4d\xca\x6e\x21\x67\xb3\x1e\xbb\xa2\x5a\xc9\x1d\xc9\xde\x40\xf0\xbb\x17\xc9\x86\xd2\xdc\x84\xe6\x67\xbe\xf9\xfe\xb1\x3d\xfd\x6e\x07\xc6\x73\x6b\x3d\x80\x3d\x8f\x41\x12\x6a\x50\xd4\x9f\x13\x81\x22\xcf\x69\xf7\x2b\xa5\x31\xbf\x43\x24\x00\x45\x43\x38\xb9\x30\x75\x75\xc7\xf3\x4e\xf8\xcf\xc4\x31\xb9\x30\xd0\x9b\x49\x64\x99\x59\x08\x0c\x40\x3f\xf9\x53\x01\x68\x83\xaf\xa0\x0a\x62\xdf\x60\x88\xf5\x37\x4e\xec\x67\x4d\x3f\x9e\x0e\x3f\xc9\x80\xb2\x3d\x96\x69\xd3\x20\x51\xce\xae\xe1\x06\xe9\xe6\xfa\xe6\x9a\x40\x2d\xa0\xf2\x35\xf5\x7d\xeb\x3b\xc7\x77\x93\x3f\x69\xda\x51\x85\x83\x30\x27\x03\x2a\xca\x9c\x57\xaf\xf0\xfa\x91\x2f\xba\xc4\xbf\x70\xdf\x4e\x2e\x1d\xf3\xf7\xf7\xe9\xa5\xc2\x77\x5b\xe2\x69\x4c\x36\xf8\x98\x49\x87\x55\xe5\x21\x0c\x03\xcb\x1e\xff\x99\xe5\x35\x8f\x9f\x8f\xb7\xeb\x44\x87\x58\x1f\x53\x17\xa6\x54\x61\x16\xd3\x2c\x12\xc4\xe0\xeb\x62\x2a\x50\xcb\x2a\xc1\x22\xe5\x0c\x99\xeb\x07\x1b\x13\xfb\x5b\xdf\x15\xb8\xa6\x3d\xe1\x87\x22\x69\x3e\x95\xd8\x55\x83\xde\xba\xe2\xda\x9f\x53\x7d\x37\x8a\xf5\xc9\xf9\x0d\xc3\x22\x55\x8e\x19\x50\x2a\xc4\xfa\xeb\x8b\x4d\xfa\xa3\xc9\x3d\x2c\x5b\xaf\xc5\x5c\x5f\xb0\x78\x1e\x38\x8e\xc1\x47\x7e\x16\x9b\x58\x2a\x14\x7c\xbf\xfd\x17\x99\x52\xff\x7f\x94\x4b\x85\x74\xcf\xce\x85\x0a\x9f\x83\xb8\xee\x8a\x0c\x2c\xf0\x37\x00\x00\xff\xff\x32\xfa\x00\xb7\x15\x02\x00\x00"),
		},
	}
	fs["/"].(*vfsgen۰DirInfo).entries = []os.FileInfo{
		fs["/demo"].(os.FileInfo),
		fs["/init"].(os.FileInfo),
	}
	fs["/demo"].(*vfsgen۰DirInfo).entries = []os.FileInfo{
		fs["/demo/blob"].(os.FileInfo),
		fs["/demo/runtimevar"].(os.FileInfo),
		fs["/demo/secrets"].(os.FileInfo),
	}
	fs["/demo/blob"].(*vfsgen۰DirInfo).entries = []os.FileInfo{
		fs["/demo/blob/demo_blob.go"].(os.FileInfo),
	}
	fs["/demo/runtimevar"].(*vfsgen۰DirInfo).entries = []os.FileInfo{
		fs["/demo/runtimevar/demo_runtimevar.go"].(os.FileInfo),
	}
	fs["/demo/secrets"].(*vfsgen۰DirInfo).entries = []os.FileInfo{
		fs["/demo/secrets/demo_secrets.go"].(os.FileInfo),
	}
	fs["/init"].(*vfsgen۰DirInfo).entries = []os.FileInfo{
		fs["/init/.dockerignore"].(os.FileInfo),
		fs["/init/.gitignore"].(os.FileInfo),
		fs["/init/Dockerfile"].(os.FileInfo),
		fs["/init/README.md"].(os.FileInfo),
		fs["/init/biomes"].(os.FileInfo),
		fs["/init/go.mod"].(os.FileInfo),
		fs["/init/main.go"].(os.FileInfo),
	}
	fs["/init/biomes"].(*vfsgen۰DirInfo).entries = []os.FileInfo{
		fs["/init/biomes/README.md"].(os.FileInfo),
		fs["/init/biomes/dev"].(os.FileInfo),
	}
	fs["/init/biomes/dev"].(*vfsgen۰DirInfo).entries = []os.FileInfo{
		fs["/init/biomes/dev/biome.json"].(os.FileInfo),
		fs["/init/biomes/dev/main.tf"].(os.FileInfo),
		fs["/init/biomes/dev/outputs.tf"].(os.FileInfo),
		fs["/init/biomes/dev/secrets.auto.tfvars"].(os.FileInfo),
		fs["/init/biomes/dev/variables.tf"].(os.FileInfo),
	}

	return fs
}()

type vfsgen۰FS map[string]interface{}

func (fs vfsgen۰FS) Open(path string) (http.File, error) {
	path = pathpkg.Clean("/" + path)
	f, ok := fs[path]
	if !ok {
		return nil, &os.PathError{Op: "open", Path: path, Err: os.ErrNotExist}
	}

	switch f := f.(type) {
	case *vfsgen۰CompressedFileInfo:
		gr, err := gzip.NewReader(bytes.NewReader(f.compressedContent))
		if err != nil {
			// This should never happen because we generate the gzip bytes such that they are always valid.
			panic("unexpected error reading own gzip compressed bytes: " + err.Error())
		}
		return &vfsgen۰CompressedFile{
			vfsgen۰CompressedFileInfo: f,
			gr:                        gr,
		}, nil
	case *vfsgen۰FileInfo:
		return &vfsgen۰File{
			vfsgen۰FileInfo: f,
			Reader:          bytes.NewReader(f.content),
		}, nil
	case *vfsgen۰DirInfo:
		return &vfsgen۰Dir{
			vfsgen۰DirInfo: f,
		}, nil
	default:
		// This should never happen because we generate only the above types.
		panic(fmt.Sprintf("unexpected type %T", f))
	}
}

// vfsgen۰CompressedFileInfo is a static definition of a gzip compressed file.
type vfsgen۰CompressedFileInfo struct {
	name              string
	modTime           time.Time
	compressedContent []byte
	uncompressedSize  int64
}

func (f *vfsgen۰CompressedFileInfo) Readdir(count int) ([]os.FileInfo, error) {
	return nil, fmt.Errorf("cannot Readdir from file %s", f.name)
}
func (f *vfsgen۰CompressedFileInfo) Stat() (os.FileInfo, error) { return f, nil }

func (f *vfsgen۰CompressedFileInfo) GzipBytes() []byte {
	return f.compressedContent
}

func (f *vfsgen۰CompressedFileInfo) Name() string       { return f.name }
func (f *vfsgen۰CompressedFileInfo) Size() int64        { return f.uncompressedSize }
func (f *vfsgen۰CompressedFileInfo) Mode() os.FileMode  { return 0444 }
func (f *vfsgen۰CompressedFileInfo) ModTime() time.Time { return f.modTime }
func (f *vfsgen۰CompressedFileInfo) IsDir() bool        { return false }
func (f *vfsgen۰CompressedFileInfo) Sys() interface{}   { return nil }

// vfsgen۰CompressedFile is an opened compressedFile instance.
type vfsgen۰CompressedFile struct {
	*vfsgen۰CompressedFileInfo
	gr      *gzip.Reader
	grPos   int64 // Actual gr uncompressed position.
	seekPos int64 // Seek uncompressed position.
}

func (f *vfsgen۰CompressedFile) Read(p []byte) (n int, err error) {
	if f.grPos > f.seekPos {
		// Rewind to beginning.
		err = f.gr.Reset(bytes.NewReader(f.compressedContent))
		if err != nil {
			return 0, err
		}
		f.grPos = 0
	}
	if f.grPos < f.seekPos {
		// Fast-forward.
		_, err = io.CopyN(ioutil.Discard, f.gr, f.seekPos-f.grPos)
		if err != nil {
			return 0, err
		}
		f.grPos = f.seekPos
	}
	n, err = f.gr.Read(p)
	f.grPos += int64(n)
	f.seekPos = f.grPos
	return n, err
}
func (f *vfsgen۰CompressedFile) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		f.seekPos = 0 + offset
	case io.SeekCurrent:
		f.seekPos += offset
	case io.SeekEnd:
		f.seekPos = f.uncompressedSize + offset
	default:
		panic(fmt.Errorf("invalid whence value: %v", whence))
	}
	return f.seekPos, nil
}
func (f *vfsgen۰CompressedFile) Close() error {
	return f.gr.Close()
}

// vfsgen۰FileInfo is a static definition of an uncompressed file (because it's not worth gzip compressing).
type vfsgen۰FileInfo struct {
	name    string
	modTime time.Time
	content []byte
}

func (f *vfsgen۰FileInfo) Readdir(count int) ([]os.FileInfo, error) {
	return nil, fmt.Errorf("cannot Readdir from file %s", f.name)
}
func (f *vfsgen۰FileInfo) Stat() (os.FileInfo, error) { return f, nil }

func (f *vfsgen۰FileInfo) NotWorthGzipCompressing() {}

func (f *vfsgen۰FileInfo) Name() string       { return f.name }
func (f *vfsgen۰FileInfo) Size() int64        { return int64(len(f.content)) }
func (f *vfsgen۰FileInfo) Mode() os.FileMode  { return 0444 }
func (f *vfsgen۰FileInfo) ModTime() time.Time { return f.modTime }
func (f *vfsgen۰FileInfo) IsDir() bool        { return false }
func (f *vfsgen۰FileInfo) Sys() interface{}   { return nil }

// vfsgen۰File is an opened file instance.
type vfsgen۰File struct {
	*vfsgen۰FileInfo
	*bytes.Reader
}

func (f *vfsgen۰File) Close() error {
	return nil
}

// vfsgen۰DirInfo is a static definition of a directory.
type vfsgen۰DirInfo struct {
	name    string
	modTime time.Time
	entries []os.FileInfo
}

func (d *vfsgen۰DirInfo) Read([]byte) (int, error) {
	return 0, fmt.Errorf("cannot Read from directory %s", d.name)
}
func (d *vfsgen۰DirInfo) Close() error               { return nil }
func (d *vfsgen۰DirInfo) Stat() (os.FileInfo, error) { return d, nil }

func (d *vfsgen۰DirInfo) Name() string       { return d.name }
func (d *vfsgen۰DirInfo) Size() int64        { return 0 }
func (d *vfsgen۰DirInfo) Mode() os.FileMode  { return 0755 | os.ModeDir }
func (d *vfsgen۰DirInfo) ModTime() time.Time { return d.modTime }
func (d *vfsgen۰DirInfo) IsDir() bool        { return true }
func (d *vfsgen۰DirInfo) Sys() interface{}   { return nil }

// vfsgen۰Dir is an opened dir instance.
type vfsgen۰Dir struct {
	*vfsgen۰DirInfo
	pos int // Position within entries for Seek and Readdir.
}

func (d *vfsgen۰Dir) Seek(offset int64, whence int) (int64, error) {
	if offset == 0 && whence == io.SeekStart {
		d.pos = 0
		return 0, nil
	}
	return 0, fmt.Errorf("unsupported Seek in directory %s", d.name)
}

func (d *vfsgen۰Dir) Readdir(count int) ([]os.FileInfo, error) {
	if d.pos >= len(d.entries) && count > 0 {
		return nil, io.EOF
	}
	if count <= 0 || count > len(d.entries)-d.pos {
		count = len(d.entries) - d.pos
	}
	e := d.entries[d.pos : d.pos+count]
	d.pos += count
	return e, nil
}
