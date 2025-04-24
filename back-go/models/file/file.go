package file

import "gin-web/utils"

var (
	FilePtr *fileInfos
	//saveFileMultiple == true：多个文件可以共享相同的 MD5 文件，不会重复存储相同内容的文件，节省存储空间。
	//saveFileMultiple == false：每个文件都单独存储，即使内容相同，也会创建物理副本，不共享 MD5 文件，保证独立性。
	SaveFileMultiple = true
	FileDiskTotal    = 50 * utils.MB // 默认50M
)

// 文件信息管理
type fileInfos struct {
	FileInfo *fileInfo           `json:"fileInfo"`
	MD5Files map[string]*md5File `json:"md5Files"`
	UsedDisk uint64              `json:"_"`
}

// md5文件信息
type md5File struct {
	File string   `json:"file"` // 原始文件路径
	Size uint64   `json:"size"` //文件大小
	MD5  string   `json:"md5"`  //文件MD5
	Ptr  []string `json:"ptr"`  // 文件引用
}

// 文件/目录信息
type fileInfo struct {
	Path       string               `json:"path"`            // 相对路径
	Name       string               `json:"name,omitempty"`  // 名字
	AbsPath    string               `json:"absPath"`         // 绝对路径
	IsDir      bool                 `json:"isDir,omitempty"` //是否是目录
	ModeTime   string               `json:"modeTime"`
	FileSize   uint64               `json:"fileSize"`  // 文件大小
	FileMD5    string               `json:"fileMd5"`   // 文件 MD5
	FileInfos  map[string]*fileInfo `json:"fileInfos"` // 子文件夹信息
	FileUpload *upload              `json:"_"`         // 上传时的临时数据
}

// 文件上传信息
type upload struct {
	Md5        string           `json:"md5"`        // 文件上传时的 MD5 值
	Size       uint64           `json:"size"`       // 文件总大小
	SliceSize  uint64           `json:"sliceSize"`  // 上传的分片大小
	Total      int              `json:"total"`      // 文件上传时分片总数
	ExistSlice map[string]int64 `json:"existSlice"` // 已上传的分片
	Token      string           `json:"token"`      // 上传时需要的 token
}
