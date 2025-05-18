package misc

import "pluto/convert"

type SingleInfo struct {
	Name      string
	Class     string
	Signature string
	Type      string
}

func PackClassInfo(name string) SingleInfo {
	return SingleInfo{
		Name:      name,
		Class:     "",
		Signature: convert.ClassToByteCodeSignature(name),
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
