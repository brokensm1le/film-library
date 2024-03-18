package http

import (
	"bytes"
	"context"
	"encoding/json"
	"film_library/internal/auth"
	mock_auth "film_library/internal/auth/mocks"
	"film_library/internal/service"
	mock_service "film_library/internal/service/mocks"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValidateActor(t *testing.T) {
	cases := []struct {
		name   string
		in     *service.Actor
		expErr error
	}{
		{
			name:   "NoName",
			in:     &service.Actor{Sex: "m", BDate: "1999-10-10"},
			expErr: fmt.Errorf("size Name should be [1;100]"),
		},
		{
			name:   "NoSex",
			in:     &service.Actor{Name: "Sasha", BDate: "1999-10-10"},
			expErr: fmt.Errorf("sex should be 'm' - male or 'f' - famale"),
		},
		{
			name:   "NoBDate",
			in:     &service.Actor{Name: "Sasha", Sex: "m"},
			expErr: fmt.Errorf("bdate should be '2000-01-01' format"),
		},
		{
			name:   "UncorrectSex",
			in:     &service.Actor{Name: "Sasha", Sex: "d", BDate: "1999-10-10"},
			expErr: fmt.Errorf("sex should be 'm' - male or 'f' - famale"),
		},
		{
			name:   "UncorrectBDate",
			in:     &service.Actor{Name: "Sasha", Sex: "f", BDate: "10-10-1999"},
			expErr: fmt.Errorf("bdate should be '2000-01-01' format"),
		},
	}
	h := ServiceHandler{}
	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			err := h.validateActor(tCase.in)
			require.Error(t, err)
			require.EqualError(t, tCase.expErr, err.Error())
		})
	}
	err := h.validateActor(&service.Actor{Name: "Sasha", Sex: "m", BDate: "1999-10-10"})
	require.NoError(t, err)
}

func TestValidateFilm(t *testing.T) {
	textError := `errorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorer
	rorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerr
	orerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerro
	rerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerror
	errorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrore
	rrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorer
	rorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerr
	orerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerro
	rerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerro
	rerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerrorerror`
	cases := []struct {
		name   string
		in     *service.Film
		expErr error
	}{
		{
			name:   "NoName",
			in:     &service.Film{RDate: "1999-10-10", Rating: 9.0, Desc: "text"},
			expErr: fmt.Errorf("size Name should be [1;150]"),
		},
		{
			name:   "NoRdate",
			in:     &service.Film{Name: "Lilo&Stich", Rating: 9.0, Desc: "text"},
			expErr: fmt.Errorf("rdate should be '2000-01-01' format"),
		},
		{
			name:   "NoRating",
			in:     &service.Film{Name: "Lilo&Stich", RDate: "1999-10-10", Desc: "text"},
			expErr: fmt.Errorf("rating should be (0;10]"),
		},
		{
			name:   "0Rating",
			in:     &service.Film{Name: "Lilo&Stich", Rating: 0, RDate: "1999-10-10", Desc: "text"},
			expErr: fmt.Errorf("rating should be (0;10]"),
		},
		{
			name:   "UncorrectBDate",
			in:     &service.Film{Name: "Lilo&Stich", Rating: 9.0, RDate: "10-10-1999", Desc: "text"},
			expErr: fmt.Errorf("rdate should be '2000-01-01' format"),
		},
		{
			name:   "LimitDesc",
			in:     &service.Film{Name: "Lilo&Stich", Rating: 9.0, RDate: "1999-10-10", Desc: textError},
			expErr: fmt.Errorf("size Name should be < 1000 symbols"),
		},
	}
	h := ServiceHandler{}
	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			err := h.validateFilm(tCase.in)
			require.Error(t, err)
			require.EqualError(t, tCase.expErr, err.Error())
		})
	}
	err := h.validateFilm(&service.Film{Name: "Lilo&Stich", Rating: 7.7, RDate: "1999-10-10", Desc: "text"})
	require.NoError(t, err)
}

