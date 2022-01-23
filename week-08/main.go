package main

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var bytes = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890") // 组成 value 的 字符
var nums = 10000                                                                     // key-value 的数量                                                              // 数据量 1W
var dataSizeArray = []int{10, 20, 50, 100, 200, 1000, 5000}                          // value的大小，单位字节

func main() {
	// 连接 Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6378",
		Password: "pwd-turato",
		DB:       0, // use default DB
	})

	valueMap := make(map[int][]string) // 字节大小:列数据的数组
	headKey := make([]string, 0, 39)   // csv 文件的列
	// 根据数据大小生成数据
	for _, dataSize := range dataSizeArray {
		var headValue []string
		writeData(dataSize, nums, rdb) //  写入数据

		result, err := rdb.Info(ctx, "memory").Result() // 执行命令： info memory
		if err != nil {
			panic(err)
		}
		// 对结果的格式进行处理：解析结果
		headValue, headKey = parseResults(result)
		valueMap[dataSize] = headValue

		err = rdb.FlushDB(ctx).Err() // 清空缓存
		if err != nil {
			panic(err)
		}
	}

	// 打印数据
	printCsv(headKey, valueMap)
}

// parseResults 解析 info memory 的结果
func parseResults(result string) ([]string, []string) {
	results := strings.Split(result, "\r\n")
	headKey := make([]string, 0, 39)
	var headValue []string
	for _, str := range results {
		tmp := strings.Split(str, ":")
		if len(tmp) < 2 {
			continue
		}
		if len(headValue) <= 39 {
			headKey = append(headKey, tmp[0])
			headValue = append(headValue, tmp[1])
		}
	}
	return headValue, headKey
}

// writeData 写入数据， dataSize：数据大小， nums：数据量
func writeData(dataSize int, nums int, rdb *redis.Client) {
	for i := 0; i < nums; i++ { // 数据的数量
		key := strconv.Itoa(i)
		value := GenerateRandomString(dataSize) // 保证每一个 key 及 value 基本不重复
		err := rdb.Set(ctx, key, value, 0).Err()
		if err != nil {
			panic(err)
		}
	}
}

// printCsv headKey 列名, valueMap 中 key-value 为一行数据，key 为第一列的值
func printCsv(headKey []string, valueMap map[int][]string) {
	var head string
	head = fmt.Sprintf("%s", "data_size")
	for _, v := range headKey {
		head = fmt.Sprintf("%s,%s", head, v)
	}
	sb := new(strings.Builder)
	sb.WriteString(head + "\n")
	for _, dataSize := range dataSizeArray {
		value := fmt.Sprintf("%d", dataSize)
		for _, s := range valueMap[dataSize] {
			value = fmt.Sprintf("%s,%s", value, s)
		}
		sb.WriteString(value + "\n")
	}
	fmt.Println("csv 格式：")
	fmt.Println(sb.String())
}

// GenerateRandomString n 随机字符串的长度，也是字节数
func GenerateRandomString(n int) string {
	result := make([]byte, 0, n)
	for i := 0; i < n; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}
