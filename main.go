package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/joho/godotenv"
)

type Response struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Data    *ResponesData `json:"data"`
}

type ResponesData struct {
	Isbn                string `json:"isbn"`
	Title               string `json:"title"`
	Author              string `json:"author"`
	Image               string `json:"image"`
	Publisher           string `json:"publisher"`
	PublishingDate      string `json:"publishingDate"` //"2024-01-19"
	Introduction        string `json:"introduction"`
	TableOfContents     string `json:"tableOfContents"`
	PublisherBookReview string `json:"publisherBookReview"`
	Price               int64  `json:"price"`
	PurchaseURL         string `json:"purchaseURL"`
}

func connectElasticSearch(CLOUD_ID, API_KEY string) (*elasticsearch.Client, error) {
	config := elasticsearch.Config{
		CloudID: CLOUD_ID,
		APIKey:  API_KEY,
	}

	es, err := elasticsearch.NewClient(config)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	fmt.Print("엘라스틱 클라이언트 : ", es)

	// Elasticsearch 서버에 핑을 보내 연결을 테스트합니다.
	res, err := es.Ping()
	if err != nil {
		fmt.Println("Elasticsearch와 연결 중 오류 발생:", err)
		return nil, err
	}
	defer res.Body.Close()

	fmt.Println("Elasticsearch 클라이언트가 성공적으로 연결되었습니다.")

	return es, nil

}

func searchIndex(es *elasticsearch.Client, indexName, fieldName, value string) ([]map[string]interface{}, error) {

	var allHits []map[string]interface{}
	//검색 쿼리 작성
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				fieldName: value,
			},
		},
	}

	// 쿼리를 JSON으로 변환합니다.
	queryJSON, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	// 검색 요청을 수행합니다.
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex(indexName),
		es.Search.WithBody(bytes.NewReader(queryJSON)),
	)
	if err != nil {
		return nil, err
	}

	// 검색 응답을 디코딩합니다.
	var searchResponse map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&searchResponse); err != nil {
		fmt.Println("검색 응답 디코딩 중 오류 발생:", err)
		return nil, err
	}

	// 히트를 추출하고 후 저장
	hits := searchResponse["hits"].(map[string]interface{})["hits"].([]interface{})
	for _, hit := range hits {
		allHits = append(allHits, hit.(map[string]interface{})["_source"].(map[string]interface{}))
	}

	return allHits, nil
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	//0. 환경변수
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	CLOUD_ID := os.Getenv("CLOUD_ID")
	API_KEY := os.Getenv("API_KEY")
	INDEX_NAME := os.Getenv("INDEX_NAME")
	FIELD_NAME := os.Getenv("FIELD_NAME")

	// 1. url path paramether로 isbn 값 받아오기
	isbn, ok := request.PathParameters["isbn"]
	if !ok {
		bodyJSON, err := json.Marshal(Response{
			Code:    400,
			Message: "isbn값이 없습니다.",
			Data:    nil,
		})
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 500}, err
		}

		// Return the response
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       string(bodyJSON),
		}, nil
	}
	fmt.Printf("파라미터 ISBN 값 : %s\n", isbn)

	//2. es cloud 연결하기
	esClient, err := connectElasticSearch(CLOUD_ID, API_KEY)
	if err != nil {
		fmt.Println("Error connecting to Elasticsearch:", err)
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	//3. isbn 값으로 검색하기
	res, err := searchIndex(esClient, INDEX_NAME, FIELD_NAME, isbn)
	if err != nil {
		fmt.Println("인덱스 검색 중 오류 발생:", err)
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	//4. 검색결과 파싱하여 response
	//4.1 res가 빈경우
	if len(res) == 0 {
		bodyJSON, err := json.Marshal(Response{
			Code:    404,
			Message: "없는 책입니다.",
			Data:    nil,
		})
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 500}, err
		}

		// Return the response
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       string(bodyJSON),
		}, nil
	}

	// 4.2 res가 존재하는경우
	bodyJSON, err := json.Marshal(Response{
		Code:    200,
		Message: "책의 상세 정보를 가져오는데 성공했습니다.",
		Data: &ResponesData{
			Isbn:                res[0]["ISBN"].(string),
			Title:               res[0]["Title"].(string),
			Author:              res[0]["Author"].(string),
			Image:               res[0]["ImageURL"].(string),
			Publisher:           res[0]["Publisher"].(string),
			PublishingDate:      res[0]["PubDate"].(string),
			Introduction:        res[0]["Introduction"].(string),
			TableOfContents:     res[0]["IndexContent"].(string),
			PublisherBookReview: res[0]["PublisherReview"].(string),
			Price:               int64(res[0]["Price"].(float64)), // float64로 처리 : escloud 의 response는 일반적으로 go에서 float64로 반환
			PurchaseURL:         res[0]["PurchaseURL"].(string),
		},
	})

	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(bodyJSON),
	}, nil

}

func main() {
	// 람다
	lambda.Start(handler)

	// //test~~~~~~~~~~~~~~~~~~~~~~~~~~
	// testEventFile, err := os.Open("test-event.json")
	// if err != nil {
	// 	log.Fatalf("Error opening test event file: %s", err)
	// }
	// defer testEventFile.Close()

	// // Decode the test event JSON
	// var testEvent events.APIGatewayProxyRequest
	// err = json.NewDecoder(testEventFile).Decode(&testEvent)
	// if err != nil {
	// 	log.Fatalf("Error decoding test event JSON: %s", err)
	// }

	// // Invoke the Lambda handler function with the test event
	// response, err := handler(context.Background(), testEvent)
	// if err != nil {
	// 	log.Fatalf("Error invoking Lambda handler: %s", err)
	// }

	// // Print the response
	// fmt.Printf("%v\n", response.StatusCode)
	// fmt.Printf("%v\n", response.Body)

}
