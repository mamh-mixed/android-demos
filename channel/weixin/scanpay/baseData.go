package scanpay

// BaseData 只是为了注入签名方便
// TODO 写一个类似 json 和 xml 的工具类，生成 queryString
type BaseData interface {
	GenSign()
}