func TestCreateActor(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUsecase, user service.Actor)
	var resp *auth.ResponseModel = &auth.ResponseModel{Status: "OK"}
	ans, _ := json.Marshal(resp)

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           service.Actor
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody []byte
	}{
		{
			name:      "OK",
			inputBody: `{"Name":"Sasha", "Sex":"m", "BDate":"1999-10-10"}`,
			inputUser: service.Actor{
				Name:  "Sasha",
				Sex:   "m",
				BDate: "1999-10-10",
			},
			mockBehavior: func(s *mock_service.MockUsecase, actor service.Actor) {
				s.EXPECT().CreateActor(&actor).Return(nil).Times(1)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: ans,
		},
		{
			name:      "Error",
			inputBody: `{"Name":"Sasha", "Sex":"m", "BDate":"1999-10-10"}`,
			inputUser: service.Actor{
				Name:  "Sasha",
				Sex:   "m",
				BDate: "1999-10-10",
			},
			mockBehavior: func(s *mock_service.MockUsecase, actor service.Actor) {
				s.EXPECT().CreateActor(&actor).Return(fmt.Errorf("error")).Times(1)
			},
			expectedStatusCode:  http.StatusInternalServerError,
			expectedRequestBody: []byte("error\n"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockService := mock_service.NewMockUsecase(c)
			mockAuth := mock_auth.NewMockUsecase(c)
			testCase.mockBehavior(mockService, testCase.inputUser)

			handler := NewServiceHandler(mockService, mockAuth)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/actor/add", bytes.NewBufferString(testCase.inputBody))
			ctx := context.WithValue(r.Context(), "tokenData", &auth.TokenData{Id: 1, Role: 1})
			r = r.WithContext(ctx)
			handler.CreateActor(w, r)

			require.Equal(t, testCase.expectedStatusCode, w.Code)
			require.Equal(t, testCase.expectedRequestBody, w.Body.Bytes())
		})
	}
}

func TestGetActor(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUsecase, name string, actor service.Actor)
	var resp *service.Actor = &service.Actor{
		Name:  "Sasha",
		Sex:   "m",
		BDate: "1999-10-10",
	}
	ans, _ := json.Marshal(resp)

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           service.Actor
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody []byte
	}{
		{
			name:      "OK",
			inputBody: `{"Name":"Sasha", "Sex":"m", "BDate":"1999-10-10"}`,
			inputUser: service.Actor{
				Name:  "Sasha",
				Sex:   "m",
				BDate: "1999-10-10",
			},
			mockBehavior: func(s *mock_service.MockUsecase, name string, actor service.Actor) {
				s.EXPECT().GetActor(name).Return(&actor, nil).Times(1)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: ans,
		},
		{
			name:      "Error",
			inputBody: ``,
			inputUser: service.Actor{},
			mockBehavior: func(s *mock_service.MockUsecase, name string, actor service.Actor) {
				s.EXPECT().GetActor(name).Return(nil, fmt.Errorf("error")).Times(1)
			},
			expectedStatusCode:  http.StatusInternalServerError,
			expectedRequestBody: []byte("error\n"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockService := mock_service.NewMockUsecase(c)
			mockAuth := mock_auth.NewMockUsecase(c)
			testCase.mockBehavior(mockService, "Sasha", testCase.inputUser)

			handler := NewServiceHandler(mockService, mockAuth)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/actor/get/Sasha", bytes.NewBufferString(""))
			ctx := context.WithValue(r.Context(), "tokenData", &auth.TokenData{Id: 1, Role: 1})
			r = r.WithContext(ctx)
			handler.GetActor(w, r)

			require.Equal(t, testCase.expectedStatusCode, w.Code)
			require.Equal(t, testCase.expectedRequestBody, w.Body.Bytes())
		})
	}
}

