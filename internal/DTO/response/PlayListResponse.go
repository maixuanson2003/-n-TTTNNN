package response

import "time"

type PlayListResponse struct {
	ID        int
	Name      string
	CreateDay time.Time
}
