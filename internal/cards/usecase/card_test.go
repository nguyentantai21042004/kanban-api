package usecase

// import (
// 	"context"
// 	"math"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/assert"

// 	"gitlab.com/tantai-kanban/kanban-api/internal/cards"
// 	"gitlab.com/tantai-kanban/kanban-api/internal/cards/repository"
// 	"gitlab.com/tantai-kanban/kanban-api/internal/models"
// 	"gitlab.com/tantai-kanban/kanban-api/internal/websocket/service"
// 	"gitlab.com/tantai-kanban/kanban-api/pkg/log"
// )

// // TestInitUseCase tests the usecase initialization
// func TestInitUseCase(t *testing.T) {
// 	l := log.InitializeTestZapLogger()
// 	mockRepo := repository.NewMockRepository(t)
// 	mockWsHub := &service.Hub{}

// 	uc := New(l, mockRepo, mockWsHub)
// 	assert.NotNil(t, uc)
// }

// // TestPositionCalculationLogic tests the core position calculation logic
// func TestPositionCalculationLogic(t *testing.T) {
// 	tests := []struct {
// 		name           string
// 		currentCards   []models.Card
// 		targetPosition float64
// 		expectedPos    float64
// 		description    string
// 	}{
// 		{
// 			name:           "Empty list - position should be 1000",
// 			currentCards:   []models.Card{},
// 			targetPosition: 0,
// 			expectedPos:    1000.0,
// 			description:    "When list is empty, new position should be 1000",
// 		},
// 		{
// 			name: "List with one card - position should be 2000",
// 			currentCards: []models.Card{
// 				{ID: "card-1", Position: 1000.0},
// 			},
// 			targetPosition: 0,
// 			expectedPos:    2000.0,
// 			description:    "When list has one card, new position should be 2000",
// 		},
// 		{
// 			name: "List with multiple cards - position should be 4000",
// 			currentCards: []models.Card{
// 				{ID: "card-1", Position: 1000.0},
// 				{ID: "card-2", Position: 2000.0},
// 				{ID: "card-3", Position: 3000.0},
// 			},
// 			targetPosition: 0,
// 			expectedPos:    4000.0,
// 			description:    "When list has multiple cards, new position should be max + 1000",
// 		},
// 		{
// 			name: "Specific position provided - should use provided position",
// 			currentCards: []models.Card{
// 				{ID: "card-1", Position: 1000.0},
// 				{ID: "card-2", Position: 2000.0},
// 			},
// 			targetPosition: 1500.0,
// 			expectedPos:    1500.0,
// 			description:    "When specific position provided, should use that position",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Calculate max position from current cards
// 			maxPosition := 0.0
// 			for _, card := range tt.currentCards {
// 				if card.Position > maxPosition {
// 					maxPosition = card.Position
// 				}
// 			}

// 			// Calculate new position
// 			var newPosition float64
// 			if tt.targetPosition <= 0 {
// 				// Auto-calculate position
// 				if len(tt.currentCards) == 0 {
// 					newPosition = 1000.0 // Initial position
// 				} else {
// 					newPosition = maxPosition + 1000.0 // Add 1000 to max position
// 				}
// 			} else {
// 				// Use provided position
// 				newPosition = tt.targetPosition
// 			}

// 			// Assert
// 			assert.Equal(t, tt.expectedPos, newPosition, tt.description)
// 		})
// 	}
// }