func TestGetActors(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUsecase, details service.DetailsParams, actor []service.Actor)
	var resp []service.Actor = []service.Actor{{
		Name:  "Sasha",
		Sex:   "m",
		BDate: "1999-10-10"},
	}
	ans, _ := json.Marshal(resp)
	var details service.DetailsParams = service.DetailsParams{Sort: "actor_name"}

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           service.Actor
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody []byte
	}{
		{
			name:      "OK",
			inputBody: `{"Name":"Sasha", "Sex":"m", "BDate":"1999-10-10"}`,
			inputUser: service.Actor{
				Name:  "Sasha",
				Sex:   "m",
				BDate: "1999-10-10",
			},
			mockBehavior: func(s *mock_service.MockUsecase, details service.DetailsParams, actors []service.Actor) {
				s.EXPECT().GetActors(&details).Return(actors, nil).Times(1)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: ans,
		},
		{
			name:      "Error",
			inputBody: ``,
			inputUser: service.Actor{},
			mockBehavior: func(s *mock_service.MockUsecase, details service.DetailsParams, actors []service.Actor) {
				s.EXPECT().GetActors(&details).Return(actors, fmt.Errorf("error")).Times(1)
			},
			expectedStatusCode:  http.StatusInternalServerError,
			expectedRequestBody: []byte("error\n"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockService := mock_service.NewMockUsecase(c)
			mockAuth := mock_auth.NewMockUsecase(c)
			testCase.mockBehavior(mockService, details, []service.Actor{testCase.inputUser})

			handler := NewServiceHandler(mockService, mockAuth)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/actor/get_all", bytes.NewBufferString(""))
			ctx := context.WithValue(r.Context(), "tokenData", &auth.TokenData{Id: 1, Role: 1})
			r = r.WithContext(ctx)
			handler.GetActors(w, r)

			require.Equal(t, testCase.expectedStatusCode, w.Code)
			require.Equal(t, testCase.expectedRequestBody, w.Body.Bytes())
		})
	}
}

func TestUpdateActor(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUsecase, name string, user service.Actor)
	var resp *auth.ResponseModel = &auth.ResponseModel{Status: "OK"}
	ans, _ := json.Marshal(resp)

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           service.Actor
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody []byte
	}{
		{
			name:      "OK",
			inputBody: `{"name":"Sasha", "sex":"m", "bdate":"1999-10-10"}`,
			inputUser: service.Actor{
				Name:  "Sasha",
				Sex:   "m",
				BDate: "1999-10-10",
			},
			mockBehavior: func(s *mock_service.MockUsecase, name string, actor service.Actor) {
				s.EXPECT().UpdateActor(name, &actor).Return(nil).Times(1)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: ans,
		},
		{
			name:      "Error",
			inputBody: `{"Name":"Sasha", "Sex":"m", "BDate":"1999-10-10"}`,
			inputUser: service.Actor{
				Name:  "Sasha",
				Sex:   "m",
				BDate: "1999-10-10",
			},
			mockBehavior: func(s *mock_service.MockUsecase, name string, actor service.Actor) {
				s.EXPECT().UpdateActor(name, &actor).Return(fmt.Errorf("error")).Times(1)
			},
			expectedStatusCode:  http.StatusInternalServerError,
			expectedRequestBody: []byte("error\n"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockService := mock_service.NewMockUsecase(c)
			mockAuth := mock_auth.NewMockUsecase(c)
			testCase.mockBehavior(mockService, "Sasha", testCase.inputUser)

			handler := NewServiceHandler(mockService, mockAuth)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/actor/update/Sasha", bytes.NewBufferString(testCase.inputBody))
			ctx := context.WithValue(r.Context(), "tokenData", &auth.TokenData{Id: 1, Role: 1})
			r = r.WithContext(ctx)
			handler.UpdateActor(w, r)

			require.Equal(t, testCase.expectedStatusCode, w.Code)
			require.Equal(t, testCase.expectedRequestBody, w.Body.Bytes())
		})
	}
}

func TestDeleteActor(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUsecase, name string)
	var resp *auth.ResponseModel = &auth.ResponseModel{Status: "OK"}
	ans, _ := json.Marshal(resp)

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           service.Actor
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody []byte
	}{
		{
			name:      "OK",
			inputBody: `{"Name":"Sasha", "Sex":"m", "BDate":"1999-10-10"}`,
			inputUser: service.Actor{
				Name:  "Sasha",
				Sex:   "m",
				BDate: "1999-10-10",
			},
			mockBehavior: func(s *mock_service.MockUsecase, name string) {
				s.EXPECT().DeleteActor(name).Return(nil).Times(1)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: ans,
		},
		{
			name:      "Error",
			inputBody: ``,
			inputUser: service.Actor{},
			mockBehavior: func(s *mock_service.MockUsecase, name string) {
				s.EXPECT().DeleteActor(name).Return(fmt.Errorf("error")).Times(1)
			},
			expectedStatusCode:  http.StatusInternalServerError,
			expectedRequestBody: []byte("error\n"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockService := mock_service.NewMockUsecase(c)
			mockAuth := mock_auth.NewMockUsecase(c)
			testCase.mockBehavior(mockService, "Sasha")

			handler := NewServiceHandler(mockService, mockAuth)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/actor/delete/Sasha", bytes.NewBufferString(""))
			ctx := context.WithValue(r.Context(), "tokenData", &auth.TokenData{Id: 1, Role: 1})
			r = r.WithContext(ctx)
			handler.DeleteActor(w, r)

			require.Equal(t, testCase.expectedStatusCode, w.Code)
			require.Equal(t, testCase.expectedRequestBody, w.Body.Bytes())
		})
	}
}

