package handler

//type createUserRequest struct {
//	Phone    string `json:"phone" binding:"required"`
//	Username string `json:"username" binding:"required,alphanumunicode"`
//	Password string `json:"password" binding:"required,min=6"`
//	Email    string `json:"email" binding:"required,email"`
//}
//
//type userResponse struct {
//	Phone       string    `json:"phone"`
//	Username    string    `json:"username"`
//	Email       string    `json:"email"`
//	AccessToken string    `json:"access_token"`
//	CreatedAt   time.Time `json:"created_at"`
//}
//
//func newUserResponse(user model.User) userResponse {
//	return userResponse{
//		Phone:     user.Phone,
//		Username:  user.Username,
//		CreatedAt: user.CreatedAt,
//	}
//}

//func (h *Handler) CreateUser(c fiber.Ctx) error {
//	var req createUserRequest
//	if err := c.Bind().JSON(&req); err != nil {
//		return c.JSON(errorResponse(err))
//	}
//
//	// 手动验证手机号格式
//	if !regexp.MustCompile(`^1[3-9]\d{9}$`).MatchString(req.Phone) {
//		return c.JSON("手机号格式不正确")
//
//	}
//
//	hashPassword, err := util.HashPassword(req.Password)
//	if err != nil {
//		return c.JSON(errorResponse(err))
//	}
//	arg := model.CreateUserParams{
//		Phone:    req.Phone,
//		Username: req.Username,
//		Password: hashPassword,
//	}
//
//	user, err := h.di.Db.CreateUser(c.Context(), arg)
//	if err != nil {
//		//如果是数据库出错
//		if pqErr, ok := err.(*pq.Error); ok {
//			switch pqErr.Code.Name() {
//			case "unique_violation":
//				return c.JSON("手机号或用户名已存在")
//			}
//		}
//		return c.JSON(errorResponse(err))
//	}
//	resp := userResponse{
//		Phone:     user.Phone,
//		Username:  user.Username,
//		CreatedAt: user.CreatedAt,
//	}
//	return c.JSON(resp)
//}

//type loginUserRequest struct {
//	Phone    string `json:"phone" binding:"required"`
//	Password string `json:"password" binding:"required,min=6"`
//}
//
//type loginUserResponse struct {
//	AccessToken string
//	User        userResponse
//}
//
//func (h *Handler) Login(c fiber.Ctx) error {
//	var req loginUserRequest
//	if err := c.Bind().JSON(&req); err != nil {
//		return c.JSON(errorResponse(err))
//	}
//
//	// 手动验证手机号格式
//	if !regexp.MustCompile(`^1[3-9]\d{9}$`).MatchString(req.Phone) {
//		return c.JSON(fiber.Map{"error": "手机号格式不正确"})
//
//	}
//
//	user, err := h.svc.Db.FindUserByPhone(c.Context(), req.Phone)
//	if err != nil {
//		if errors.Is(err, sql.ErrNoRows) {
//			return c.JSON("账号或密码错误")
//
//		}
//		return c.JSON(errorResponse(err))
//	}
//	err = util.CheckPassword(req.Password, user.Password)
//	if err != nil {
//		return c.JSON(errorResponse(err))
//	}
//
//	//create access token for user
//	accessToken, err := h.svc.TokenMaker.CreateToken(
//		uint64(user.ID),
//		h.svc.Config.AccessTokenDuration,
//	)
//	if err != nil {
//		return c.JSON(errorResponse(err))
//	}
//	rsp := loginUserResponse{
//		AccessToken: accessToken,
//		User:        newUserResponse(user),
//	}
//	return c.JSON(rsp)
//}
//
//type updateUserRequest struct {
//	Id       int64  `json:"id" binding:"required,number"`
//	Username string `json:"username" binding:"required,alphanumunicode"`
//}
//
//func (h *Handler) Update(c fiber.Ctx) error {
//	var req updateUserRequest
//
//	if err := c.Bind().JSON(req); err != nil {
//		return c.JSON(err)
//	}
//
//	payload := c.Locals(middleware.AuthorizationPayloadKey).(*token.Payload)
//	if int64(payload.UserId) != req.Id {
//		return c.JSON("参数异常")
//	}
//
//	arg := model.UpdateUserInfoParams{
//		ID:       req.Id,
//		Username: req.Username,
//	}
//	result, err := h.svc.Db.UpdateUserInfo(c.Context(), arg)
//	if err != nil {
//		if pqErr, ok := err.(*pq.Error); ok {
//			switch pqErr.Code.Name() {
//			case "unique_violation":
//				return c.JSON("用户名已存在")
//			}
//		}
//		return c.JSON(err)
//	}
//
//	rowsAffected, err := result.RowsAffected()
//	if err != nil {
//		return c.JSON(err.Error())
//	}
//
//	return c.JSON(fiber.Map{
//		"rowsAffected": rowsAffected,
//	})
//}
