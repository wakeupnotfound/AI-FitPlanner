package request

// 用户注册请求
type RegisterRequest struct {
	Username        string `json:"username" binding:"required,min=3,max=20,alphanum"`
	Email           string `json:"email" binding:"required,email,email_format"`
	Password        string `json:"password" binding:"required,min=8,max=20,password_strength"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

// 用户登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"` // 用户名或邮箱
	Password string `json:"password" binding:"required"`
}

// 刷新token请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// 更新用户信息请求
type UpdateUserRequest struct {
	Username string `json:"username" binding:"omitempty,min=3,max=20,alphanum"`
	Nickname string `json:"nickname" binding:"omitempty,min=1,max=50"`
	Phone    string `json:"phone" binding:"omitempty,e164"`
	Avatar   string `json:"avatar" binding:"omitempty,avatar"`
}

// 更新密码请求
type UpdatePasswordRequest struct {
	OldPassword     string `json:"old_password" binding:"required,min=1,max=100"`
	NewPassword     string `json:"new_password" binding:"required,min=8,max=20,password_strength"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword"`
}