func TestSearchActor(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUsecase, pattern string, actor []string)
	var resp []string = []string{"Sasha"}
	ans, _ := json.Marshal(resp)

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           service.Actor
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody []byte
	}{
		{
			name:      "OK",
			inputBody: `{"Name":"Sasha", "Sex":"m", "BDate":"1999-10-10"}`,
			inputUser: service.Actor{
				Name:  "Sasha",
				Sex:   "m",
				BDate: "1999-10-10",
			},
			mockBehavior: func(s *mock_service.MockUsecase, pattern string, actors []string) {
				s.EXPECT().SearchActor(pattern).Return(actors, nil).Times(1)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: ans,
		},
		{
			name:      "Error",
			inputBody: ``,
			inputUser: service.Actor{},
			mockBehavior: func(s *mock_service.MockUsecase, pattern string, actors []string) {
				s.EXPECT().SearchActor(pattern).Return(actors, fmt.Errorf("error")).Times(1)
			},
			expectedStatusCode:  http.StatusInternalServerError,
			expectedRequestBody: []byte("error\n"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockService := mock_service.NewMockUsecase(c)
			mockAuth := mock_auth.NewMockUsecase(c)
			testCase.mockBehavior(mockService, "asha", []string{testCase.inputUser.Name})

			handler := NewServiceHandler(mockService, mockAuth)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/actor/search/asha", bytes.NewBufferString(""))
			ctx := context.WithValue(r.Context(), "tokenData", &auth.TokenData{Id: 1, Role: 1})
			r = r.WithContext(ctx)
			handler.SearchActor(w, r)

			require.Equal(t, testCase.expectedStatusCode, w.Code)
			require.Equal(t, testCase.expectedRequestBody, w.Body.Bytes())
		})
	}
}

func TestCreateFilm(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUsecase, film service.Film)
	var resp *auth.ResponseModel = &auth.ResponseModel{Status: "OK"}
	ans, _ := json.Marshal(resp)

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           service.Film
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody []byte
	}{
		{
			name:      "OK",
			inputBody: `{"name":"Sasha", "rating":5.5, "RDate":"1999-10-10", "Desc":"nice film"}`,
			inputUser: service.Film{
				Name:   "Sasha",
				Rating: 5.5,
				RDate:  "1999-10-10",
				Desc:   "nice film",
			},
			mockBehavior: func(s *mock_service.MockUsecase, film service.Film) {
				s.EXPECT().CreateFilm(&film).Return(nil).Times(1)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: ans,
		},
		{
			name:      "Error",
			inputBody: `{"Name":"Sasha", "Rating":5.5, "RDate":"1999-10-10", "Desc":"nice film"}`,
			inputUser: service.Film{
				Name:   "Sasha",
				Rating: 5.5,
				RDate:  "1999-10-10",
				Desc:   "nice film",
			},
			mockBehavior: func(s *mock_service.MockUsecase, film service.Film) {
				s.EXPECT().CreateFilm(&film).Return(fmt.Errorf("error")).Times(1)
			},
			expectedStatusCode:  http.StatusInternalServerError,
			expectedRequestBody: []byte("error\n"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockService := mock_service.NewMockUsecase(c)
			mockAuth := mock_auth.NewMockUsecase(c)
			testCase.mockBehavior(mockService, testCase.inputUser)

			handler := NewServiceHandler(mockService, mockAuth)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/film/add", bytes.NewBufferString(testCase.inputBody))
			ctx := context.WithValue(r.Context(), "tokenData", &auth.TokenData{Id: 1, Role: 1})
			r = r.WithContext(ctx)
			handler.CreateFilm(w, r)

			require.Equal(t, testCase.expectedStatusCode, w.Code)
			require.Equal(t, testCase.expectedRequestBody, w.Body.Bytes())
		})
	}
}

