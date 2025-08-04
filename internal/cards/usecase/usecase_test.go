package usecase

// import (
// 	"testing"
// 	"time"

// 	"gitlab.com/tantai-kanban/kanban-api/internal/cards"
// 	"gitlab.com/tantai-kanban/kanban-api/internal/cards/repository"
// 	"gitlab.com/tantai-kanban/kanban-api/pkg/log"
// )

// type mockDeps struct {
// 	mockRepo *repository.MockRepository
// }

// func initUseCase(t *testing.T, mockTime time.Time) (cards.UseCase, mockDeps) {
// 	t.Helper()
// 	l := log.InitializeTestZapLogger()

// 	mockRepo := repository.NewMockRepository(t)

// 	uc := implUsecase{
// 		l:     l,
// 		repo:  mockRepo,
// 		clock: func() time.Time { return mockTime },
// 	}

// 	return &uc, mockDeps{
// 		mockRepo: mockRepo,
// 	}
// }
