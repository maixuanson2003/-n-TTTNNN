package gemini

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"ten_module/internal/Config"

	"github.com/go-resty/resty/v2"
)

type GeminiClient struct {
	apiKey string
	client *resty.Client
}

var Chat *GeminiClient

func InitGeminiClient() {
	env := Config.GetEnvConfig()
	log.Print(env.GeminiAiKey())
	Chat = &GeminiClient{
		apiKey: env.GeminiAiKey(),
		client: resty.New(),
	}
}

func (g *GeminiClient) GenerateText(prompt string) (string, error) {
	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=" + g.apiKey
	body := map[string]interface{}{
		"contents": []map[string]interface{}{
			{"parts": []map[string]string{
				{"text": prompt},
			}},
		},
	}

	var response map[string]interface{}
	resp, err := g.client.R().
		SetBody(body).
		SetResult(&response).
		Post(url)

	if err != nil {
		return "", err
	}

	if !resp.IsSuccess() {
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode(), resp.String())
	}

	text := response["candidates"].([]interface{})[0].(map[string]interface{})["content"].(map[string]interface{})["parts"].([]interface{})[0].(map[string]interface{})["text"].(string)
	return text, nil
}

type MusicQuery struct {
	Genre    string `json:"genre"`
	Artist   string `json:"artist"`
	Song     string `json:"song"`
	Album    string `json:"album"`
	Intent   string `json:"intent"`
	Keywords string `json:"keywords"`

	// Mới thêm
	TimeRange string `json:"time_range"`
	SortBy    string `json:"sort_by"`
}

func getSimplePrompt() string {
	return `
Bạn là trợ lý âm nhạc. Phân tích câu người dùng và trả về JSON:

{
  "genre": "thể loại nhạc (pop, rap, ballad, rock, edm, vpop...)",
  "artist": "tên nghệ sĩ nếu có",
  "song": "tên bài hát cụ thể nếu có",
  "album": "tên album nếu có",
  "intent": "play, search, recommend, top",
  "keywords": "từ khóa tìm kiếm chung",
  "time_range": "today, week, month, year - nếu có yêu cầu về thời gian",
  "sort_by": "most_played, latest, top, popular - cách sắp xếp"
}

Quy tắc phân tích:
- "hôm nay", "today" → time_range: "today"
- "tuần này", "this week" → time_range: "week"  
- "tháng này", "this month" → time_range: "month"
- "năm nay", "this year" → time_range: "year"
- "mới nhất", "latest", "vừa ra" → sort_by: "latest"
- "hot nhất", "phổ biến", "trending" → sort_by: "popular"
- "hay nhất", "top", "xuất sắc" → sort_by: "top"
- "nghe nhiều", "được nghe nhiều" → sort_by: "most_played"

Chỉ trả về JSON, không giải thích.
`
}

// Prompt cho tìm kiếm nhạc
func getSearchPrompt() string {
	return `
Phân tích yêu cầu tìm kiếm nhạc và trả về JSON:

Ví dụ:
- "Tìm bài hát ballad mới nhất" → {"genre": "ballad", "intent": "search", "sort_by": "latest"}
- "Nghe Sơn Tùng MTP" → {"artist": "Sơn Tùng MTP", "intent": "play"}
- "Tìm bài Lạc Trôi" → {"song": "Lạc Trôi", "intent": "search"}
- "Album Hoàng Thùy Linh hot nhất" → {"artist": "Hoàng Thùy Linh", "intent": "search", "sort_by": "popular"}
- "Nhạc pop tuần này" → {"genre": "pop", "intent": "search", "time_range": "week"}
- "Bài hát được nghe nhiều nhất tháng này" → {"intent": "search", "time_range": "month", "sort_by": "most_played"}

{
  "genre": "",
  "artist": "",
  "song": "",
  "album": "",
  "intent": "play/search/recommend",
  "keywords": "",
  "time_range": "today/week/month/year",
  "sort_by": "most_played/latest/top/popular"
}

Chỉ trả về JSON.
`
}

