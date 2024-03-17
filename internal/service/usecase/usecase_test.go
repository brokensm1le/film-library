package usecase

import (
	"film_library/internal/service"
	mock_service "film_library/internal/service/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestActor(t *testing.T) {
	ctr := gomock.NewController(t)
	defer ctr.Finish()

	repo := mock_service.NewMockRepository(ctr)
	in := service.Actor{Name: "Sasha", Sex: "m", BDate: "1999-10-10"}
	detail := service.DetailsParams{Sort: "Name"}

	repo.EXPECT().CreateActor(&in).Return(nil).Times(1)
	repo.EXPECT().GetActor("Sasha").Return(&in, nil).Times(1)
	repo.EXPECT().UpdateActor("Sasha", &in).Return(nil).Times(1)
	repo.EXPECT().DeleteActor("Sasha").Return(nil).Times(1)
	repo.EXPECT().SearchActor("Sasha").Return([]string{"Sasha1", "Sasha2"}, nil).Times(1)
	repo.EXPECT().GetActors(&detail).Return([]service.Actor{in}, nil).Times(1)
	useCase := NewServiceUsecase(repo)
	err := useCase.CreateActor(&in)
	require.NoError(t, err)
	err = useCase.UpdateActor("Sasha", &in)
	require.NoError(t, err)
	err = useCase.DeleteActor("Sasha")
	require.NoError(t, err)
	resp, err := useCase.SearchActor("Sasha")
	require.NoError(t, err)
	require.Equal(t, resp, []string{"Sasha1", "Sasha2"})
	actors, err := useCase.GetActors(&detail)
	require.NoError(t, err)
	require.Equal(t, actors, []service.Actor{in})
	actor, err := useCase.GetActor("Sasha")
	require.NoError(t, err)
	require.Equal(t, actor, &in)
}

func TestFilm(t *testing.T) {
	ctr := gomock.NewController(t)
	defer ctr.Finish()

	repo := mock_service.NewMockRepository(ctr)
	in := service.Film{Name: "Sasha", Rating: 7.7, RDate: "1999-10-10", Desc: "nice file, klyanus`"}
	detail := service.DetailsParams{Sort: "Name"}

	repo.EXPECT().CreateFilm(&in).Return(nil).Times(1)
	repo.EXPECT().GetFilm("Rocky").Return(&in, nil).Times(1)
	repo.EXPECT().UpdateFilm("Rocky", &in).Return(nil).Times(1)
	repo.EXPECT().DeleteFilm("Rocky").Return(nil).Times(1)
	repo.EXPECT().SearchFilms("Rocky").Return([]string{"Rocky 1", "Rocky 2"}, nil).Times(1)
	repo.EXPECT().GetFilms(&detail).Return([]service.Film{in}, nil).Times(1)
	useCase := NewServiceUsecase(repo)
	err := useCase.CreateFilm(&in)
	require.NoError(t, err)
	err = useCase.UpdateFilm("Rocky", &in)
	require.NoError(t, err)
	err = useCase.DeleteFilm("Rocky")
	require.NoError(t, err)
	resp, err := useCase.SearchFilms("Rocky")
	require.NoError(t, err)
	require.Equal(t, resp, []string{"Rocky 1", "Rocky 2"})
	films, err := useCase.GetFilms(&detail)
	require.NoError(t, err)
	require.Equal(t, films, []service.Film{in})

	film, err := useCase.GetFilm("Rocky")
	require.NoError(t, err)
	require.Equal(t, film, &in)
}

func TestOther(t *testing.T) {
	ctr := gomock.NewController(t)
	defer ctr.Finish()

	repo := mock_service.NewMockRepository(ctr)
	actors := []string{"Milla Jovovich", "Mark Zakharov", "Cameron Diaz", "John Travolta", "James cameron", "Kate Winslet"}
	films := []string{"Forrest Gump", "The Shawshank Redemption", "The Social Network", "Pulp Fiction", "The King's Speech", "Dead Poets Society"}
	repo.EXPECT().AddFilmsByActor(&service.AddFilmsByActorParams{Films: films, Actor: actors[0]}).Return(nil).Times(1)
	repo.EXPECT().AddActorsByFilm(&service.AddActorsByFilmParams{Film: films[0], Actors: actors}).Return(nil).Times(1)
	repo.EXPECT().DeleteActorFilm(&service.DeleteActorFilmParams{Film: films[0], Actor: actors[0]}).Return(nil).Times(1)

	useCase := NewServiceUsecase(repo)
	err := useCase.AddActorsByFilm(&service.AddActorsByFilmParams{Film: films[0], Actors: actors})
	require.NoError(t, err)
	err = useCase.AddFilmsByActor(&service.AddFilmsByActorParams{Films: films, Actor: actors[0]})
	require.NoError(t, err)
	err = useCase.DeleteActorFilm(&service.DeleteActorFilmParams{Film: films[0], Actor: actors[0]})
	require.NoError(t, err)
}
