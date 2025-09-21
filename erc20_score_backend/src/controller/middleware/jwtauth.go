package middleware

type Login struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type Register struct {
	Login
	Email string `json:"email" binding:"required"`
}

// JWTAuthMiddleware 基于JWT的认证中间件
// func JWTAuthMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		result := base.Result{}
// 		token := c.Request.Header.Get("Authorization")
// 		if token == "" {
// 			c.JSON(200, result.Fail(errorcode.Token_missing))
// 			c.Abort() //c.Abort() 通过设置 c.index = abortIndex（默认值为 63）来中断当前挂起的请求处理链
// 			return
// 		}

// 		config := config.GetConfig()
// 		key := config.Userpassword.Key
// 		jwtKey := []byte(key)

// 		jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
// 			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 				return nil, jwt.ErrInvalidKeyType
// 			}
// 			return jwtKey, nil
// 		})

// 		if err != nil {
// 			c.JSON(200, result.Fail(errorcode.Token_invalid))
// 			c.Abort()
// 			return
// 		}

// 		if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok && jwtToken.Valid {
// 			logger.Log.Printf("claims is :%v", claims)
// 			c.Set("userId", claims["userId"]) // 将claims信息设置到context中，后续可以通过c.Get("user")获取到用户信息
// 			c.Set("username", claims["username"])
// 			c.Next() // 继续执行后续的请求处理函数
// 		} else {
// 			c.JSON(200, result.Fail(errorcode.Token_invalid))
// 			c.Abort()
// 			return
// 		}
// 	}
// }

// func UserRegister(c *gin.Context) {

// 	result := base.Result{}
// 	var userRegister Register
// 	if err := c.ShouldBindJSON(&userRegister); err != nil {
// 		c.JSON(200, result.FailWeb(400, err.Error()))
// 		return
// 	}

// 	// 加密密码
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRegister.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		c.JSON(200, result.FailWeb(http.StatusInternalServerError, err.Error()))
// 		return
// 	}
// 	userRegister.Password = string(hashedPassword)

// 	user := &dao.User{
// 		Username: userRegister.Username,
// 		Password: userRegister.Password,
// 		Email:    userRegister.Email,
// 	}

// 	db := dao.InitDb()
// 	err = db.Transaction(func(tx *gorm.DB) error {

// 		err := dao.UserInsert(db, user)
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		c.JSON(200, result.FailCommon(errorcode.Operation_error, err.Error()))
// 		return
// 	}

// 	logger.Log.Printf("---------------------------user:%v", user)

// 	c.JSON(200, result.Sucess())

// }

// func UserLogin(c *gin.Context) {

// 	var login Login
// 	if err := c.ShouldBindJSON(&login); err != nil {
// 		c.JSON(400, gin.H{"error": err.Error()})
// 		return
// 	}
// 	result := base.Result{}

// 	db := dao.InitDb()
// 	storedUser := dao.UserGet(db, login.Username)

// 	//查无此人
// 	if storedUser.ID == 0 {
// 		c.JSON(200, result.Fail(errorcode.Login_fail))
// 		return
// 	}

// 	// 验证密码
// 	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(login.Password)); err != nil {
// 		c.JSON(200, result.Fail(errorcode.Login_fail))
// 		return
// 	}

// 	// 生成 JWT
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"userId":   storedUser.ID,
// 		"username": storedUser.Username,
// 		"exp":      time.Now().Add(time.Hour * 1).Unix(),
// 	})

// 	config := config.GetConfig()
// 	key := config.Userpassword.Key

// 	//生成token
// 	tokenString, _ := token.SignedString([]byte(key))
// 	tokenMap := map[string]string{
// 		"token": tokenString,
// 	}

// 	c.JSON(200, result.SucessData(tokenMap))

// }
