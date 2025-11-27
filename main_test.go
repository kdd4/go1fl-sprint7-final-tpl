package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		count   int
	}{
		{"/cafe?city=moscow&count=0", 0},
		{"/cafe?city=tula&count=1", 1},
		{"/cafe?city=moscow&count=2", 2},
		{"/cafe?city=moscow&count=100", min(100, len(cafeList["moscow"]))},
	}

	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)

		responseBody := response.Body.String()

		cafeCount := 0

		if responseBody != "" {
			splitedCafe := strings.Split(responseBody, ",")
			cafeCount = len(splitedCafe)
		}

		assert.Equal(t, v.count, cafeCount)
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		search    string
		wantCount int
	}{
		{"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
		{"СытыЙ", 1},
		{"КОФЕ", 2},
	}

	for _, v := range requests {
		response := httptest.NewRecorder()
		target := fmt.Sprintf("/cafe?city=moscow&search=%s", v.search)
		req := httptest.NewRequest("GET", target, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)

		responseBody := response.Body.String()

		splitedCafe := []string{}
		cafeCount := 0

		if responseBody != "" {
			splitedCafe = strings.Split(responseBody, ",")
			cafeCount = len(splitedCafe)
		}

		assert.Equal(t, v.wantCount, cafeCount, v.search)

		for _, cafe := range splitedCafe {
			isCurrect := strings.Contains(strings.ToLower(cafe), strings.ToLower(v.search))
			assert.True(t, isCurrect)
		}
	}
}
