package alicloud

import (
	"fmt"
)

type ProxyStruct struct {
	ProxyId         string
	ProxyType       string
	Weight          int
	HTTPProxyConfig struct {
		Address  string
		Port     int
		User     string
		Password string
	}
	Socks5ProxyConfig struct {
		Address  string
		Port     int
		User     string
		Password string
	}
	SSHProxyConfig struct {
		Address  string
		Port     int
		User     string
		Password string
	}
}

func getPasswordFromProxies(proxies []ProxyStruct, proxy map[string]interface{}) string {
	for _, v := range proxies {
		if v.ProxyType == proxy["ProxyType"] {
			if v.ProxyType == "HTTPProxy" {
				l1 := fmt.Sprintf("%v %v %v", v.HTTPProxyConfig.User, v.HTTPProxyConfig.Address, v.HTTPProxyConfig.Port)
				l2 := fmt.Sprintf("%v %v %v", proxy["HTTPProxyConfig"].(map[string]interface{})["User"],
					proxy["HTTPProxyConfig"].(map[string]interface{})["Address"],
					proxy["HTTPProxyConfig"].(map[string]interface{})["Port"])
				if l1 == l2 {
					return v.HTTPProxyConfig.Password
				}
			} else if v.ProxyType == "Socks5Proxy" {
				l1 := fmt.Sprintf("%v %v %v", v.Socks5ProxyConfig.User, v.Socks5ProxyConfig.Address, v.Socks5ProxyConfig.Port)
				l2 := fmt.Sprintf("%v %v %v", proxy["Socks5ProxyConfig"].(map[string]interface{})["User"],
					proxy["Socks5ProxyConfig"].(map[string]interface{})["Address"],
					proxy["Socks5ProxyConfig"].(map[string]interface{})["Port"])
				if l1 == l2 {
					return v.Socks5ProxyConfig.Password
				}
			} else if v.ProxyType == "SSHProxy" {
				l1 := fmt.Sprintf("%v %v %v", v.SSHProxyConfig.User, v.SSHProxyConfig.Address, v.SSHProxyConfig.Port)
				l2 := fmt.Sprintf("%v %v %v", proxy["SSHProxyConfig"].(map[string]interface{})["User"],
					proxy["SSHProxyConfig"].(map[string]interface{})["Address"],
					proxy["SSHProxyConfig"].(map[string]interface{})["Port"])
				if l1 == l2 {
					return v.SSHProxyConfig.Password
				}
			}
		}
	}

	return ""
}

func isProxyEqual(source ProxyStruct, target ProxyStruct) bool {
	if source.ProxyType == target.ProxyType && source.Weight == target.Weight {
		if target.ProxyType == "HTTPProxy" {
			if source.HTTPProxyConfig.Address == target.HTTPProxyConfig.Address &&
				source.HTTPProxyConfig.Password == target.HTTPProxyConfig.Password &&
				source.HTTPProxyConfig.Port == target.HTTPProxyConfig.Port &&
				source.HTTPProxyConfig.User == target.HTTPProxyConfig.User {
				return true
			}
		} else if target.ProxyType == "Socks5Proxy" {
			if source.Socks5ProxyConfig.Address == target.Socks5ProxyConfig.Address &&
				source.Socks5ProxyConfig.Password == target.Socks5ProxyConfig.Password &&
				source.Socks5ProxyConfig.Port == target.Socks5ProxyConfig.Port &&
				source.Socks5ProxyConfig.User == target.Socks5ProxyConfig.User {
				return true
			}
		} else if target.ProxyType == "SSHProxy" {
			if source.SSHProxyConfig.Address == target.SSHProxyConfig.Address &&
				source.SSHProxyConfig.Password == target.SSHProxyConfig.Password &&
				source.SSHProxyConfig.Port == target.SSHProxyConfig.Port &&
				source.SSHProxyConfig.User == target.SSHProxyConfig.User {
				return true
			}
		}
	}

	return false
}

func isProxyChange(source ProxyStruct, target ProxyStruct) bool {
	if source.ProxyId == target.ProxyId {
		return !isProxyEqual(source, target)
	}

	return false
}

func isProxyContain(proxies []ProxyStruct, target ProxyStruct) bool {
	for _, source := range proxies {
		if isProxyEqual(source, target) {
			return true
		}
	}

	return false
}

func isProxyDifference(proxies []ProxyStruct, target ProxyStruct) string {
	for _, source := range proxies {
		if isProxyEqual(source, target) {
			return "equal"
		} else if isProxyChange(source, target) {
			return source.ProxyId
		}
	}

	return "diff"
}

func compareProxies(oldObjects []ProxyStruct, newObjects []ProxyStruct) ([]ProxyStruct, []ProxyStruct, []ProxyStruct) {
	add := make([]ProxyStruct, 0)
	remove := make([]ProxyStruct, 0)
	update := make([]ProxyStruct, 0)

	nSize := len(newObjects)
	oSize := len(oldObjects)

	if nSize == 0 && oSize == 0 {
		return add, remove, update
	} else if nSize == 0 && oSize > 0 {
		return add, oldObjects, update
	} else if nSize > 0 && oSize == 0 {
		return newObjects, remove, update
	}

	for _, n := range newObjects {
		r := isProxyDifference(oldObjects, n)
		if r == "diff" {
			add = append(add, n)
		} else if r != "equal" {
			n.ProxyId = r
			update = append(update, n)
		}
	}

	for _, o := range oldObjects {
		r := isProxyDifference(newObjects, o)
		if r == "diff" {
			remove = append(remove, o)
		}
	}

	return add, remove, update
}
