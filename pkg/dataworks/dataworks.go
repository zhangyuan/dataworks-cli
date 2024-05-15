package dataworks

const (
	NOMARL_USE_TYPE = "NORMAL"
)

func GetFileExt(fileType int32) string {
	// ODPS SQL
	if fileType == 10 {
		return "sql"
	}
	// 数据集成
	if fileType == 23 {
		return "json"
	}
	return ""
}
