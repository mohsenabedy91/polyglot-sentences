package requests

type Header struct {
	UserID uint64 `header:"userID" binding:"required,number"`
	JTI    string `header:"jti" binding:"required,uuid"`
	EXP    int64  `header:"exp" binding:"required"`
}