func TestGetFilm(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUsecase, name string, actor service.Film)
	var resp *service.Film = &service.Film{
		Name:   "Sasha",
		Rating: 5.5,
		RDate:  "1999-10-10",
		Desc:   "nice film",
	}
	ans, _ := json.Marshal(resp)

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           service.Film
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody []byte
	}{
		{
			name:      "OK",
			inputBody: ``,
			inputUser: service.Film{
				Name:   "Sasha",
				Rating: 5.5,
				RDate:  "1999-10-10",
				Desc:   "nice film",
			},
			mockBehavior: func(s *mock_service.MockUsecase, name string, film service.Film) {
				s.EXPECT().GetFilm(name).Return(&film, nil).Times(1)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: ans,
		},
		{
			name:      "Error",
			inputBody: ``,
			inputUser: service.Film{},
			mockBehavior: func(s *mock_service.MockUsecase, name string, actor service.Film) {
				s.EXPECT().GetFilm(name).Return(nil, fmt.Errorf("error")).Times(1)
			},
			expectedStatusCode:  http.StatusInternalServerError,
			expectedRequestBody: []byte("error\n"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockService := mock_service.NewMockUsecase(c)
			mockAuth := mock_auth.NewMockUsecase(c)
			testCase.mockBehavior(mockService, "Sasha", testCase.inputUser)

			handler := NewServiceHandler(mockService, mockAuth)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/film/get/Sasha", bytes.NewBufferString(""))
			ctx := context.WithValue(r.Context(), "tokenData", &auth.TokenData{Id: 1, Role: 1})
			r = r.WithContext(ctx)
			handler.GetFilm(w, r)

			require.Equal(t, testCase.expectedStatusCode, w.Code)
			require.Equal(t, testCase.expectedRequestBody, w.Body.Bytes())
		})
	}
}

func TestGetFilms(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUsecase, details service.DetailsParams, actor []service.Film)
	var resp []service.Film = []service.Film{{
		Name:   "Sasha",
		Rating: 5.5,
		RDate:  "1999-10-10",
		Desc:   "nice film"},
	}
	ans, _ := json.Marshal(resp)
	var details service.DetailsParams = service.DetailsParams{Sort: "rating"}

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           service.Film
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody []byte
	}{
		{
			name:      "OK",
			inputBody: `{"name":"Sasha", "rating":5.5, "RDate":"1999-10-10", "Desc":"nice film"}`,
			inputUser: service.Film{
				Name:   "Sasha",
				Rating: 5.5,
				RDate:  "1999-10-10",
				Desc:   "nice film",
			},
			mockBehavior: func(s *mock_service.MockUsecase, details service.DetailsParams, actors []service.Film) {
				s.EXPECT().GetFilms(&details).Return(actors, nil).Times(1)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: ans,
		},
		{
			name:      "Error",
			inputBody: ``,
			inputUser: service.Film{},
			mockBehavior: func(s *mock_service.MockUsecase, details service.DetailsParams, actors []service.Film) {
				s.EXPECT().GetFilms(&details).Return(actors, fmt.Errorf("error")).Times(1)
			},
			expectedStatusCode:  http.StatusInternalServerError,
			expectedRequestBody: []byte("error\n"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockService := mock_service.NewMockUsecase(c)
			mockAuth := mock_auth.NewMockUsecase(c)
			testCase.mockBehavior(mockService, details, []service.Film{testCase.inputUser})

			handler := NewServiceHandler(mockService, mockAuth)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/film/get_all", bytes.NewBufferString(""))
			ctx := context.WithValue(r.Context(), "tokenData", &auth.TokenData{Id: 1, Role: 1})
			r = r.WithContext(ctx)
			handler.GetFilms(w, r)

			require.Equal(t, testCase.expectedStatusCode, w.Code)
			require.Equal(t, testCase.expectedRequestBody, w.Body.Bytes())
		})
	}
}