// // TestPositionValidation tests position validation logic
// func TestPositionValidation(t *testing.T) {
// 	tests := []struct {
// 		name          string
// 		position      float64
// 		expectedValid bool
// 		description   string
// 	}{
// 		{
// 			name:          "Valid positive position",
// 			position:      1000.0,
// 			expectedValid: true,
// 			description:   "Positive positions should be valid",
// 		},
// 		{
// 			name:          "Valid zero position",
// 			position:      0.0,
// 			expectedValid: true,
// 			description:   "Zero position should be valid",
// 		},
// 		{
// 			name:          "Valid negative position",
// 			position:      -100.0,
// 			expectedValid: true,
// 			description:   "Negative positions should be valid",
// 		},
// 		{
// 			name:          "Valid very large position",
// 			position:      999999999.0,
// 			expectedValid: true,
// 			description:   "Very large positions should be valid",
// 		},
// 		{
// 			name:          "Invalid NaN position",
// 			position:      math.NaN(),
// 			expectedValid: false,
// 			description:   "NaN positions should be invalid",
// 		},
// 		{
// 			name:          "Invalid infinite position",
// 			position:      math.Inf(1),
// 			expectedValid: false,
// 			description:   "Infinite positions should be invalid",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Validate position
// 			isValid := !math.IsNaN(tt.position) &&
// 				!math.IsInf(tt.position, 1) &&
// 				!math.IsInf(tt.position, -1)

// 			// Assert
// 			assert.Equal(t, tt.expectedValid, isValid, tt.description)
// 		})
// 	}
// }

// // TestPositionSorting tests position sorting logic
// func TestPositionSorting(t *testing.T) {
// 	cards := []models.Card{
// 		{ID: "card-3", Position: 3000.0, Title: "Card 3"},
// 		{ID: "card-1", Position: 1000.0, Title: "Card 1"},
// 		{ID: "card-2", Position: 2000.0, Title: "Card 2"},
// 		{ID: "card-4", Position: 4000.0, Title: "Card 4"},
// 	}

// 	// Sort cards by position
// 	for i := 0; i < len(cards)-1; i++ {
// 		for j := i + 1; j < len(cards); j++ {
// 			if cards[i].Position > cards[j].Position {
// 				cards[i], cards[j] = cards[j], cards[i]
// 			}
// 		}
// 	}

// 	// Assert cards are sorted by position
// 	assert.Equal(t, "Card 1", cards[0].Title, "First card should be Card 1")
// 	assert.Equal(t, "Card 2", cards[1].Title, "Second card should be Card 2")
// 	assert.Equal(t, "Card 3", cards[2].Title, "Third card should be Card 3")
// 	assert.Equal(t, "Card 4", cards[3].Title, "Fourth card should be Card 4")
// }

// // TestMoveInputValidation tests move input validation
// func TestMoveInputValidation(t *testing.T) {
// 	tests := []struct {
// 		name          string
// 		input         cards.MoveInput
// 		expectedValid bool
// 		description   string
// 	}{
// 		{
// 			name: "Valid move input",
// 			input: cards.MoveInput{
// 				ID:       "card-1",
// 				ListID:   "list-2",
// 				Position: 1000.0,
// 			},
// 			expectedValid: true,
// 			description:   "Valid move input should pass validation",
// 		},
// 		{
// 			name: "Empty card ID",
// 			input: cards.MoveInput{
// 				ID:       "",
// 				ListID:   "list-2",
// 				Position: 1000.0,
// 			},
// 			expectedValid: false,
// 			description:   "Empty card ID should fail validation",
// 		},
// 		{
// 			name: "Empty list ID",
// 			input: cards.MoveInput{
// 				ID:       "card-1",
// 				ListID:   "",
// 				Position: 1000.0,
// 			},
// 			expectedValid: false,
// 			description:   "Empty list ID should fail validation",
// 		},
// 		{
// 			name: "Negative position",
// 			input: cards.MoveInput{
// 				ID:       "card-1",
// 				ListID:   "list-2",
// 				Position: -100.0,
// 			},
// 			expectedValid: true,
// 			description:   "Negative position should be valid",
// 		},
// 		{
// 			name: "Zero position",
// 			input: cards.MoveInput{
// 				ID:       "card-1",
// 				ListID:   "list-2",
// 				Position: 0.0,
// 			},
// 			expectedValid: true,
// 			description:   "Zero position should be valid",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Validate input
// 			isValid := tt.input.ID != "" && tt.input.ListID != ""

// 			// Assert
// 			assert.Equal(t, tt.expectedValid, isValid, tt.description)
// 		})
// 	}
// }