// Prompt cho gợi ý nhạc
func getRecommendPrompt() string {
	return `
Phân tích yêu cầu gợi ý nhạc và trả về JSON:

Tập trung vào:
- Thể loại nhạc
- Nghệ sĩ yêu thích  
- Bài hát tương tự
- Album hay
- Thời gian phát hành
- Cách sắp xếp

Ví dụ:
- "Gợi ý nhạc pop mới nhất" → {"genre": "pop", "intent": "recommend", "sort_by": "latest"}
- "Nhạc giống như Đen Vâu" → {"artist": "Đen Vâu", "intent": "recommend"}
- "Tương tự bài Nơi Này Có Anh" → {"song": "Nơi Này Có Anh", "intent": "recommend"}
- "Gợi ý nhạc hot tuần này" → {"intent": "recommend", "time_range": "week", "sort_by": "popular"}
- "Nhạc hay nhất năm nay" → {"intent": "recommend", "time_range": "year", "sort_by": "top"}
- "Nhạc được nghe nhiều tháng này" → {"intent": "recommend", "time_range": "month", "sort_by": "most_played"}

{
  "genre": "",
  "artist": "",
  "song": "",
  "album": "",
  "intent": "recommend",
  "keywords": "",
  "time_range": "today/week/month/year",
  "sort_by": "most_played/latest/top/popular"
}

Chỉ trả về JSON.
`
}

// Prompt cho top nhạc
func getTopPrompt() string {
	return `
Phân tích yêu cầu tìm top nhạc và trả về JSON:

Tập trung vào:
- Top bài hát theo thể loại
- Top nghệ sĩ
- Top album
- Thời gian (hôm nay, tuần, tháng, năm)
- Tiêu chí sắp xếp

Ví dụ:
- "Top 10 bài hát hay nhất" → {"intent": "top", "sort_by": "top"}
- "Top nhạc pop tuần này" → {"genre": "pop", "intent": "top", "time_range": "week"}
- "BXH Vpop tháng này" → {"genre": "vpop", "intent": "top", "time_range": "month"}
- "Top bài hát được nghe nhiều nhất" → {"intent": "top", "sort_by": "most_played"}
- "Nhạc trending hôm nay" → {"intent": "top", "time_range": "today", "sort_by": "popular"}

{
  "genre": "",
  "artist": "",
  "song": "",
  "album": "",
  "intent": "top",
  "keywords": "",
  "time_range": "today/week/month/year",
  "sort_by": "most_played/latest/top/popular"
}

Chỉ trả về JSON.
`
}

func ExtractMusicInfo(userInput string) (*MusicQuery, error) {
	// Chọn prompt dựa trên từ khóa
	var prompt string
	input := strings.ToLower(userInput)

	if strings.Contains(input, "top") || strings.Contains(input, "bxh") ||
		strings.Contains(input, "bảng xếp hạng") || strings.Contains(input, "trending") ||
		strings.Contains(input, "thịnh hành") {
		prompt = getTopPrompt()
	} else if strings.Contains(input, "tìm") || strings.Contains(input, "search") {
		prompt = getSearchPrompt()
	} else if strings.Contains(input, "gợi ý") || strings.Contains(input, "recommend") ||
		strings.Contains(input, "muốn nghe") || strings.Contains(input, "đề xuất") {
		prompt = getRecommendPrompt()
	} else {
		prompt = getSimplePrompt()
	}

	fullPrompt := prompt + "\n\nCâu người dùng: " + userInput

	text, err := Chat.GenerateText(fullPrompt)
	if err != nil {
		return nil, err
	}

	// Clean JSON response - xử lý nhiều format khác nhau
	cleanText := strings.TrimSpace(text)

	// Loại bỏ markdown code blocks
	cleanText = strings.TrimPrefix(cleanText, "```json")
	cleanText = strings.TrimPrefix(cleanText, "```JSON")
	cleanText = strings.TrimPrefix(cleanText, "```")
	cleanText = strings.TrimSuffix(cleanText, "```")

	// Loại bỏ text thừa trước và sau JSON
	if startIdx := strings.Index(cleanText, "{"); startIdx != -1 {
		cleanText = cleanText[startIdx:]
	}
	if endIdx := strings.LastIndex(cleanText, "}"); endIdx != -1 {
		cleanText = cleanText[:endIdx+1]
	}

	cleanText = strings.TrimSpace(cleanText)
	log.Printf("Cleaned JSON Response: %s", cleanText)

	var query MusicQuery
	if err := json.Unmarshal([]byte(cleanText), &query); err != nil {
		log.Printf("JSON Parse Error: %v", err)
		log.Printf("Raw response: %s", text)
		return nil, fmt.Errorf("lỗi parse JSON: %w", err)
	}

	// Post-processing: normalize và validate data
	query = normalizeQuery(query)

	log.Printf("Final parsed query: %+v", query)
	return &query, nil
}

