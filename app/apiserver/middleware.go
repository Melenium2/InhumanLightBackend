package apiserver

import (
	"context"
	"net/http"
	"strings"

	"github.com/felixge/httpsnoop"
	"github.com/inhumanLightBackend/app/apiserver/apierrors"
	"github.com/inhumanLightBackend/app/utils/jwtHelper"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

type (
	ctxKey      int8
	HttpReqInfo struct {
		method    string
		uri       string
		refer     string
		ipaddr    string
		code      int
		size      int64
		duration  float64
		userAgent string
	}
)

const (
	// Context key for authorize user in the system
	ctxUserKey ctxKey = iota
)

var (
	// Init new logrus instance
	logger = logrus.New()
)

// Init logrus date dormat
func init() {
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp: true,
	})
}

// Authenticate user in the system
func authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := getToken(r)
		if err != nil {
			sendError(w, r, http.StatusUnauthorized, err)
			return
		}

		claims, err := jwtHelper.Validate(token)
		if err != nil || claims.Type == "refresh" {
			sendError(w, r, http.StatusUnauthorized, apierrors.ErrNotAuthenticated)
			return
		}

		ctx := context.WithValue(r.Context(), ctxUserKey, map[string]interface{}{
			"id":     claims.UserId,
			"access": claims.Access,
		})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Loggin request recived by API
func logging(next http.Handler) http.Handler {

	ipAddrFromRemoteAddr := func(s string) string {
		idx := strings.LastIndex(s, ":")
		if idx == -1 {
			return s
		}

		return s[:idx]
	}

	// getRemoteAddr returns ip address of the client making the request,
	// taking into account http proxies
	getRemoteAddr := func(r *http.Request) string {
		header := r.Header
		hRealIp := header.Get("X-Real-IP")
		hForwardedFor := header.Get("X-Forwarded-For")
		if hRealIp == "" && hForwardedFor == "" {
			return ipAddrFromRemoteAddr(r.RemoteAddr)
		} 
		if hForwardedFor != "" {
			parts := strings.Split(hForwardedFor, ",")
			for i, p := range parts {
				parts[i] = strings.TrimSpace(p)
			}

			return parts[0]
		}

		return hRealIp
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metrics := httpsnoop.CaptureMetrics(next, w, r)

		reqInfo := &HttpReqInfo{
			method: r.Method,
			uri: r.URL.String(),
			refer: r.Header.Get("Refer"),
			userAgent: r.Header.Get("User-Agent"),
			ipaddr: getRemoteAddr(r),
			code: metrics.Code,
			size: metrics.Written,
			duration: metrics.Duration.Seconds(),
		}
		
		logger.WithFields(logrus.Fields{
			"IP": reqInfo.ipaddr,
			"Method": reqInfo.method,
			"URI": reqInfo.uri,
			"Code": reqInfo.code,
			"Duration": reqInfo.duration,
			"Refer": reqInfo.refer,
			"User-agent": reqInfo.userAgent,
			"Length": reqInfo.size, 
		}).Info("Request")
	})
}

// Get auth token from header
func getToken(r *http.Request) (string, error) {
	header := r.Header.Get("Authentication")
	if header == "" {
		return "", apierrors.ErrNotAuthenticated
	}

	splitedToken := strings.Split(header, " ")
	if len(splitedToken) != 2 {
		return "", apierrors.ErrNotAuthenticated
	}

	token := splitedToken[1]
	if token == "" {
		return "", apierrors.ErrNotAuthenticated
	}

	return token, nil
}

// Get user map in context of request
func userContextMap(ctx interface{}) map[string]string {
	return cast.ToStringMapString(ctx)
}
