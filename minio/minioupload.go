package minio

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"jvmdump4k8s/config"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

func putUrl(url string, downUrl string) {
	data := map[string]interface{}{
		"information": map[string]interface{}{
			"BUILD_PACKAGE_DOWNLOAD_ADDRESS": downUrl,
		},
	}

	// 将数据编码为 JSON 格式
	payload, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	// 发送请求
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// 添加 Token 认证头部，如果需要的话
	// req.Header.Set("Authorization", "Token your_token_here")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)
	//fmt.Println("Response Body:")
	//buf := new(bytes.Buffer)
	//buf.ReadFrom(resp.Body)
	//fmt.Println(buf.String())
}

func getSuffixNew(filename string) string {
	re := regexp.MustCompile(`[^/]+$`)

	// 查找匹配的字符串
	match := re.FindString(filename)

	if match != "" {
		return match
	} else {
		return ""
	}
}

func getUUIDFile(filePath string) (string, error) {
	// 获取文件的后缀
	fileExt := getSuffixNew(filePath)
	// 如果文件没有后缀，可以根据实际需求处理，这里假设有后缀
	if fileExt == "" {
		log.Fatal("文件没有后缀")
	}
	currentTime := time.Now()

	// 获取年份和月份
	year := currentTime.Year()
	month := currentTime.Month()

	// 组装新的文件名
	newFileName := fmt.Sprintf("%d/%d/%s", year, month, fileExt)
	//fmt.Println(newFileName)
	return newFileName, nil
}

func Upload(file string, podName string) string {

	var bucketName = config.GlobalConfig.MinioBucket
	var accessKey = config.GlobalConfig.MinioAccessKey
	var accessSecret = config.GlobalConfig.MinioSecretKey
	var apiHost = config.GlobalConfig.MinioApiHost
	fmt.Printf("开始上传minio OSS accessKey=%s bucketName=%s apihost=%s \n", accessKey, bucketName, apiHost)

	//var filename = filepath.Base(file) //获取文件名
	//var ext = path.Ext(file)           //获取扩展名
	//var objectName = "/" + filename + util.FormartdateNow() + ext

	///

	// Initialize minio client object.
	minioClient, err := minio.New(apiHost, &minio.Options{
		Creds: credentials.NewStaticV4(accessKey, accessSecret, ""),
	})
	if err != nil {
		log.Fatalln(err)
	}

	// 打开本地文件
	fileReader, err := os.Open(file)
	if err != nil {
		log.Fatalln(err)
	}
	defer fileReader.Close()
	fileInfo, err := fileReader.Stat()
	if err != nil {
		log.Fatalln(err)
	}
	// 5. 上传文件
	objectName := filepath.Base(file) // 使用文件名作为对象名称
	currentTime := time.Now()

	// 获取年份和月份
	year := currentTime.Year()
	month := currentTime.Month()
	newFileName := fmt.Sprintf("%d/%d/%s-%s", year, month, podName, objectName)
	_, err = minioClient.PutObject(
		context.Background(),
		bucketName,
		newFileName,
		fileReader,
		fileInfo.Size(),
		minio.PutObjectOptions{ContentType: "application/octet-stream"},
	)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("文件 %s 成功上传到存储桶 %s", objectName, newFileName)

	//if err != nil {
	//	fmt.Println("minio文件上传发生错误", err)
	//	os.Exit(-1)
	//	return ""
	//}
	url := apiHost + "/" + bucketName + "/" + newFileName
	fmt.Printf("上传成功 %s\n", url)

	return url
}
