package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AlehaWP/YaPracticum.git/internal/global"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type UrlsMock struct {
	mock.Mock
}

func (m *UrlsMock) SaveURL(url []byte, baseURL, userID string) string {
	args := m.Called(url, baseURL, userID)
	return args.String(0)
}

func (m *UrlsMock) GetURL(id string) (string, error) {
	args := m.Called(id)
	return args.String(0), args.Error(1)
}

// func (m *UrlsMock) Get() map[string][]string {
// 	return nil
// }

// func (m *UrlsMock) ToSet() *map[string][]string {
// 	return nil
// }

func (m *UrlsMock) FindUser(string) bool {
	return false
}

func (m *UrlsMock) CreateUser() (string, error) {
	return "", nil
}

func (m *UrlsMock) GetUserURLs(string) []global.URLs {
	return nil
}

/*type OptsMock struct {
	mock.Mock
}

func (o *OptsMock) ServAddr() string{
	args := o.Called()
    return args.String()
}

func (o *OptsMock)  RespBaseURL() string{
	args := o.Called()
    return args.String()
}

func (o *OptsMock) RepoFileName() string{
	args := o.Called()
    return args.String()
}
*/

func TestHandlerUrlGet(t *testing.T) {
	dataTests := map[string]map[string]interface{}{
		"test1": {
			"reqID":       "123123asdasd",
			"result":      "www.example.com",
			"resStatus":   307,
			"mockReturn1": "www.example.com",
		},
		"test2": {
			"reqID":       "123123",
			"result":      "",
			"resStatus":   400,
			"mockReturn1": "",
			"mockReturn2": errors.New("not found"),
		},
	}

	repoMock := new(UrlsMock)

	NewHandlers(repoMock, "http://someurl")
	handler := http.HandlerFunc(HandlerURLGet)

	for key, value := range dataTests {
		log.Println("start test", key)
		var err error
		if value["mockReturn2"] != nil {
			err = value["mockReturn2"].(error)
		}
		repoMock.On("GetURL", value["reqID"].(string)).Return(value["mockReturn1"].(string), err)
		r := httptest.NewRequest("GET", "/"+value["reqID"].(string), strings.NewReader(""))
		w := httptest.NewRecorder()
		ctx := context.WithValue(context.Background(), global.CtxString("url_id"), value["reqID"].(string))
		handler.ServeHTTP(w, r.WithContext(ctx))

		res := w.Result()
		assert.Equal(t, value["resStatus"].(int), res.StatusCode, "Не верный код ответа GET")
		assert.Equal(t, w.Header().Get("Location"), value["result"].(string), "Не верный ответ GET")
		defer res.Body.Close()
	}
}

func TestHandlerUrlPost(t *testing.T) {
	repoMock := new(UrlsMock)
	repoMock.On("SaveURL", []byte("www.example.com"), "http://baseURL/", "asdasd").Return("http://baseURL/123123asdasd")

	NewHandlers(repoMock, "http://baseURL")
	handler := http.HandlerFunc(HandlerURLPost)
	r := httptest.NewRequest("POST", "http://localhost:8080", strings.NewReader("www.example.com"))
	w := httptest.NewRecorder()

	ctx := context.WithValue(context.Background(), global.CtxString("UserID"), "asdasd")
	handler.ServeHTTP(w, r.WithContext(ctx))

	res := w.Result()
	b, _ := io.ReadAll(res.Body)
	defer res.Body.Close()
	assert.Equal(t, 201, res.StatusCode, "Не верный код ответа POST")
	assert.Equal(t, "http://baseURL/123123asdasd", string(b), "Не верный ответ POST")

}

func TestHandlerApiUrlPost(t *testing.T) {
	str := &struct {
		URL string
	}{
		URL: "www.example.com",
	}
	bOut, err := json.Marshal(str)
	if err != nil {
		t.Error("Ошибка серилизации")
	}
	repoMock := new(UrlsMock)
	repoMock.On("SaveURL", []byte("www.example.com"), "http://baseURL/", "aasdasdSQW").Return("http://baseURL/123123asdasd")
	NewHandlers(repoMock, "http://baseURL")
	handler := http.HandlerFunc(HandlerAPIURLPost)
	r := httptest.NewRequest("POST", "http://localhost:8080", bytes.NewBuffer(bOut))
	w := httptest.NewRecorder()

	ctx := context.WithValue(context.Background(), global.CtxString("UserID"), "aasdasdSQW")
	handler.ServeHTTP(w, r.WithContext(ctx))
	res := w.Result()
	b, _ := io.ReadAll(res.Body)
	defer res.Body.Close()
	assert.Equal(t, 201, res.StatusCode, "Не верный код ответа POST")
	assert.Equal(t, `{"result":"http://baseURL/123123asdasd"}`, string(b), "Не верный ответ POST")

}
