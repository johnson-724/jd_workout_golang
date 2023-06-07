package auth

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"jd_workout_golang/app/models"
	repo "jd_workout_golang/app/repositories/user"
	"jd_workout_golang/app/services/jwtHelper"
	db "jd_workout_golang/lib/database"
	google "jd_workout_golang/lib/google"
	"net/http"
)

// LoginAction logs in a user with the provided email and password,
// and generates a JWT token for the user.
//
// @Summary Login user
// @Description Logs in a user with the provided email and password, and generates a JWT token for the user
// @Tags Auth
// @Accept x-www-form-urlencoded
// @Produce json
// @Param email formData string true "User email"
// @Param password formData string true "User password"
// @Success 200 {string} string "{'message': 'login success', 'token': 'JWT token', 'reset': '0'}"
// @Failure 422 {string} string "{'message': '缺少必要欄位', 'error': 'error message'}"
// @Failure 422 {string} string "{'message': '帳號或密碼錯誤', 'error': 'error message'}"
// @Router /login [post]
func LoginAction(c *gin.Context) {
	loginForm := struct {
		Email    string `json:"email" form:"email" binding:"required"`
		Password string `json:"password" form:"password" binding:"required"`
	}{}

	if err := c.ShouldBind(&loginForm); err != nil {
		c.JSON(422, gin.H{
			"message": "缺少必要欄位",
			"error":   err.Error(),
		})

		return
	}

	user := models.User{
		Email:    loginForm.Email,
		Password: loginForm.Password,
	}

	if _, err := validateLogin(&user, db.Connection); err != nil {
		c.JSON(422, gin.H{
			"message": "帳號或密碼錯誤",
			"error":   err.Error(),
		})

		return
	}

	if user.EmailVerified == 0 {
		c.JSON(422, gin.H{
			"message": "信箱尚未驗證",
		})

		return
	}

	token, _ := jwtHelper.GenerateToken(&user)

	c.JSON(200, gin.H{
		"message": "login success",
		"token":   token,
		"reset":   user.ResetPassword,
	})
}

func validateLogin(user *models.User, db *gorm.DB) (bool, error) {
	record := models.User{}

	result := db.Where("email = ?", user.Email).First(&record)

	if result.Error != nil || result.RowsAffected == 0 {
		return false, result.Error
	}

	error := bcrypt.CompareHashAndPassword([]byte(record.Password), []byte(user.Password))

	if error != nil {
		return false, error
	}

	*user = record

	return true, nil
}

// LoginWithGoogleAction logs in a user with the provided google token,
// and generates a JWT token for the user.
func LoginWithGoogleAction(c *gin.Context) {
	redirectURL := google.CreateGoogleOAuthURL()
	c.Redirect(http.StatusSeeOther, redirectURL) // http.StatusSeeOther 為 303
}

// LoginWithAuthkAction logs in a user with the provided google token,
func LoginWithGoogleAuthAction(c *gin.Context) {
	token, err := google.GetAccessToken(c.Query("code"))

	if err != nil {
		c.JSON(422, gin.H{
			"message": "登入失敗",
			"error":   err.Error(),
		})

		return
	}

	userInfo, err := google.GetUserInfo(token)

	if err != nil {
		c.JSON(422, gin.H{
			"message": "登入失敗, 無法取得使用者授權資訊",
			"error":   err.Error(),
		})

		c.Abort()

		return
	}

	validateOnGoogle(userInfo, c)
}

type LoginWithGoogleAccessTokenForm struct {
	Token string `json:"token" form:"token" binding:"required"`
}

// Login with google oauth2 access token
// and generates a JWT token for the user.
//
// @Summary Login user with google oauth2 access token
// @Description Logs in a user with the google oauth2 access token, and generates a JWT token for the user
// @Tags Auth
// @Accept json
// @Produce json
// @Param LoginWithGoogleAccessTokenForm body LoginWithGoogleAccessTokenForm true "LoginWithGoogleAccessTokenForm"
// @Success 200 {string} string "{'message': 'login success', 'token': 'JWT token'}"
// @Failure 422 {string} string "{'message': '缺少必要欄位', 'error': 'error message'}"
// @Router /login/google/access-token [post]
func LoginWithGoogleAccessTokenAction(c *gin.Context) {
	loginForm := LoginWithGoogleAccessTokenForm{}

	if err := c.ShouldBind(&loginForm); err != nil {
		c.JSON(422, gin.H{
			"message": "缺少必要欄位",
			"error":   err.Error(),
		})

		c.Abort()

		return
	}

	userInfo, err := google.GetUserInfoWithAccessToken(loginForm.Token)

	if err != nil {
		c.JSON(422, gin.H{
			"message": "登入失敗, 無法取得使用者授權資訊",
			"error":   err.Error(),
		})

		c.Abort()

		return
	}

	validateOnGoogle(userInfo, c)
}

func bindUserWithThirdPartyAccount(thirdPartyInfo *google.UserInfo, user *models.User) {
	// create user with third party account
	if user.ID == 0 {
		user.Username = thirdPartyInfo.Name
		user.Email = thirdPartyInfo.Email
		user.EmailVerified = 1

		repo.Create(user)

		return
	}

	// user exist and email not verified
	if user.ID != 0 && user.EmailVerified != 1 {
		user.EmailVerified = 1

		repo.Update(user)

		return
	}
}

func validateOnGoogle(userInfo *google.UserInfo, c *gin.Context) {
	user, err := repo.GetUserByEmail(userInfo.Email)

	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(422, gin.H{
			"message": "登入失敗, 無法取得使用者資訊",
			"error":   err.Error(),
		})

		c.Abort()

		return
	}

	bindUserWithThirdPartyAccount(userInfo, user)

	jwtToken, _ := jwtHelper.GenerateToken(user)

	c.JSON(200, gin.H{
		"message": "login success",
		"token":   jwtToken,
	})
}
