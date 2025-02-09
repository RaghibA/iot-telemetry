package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/RaghibA/iot-telemetry/authn-service/internal/db"
	"github.com/RaghibA/iot-telemetry/authn-service/internal/models"
	"github.com/RaghibA/iot-telemetry/authn-service/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type RegisterBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type LoginBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func emailInUse(email string) (bool, error) {
	var ct int64

	err := db.IotDb.Db.Model(&models.User{}).Where("email = ?", email).Count(&ct).Error
	if err != nil {
		return false, err
	}
	return ct > 0, nil
}

func usernameInUse(username string) (bool, error) {
	var ct int64

	err := db.IotDb.Db.Model(&models.User{}).Where("username = ?", username).Count(&ct).Error
	if err != nil {
		return false, err
	}
	return ct > 0, nil
}

func HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Health OK",
	})
}

func RegisterUserHandler(c *gin.Context) {
	var userBody RegisterBody
	validate := validator.New()

	if c.Bind(&userBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON Body",
			"code":  400001,
		})
		return
	}

	err := validate.Var(userBody.Email, "required,email")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email",
			"code":  400002,
		})
		return
	}

	err = validate.Var(userBody.Username, "required,min=5")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Username must be at least 5 characters",
			"code":  400003,
		})
		return
	}

	err = validate.Var(userBody.Password, "required,min=8")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Password must be at least 8 characters",
			"code":  400004,
		})
		return
	}

	exists, err := emailInUse(userBody.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Service Error",
			"code":  500002,
		})
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Email in use",
			"code":  409001,
		})
		return
	}

	exists, err = usernameInUse(userBody.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Service Error",
			"code":  500002,
		})
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Username in use",
			"code":  409002,
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userBody.Password), 14)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
			"code":  500003,
		})
		return
	}

	// issue api key for device & store hashed key in db
	apiKey, err := utils.GenerateAPIKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
			"code":  500005,
		})
		return
	}

	hashedKey, err := utils.HashAPIKey(apiKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
			"code":  500006,
		})
	}

	user := models.User{
		Username: userBody.Username,
		Password: hashedPassword,
		Email:    userBody.Email,
		UserID:   uuid.New().String(),
		APIKey:   hashedKey,
	}

	err = db.IotDb.Db.Create(&user).Error
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
			"code":    500004,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":  "Account Created",
		"username": user.Username,
		"email":    user.Email,
		"apiKey":   apiKey,
	})
}

func LoginHandler(c *gin.Context) {
	var loginCreds LoginBody

	if c.Bind(&loginCreds) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Login Failed",
			"code":    400005,
		})
		return
	}

	// find associated account
	var account models.User
	err := db.IotDb.Db.Where("username = ?", loginCreds.Username).First(&account).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
			"code":    500002,
		})
		return
	}

	if bcrypt.CompareHashAndPassword(account.Password, []byte(loginCreds.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "login failed",
			"code":    401001,
		})
		return
	}

	accessToken, err := utils.GenerateJWT(true, account, time.Now().Add(time.Hour*24*7).Unix(), time.Now().Unix())
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
			"code":    500003,
		})
		return
	}

	refreshToken, err := utils.GenerateJWT(false, account, time.Now().Add(time.Minute*30).Unix(), time.Now().Unix())
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
			"code":    500003,
		})
		return
	}

	c.SetCookie("refresh_token", refreshToken, 60*60*24*7, "/", "", true, true) // 7 days
	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
	})
}

func RefreshHandler(c *gin.Context) {
	var user models.User
	userId, ok := c.Get("userID")
	log.Println(userId)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
			"code":    500004,
		})
		return
	}

	err := db.IotDb.Db.Where("user_id = ?", userId).First(&user).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "unable to find user",
			"code":    404001,
		})
		return
	}

	accessToken, err := utils.GenerateJWT(true, user, time.Now().Add(time.Minute*15).Unix(), time.Now().Unix())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
			"code":    500005,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken": accessToken,
	})
}

func LogoutHandler(c *gin.Context) {
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "logged out succesfully",
	})
}

func DeactivateHandler(c *gin.Context) {
	userId, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "user not found",
			"code":    404002,
		})
		return
	}

	var user models.User
	err := db.IotDb.Db.Where("user_id = ?", userId).First(&user).Error
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
			"code":    500006,
		})
		return
	}

	err = db.IotDb.Db.Where("user_id = ?", userId).Delete(&user).Error
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
			"code":    500007,
		})
		return
	}

	//! Make some cascading delete for all user resources

	c.JSON(http.StatusNoContent, gin.H{})
}

func GenerateAPIKeyHandler(c *gin.Context) {
	apiKey, err := utils.GenerateAPIKey()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
			"code":  500001,
		})
		return
	}

	hashedKey, err := utils.HashAPIKey(apiKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
			"code":  500002,
		})
		return
	}

	userId, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
			"code":  500003,
		})
		return
	}

	err = db.IotDb.Db.Model(&models.User{}).Where("user_id = ?", userId).Update("api_key", hashedKey).Error
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
			"code":  500004,
		})
		return
	}

	err = db.IotDb.Db.Model(&models.KafkaACL{}).Where("user_id = ?", userId).Update("api_key", hashedKey).Error
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
			"code":  500004,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"apiKey": apiKey,
	})
}