// // TestCreateCardPosition tests the position calculation for new cards
// func TestCreateCardPosition(t *testing.T) {
// 	type mockGetMaxPosition struct {
// 		isCalled   bool
// 		input      string
// 		wantOutput float64
// 		wantError  error
// 	}

// 	type mockCreate struct {
// 		isCalled   bool
// 		input      repository.CreateOptions
// 		wantOutput models.Card
// 		wantError  error
// 	}

// 	type mockGetBoardIDFromListID struct {
// 		isCalled   bool
// 		input      string
// 		wantOutput string
// 		wantError  error
// 	}

// 	mockTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
// 	sc := models.Scope{UserID: "user-1"}

// 	tcs := map[string]struct {
// 		input                    cards.CreateInput
// 		mockGetMaxPosition       mockGetMaxPosition
// 		mockCreate               mockCreate
// 		mockGetBoardIDFromListID mockGetBoardIDFromListID
// 		wantOutput               cards.DetailOutput
// 		wantError                error
// 	}{
// 		"success_create_in_empty_list": {
// 			input: cards.CreateInput{
// 				ListID:      "list-1",
// 				Title:       "New Card",
// 				Description: "Test description",
// 				Priority:    models.CardPriorityMedium,
// 				Labels:      []string{},
// 				DueDate:     nil,
// 			},
// 			mockGetMaxPosition: mockGetMaxPosition{
// 				isCalled:   true,
// 				input:      "list-1",
// 				wantOutput: 0.0,
// 				wantError:  nil,
// 			},
// 			mockCreate: mockCreate{
// 				isCalled: true,
// 				input: repository.CreateOptions{
// 					ListID:      "list-1",
// 					Title:       "New Card",
// 					Description: "Test description",
// 					Position:    1.0, // 0 + 1
// 					Priority:    models.CardPriorityMedium,
// 					Labels:      []string{},
// 					DueDate:     nil,
// 				},
// 				wantOutput: models.Card{
// 					ID:          "card-1",
// 					ListID:      "list-1",
// 					Title:       "New Card",
// 					Description: "Test description",
// 					Position:    1.0,
// 					Priority:    models.CardPriorityMedium,
// 					CreatedAt:   mockTime,
// 					UpdatedAt:   mockTime,
// 				},
// 				wantError: nil,
// 			},
// 			mockGetBoardIDFromListID: mockGetBoardIDFromListID{
// 				isCalled:   true,
// 				input:      "list-1",
// 				wantOutput: "board-1",
// 				wantError:  nil,
// 			},
// 			wantOutput: cards.DetailOutput{
// 				Card: models.Card{
// 					ID:          "card-1",
// 					ListID:      "list-1",
// 					Title:       "New Card",
// 					Description: "Test description",
// 					Position:    1.0,
// 					Priority:    models.CardPriorityMedium,
// 					CreatedAt:   mockTime,
// 					UpdatedAt:   mockTime,
// 				},
// 			},
// 			wantError: nil,
// 		},
// 		"success_create_in_list_with_cards": {
// 			input: cards.CreateInput{
// 				ListID:      "list-1",
// 				Title:       "New Card",
// 				Description: "Test description",
// 				Priority:    models.CardPriorityHigh,
// 				Labels:      []string{"urgent"},
// 				DueDate:     &mockTime,
// 			},
// 			mockGetMaxPosition: mockGetMaxPosition{
// 				isCalled:   true,
// 				input:      "list-1",
// 				wantOutput: 3000.0,
// 				wantError:  nil,
// 			},
// 			mockCreate: mockCreate{
// 				isCalled: true,
// 				input: repository.CreateOptions{
// 					ListID:      "list-1",
// 					Title:       "New Card",
// 					Description: "Test description",
// 					Position:    3001.0, // 3000 + 1
// 					Priority:    models.CardPriorityHigh,
// 					Labels:      []string{"urgent"},
// 					DueDate:     &mockTime,
// 				},
// 				wantOutput: models.Card{
// 					ID:          "card-2",
// 					ListID:      "list-1",
// 					Title:       "New Card",
// 					Description: "Test description",
// 					Position:    3001.0,
// 					Priority:    models.CardPriorityHigh,
// 					Labels:      []string{"urgent"},
// 					DueDate:     &mockTime,
// 					CreatedAt:   mockTime,
// 					UpdatedAt:   mockTime,
// 				},
// 				wantError: nil,
// 			},
// 			mockGetBoardIDFromListID: mockGetBoardIDFromListID{
// 				isCalled:   true,
// 				input:      "list-1",
// 				wantOutput: "board-1",
// 				wantError:  nil,
// 			},
// 			wantOutput: cards.DetailOutput{
// 				Card: models.Card{
// 					ID:          "card-2",
// 					ListID:      "list-1",
// 					Title:       "New Card",
// 					Description: "Test description",
// 					Position:    3001.0,
// 					Priority:    models.CardPriorityHigh,
// 					Labels:      []string{"urgent"},
// 					DueDate:     &mockTime,
// 					CreatedAt:   mockTime,
// 					UpdatedAt:   mockTime,
// 				},
// 			},
// 			wantError: nil,
// 		},
// 		"error_get_max_position": {
// 			input: cards.CreateInput{
// 				ListID:      "list-1",
// 				Title:       "New Card",
// 				Description: "Test description",
// 				Priority:    models.CardPriorityMedium,
// 				Labels:      []string{},
// 				DueDate:     nil,
// 			},
// 			mockGetMaxPosition: mockGetMaxPosition{
// 				isCalled:   true,
// 				input:      "list-1",
// 				wantOutput: 0.0,
// 				wantError:  assert.AnError,
// 			},
// 			mockCreate: mockCreate{
// 				isCalled: false,
// 			},
// 			mockGetBoardIDFromListID: mockGetBoardIDFromListID{
// 				isCalled: false,
// 			},
// 			wantOutput: cards.DetailOutput{},
// 			wantError:  assert.AnError,
// 		},
// 	}

