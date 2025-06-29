package gemini

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"ten_module/internal/Config"
	"ten_module/internal/repository"

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
	Genre     string `json:"genre"`
	Artist    string `json:"artist"`
	Intent    string `json:"intent"`
	Keywords  string `json:"keywords"`
	Country   string `json:"country"`
	TimeRange string `json:"time_range"`
	SortBy    string `json:"sort_by"`
}

func getSimplePrompt() string {
	Crepo, _ := repository.CountryRepo.FindAll()
	songtype, _ := repository.SongTypeRepo.FindAll()
	Art, _ := repository.ArtistRepo.FindAll()
	var countryNames []string
	for _, c := range Crepo {
		countryNames = append(countryNames, fmt.Sprintf(`"%s"`, c.CountryName))
	}

	var songTypes []string
	for _, s := range songtype {
		songTypes = append(songTypes, fmt.Sprintf(`"%s"`, s.Type))
	}

	var artistNames []string
	for _, a := range Art {
		artistNames = append(artistNames, fmt.Sprintf(`"%s"`, a.Name))
	}

	countriesStr := strings.Join(countryNames, ", ")
	songTypesStr := strings.Join(songTypes, ", ")
	artistsStr := strings.Join(artistNames, ", ")
	return fmt.Sprintf(`Bạn là trợ lý âm nhạc. Phân tích câu người dùng và trả về JSON:

{
  "genre": "thể loại nhạc (ví dụ: %s)",
  "artist": "tên nghệ sĩ nếu có (ví dụ: %s)",
  "intent": "play, search, recommend, top",
  "keywords": "từ khóa tìm kiếm chung",
  "country": "tên đất nước (ví dụ: %s)",
  "time_range": "week, month - nếu có yêu cầu về thời gian",
  "sort_by": "most_played, latest, top, popular - cách sắp xếp"
}

Quy tắc phân tích:

1. Ý định (intent):
- Nếu người dùng yêu cầu phát nhạc: intent = "play"
- Nếu người dùng yêu cầu tìm kiếm thông tin (tên bài, nghệ sĩ...): intent = "search"
- Nếu người dùng hỏi gợi ý ("phù hợp", "nên nghe gì", "gợi ý", "recommend"): intent = "recommend"
- Nếu người dùng hỏi nhạc top, hay nhất, xếp hạng: intent = "top"

2. Tên nghệ sĩ (artist):
- So khớp chuỗi trong câu với danh sách nghệ sĩ (%s)
- Nếu có khớp thì điền vào trường "artist"

3. Thể loại nhạc (genre):
- So khớp các từ với danh sách thể loại (%s)
- Nếu không khớp rõ ràng, nhưng ngữ cảnh như "buồn", "lái xe", "tập gym", "ngủ", "làm việc", "thư giãn","đám cưới","học tập" thì suy luận genre phù hợp:
  - "lái xe" → ["rock", "edm", "lofi"]
  - "tập gym" → ["edm", "hiphop", "remix"]
  - "thư giãn" → ["lofi", "acoustic", "jazz"]
  - "ngủ" → ["instrumental", "ambient", "piano"]
  - "làm việc" → ["lofi", "instrumental"]
  - "tiệc tùng" → ["edm", "dance", "pop"]

4. Quốc gia (country):
- Nếu có đề cập tên quốc gia hoặc vùng nhạc (ví dụ: "nhạc Hàn", "Kpop", "USUK", "Vpop") → country = tương ứng trong danh sách (%s)

5. Từ khóa:
- Nếu người dùng dùng từ không xác định rõ ràng (ví dụ: "bài buồn", "vui vẻ", "đi bar", "cảm xúc") thì đưa vào "keywords"

6. Thời gian và sắp xếp:
- "hôm nay", "today" → time_range: "today"
- "tuần này", "this week" → time_range: "week"
- "tháng này", "this month" → time_range: "month"
- "năm nay", "this year" → time_range: "year"
- "mới nhất", "latest", "vừa ra" → sort_by: "latest"
- "hot nhất", "phổ biến", "trending" → sort_by: "popular"
- "hay nhất", "top", "xuất sắc" → sort_by: "top"
- "nghe nhiều", "được nghe nhiều" → sort_by: "most_played"

Chỉ trả về JSON, không giải thích.
`, songTypesStr, artistsStr, countriesStr, artistsStr, songTypesStr, countriesStr)
}

func ExtractMusicInfo(userInput string) (*MusicQuery, error) {
	var prompt string
	prompt = getSimplePrompt()

	fullPrompt := prompt + "\n\nCâu người dùng: " + userInput

	text, err := Chat.GenerateText(fullPrompt)
	if err != nil {
		return nil, err
	}

	cleanText := strings.TrimSpace(text)

	cleanText = strings.TrimPrefix(cleanText, "```json")
	cleanText = strings.TrimPrefix(cleanText, "```JSON")
	cleanText = strings.TrimPrefix(cleanText, "```")
	cleanText = strings.TrimSuffix(cleanText, "```")

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

	log.Printf("Final parsed query: %+v", query)
	return &query, nil
}