func TestUpdateFilm(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUsecase, name string, user service.Film)
	var resp *auth.ResponseModel = &auth.ResponseModel{Status: "OK"}
	ans, _ := json.Marshal(resp)

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           service.Film
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody []byte
	}{
		{
			name:      "OK",
			inputBody: `{"name":"Sasha", "rating":5.5, "RDate":"1999-10-10", "Desc":"nice film"}`,
			inputUser: service.Film{
				Name:   "Sasha",
				Rating: 5.5,
				RDate:  "1999-10-10",
				Desc:   "nice film",
			},
			mockBehavior: func(s *mock_service.MockUsecase, name string, actor service.Film) {
				s.EXPECT().UpdateFilm(name, &actor).Return(nil).Times(1)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: ans,
		},
		{
			name:      "Error",
			inputBody: `{"name":"Sasha", "rating":5.5, "RDate":"1999-10-10", "Desc":"nice film"}`,
			inputUser: service.Film{
				Name:   "Sasha",
				Rating: 5.5,
				RDate:  "1999-10-10",
				Desc:   "nice film",
			},
			mockBehavior: func(s *mock_service.MockUsecase, name string, actor service.Film) {
				s.EXPECT().UpdateFilm(name, &actor).Return(fmt.Errorf("error")).Times(1)
			},
			expectedStatusCode:  http.StatusInternalServerError,
			expectedRequestBody: []byte("error\n"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockService := mock_service.NewMockUsecase(c)
			mockAuth := mock_auth.NewMockUsecase(c)
			testCase.mockBehavior(mockService, "Sasha", testCase.inputUser)

			handler := NewServiceHandler(mockService, mockAuth)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/api/actor/update/Sasha", bytes.NewBufferString(testCase.inputBody))
			ctx := context.WithValue(r.Context(), "tokenData", &auth.TokenData{Id: 1, Role: 1})
			r = r.WithContext(ctx)
			handler.UpdateFilm(w, r)

			require.Equal(t, testCase.expectedStatusCode, w.Code)
			require.Equal(t, testCase.expectedRequestBody, w.Body.Bytes())
		})
	}
}

func TestDeleteFilm(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUsecase, name string)
	var resp *auth.ResponseModel = &auth.ResponseModel{Status: "OK"}
	ans, _ := json.Marshal(resp)

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           service.Film
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody []byte
	}{
		{
			name:      "OK",
			inputBody: `{"name":"Sasha", "rating":5.5, "RDate":"1999-10-10", "Desc":"nice film"}`,
			inputUser: service.Film{
				Name:   "Sasha",
				Rating: 5.5,
				RDate:  "1999-10-10",
				Desc:   "nice film",
			},
			mockBehavior: func(s *mock_service.MockUsecase, name string) {
				s.EXPECT().DeleteFilm(name).Return(nil).Times(1)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: ans,
		},
		{
			name:      "Error",
			inputBody: ``,
			inputUser: service.Film{},
			mockBehavior: func(s *mock_service.MockUsecase, name string) {
				s.EXPECT().DeleteFilm(name).Return(fmt.Errorf("error")).Times(1)
			},
			expectedStatusCode:  http.StatusInternalServerError,
			expectedRequestBody: []byte("error\n"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockService := mock_service.NewMockUsecase(c)
			mockAuth := mock_auth.NewMockUsecase(c)
			testCase.mockBehavior(mockService, "Sasha")

			handler := NewServiceHandler(mockService, mockAuth)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/film/delete/Sasha", bytes.NewBufferString(""))
			ctx := context.WithValue(r.Context(), "tokenData", &auth.TokenData{Id: 1, Role: 1})
			r = r.WithContext(ctx)
			handler.DeleteFilm(w, r)

			require.Equal(t, testCase.expectedStatusCode, w.Code)
			require.Equal(t, testCase.expectedRequestBody, w.Body.Bytes())
		})
	}
}