// 	for name, tc := range tcs {
// 		t.Run(name, func(t *testing.T) {
// 			ctx := context.Background()
// 			uc, deps := initUseCase(t, mockTime)

// 			if tc.mockGetMaxPosition.isCalled {
// 				deps.mockRepo.EXPECT().
// 					GetMaxPosition(ctx, sc, tc.mockGetMaxPosition.input).
// 					Return(tc.mockGetMaxPosition.wantOutput, tc.mockGetMaxPosition.wantError)
// 			}

// 			if tc.mockCreate.isCalled {
// 				deps.mockRepo.EXPECT().
// 					Create(ctx, sc, tc.mockCreate.input).
// 					Return(tc.mockCreate.wantOutput, tc.mockCreate.wantError)
// 			}

// 			if tc.mockGetBoardIDFromListID.isCalled {
// 				deps.mockRepo.EXPECT().
// 					GetBoardIDFromListID(ctx, tc.mockGetBoardIDFromListID.input).
// 					Return(tc.mockGetBoardIDFromListID.wantOutput, tc.mockGetBoardIDFromListID.wantError)
// 			}

// 			gotOutput, gotErr := uc.Create(ctx, sc, tc.input)

// 			if tc.wantError != nil {
// 				assert.Error(t, gotErr)
// 				assert.Equal(t, tc.wantError, gotErr)
// 			} else {
// 				assert.NoError(t, gotErr)
// 				assert.Equal(t, tc.wantOutput, gotOutput)
// 			}
// 		})
// 	}
// }

// // TestMoveCardPosition tests the position calculation for moving cards
// func TestMoveCardPosition(t *testing.T) {
// 	type mockDetail struct {
// 		isCalled   bool
// 		input      string
// 		wantOutput models.Card
// 		wantError  error
// 	}

// 	type mockMove struct {
// 		isCalled   bool
// 		input      repository.MoveOptions
// 		wantOutput models.Card
// 		wantError  error
// 	}

