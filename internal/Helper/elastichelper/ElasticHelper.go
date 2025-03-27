package elastichelper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"ten_module/internal/DTO/response"
	entity "ten_module/internal/Entity"

	"github.com/elastic/go-elasticsearch/v8"
)

type ElasticHelper struct {
}

var ElasticHelpers *ElasticHelper

func InitElasticHelpers() {
	ElasticHelpers = &ElasticHelper{}
}

func (ElasticHelp *ElasticHelper) CreateIndexElastic(ContentSearch string) error {
	cfg := elasticsearch.Config{
		Addresses: []string{"http://127.0.0.1:9200"}, // ⚡ Đổi localhost thành IPv4
	}

	ElasticSearch, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Print(err)
		return err
	}
	Response, errorToCreateIndex := ElasticSearch.Indices.Create(ContentSearch)
	if errorToCreateIndex != nil {
		log.Print(errorToCreateIndex)
		return errorToCreateIndex
	}
	fmt.Print(Response)
	defer Response.Body.Close()
	return nil
}
func (ElasticHelp *ElasticHelper) InsertDataSongToIndex(Index string, Data []entity.Song) error {
	cfg := elasticsearch.Config{
		Addresses: []string{"http://127.0.0.1:9200"}, // ⚡ Đổi localhost thành IPv4
	}

	ElasticSearch, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Print(err)
		return err
	}
	DataSongRespone := []response.SongResponse{}
	for _, SongItem := range Data {
		DataSongRespone = append(DataSongRespone, response.SongResponse{
			ID:           SongItem.ID,
			NameSong:     SongItem.NameSong,
			Description:  SongItem.Description,
			ReleaseDay:   SongItem.ReleaseDay,
			CreateDay:    SongItem.CreateDay,
			UpdateDay:    SongItem.UpdateDay,
			Point:        SongItem.Point,
			LikeAmount:   SongItem.LikeAmount,
			Status:       SongItem.Status,
			CountryId:    SongItem.CountryId,
			ListenAmout:  SongItem.ListenAmout,
			AlbumId:      SongItem.AlbumId,
			SongResource: SongItem.SongResource,
		})
	}
	for _, SongResponseItem := range DataSongRespone {
		SongAdd, err := json.Marshal(SongResponseItem)
		if err != nil {
			log.Print(err)
			return err
		}
		res, errorToHanldeRequest := ElasticSearch.Index(Index, bytes.NewReader(SongAdd))
		if errorToHanldeRequest != nil {
			log.Print(errorToHanldeRequest)
			return errorToHanldeRequest
		}
		defer res.Body.Close()
		if res.IsError() {
			log.Printf("Lỗi từ Elasticsearch: %s", res.String())
		} else {
			log.Printf("Thêm bài hát thành công: %s", "ok")
		}
	}
	return nil

}

type ArtistElastic struct {
	ID          int
	Name        string
	BirthDay    string
	Description string
}

func (ElasticHelp *ElasticHelper) InsertDataArtistToIndex(Index string, Data []entity.Artist) error {
	cfg := elasticsearch.Config{
		Addresses: []string{"http://127.0.0.1:9200"}, // ⚡ Đổi localhost thành IPv4
	}

	ElasticSearch, err := elasticsearch.NewClient(cfg)

	if err != nil {
		log.Print(err)
		return err
	}
	ArtistElasticArray := []ArtistElastic{}
	for _, ArtistItem := range Data {
		ArtistElasticArray = append(ArtistElasticArray, ArtistElastic{
			ID:          ArtistItem.ID,
			Name:        ArtistItem.Name,
			BirthDay:    ArtistItem.BirthDay,
			Description: ArtistItem.Description,
		})
	}
	for _, ArtistItem := range ArtistElasticArray {
		ArtistAdd, errs := json.Marshal(ArtistItem)
		if errs != nil {
			log.Print(errs)
			return errs
		}
		res, errorToHanldeRequest := ElasticSearch.Index(Index, bytes.NewReader(ArtistAdd))
		if errorToHanldeRequest != nil {
			log.Print(errorToHanldeRequest)
			return errorToHanldeRequest
		}
		defer res.Body.Close()
		if res.IsError() {
			log.Printf("Lỗi từ Elasticsearch: %s", res.String())
		} else {
			log.Printf("Thêm bài hát thành công: %s", "ok")
		}
	}
	return nil
}

type SearchResult struct {
	Hits struct {
		Hits []struct {
			Source response.SongResponse `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

func (ElasticHelp *ElasticHelper) SearchSong(Keyword string) ([]response.SongResponse, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{"http://127.0.0.1:9200"}, // ⚡ Đổi localhost thành IPv4
	}

	ElasticSearch, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	query := fmt.Sprintf(`{
		"query": {
			"wildcard": {
				"NameSong": { 
				"value": "*%s*"
				}
			}
		}
	}`, Keyword)
	res, errorToSearch := ElasticSearch.Search(
		ElasticSearch.Search.WithIndex("song"),
		ElasticSearch.Search.WithBody(strings.NewReader(query)),
		ElasticSearch.Search.WithPretty(),
	)
	if errorToSearch != nil {
		log.Print(errorToSearch)
		return nil, errorToSearch
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	var search SearchResult
	json.Unmarshal(body, &search)
	Result := search.Hits.Hits
	SongResponse := []response.SongResponse{}
	for _, DocumentItem := range Result {
		SongResponse = append(SongResponse, DocumentItem.Source)
	}
	return SongResponse, nil

}

type searchArtistResult struct {
	Hits struct {
		Hits []struct {
			Source response.ArtistResponse `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

func (ElasticHelp *ElasticHelper) SearchArtist(Keyword string) ([]response.ArtistResponse, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{"http://127.0.0.1:9200"}, // ⚡ Đổi localhost thành IPv4
	}

	ElasticSearch, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	query := fmt.Sprintf(`{
		"query": {
			"wildcard": {
				"Name": {
				  "value": "*%s*"
				}
			}
		}
	}`, Keyword)
	res, errorToSearch := ElasticSearch.Search(
		ElasticSearch.Search.WithIndex("artist"),
		ElasticSearch.Search.WithBody(strings.NewReader(query)),
		ElasticSearch.Search.WithPretty(),
	)
	if errorToSearch != nil {
		log.Print(errorToSearch)
		return nil, errorToSearch
	}
	defer res.Body.Close()
	body, errorToRead := io.ReadAll(res.Body)
	if errorToRead != nil {
		log.Print(errorToRead)
		return nil, errorToRead
	}
	var search searchArtistResult
	json.Unmarshal(body, &search)
	Result := search.Hits.Hits
	ArtistResponse := []response.ArtistResponse{}
	for _, DocumentItem := range Result {
		ArtistResponse = append(ArtistResponse, DocumentItem.Source)
	}
	return ArtistResponse, nil

}
