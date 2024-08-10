package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Input  string `json:"input"`
	Output string `json:"output"`
	Src    string `json:"src"`
	Dst    string `json:"dst"`
	Count  int    `json:"count"`
}

func parseHexString(hexString string) ([]byte, error) {
	hexStrs := strings.Split(hexString, " ")
	bytes := make([]byte, len(hexStrs))
	for i, h := range hexStrs {
		var b byte
		_, err := fmt.Sscanf(h, "%02X", &b)
		if err != nil {
			return nil, fmt.Errorf("invalid hex string: %s", h)
		}
		bytes[i] = b
	}
	return bytes, nil
}

func replaceBytes(config Config) error {
	// 读取输入文件
	data, err := os.ReadFile(config.Input)
	if err != nil {
		return fmt.Errorf("failed to read input file: %v", err)
	}

	// 解析源和目标字节序列
	srcBytes, err := parseHexString(config.Src)
	if err != nil {
		return fmt.Errorf("failed to parse src: %v", err)
	}

	dstBytes, err := parseHexString(config.Dst)
	if err != nil {
		return fmt.Errorf("failed to parse dst: %v", err)
	}

	// 检查源和目标字节序列的长度是否相同
	if len(srcBytes) != len(dstBytes) {
		return fmt.Errorf("source and destination sequences must be of the same length")
	}

	// 计数替换的次数
	replacements := 0

	for i := 0; i <= len(data)-len(srcBytes); i++ {
		if string(data[i:i+len(srcBytes)]) == string(srcBytes) {
			copy(data[i:i+len(dstBytes)], dstBytes)
			replacements++

			// 如果 count 不为 0 并且替换次数达到 count，停止替换
			if config.Count > 0 && replacements >= config.Count {
				break
			}

			// 跳过已经替换的部分以避免重叠替换
			i += len(dstBytes) - 1
		}
	}

	if replacements == 0 {
		return fmt.Errorf("pattern not found in the file")
	}

	// 将修改后的数据写入输出文件
	err = os.WriteFile(config.Output, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write output file: %v", err)
	}

	fmt.Printf("Replaced %d occurrences.\n", replacements)
	return nil
}

func main() {
	// 定义命令行标志
	configPath := flag.String("c", "", "Path to the configuration file")
	configPathAlt := flag.String("config", "", "Path to the configuration file (alternative flag)")

	// 解析命令行标志
	flag.Parse()

	// 确定配置文件路径
	configFile := *configPath
	if configFile == "" {
		configFile = *configPathAlt
	}

	// 检查是否提供了配置文件路径
	if configFile == "" {
		fmt.Println("Usage: ./rehex -c config.json")
		os.Exit(1)
	}

	// 读取配置文件
	configData, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Printf("Failed to read config file: %v\n", err)
		os.Exit(1)
	}

	// 解析配置文件
	var config Config
	err = json.Unmarshal(configData, &config)
	if err != nil {
		fmt.Printf("Failed to parse config file: %v\n", err)
		os.Exit(1)
	}

	// 执行字节替换
	err = replaceBytes(config)
	if err != nil {
		fmt.Printf("Error during byte replacement: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Replacement successful!")
}
