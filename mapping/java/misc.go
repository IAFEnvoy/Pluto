package java

import "strings"

type SingleInfo struct {
	Name      string
	Class     string
	Signature string
	Type      string
}

func PackClassInfo(name string) SingleInfo {
	return SingleInfo{
		Name:      PathToClassName(name),
		Class:     name,
		Signature: ClassToByteCodeSignature(name),
		Type:      "class",
	}
}

func PackMethodInfo(name, class, signature string) SingleInfo {
	return SingleInfo{
		Name:      name,
		Class:     class,
		Signature: signature,
		Type:      "method",
	}
}

func PackFieldInfo(name, class, signature string) SingleInfo {
	return SingleInfo{
		Name:      name,
		Class:     class,
		Signature: signature,
		Type:      "field",
	}
}

func BuildMapping(mapping *map[SingleInfo]SingleInfo) *Mappings {
	result := Mappings{
		AllMapping:  make(map[SingleInfo]SingleInfo, len(*mapping)),
		NotchByName: make(map[string][]SingleInfo),
		NamedByName: make(map[string][]SingleInfo),
	}
	for k, v := range *mapping {
		result.AllMapping[k] = v
		result.NotchByName[k.Name] = append(result.NotchByName[k.Name], k)
		result.NamedByName[v.Name] = append(result.NamedByName[v.Name], v)
	}
	return &result
}

func PathToClassName(path string) string {
	split := strings.Split(path, "/")
	return split[len(split)-1]
}
