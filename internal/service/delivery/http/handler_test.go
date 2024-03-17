package http

import (
	"film_library/internal/service"
	"fmt"
	"github.com/stretchr/testify/require"
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
