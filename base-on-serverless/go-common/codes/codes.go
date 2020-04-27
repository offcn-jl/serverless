/*
   @Time : 2020/4/26 4:28 下午
   @Author : ShadowWalker
   @Email : master@rebeta.cn
   @File : codes
   @Software: GoLand
*/

package codes

const (
	// 授权错误
	NotCertifiedCORS = 13001
)

var errorText = map[int]string{
	// 授权错误
	NotCertifiedCORS: "Not Certified ( CORS )",
}

func ErrorText(code int) string {
	return errorText[code]
}
