package main

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	dump(r)

	w.Header().Set("Content-Type", "application/json")

	claims, err := parseToken(r)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("cognito:groups: %v\n", claims["cognito:groups"])

	groups := parseGroups(claims)

	if len(groups) == 0 {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	token := newToken(groups)
	write(w, fmt.Sprintf("{\"authorizationToken\":\"%s\"}", token))
}

func parseGroups(claims map[string]any) []string {
	groupsRaw, ok := claims["cognito:groups"]
	if !ok {
		return []string{}
	}

	groupsSlice, ok := groupsRaw.([]interface{})
	if !ok {
		panic("cognito:groups claim is not a list")
	}

	var groups []string
	for _, val := range groupsSlice {
		switch v := val.(type) {
		case string:
			groups = append(groups, v)
		default:
			// fallback: stringify non-string values
			groups = append(groups, fmt.Sprintf("%v", v))
		}
	}
	return groups
}

// DEBUG ONLY - we are _not_ verifying the signature!
func parseToken(r *http.Request) (map[string]any, error) {
	emptyClaims := map[string]any{}
	h := r.Header.Get("authorization")
	if !strings.HasPrefix(h, "Bearer") {
		return emptyClaims, fmt.Errorf("no bearer token found")
	}
	tokenString := strings.TrimPrefix(h, "Bearer ")

	// Parse the token without verifying the signature
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return emptyClaims, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		debug(fmt.Sprintf("id token: %v", asJson(claims)))
		return claims, nil
	} else {
		return emptyClaims, fmt.Errorf("failed to extract claims")
	}
}

func dump(r *http.Request) {
	fmt.Println("--------------------------------------------------------------------------------")
	dr, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Printf("Dump error: %v", err)
	} else {
		fmt.Println(string(dr))
	}
	fmt.Println("--------------------------------------------------------------------------------")
	/*
		debug(fmt.Sprintf("%s %s (%v) %s\n", r.Method, r.URL, r.URL.Query(), r.Proto))
		fmt.Printf("Headers:\n")
		for k, v := range r.Header {
			fmt.Printf("%v: %v\n", k, v)
		}
	*/
}

func newToken(groups []string) string {
	secret := []byte("my-secret-key")

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":    "1234567890",
		"name":   "Pontus Example",
		"groups": groups,
		"iat":    time.Now().Unix(),
		"exp":    time.Now().Add(time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		fail("could not sign token")
		log.Fatal(err)
	}

	debug(fmt.Sprintf("response token: %s", tokenString))
	return tokenString
}

func main() {
	debug("users server listening at :8082")
	http.HandleFunc("/", defaultHandler)
	log.Fatal(http.ListenAndServe(":8082", nil))
}

func write(w http.ResponseWriter, s string) {
	_, err := w.Write([]byte(s))
	if err != nil {
		log.Fatal(err)
	}
}

func debug(s string) {
	fmt.Printf("%s\t[debug] %s\n", time.Now().Format("2006-01-02T15:04:05Z07:00"), s)
}

func fail(s string) {
	fmt.Printf("%s\t[error] %s\n", time.Now().Format("2006-01-02T15:04:05Z07:00"), s)
}

func asJson(m map[string]interface{}) string {
	jsonBytes, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}
	return string(jsonBytes)
}
