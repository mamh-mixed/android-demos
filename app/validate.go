package app

import (
	"github.com/CardInfoLink/quickpay/model"
)

// requestDataValidate 请求数据验证
func requestDataValidate(req *reqParams) *model.AppResult {
	if req.UserName != "" {
		if len(req.UserName) > 50 {
			return model.NewAppResult("FAIL", "username 长度过长")
		}
	}

	if req.OrderNum != "" {
		if len(req.OrderNum) > 50 {
			return model.NewAppResult("FAIL", "orderNum 长度过长")
		}
	}

	if req.Date != "" {
		if len(req.Date) > 10 {
			return model.NewAppResult("FAIL", "date 长度过长")
		}
	}

	if req.Month != "" {
		if len(req.Month) > 6 {
			return model.NewAppResult("FAIL", "month 长度过长")
		}
	}

	if req.Status != "" {
		if len(req.Status) > 20 {
			return model.NewAppResult("FAIL", "status 长度过长")
		}
	}

	// if req.Index != "" {
	// 	if len(req.Index) > 10 {
	// 		return model.NewAppResult("FAIL", "index 长度过长")
	// 	}
	// }

	if req.Password != "" {
		if len(req.Password) > 50 {
			return model.NewAppResult("FAIL", "password 长度过长")
		}
	}

	if req.OldPassword != "" {
		if len(req.OldPassword) > 50 {
			return model.NewAppResult("FAIL", "oldpassword 长度过长")
		}
	}

	if req.NewPassword != "" {
		if len(req.NewPassword) > 50 {
			return model.NewAppResult("FAIL", "newpassword 长度过长")
		}
	}

	if req.Transtime != "" {
		if len(req.Transtime) > 30 {
			return model.NewAppResult("FAIL", "transtime 长度过长")
		}
	}

	if req.Province != "" {
		if len(req.Province) > 50 {
			return model.NewAppResult("FAIL", "province 长度过长")
		}
	}

	if req.City != "" {
		if len(req.City) > 50 {
			return model.NewAppResult("FAIL", "city 长度过长")
		}
	}

	if req.BankOpen != "" {
		var r = []rune(req.BankOpen)
		if len(r) > 100 {
			return model.NewAppResult("FAIL", "bank_open 长度过长")
		}
	}

	if req.Payee != "" {
		var r = []rune(req.Payee)
		if len(r) > 100 {
			return model.NewAppResult("FAIL", "payee 长度过长")
		}
	}

	if req.BranchBank != "" {
		var r = []rune(req.BranchBank)
		if len(r) > 200 {
			return model.NewAppResult("FAIL", "branch_bank 长度过长")
		}
	}

	if req.BankNo != "" {
		if len(req.BankNo) > 50 {
			return model.NewAppResult("FAIL", "bankNo 长度过长")
		}
	}

	if req.PayeeCard != "" {
		if len(req.PayeeCard) > 50 {
			return model.NewAppResult("FAIL", "payee_card 长度过长")
		}
	}

	if req.PhoneNum != "" {
		if len(req.PhoneNum) > 20 {
			return model.NewAppResult("FAIL", "phone_num 长度过长")
		}
	}

	if req.Email != "" {
		if len(req.Email) > 50 {
			return model.NewAppResult("FAIL", "email 长度过长")
		}
	}

	if req.Sign != "" {
		if len(req.Sign) > 32 {
			return model.NewAppResult("FAIL", "sign 长度过长")
		}
	}

	return nil
}
