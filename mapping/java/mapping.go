package java

import (
	"sort"
	"strings"
)

type Mappings struct {
	AllMapping  map[SingleInfo]SingleInfo
	NotchByName map[string][]SingleInfo
	NamedByName map[string][]SingleInfo
}

type InfoForNetwork struct {
	Notch      SingleInfo `json:"notch"`
	Named      SingleInfo `json:"named"`
	Translated SingleInfo `json:"translated,omitzero"`
}

type searchResult struct {
	info InfoForNetwork
	// 匹配类型权重: 完全匹配键=3, 前缀匹配键=2, 包含匹配键=1
	matchType int
	// 类型权重: class=3, method=2, field=1, 其他=0
	typeWeight int
	// 匹配名称权重: Named匹配=2, Notch匹配=1
	nameType int
}

func (m *Mappings) Search(keyword string, maxCount int) []InfoForNetwork {
	if maxCount <= 0 {
		return []InfoForNetwork{}
	}

	results := make([]searchResult, 0)
	keyword = strings.ToLower(keyword)
	seen := make(map[string]struct{}) // 用于去重

	// 搜索NotchByName
	for name, infos := range m.NotchByName {
		if !strings.Contains(strings.ToLower(name), keyword) {
			continue
		}
		matchType := getMatchType(name, keyword)
		for _, info := range infos {
			if named, exists := m.AllMapping[info]; exists {
				key := info.Name + "|" + named.Name
				if _, ok := seen[key]; !ok {
					seen[key] = struct{}{}
					results = append(results, searchResult{
						info: InfoForNetwork{
							Notch: info,
							Named: named,
						},
						typeWeight: getTypeWeight(info.Type),
						nameType:   1, // Notch匹配
						matchType:  matchType,
					})
				}
			}
		}
	}

	// 搜索NamedByName
	for name, infos := range m.NamedByName {
		if !strings.Contains(strings.ToLower(name), keyword) {
			continue
		}
		matchType := getMatchType(name, keyword)
		for _, info := range infos {
			// 查找对应的Notch
			var notch SingleInfo
			found := false
			for n, named := range m.AllMapping {
				if named == info {
					notch = n
					found = true
					break
				}
			}
			if found {
				key := notch.Name + "|" + info.Name
				if _, ok := seen[key]; !ok {
					seen[key] = struct{}{}
					results = append(results, searchResult{
						info: InfoForNetwork{
							Notch: notch,
							Named: info,
						},
						typeWeight: getTypeWeight(info.Type),
						nameType:   2, // Named匹配
						matchType:  matchType,
					})
				}
			}
		}
	}

	// 按相关性排序
	sort.Slice(results, func(i, j int) bool {
		a, b := results[i], results[j]
		// 按匹配类型排序
		if a.matchType != b.matchType {
			return a.matchType > b.matchType
		}
		// 按类型权重排序
		if a.typeWeight != b.typeWeight {
			return a.typeWeight > b.typeWeight
		}
		// 名称类型排序(Named > Notch)
		if a.nameType != b.nameType {
			return a.nameType > b.nameType
		}
		// 按名称字母顺序
		return a.info.Notch.Name < b.info.Notch.Name
	})

	// 截取前maxCount个结果
	if len(results) > maxCount {
		results = results[:maxCount]
	}

	// 转换为最终结果
	final := make([]InfoForNetwork, len(results))
	for i, res := range results {
		final[i] = res.info
	}

	return final
}

// 判断匹配类型并返回权重
func getMatchType(name, keyword string) int {
	nameLower := strings.ToLower(name)
	keywordLower := strings.ToLower(keyword)
	if nameLower == keywordLower {
		return 3 // 完全匹配
	}
	if strings.HasPrefix(nameLower, keywordLower) {
		return 2 // 前缀匹配
	}
	return 1 // 包含匹配
}

// 获取Type的权重
func getTypeWeight(t string) int {
	switch strings.ToLower(t) {
	case "class":
		return 3
	case "method":
		return 2
	case "field":
		return 1
	default:
		return 0
	}
}
