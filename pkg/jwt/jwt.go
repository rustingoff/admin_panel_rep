package jwt

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"github.com/rustingoff/admin_panel_rep/internal/model"
	"github.com/rustingoff/admin_panel_rep/pkg/redis"
	"github.com/twinj/uuid"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func init() {
	var _ = godotenv.Load("./token.env")
}

func CreateToken(userid uint) (*model.TokenDetails, error) {
	td := &model.TokenDetails{}
	td.AtExpires = time.Now().Add(time.Hour * 8).Unix()
	td.AccessUuid = uuid.NewV4().String()

	var err error

	acs := os.Getenv("ACCESS_SECRET")

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = userid
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(acs))
	if err != nil {
		return nil, err
	}

	return td, nil
}

func CreateAuth(userid uint, td *model.TokenDetails) error {
	at := time.Unix(td.AtExpires, 0)
	now := time.Now()

	errAccess := redis.Client.Set(context.Background(), td.AccessUuid, strconv.Itoa(int(userid)), at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	return nil
}

func ExtractToken(r *http.Request) string {
	bearToken, err := r.Cookie("Authorization")

	if err != nil {
		log.Printf("extract token error: %v", err)
		return ""
	}

	return bearToken.Value
}

func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		acs := os.Getenv("ACCESS_SECRET")
		return []byte(acs), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func TokenValid(r *http.Request) error {
	token, err := VerifyToken(r)
	if err != nil {
		return err
	}

	if !token.Valid {
		return err
	}
	return nil
}

func ExtractTokenMetadata(r *http.Request) (*model.AccessDetails, error) {
	token, err := VerifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}
		return &model.AccessDetails{
			AccessUuid: accessUuid,
			UserId:     userId,
		}, nil
	}

	return nil, err
}

func FetchAuth(authD *model.AccessDetails) (uint64, error) {
	userid, err := redis.Client.Get(context.Background(), authD.AccessUuid).Result()
	if err != nil {
		return 0, err
	}
	userID, _ := strconv.ParseUint(userid, 10, 64)
	return userID, nil
}

func DeleteAuth(givenUuid string) (int64, error) {
	deleted, err := redis.Client.Del(context.Background(), givenUuid).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}
