package adapter

import(
    "strings"
)

type Abstract struct {
    pathPrefix string
    pathSeparator string
}

// 设置前缀
func (at *Abstract) SetPathPrefix(prefix string) error {
    at.pathPrefix = strings.TrimSuffix(prefix, "\\/") + at.pathSeparator
}

// 获取前缀
func (at *Abstract) GetPathPrefix() string {
    return at.pathPrefix
}

// 添加前缀
func (at *Abstract) ApplyPathPrefix(path string) string {
    return at.GetPathPrefix() + strings.TrimPrefix(path, "\\/")
}

// 移除前缀
func (at *Abstract) RemovePathPrefix(path string) string {
    prefix := at.GetPathPrefix()
    return strings.TrimPrefix(path, prefix)
}
