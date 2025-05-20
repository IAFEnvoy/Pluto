package java

import (
	"fmt"
	"strings"
)

func ObfuscateMethodSignature(signature string, obfuscationMap map[string]string) string {
	paramStart := strings.Index(signature, "(")
	paramEnd := strings.Index(signature, ")")
	if paramStart == -1 || paramEnd == -1 || paramStart > paramEnd {
		return signature // 无效签名格式，保持原样
	}
	paramTypes := signature[paramStart+1 : paramEnd]
	returnType := signature[paramEnd+1:]
	processedParams := ObfuscateTypeSignature(paramTypes, obfuscationMap)
	processedReturn := processType(returnType, obfuscationMap)
	return fmt.Sprintf("(%s)%s", processedParams, processedReturn)
}

func ObfuscateTypeSignature(signature string, obfuscationMap map[string]string) string {
	var result strings.Builder
	i := 0
	for i < len(signature) {
		switch signature[i] {
		case '[': // 数组类型
			result.WriteByte('[')
			i++
		case 'L': // 对象类型 (类或接口)
			// 查找类型结束符 ';'
			end := i + 1
			for end < len(signature) && signature[end] != ';' {
				end++
			}
			if end < len(signature) {
				classSig := signature[i : end+1]
				// 查找混淆映射
				if obfuscated, exists := obfuscationMap[classSig]; exists {
					result.WriteString(obfuscated)
				} else {
					result.WriteString(classSig) // 未找到映射，保持原样
				}
				i = end + 1
			} else {
				// 无效类型格式，保持原样
				result.WriteString(signature[i:])
				i = len(signature)
			}
		default: // 基本类型 (Z, B, C, S, I, J, F, D, V)
			result.WriteByte(signature[i])
			i++
		}
	}
	return result.String()
}

// 处理单个类型
func processType(t string, obfuscationMap map[string]string) string {
	if len(t) == 0 {
		return ""
	}
	if t[0] == '[' { // 数组类型
		return "[" + processType(t[1:], obfuscationMap)
	}
	if t[0] == 'L' { // 对象类型
		if obfuscated, exists := obfuscationMap[t]; exists {
			return obfuscated
		}
	}
	return t // 基本类型或未找到映射
}
