package main

type FileDetails struct {
	TID          int    `json:"tid"`
	FileName     string `json:"filename"`
	Path         string `json:"path"`
	Size         int64  `json:"size"`
	LastModified string `json:"last_modified"`
	IsDirectory  bool   `json:"is_directory"`
}

type CopyRequest struct {
	Src string `json:"src"`
	Dst string `json:"dst"`
}

type ReadRequest struct {
	Src string `json:"src"`
}