// Hàm normalize và validate dữ liệu
func normalizeQuery(query MusicQuery) MusicQuery {
	// Normalize intent
	intent := strings.ToLower(strings.TrimSpace(query.Intent))
	switch intent {
	case "play", "nghe":
		query.Intent = "play"
	case "search", "tìm", "tim":
		query.Intent = "search"
	case "recommend", "gợi ý", "goi y", "đề xuất", "de xuat":
		query.Intent = "recommend"
	case "top", "bxh", "trending":
		query.Intent = "top"
	default:
		if query.Song != "" {
			query.Intent = "play" // Nếu có tên bài hát cụ thể thì default là play
		} else {
			query.Intent = "search" // Default là search
		}
	}

	// Normalize time_range
	timeRange := strings.ToLower(strings.TrimSpace(query.TimeRange))
	switch timeRange {
	case "today", "hôm nay", "hom nay":
		query.TimeRange = "today"
	case "week", "tuần này", "tuan nay", "this week":
		query.TimeRange = "week"
	case "month", "tháng này", "thang nay", "this month":
		query.TimeRange = "month"
	case "year", "năm nay", "nam nay", "this year":
		query.TimeRange = "year"
	default:
		query.TimeRange = ""
	}

	// Normalize sort_by
	sortBy := strings.ToLower(strings.TrimSpace(query.SortBy))
	switch sortBy {
	case "most_played", "nghe nhiều", "nghe nhieu", "được nghe nhiều", "duoc nghe nhieu":
		query.SortBy = "most_played"
	case "latest", "mới nhất", "moi nhat", "vừa ra", "vua ra", "newest":
		query.SortBy = "latest"
	case "top", "hay nhất", "hay nhat", "xuất sắc", "xuat sac", "best":
		query.SortBy = "top"
	case "popular", "hot nhất", "hot nhat", "phổ biến", "pho bien", "trending":
		query.SortBy = "popular"
	default:
		// Set default sort based on intent
		switch query.Intent {
		case "top":
			query.SortBy = "top"
		case "play":
			query.SortBy = "most_played"
		default:
			query.SortBy = "latest"
		}
	}

	// Normalize genre
	genre := strings.ToLower(strings.TrimSpace(query.Genre))
	switch genre {
	case "vpop", "v-pop", "việt pop", "viet pop":
		query.Genre = "vpop"
	case "kpop", "k-pop", "hàn quốc", "han quoc", "korean pop":
		query.Genre = "kpop"
	case "usuk", "us-uk", "âu mỹ", "au my", "english":
		query.Genre = "âu mỹ"
	case "ballad", "balad":
		query.Genre = "ballad"
	case "rap", "hiphop", "hip-hop", "hip hop":
		query.Genre = "rap"
	case "pop":
		query.Genre = "pop"
	case "rock":
		query.Genre = "rock"
	case "edm", "electronic", "điện tử", "dien tu":
		query.Genre = "EDM"
	case "jazz":
		query.Genre = "jazz"
	case "blues":
		query.Genre = "blues"
	case "country":
		query.Genre = "country"
	case "folk", "dân ca", "dan ca":
		query.Genre = "folk"
	}

	// Clean empty strings
	if strings.TrimSpace(query.Genre) == "" {
		query.Genre = ""
	}
	if strings.TrimSpace(query.Artist) == "" {
		query.Artist = ""
	}
	if strings.TrimSpace(query.Song) == "" {
		query.Song = ""
	}
	if strings.TrimSpace(query.Album) == "" {
		query.Album = ""
	}
	if strings.TrimSpace(query.Keywords) == "" {
		query.Keywords = ""
	}

	return query
}