func TestSearchFilm(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUsecase, pattern string, film []string)
	var resp []string = []string{"Sasha"}
	ans, _ := json.Marshal(resp)

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           service.Actor
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody []byte
	}{
		{
			name:      "OK",
			inputBody: `{"Name":"Sasha", "Sex":"m", "BDate":"1999-10-10"}`,
			inputUser: service.Actor{
				Name:  "Sasha",
				Sex:   "m",
				BDate: "1999-10-10",
			},
			mockBehavior: func(s *mock_service.MockUsecase, pattern string, actors []string) {
				s.EXPECT().SearchFilms(pattern).Return(actors, nil).Times(1)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: ans,
		},
		{
			name:      "Error",
			inputBody: ``,
			inputUser: service.Actor{},
			mockBehavior: func(s *mock_service.MockUsecase, pattern string, actors []string) {
				s.EXPECT().SearchFilms(pattern).Return(actors, fmt.Errorf("error")).Times(1)
			},
			expectedStatusCode:  http.StatusInternalServerError,
			expectedRequestBody: []byte("error\n"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockService := mock_service.NewMockUsecase(c)
			mockAuth := mock_auth.NewMockUsecase(c)
			testCase.mockBehavior(mockService, "asha", []string{testCase.inputUser.Name})

			handler := NewServiceHandler(mockService, mockAuth)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/film/search/asha", bytes.NewBufferString(""))
			ctx := context.WithValue(r.Context(), "tokenData", &auth.TokenData{Id: 1, Role: 1})
			r = r.WithContext(ctx)
			handler.SearchFilms(w, r)

			require.Equal(t, testCase.expectedStatusCode, w.Code)
			require.Equal(t, testCase.expectedRequestBody, w.Body.Bytes())
		})
	}
}

func TestRole(t *testing.T) {
	t.Run("UpdateErrRole", func(t *testing.T) {
		c := gomock.NewController(t)
		defer c.Finish()

		mockService := mock_service.NewMockUsecase(c)
		mockAuth := mock_auth.NewMockUsecase(c)

		handler := NewServiceHandler(mockService, mockAuth)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/film/search/asha", bytes.NewBufferString(""))
		ctx := context.WithValue(r.Context(), "tokenData", &auth.TokenData{Id: 1, Role: 0})
		r = r.WithContext(ctx)
		handler.UpdateFilm(w, r)

		require.Equal(t, http.StatusForbidden, w.Code)

		w = httptest.NewRecorder()
		r = httptest.NewRequest("GET", "/film/search/asha", bytes.NewBufferString(""))
		ctx = context.WithValue(r.Context(), "tokenData", &auth.TokenData{Id: 1, Role: 0})
		r = r.WithContext(ctx)
		handler.UpdateActor(w, r)

		require.Equal(t, http.StatusForbidden, w.Code)
	})
	t.Run("AddErrRole", func(t *testing.T) {
		c := gomock.NewController(t)
		defer c.Finish()

		mockService := mock_service.NewMockUsecase(c)
		mockAuth := mock_auth.NewMockUsecase(c)

		handler := NewServiceHandler(mockService, mockAuth)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/film/search/asha", bytes.NewBufferString(""))
		ctx := context.WithValue(r.Context(), "tokenData", &auth.TokenData{Id: 1, Role: 0})
		r = r.WithContext(ctx)
		handler.CreateActor(w, r)

		require.Equal(t, http.StatusForbidden, w.Code)

		w = httptest.NewRecorder()
		r = httptest.NewRequest("GET", "/film/search/asha", bytes.NewBufferString(""))
		ctx = context.WithValue(r.Context(), "tokenData", &auth.TokenData{Id: 1, Role: 0})
		r = r.WithContext(ctx)
		handler.CreateFilm(w, r)

		require.Equal(t, http.StatusForbidden, w.Code)
	})
	t.Run("DeleteErrRole", func(t *testing.T) {
		c := gomock.NewController(t)
		defer c.Finish()

		mockService := mock_service.NewMockUsecase(c)
		mockAuth := mock_auth.NewMockUsecase(c)

		handler := NewServiceHandler(mockService, mockAuth)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/film/search/asha", bytes.NewBufferString(""))
		ctx := context.WithValue(r.Context(), "tokenData", &auth.TokenData{Id: 1, Role: 0})
		r = r.WithContext(ctx)
		handler.DeleteActor(w, r)

		require.Equal(t, http.StatusForbidden, w.Code)

		w = httptest.NewRecorder()
		r = httptest.NewRequest("GET", "/film/search/asha", bytes.NewBufferString(""))
		ctx = context.WithValue(r.Context(), "tokenData", &auth.TokenData{Id: 1, Role: 0})
		r = r.WithContext(ctx)
		handler.DeleteFilm(w, r)

		require.Equal(t, http.StatusForbidden, w.Code)
	})
}