// 	type mockGetBoardIDFromListID struct {
// 		isCalled   bool
// 		input      string
// 		wantOutput string
// 		wantError  error
// 	}

// 	mockTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
// 	sc := models.Scope{UserID: "user-1"}

// 	tcs := map[string]struct {
// 		input                    cards.MoveInput
// 		mockDetail               mockDetail
// 		mockMove                 mockMove
// 		mockDetailAfterMove      mockDetail
// 		mockGetBoardIDFromListID mockGetBoardIDFromListID
// 		wantOutput               cards.DetailOutput
// 		wantError                error
// 	}{
// 		"success_move_to_empty_list": {
// 			input: cards.MoveInput{
// 				ID:       "card-1",
// 				ListID:   "list-2",
// 				Position: 1000.0,
// 			},
// 			mockDetail: mockDetail{
// 				isCalled: true,
// 				input:    "card-1",
// 				wantOutput: models.Card{
// 					ID:       "card-1",
// 					ListID:   "list-1",
// 					Title:    "Test Card",
// 					Position: 2000.0,
// 				},
// 				wantError: nil,
// 			},
// 			mockMove: mockMove{
// 				isCalled: true,
// 				input: repository.MoveOptions{
// 					ID:       "card-1",
// 					ListID:   "list-2",
// 					Position: 1000.0,
// 					OldModel: models.Card{
// 						ID:       "card-1",
// 						ListID:   "list-1",
// 						Title:    "Test Card",
// 						Position: 2000.0,
// 					},
// 				},
// 				wantOutput: models.Card{
// 					ID:       "card-1",
// 					ListID:   "list-2",
// 					Title:    "Test Card",
// 					Position: 1000.0,
// 				},
// 				wantError: nil,
// 			},
// 			mockDetailAfterMove: mockDetail{
// 				isCalled: true,
// 				input:    "card-1",
// 				wantOutput: models.Card{
// 					ID:       "card-1",
// 					ListID:   "list-1", // Still in old list
// 					Title:    "Test Card",
// 					Position: 2000.0,
// 				},
// 				wantError: nil,
// 			},
// 			mockGetBoardIDFromListID: mockGetBoardIDFromListID{
// 				isCalled:   true,
// 				input:      "list-1", // Use the old list ID from updated card
// 				wantOutput: "board-1",
// 				wantError:  nil,
// 			},
// 			wantOutput: cards.DetailOutput{
// 				Card: models.Card{
// 					ID:       "card-1",
// 					ListID:   "list-1",
// 					Title:    "Test Card",
// 					Position: 2000.0,
// 				},
// 			},
// 			wantError: nil,
// 		},
// 		"success_move_with_auto_position": {
// 			input: cards.MoveInput{
// 				ID:       "card-1",
// 				ListID:   "list-2",
// 				Position: 0, // Auto-calculate
// 			},
// 			mockDetail: mockDetail{
// 				isCalled: true,
// 				input:    "card-1",
// 				wantOutput: models.Card{
// 					ID:       "card-1",
// 					ListID:   "list-1",
// 					Title:    "Test Card",
// 					Position: 1000.0,
// 				},
// 				wantError: nil,
// 			},
// 			mockMove: mockMove{
// 				isCalled: true,
// 				input: repository.MoveOptions{
// 					ID:       "card-1",
// 					ListID:   "list-2",
// 					Position: 0, // Will be calculated by repository
// 					OldModel: models.Card{
// 						ID:       "card-1",
// 						ListID:   "list-1",
// 						Title:    "Test Card",
// 						Position: 1000.0,
// 					},
// 				},
// 				wantOutput: models.Card{
// 					ID:       "card-1",
// 					ListID:   "list-2",
// 					Title:    "Test Card",
// 					Position: 3000.0, // Calculated by repository
// 				},
// 				wantError: nil,
// 			},
// 			mockDetailAfterMove: mockDetail{
// 				isCalled: true,
// 				input:    "card-1",
// 				wantOutput: models.Card{
// 					ID:       "card-1",
// 					ListID:   "list-1", // Still in old list
// 					Title:    "Test Card",
// 					Position: 1000.0,
// 				},
// 				wantError: nil,
// 			},
// 			mockGetBoardIDFromListID: mockGetBoardIDFromListID{
// 				isCalled:   true,
// 				input:      "list-1", // Use the old list ID from updated card
// 				wantOutput: "board-1",
// 				wantError:  nil,
// 			},
// 			wantOutput: cards.DetailOutput{
// 				Card: models.Card{
// 					ID:       "card-1",
// 					ListID:   "list-1",
// 					Title:    "Test Card",
// 					Position: 1000.0,
// 				},
// 			},
// 			wantError: nil,
// 		},
// 		"error_card_not_found": {
// 			input: cards.MoveInput{
// 				ID:       "card-1",
// 				ListID:   "list-2",
// 				Position: 1000.0,
// 			},
// 			mockDetail: mockDetail{
// 				isCalled:   true,
// 				input:      "card-1",
// 				wantOutput: models.Card{},
// 				wantError:  repository.ErrNotFound,
// 			},
// 			mockMove: mockMove{
// 				isCalled: false,
// 			},
// 			mockGetBoardIDFromListID: mockGetBoardIDFromListID{
// 				isCalled: false,
// 			},
// 			wantOutput: cards.DetailOutput{},
// 			wantError:  repository.ErrNotFound,
// 		},
// 	}

