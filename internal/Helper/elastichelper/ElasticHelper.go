package elastichelper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"ten_module/internal/DTO/response"
	entity "ten_module/internal/Entity"
	"ten_module/internal/service/songservice"

	"github.com/elastic/go-elasticsearch/v8"
)

func CreateIndexElastic(ContentSearch string) error {
	ElasticSearch, err := elasticsearch.NewDefaultClient()
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
func InsertDataSongToIndex(Index string, Data []entity.Song) error {
	ElasticSearch, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Print(err)
		return err
	}
	DataSongRespone := []response.SongResponse{}
	for _, SongItem := range Data {
		DataSongRespone = append(DataSongRespone, songservice.SongEntityMapToSongResponse(SongItem))
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
func InsertDataArtistToIndex(Index string, Data []entity.Artist) {

}
