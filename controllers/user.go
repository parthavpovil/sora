package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/parthavpovil/sora/database"
	"github.com/parthavpovil/sora/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var secret = []byte(os.Getenv("JWT_SECRET"))
type UserReq struct{
	 Username string `json:"username"`
        Password string `json:"password"`
	}


func SignUp() gin.HandlerFunc{
	return func( c *gin.Context){
		var userreq UserReq
		var user models.User
		err :=c.BindJSON(&userreq)
		if err !=nil{
			c.JSON(http.StatusBadRequest,gin.H{
				"error":err.Error(),
			})
			return 
		}
		err =database.DB.Where("username = ?",userreq.Username).First(&models.User{}).Error

		if err ==nil{
			c.JSON(http.StatusConflict,gin.H{
				"message":"user already exist",
			})
			return 
		}else if errors.Is(err, gorm.ErrRecordNotFound){
			user.Username=userreq.Username
			user.Password_hash,err=hashPass(userreq.Password)
			if err !=nil{
				c.JSON(http.StatusInternalServerError,gin.H{
					"error":"password hasing failed",
				})
				return 
			}
			err := database.DB.Create(&user).Error
			if err !=nil{
				c.JSON(http.StatusInternalServerError,gin.H{
					"error":"error creating account",
				})
				return 
			}
			c.JSON(http.StatusOK,gin.H{
				"message ":"user created",
			})
		}else{
				c.JSON(http.StatusInternalServerError,gin.H{
					"error":"error in db",
				})
				return 
		}
	}
}

func hashPass(pass string)(string,error){
	hashedpass, err :=bcrypt.GenerateFromPassword([]byte(pass),bcrypt.DefaultCost)
	if err !=nil{
		return "",err
	}
	return string(hashedpass),nil
}

func verifypass(hasedPass string,userPass string)error{
	err :=bcrypt.CompareHashAndPassword([]byte(hasedPass),[]byte(userPass))
	return  err
	
}

func generateToken(userId int,userName string)(string,error){
	token:=jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
		"userId":userId,
		"username":userName,
		"expire":time.Now().Add(time.Hour*24).Unix(),
	})

	tokenString,err :=token.SignedString(secret)

	return tokenString,err

}

func Login() gin.HandlerFunc{
	return func(c *gin.Context) {
		var user UserReq
		var retrived models.User

		err :=c.BindJSON(&user)
		if err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{
				"error":err.Error(),
			})
			return
		}
		fmt.Println(user)
		result :=database.DB.Where("username=?",user.Username).First(&retrived)

		if result.Error !=nil{
			c.JSON(http.StatusUnauthorized,gin.H{
				"error":"user not found",
				"message":result.Error,
			})
			return 
		}
		err =verifypass(retrived.Password_hash,user.Password)
		if err !=nil{
			c.JSON(http.StatusUnauthorized,gin.H{
				"error":"user/password wrong",
			})
			return 
		}
		token,err :=generateToken(retrived.Id,retrived.Username)
		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{
				"error":"error geratig token",
			})
			return 
		}
		c.JSON(http.StatusOK,gin.H{
			"token":token,
		})
	}
}