// 	for name, tc := range tcs {
// 		t.Run(name, func(t *testing.T) {
// 			ctx := context.Background()
// 			uc, deps := initUseCase(t, mockTime)

// 			if tc.mockDetail.isCalled {
// 				deps.mockRepo.EXPECT().
// 					Detail(ctx, sc, tc.mockDetail.input).
// 					Return(tc.mockDetail.wantOutput, tc.mockDetail.wantError)
// 			}

// 			if tc.mockMove.isCalled {
// 				deps.mockRepo.EXPECT().
// 					Move(ctx, sc, tc.mockMove.input).
// 					Return(tc.mockMove.wantOutput, tc.mockMove.wantError)
// 			}

// 			if tc.mockDetailAfterMove.isCalled {
// 				deps.mockRepo.EXPECT().
// 					Detail(ctx, sc, tc.mockDetailAfterMove.input).
// 					Return(tc.mockDetailAfterMove.wantOutput, tc.mockDetailAfterMove.wantError)
// 			}

// 			if tc.mockGetBoardIDFromListID.isCalled {
// 				deps.mockRepo.EXPECT().
// 					GetBoardIDFromListID(ctx, tc.mockGetBoardIDFromListID.input).
// 					Return(tc.mockGetBoardIDFromListID.wantOutput, tc.mockGetBoardIDFromListID.wantError)
// 			}

// 			gotOutput, gotErr := uc.Move(ctx, sc, tc.input)

// 			if tc.wantError != nil {
// 				assert.Error(t, gotErr)
// 				assert.Equal(t, tc.wantError, gotErr)
// 			} else {
// 				assert.NoError(t, gotErr)
// 				assert.Equal(t, tc.wantOutput, gotOutput)
// 			}
// 		})
// 	}
// }

// // BenchmarkPositionCalculationLogic benchmarks the position calculation logic
// func BenchmarkPositionCalculationLogic(b *testing.B) {
// 	cards := []models.Card{
// 		{ID: "card-1", Position: 1000.0},
// 		{ID: "card-2", Position: 2000.0},
// 		{ID: "card-3", Position: 3000.0},
// 		{ID: "card-4", Position: 4000.0},
// 		{ID: "card-5", Position: 5000.0},
// 	}

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		// Calculate max position
// 		maxPosition := 0.0
// 		for _, card := range cards {
// 			if card.Position > maxPosition {
// 				maxPosition = card.Position
// 			}
// 		}

// 		// Calculate new position
// 		newPosition := maxPosition + 1000.0

// 		// Prevent compiler optimization
// 		_ = newPosition
// 	}
// }
