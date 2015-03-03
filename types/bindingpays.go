package types

// Binding 绑定关系
type Binding struct {
	cardNum   string
	bindingId string
	// ...
}

// BindingCreateIn 绑定支付请求
type BindingCreateIn struct {
	merId string // 商户 Id
	// ...
}

// BindingCreateOut 绑定支付响应
type BindingCreateOut struct {
	merId string // 商户 Id
	// ...
}
