package cloudscaffolder

import (
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
)

func (i *Impl) GetAllVm(w rest.ResponseWriter, r *rest.Request) {
	//token := r.Header.Get("Authorization")
	//log.Println(token)
	vm := GetAllVm(i)
	w.WriteJson(&vm)
}

func CheckJwt() (jwt.MapClaims, error) {
	hmacSecret := []byte("secret key")
	tokenString := "aaaa"
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSecret, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func Serve(i *Impl) {
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)

	api.Use(&rest.CorsMiddleware{
		RejectNonCorsRequests: false,
		OriginValidator: func(origin string, request *rest.Request) bool {
			return origin == "http://104.154.29.8:8080"
		},
		AllowedMethods:                []string{"GET", "POST", "PUT"},
		AllowedHeaders:                []string{"Content-Type", "Authorization"},
		AccessControlAllowCredentials: true,
		AccessControlMaxAge:           3600,
	})

	router, err := rest.MakeRouter(
		rest.Get("/vms", i.GetAllVm),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":4000", api.MakeHandler()))
}
