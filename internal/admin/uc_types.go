package admin

type DashboardInput struct {
	Period string // 7d | 30d | 90d
}

type DashboardUsers struct {
	Total  int64   `json:"total"`
	Active int64   `json:"active"`
	Growth float64 `json:"growth"`
}

type DashboardBoards struct {
	Total  int64 `json:"total"`
	Active int64 `json:"active"`
}

type DashboardCards struct {
	Total     int64 `json:"total"`
	Completed int64 `json:"completed"`
	Overdue   int64 `json:"overdue"`
}

type DashboardOutput struct {
	Users    DashboardUsers  `json:"users"`
	Boards   DashboardBoards `json:"boards"`
	Cards    DashboardCards  `json:"cards"`
	Activity []ActivityPoint `json:"activity"`
}

type ActivityPoint struct {
	Date           string `json:"date"`
	CardsCreated   int64  `json:"cards_created"`
	CardsCompleted int64  `json:"cards_completed"`
}

// Admin Users

type UsersInput struct {
	Search  string
	Page    int
	PerPage int
}

type RoleItem struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Alias string `json:"alias"`
}

type UserItem struct {
	ID          string   `json:"id"`
	Email       string   `json:"email"`
	FullName    string   `json:"full_name"`
	Role        RoleItem `json:"role"`
	IsActive    bool     `json:"is_active"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
	LastLoginAt *string  `json:"last_login_at"`
}

type UsersOutput struct {
	Items []UserItem `json:"items"`
	Meta  struct {
		Count       int64 `json:"count"`
		CurrentPage int   `json:"current_page"`
		PerPage     int   `json:"per_page"`
		Total       int64 `json:"total"`
		TotalPages  int   `json:"total_pages"`
	} `json:"meta"`
}

type CreateUserInput struct {
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	RoleID   string `json:"role_id"`
	Password string `json:"password"`
}

type UpdateUserInput struct {
	FullName  *string `json:"full_name"`
	RoleID    *string `json:"role_id"`
	RoleAlias *string `json:"role_alias"`
	IsActive  *bool   `json:"is_active"`
}

// Monitoring / Health

type HealthOutput struct {
	APIStatus               string  `json:"api_status"` // healthy | degraded | down
	ResponseTimeMs          float64 `json:"response_time_ms"`
	UptimePercentage        float64 `json:"uptime_percentage"`
	WebsocketConnections    int     `json:"websocket_connections"`
	WebsocketMessagesPerSec float64 `json:"websocket_messages_per_sec"`
	WebsocketAvgLatencyMs   float64 `json:"websocket_avg_latency_ms"`
	CheckedAt               string  `json:"checked_at"`
}
