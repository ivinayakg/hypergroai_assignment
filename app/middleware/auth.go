package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ivinayakg/hypergroai_assignment/helpers"
	"github.com/ivinayakg/hypergroai_assignment/models"
	"github.com/ivinayakg/hypergroai_assignment/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

type userAuth string
type userAuthCache string

const UserAuthKey userAuth = "User"
const UserAuthCacheKey userAuthCache = "user:email:%v"

func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(tokenHeader) < 2 {
			errMsg := "Authentication error!, Provide valid auth token"
			helpers.SendJSONError(&w, http.StatusForbidden, errMsg)
			log.Println(errMsg)
			return
		}
		token := tokenHeader[1]

		verifyUserData, err := utils.VerifyJwt(token)
		if err != nil {
			errMsg := err.Error()
			helpers.SendJSONError(&w, http.StatusForbidden, errMsg)
			log.Println(errMsg)
			return
		}

		user := &models.User{}

		userAuthCacheKeyValue := fmt.Sprintf(string(UserAuthCacheKey), (*verifyUserData)["email"])
		err = helpers.Redis.GetJSON(userAuthCacheKeyValue, user)
		if err != nil {
			user, err = models.GetUser((*verifyUserData)["email"])
			if err != nil {
				errMsg := err.Error()
				if err != mongo.ErrNoDocuments {
					errMsg = "Authentication error!"
				}
				helpers.SendJSONError(&w, http.StatusForbidden, errMsg)
				log.Println(errMsg)
				return
			}
			helpers.Redis.SetJSON(userAuthCacheKeyValue, user, time.Until(time.Now().Add(time.Second*5000)))
		}

		user.Token = token

		c := context.WithValue(r.Context(), UserAuthKey, user)
		next.ServeHTTP(w, r.WithContext(c))
	})
}
