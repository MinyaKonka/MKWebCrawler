package cookie

import (
	"code.google.com/p/go.net/publicsuffix"
	"net/http"
	"net/http/cookiejar"
)

// 创建http.CookieJar类型的值
func NewCookiejar() http.CookieJar {
	options := &cookiejar.Options{PublicSuffixList: &mk_publicSuffixList{}}

	jar, _ := cookiejar.New(options)

	return jar
}

// cookiejar.PublicSuffixList接口实现类型
type mk_publicSuffixList struct{}

func (list *mk_publicSuffixList) PublicSuffix(domain string) string {
	suffix, _ := list.PublicSuffix(domain)
	return suffix
}

func (list *mk_publicSuffixList) String() string {
	return "Web crawler - public suffix list (rev 1.0) power by 'code.google.com/p/go.net/publicsuffix'"
}
