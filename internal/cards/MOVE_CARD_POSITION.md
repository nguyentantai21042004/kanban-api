# Move Card Position Calculation

## Tổng quan

API move card đã được cải thiện với cơ chế tính toán lại position của các card trong cùng list để đảm bảo thứ tự chính xác và tránh xung đột position.

## Cơ chế hoạt động

### 1. Tính toán Position mới

Khi move card, hệ thống sẽ:

- **Nếu position <= 0**: Tự động tính toán position phù hợp
- **Nếu position > 0**: Sử dụng position được chỉ định

### 2. Logic tính toán position

```go
func (r implRepository) calculateNewPosition(ctx context.Context, listID string, targetPosition float64) (float64, error) {
    // Lấy tất cả card trong list
    cards, err := dbmodels.Cards(
        dbmodels.CardWhere.ListID.EQ(listID),
    ).All(ctx, r.database)
    
    if len(cards) == 0 {
        return 1000.0, nil // List trống
    }
    
    // Sắp xếp cards theo position hiện tại
    util.Sort(cards, func(a, b *dbmodels.Card) bool {
        posA := 0.0
        posB := 0.0
        if a.Position.Big != nil {
            posA, _ = a.Position.Big.Float64()
        }
        if b.Position.Big != nil {
            posB, _ = b.Position.Big.Float64()
        }
        return posA < posB
    })
    
    // Tìm vị trí phù hợp
    for i, card := range cards {
        cardPos := 0.0
        if card.Position.Big != nil {
            cardPos, _ = card.Position.Big.Float64()
        }
        
        if targetPosition <= cardPos {
            if i == 0 {
                return cardPos / 2.0, nil // Chèn vào đầu
            }
            // Chèn vào giữa 2 card
            prevCard := cards[i-1]
            prevPos := 0.0
            if prevCard.Position.Big != nil {
                prevPos, _ = prevCard.Position.Big.Float64()
            }
            return (prevPos + cardPos) / 2.0, nil
        }
    }
    
    // Chèn vào cuối
    lastCard := cards[len(cards)-1]
    lastPos := 0.0
    if lastCard.Position.Big != nil {
        lastPos, _ = lastCard.Position.Big.Float64()
    }
    return lastPos + 1000.0, nil
}
```

### 3. Tính toán lại position của các card khác

Khi move card trong cùng list, hệ thống sẽ tính toán lại position của tất cả card khác:

```go
func (r implRepository) recalculatePositions(ctx context.Context, listID string, excludeCardID string) error {
    // Lấy tất cả card trong list (trừ card đang move)
    cards, err := dbmodels.Cards(
        dbmodels.CardWhere.ListID.EQ(listID),
        dbmodels.CardWhere.ID.NEQ(excludeCardID),
    ).All(ctx, r.database)
    
    // Sắp xếp cards theo position hiện tại
    util.Sort(cards, func(a, b *dbmodels.Card) bool {
        posA := 0.0
        posB := 0.0
        if a.Position.Big != nil {
            posA, _ = a.Position.Big.Float64()
        }
        if b.Position.Big != nil {
            posB, _ = b.Position.Big.Float64()
        }
        return posA < posB
    })
    
    // Tính toán lại position với khoảng cách 1000
    position := 1000.0
    for _, card := range cards {
        card.Position = types.Decimal{Big: decimal.New(int64(position), 0)}
        _, err := card.Update(ctx, r.database, boil.Whitelist(dbmodels.CardColumns.Position))
        if err != nil {
            return err
        }
        position += 1000.0
    }
    
    return nil
}
```

## Các trường hợp sử dụng

### 1. Move card trong cùng list

```go
// Card A: position 1000
// Card B: position 2000  
// Card C: position 3000

// Move Card B đến position 1500
// Kết quả:
// Card A: position 1000
// Card B: position 1500 (mới)
// Card C: position 2000 (được tính toán lại)
```

### 2. Move card sang list khác

```go
// List 1: Card A (1000), Card B (2000)
// List 2: Card C (1000)

// Move Card A sang List 2, position 1500
// Kết quả:
// List 1: Card B (1000 - được tính toán lại)
// List 2: Card C (1000), Card A (1500)
```

### 3. Move card với position tự động

```go
// List: Card A (1000), Card B (2000)

// Move Card C với position = 0 (tự động)
// Kết quả:
// List: Card C (500), Card A (1000), Card B (2000)
```

## Lợi ích

1. **Đảm bảo thứ tự chính xác**: Cards luôn được sắp xếp theo position tăng dần
2. **Tránh xung đột position**: Không có 2 card cùng position
3. **Khoảng cách hợp lý**: Position cách nhau 1000 đơn vị, dễ dàng chèn card mới
4. **Performance tốt**: Chỉ tính toán lại khi cần thiết
5. **Backward compatibility**: Vẫn hỗ trợ position cũ

## API Endpoint

```
POST /api/v1/cards/move
```

### Request Body

```json
{
  "id": "card-uuid",
  "list_id": "list-uuid", 
  "position": 1500.0
}
```

### Response

```json
{
  "card": {
    "id": "card-uuid",
    "list_id": "list-uuid",
    "position": 1500.0,
    "title": "Card Title",
    "description": "Card Description",
    "priority": "medium",
    "due_date": "2024-01-01T00:00:00Z",
    "labels": [],
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

## Lưu ý

- Position được lưu dưới dạng `NUMERIC(10,5)` trong database
- Khoảng cách mặc định giữa các position là 1000
- Khi position = 0, hệ thống tự động tính toán position phù hợp
- Activity log được tạo để theo dõi thay đổi position
- WebSocket event được broadcast để cập nhật real-time 