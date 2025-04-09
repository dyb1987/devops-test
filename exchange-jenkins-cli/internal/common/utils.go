package common

import "os"

func ReadXml(filePath string) (string, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 文件不能大于1MB
	info, _ := file.Stat()
	if info.Size() > 1024*1024 {
		return "", ErrFileSize
	}

	// 读取文件的前 10 个字节，用于判断 XML 声明
	buffer := make([]byte, 10)
	n, err := file.Read(buffer)
	if err != nil {
		return "", err
	}

	// 检查是否以 <?xml 开头
	if string(buffer[:n])[:5] != "<?xml" {
		return "", ErrIsNotXml
	}

	// 返回模板数据
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
