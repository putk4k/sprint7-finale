package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func TestCafeWhenNotOk(t *testing.T) {
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
		responce := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)

		handler.ServeHTTP(responce, req)
		assert.Equal(t, v.status, responce.Code)
		assert.Equal(t, v.message, strings.TrimSpace(responce.Body.String()))
	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		count int
		want  int
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{100, 3},
	}

	for _, v := range requests {
		responce := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/cafe?city=tula&count="+strconv.Itoa(v.count), nil)

		handler.ServeHTTP(responce, req)
		require.Equal(t, http.StatusOK, responce.Code)

		cafes := strings.FieldsFunc(responce.Body.String(), func(r rune) bool {
			return r == ','
		})
		assert.Len(t, cafes, v.want)

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
	}

	for _, v := range requests {
		responce := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/cafe?city=moscow&search="+v.search, nil)

		handler.ServeHTTP(responce, req)
		require.Equal(t, http.StatusOK, responce.Code)
		cafes := strings.FieldsFunc(responce.Body.String(), func(r rune) bool {
			return r == ','
		})
		assert.Len(t, cafes, v.wantCount)
		for _, cafe := range cafes {
			assert.Contains(t, strings.ToLower(cafe), strings.ToLower(v.search))
		}
	}
}
