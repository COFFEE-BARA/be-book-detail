
# 📚 Checkbara 서비스 소개

![image](https://github.com/COFFEE-BARA/crawler-kyobo-isbn/assets/65851554/aeecb4be-6f02-4bfd-a6e5-88df1ef87c8a)

![image](https://github.com/COFFEE-BARA/crawler-kyobo-isbn/assets/65851554/9ad000e0-5c95-4799-9582-2c49ecd5232b)


<br/>
<br/>


## 📹 시연영상




| ![챗봇](https://github.com/COFFEE-BARA/crawler-kyobo-isbn/assets/65851554/a28f0cca-ae1f-46b9-a087-6ea18216bd9d) | ![__-ezgif com-resize](https://github.com/COFFEE-BARA/crawler-kyobo-isbn/assets/65851554/699d09a6-8691-4ca0-bab5-575340a3c34d) | ![통계](https://github.com/COFFEE-BARA/crawler-kyobo-isbn/assets/65851554/d4ac2392-57ff-406f-b03e-334c932d56ce) |
| --- | --- | --- |
| AI 책 추천 | 책 검색 & 서점 재고 확인 & 대출 가능 도서관 확인 | 키바나 통계보기 |


<br/>
<br/>


## 📡 발전 방향

![image](https://github.com/COFFEE-BARA/crawler-kyobo-isbn/assets/65851554/d18f754c-773d-43b9-b961-7931703aaeb8)


<br/>
<br/>


# 👥 팀원 소개

| <img width="165" alt="suwha" src="https://github.com/COFFEE-BARA/be-bookstore-stock/assets/72396865/19e01fac-5384-4ec7-98f1-9e1e613429b4"> | <img width="165" alt="yoonju" src="https://github.com/COFFEE-BARA/be-bookstore-stock/assets/72396865/fb0a14c6-2d02-4105-962e-4565663817cc"> | <img width="165" alt="yugyeong" src="https://github.com/COFFEE-BARA/be-bookstore-stock/assets/72396865/90b7268d-92e5-43d1-9da8-ae48afd9e8c1"> | <img width="165" alt="dayeon" src="https://github.com/COFFEE-BARA/be-bookstore-stock/assets/72396865/f19e65e6-0856-4b6a-a355-993ce83ddcb7"> |
| --- | --- | --- | --- |
| 🐼[유수화](https://github.com/YuSuhwa-ve)🐼 | 🐱[송윤주](https://github.com/raminicano)🐱 | 🐶[현유경](https://github.com/yugyeongh)🐶 | 🐤[양다연](https://github.com/dayeon1201)🐤 |
| Server / Data / BE | AI / Data / BE | Infra / BE / FE | BE / FE |



<br/>
<br/>


# ⚒️ 전체 아키텍처

![image](https://github.com/COFFEE-BARA/crawler-kyobo-isbn/assets/65851554/04a7f4e8-2c00-4084-88e2-e7ffd273187a)





<br/>
<br/>

# 🍭 Elastic stack index 구조도

## book-index

<details>
<summary>book-index mapping</summary>
<div markdown="1">

```
// book-index mapping

{
  "mappings": {
    "properties": {
      "Author": {
        "type": "text",
        "fields": {
          "keyword": {
            "type": "text",
            "analyzer": "author_analyzer"
          },
          "partial": {
            "type": "text",
            "analyzer": "edge_ngram_analyzer"
          }
        }
      },
      "DetailCategory": {
        "type": "keyword"
      },
      "ISBN": {
        "type": "keyword"
      },
      "ImageURL": {
        "type": "keyword"
      },
      "IndexContent": {
        "type": "text"
      },
      "Introduction": {
        "type": "text"
      },
      "MiddleCategory": {
        "type": "keyword"
      },
      "Price": {
        "type": "integer"
      },
      "PubDate": {
        "type": "date",
        "format": "yyyy-MM-dd"
      },
      "Publisher": {
        "type": "keyword"
      },
      "PublisherReview": {
        "type": "text"
      },
      "PurchaseURL": {
        "type": "keyword"
      },
      "Search": {
        "type": "text"
      },
      "Title": {
        "type": "text",
        "analyzer": "title_analyzer"
      },
      "Vector": {
        "type": "dense_vector",
        "dims": 768,
        "index": true,
        "similarity": "cosine"
      },
      "document": {
        "type": "object"
      },
      "id": {
        "type": "text",
        "fields": {
          "keyword": {
            "type": "keyword",
            "ignore_above": 256
          }
        }
      },
      "index": {
        "type": "text",
        "fields": {
          "keyword": {
            "type": "keyword",
            "ignore_above": 256
          }
        }
      },
      "pipeline": {
        "type": "text",
        "fields": {
          "keyword": {
            "type": "keyword",
            "ignore_above": 256
          }
        }
      }
    }
  }
}
```

</div>
</details>


<details>
<summary> book-index settings </summary>
<div markdown="1">

```
//book -index settings
{
  "settings": {
    "index": {
      "routing": {
        "allocation": {
          "include": {
            "_tier_preference": "data_content"
          }
        }
      },
      "number_of_shards": "1",
      "provided_name": "book-index",
      "creation_date": "1708182319595",
      "analysis": {
        "filter": {
          "lowercase_filter": {
            "type": "lowercase"
          },
          "edge_ngram_filter": {
            "type": "edge_ngram",
            "min_gram": "1",
            "max_gram": "10"
          }
        },
        "analyzer": {
          "edge_ngram_analyzer": {
            "filter": [
              "edge_ngram_filter",
              "lowercase_filter"
            ],
            "type": "custom",
            "tokenizer": "nori_tokenizer_mine"
          },
          "author_analyzer": {
            "filter": [
              "lowercase_filter"
            ],
            "type": "custom",
            "tokenizer": "keyword"
          },
          "title_analyzer": {
            "filter": [
              "nori_readingform",
              "lowercase_filter",
              "nori_part_of_speech"
            ],
            "type": "custom",
            "tokenizer": "nori_tokenizer_mine"
          }
        },
        "tokenizer": {
          "nori_tokenizer_mine": {
            "type": "nori_tokenizer",
            "decompound_mode": "mixed"
          }
        }
      },
      "number_of_replicas": "2",
      "uuid": "okUbOg_pTJKVG2WO7e3rYQ",
      "version": {
        "created": "8500003"
      }
    }
  },
  "defaults": {
    "index": {
      "flush_after_merge": "512mb",
      "time_series": {
        "end_time": "9999-12-31T23:59:59.999Z",
        "start_time": "-9999-01-01T00:00:00Z",
        "es87tsdb_codec": {
          "enabled": "true"
        }
      },
      "final_pipeline": "_none",
      "max_inner_result_window": "100",
      "unassigned": {
        "node_left": {
          "delayed_timeout": "1m"
        }
      },
      "max_terms_count": "65536",
      "rollup": {
        "source": {
          "name": "",
          "uuid": ""
        }
      },
      "lifecycle": {
        "prefer_ilm": "true",
        "rollover_alias": "",
        "origination_date": "-1",
        "name": "",
        "parse_origination_date": "false",
        "step": {
          "wait_time_threshold": "12h"
        },
        "indexing_complete": "false"
      },
      "mode": "standard",
      "routing_partition_size": "1",
      "force_memory_term_dictionary": "false",
      "max_docvalue_fields_search": "100",
      "merge": {
        "scheduler": {
          "max_thread_count": "1",
          "auto_throttle": "true",
          "max_merge_count": "6"
        },
        "policy": {
          "merge_factor": "32",
          "floor_segment": "2mb",
          "max_merge_at_once_explicit": "30",
          "max_merge_at_once": "10",
          "max_merged_segment": "0b",
          "expunge_deletes_allowed": "10.0",
          "segments_per_tier": "10.0",
          "type": "UNSET",
          "deletes_pct_allowed": "20.0"
        }
      },
      "max_refresh_listeners": "1000",
      "max_regex_length": "1000",
      "load_fixed_bitset_filters_eagerly": "true",
      "number_of_routing_shards": "1",
      "write": {
        "wait_for_active_shards": "1"
      },
      "verified_before_close": "false",
      "mapping": {
        "coerce": "false",
        "nested_fields": {
          "limit": "50"
        },
        "depth": {
          "limit": "20"
        },
        "field_name_length": {
          "limit": "9223372036854775807"
        },
        "total_fields": {
          "limit": "1000"
        },
        "nested_objects": {
          "limit": "10000"
        },
        "ignore_malformed": "false",
        "dimension_fields": {
          "limit": "21"
        }
      },
      "source_only": "false",
      "soft_deletes": {
        "enabled": "true",
        "retention": {
          "operations": "0"
        },
        "retention_lease": {
          "period": "12h"
        }
      },
      "max_script_fields": "32",
      "query": {
        "default_field": [
          "*"
        ],
        "parse": {
          "allow_unmapped_fields": "true"
        }
      },
      "format": "0",
      "frozen": "false",
      "sort": {
        "missing": [],
        "mode": [],
        "field": [],
        "order": []
      },
      "priority": "1",
      "routing_path": [],
      "version": {
        "compatibility": "8500003"
      },
      "codec": "default",
      "max_rescore_window": "10000",
      "bloom_filter_for_id_field": {
        "enabled": "true"
      },
      "max_adjacency_matrix_filters": "100",
      "analyze": {
        "max_token_count": "10000"
      },
      "gc_deletes": "60s",
      "top_metrics_max_size": "10",
      "optimize_auto_generated_id": "true",
      "max_ngram_diff": "1",
      "hidden": "false",
      "translog": {
        "flush_threshold_age": "1m",
        "generation_threshold_size": "64mb",
        "flush_threshold_size": "10gb",
        "sync_interval": "5s",
        "retention": {
          "size": "-1",
          "age": "-1"
        },
        "durability": "REQUEST"
      },
      "auto_expand_replicas": "false",
      "fast_refresh": "false",
      "recovery": {
        "type": ""
      },
      "requests": {
        "cache": {
          "enable": "true"
        }
      },
      "data_path": "",
      "highlight": {
        "max_analyzed_offset": "1000000",
        "weight_matches_mode": {
          "enabled": "true"
        }
      },
      "look_back_time": "2h",
      "routing": {
        "rebalance": {
          "enable": "all"
        },
        "allocation": {
          "disk": {
            "watermark": {
              "ignore": "false"
            }
          },
          "enable": "all",
          "total_shards_per_node": "-1"
        }
      },
      "search": {
        "slowlog": {
          "level": "TRACE",
          "threshold": {
            "fetch": {
              "warn": "-1",
              "trace": "-1",
              "debug": "-1",
              "info": "-1"
            },
            "query": {
              "warn": "-1",
              "trace": "-1",
              "debug": "-1",
              "info": "-1"
            }
          }
        },
        "idle": {
          "after": "30s"
        },
        "throttled": "false"
      },
      "fielddata": {
        "cache": "node"
      },
      "look_ahead_time": "2h",
      "default_pipeline": "_none",
      "max_slices_per_scroll": "1024",
      "shard": {
        "check_on_startup": "false"
      },
      "xpack": {
        "watcher": {
          "template": {
            "version": ""
          }
        },
        "version": "",
        "ccr": {
          "following_index": "false"
        }
      },
      "percolator": {
        "map_unmapped_fields_as_text": "false"
      },
      "allocation": {
        "max_retries": "5",
        "existing_shards_allocator": "gateway_allocator"
      },
      "refresh_interval": "1s",
      "indexing": {
        "slowlog": {
          "reformat": "true",
          "threshold": {
            "index": {
              "warn": "-1",
              "trace": "-1",
              "debug": "-1",
              "info": "-1"
            }
          },
          "source": "1000",
          "level": "TRACE"
        }
      },
      "compound_format": "1gb",
      "blocks": {
        "metadata": "false",
        "read": "false",
        "read_only_allow_delete": "false",
        "read_only": "false",
        "write": "false"
      },
      "max_result_window": "10000",
      "store": {
        "stats_refresh_interval": "10s",
        "type": "",
        "fs": {
          "fs_lock": "native"
        },
        "preload": [],
        "snapshot": {
          "snapshot_name": "",
          "index_uuid": "",
          "cache": {
            "prewarm": {
              "enabled": "true"
            },
            "enabled": "true",
            "excluded_file_types": []
          },
          "repository_uuid": "",
          "uncached_chunk_size": "-1b",
          "delete_searchable_snapshot": "false",
          "index_name": "",
          "partial": "false",
          "blob_cache": {
            "metadata_files": {
              "max_length": "64kb"
            }
          },
          "repository_name": "",
          "snapshot_uuid": ""
        }
      },
      "queries": {
        "cache": {
          "enabled": "true"
        }
      },
      "shard_limit": {
        "group": "normal"
      },
      "warmer": {
        "enabled": "true"
      },
      "downsample": {
        "origin": {
          "name": "",
          "uuid": ""
        },
        "source": {
          "name": "",
          "uuid": ""
        },
        "status": "unknown"
      },
      "override_write_load_forecast": "0.0",
      "max_shingle_diff": "3",
      "query_string": {
        "lenient": "false"
      }
    }
  }
}
```


</div>
</details>



<br/>

## 데이터 관리전략


![image](https://github.com/COFFEE-BARA/crawler-kyobo-isbn/assets/65851554/10c4db70-82ec-4219-a9cb-77a2fe11e69b)

<br/>
<br/>


# 🏆 Tech Stack


## Programming language
<img src="https://img.shields.io/badge/go-00ADD8?style=for-the-badge&logo=go&logoColor=white"/>

<br/>

## DB

<img src="https://img.shields.io/badge/elastic-005571?style=for-the-badge&logo=elastic&logoColor=white">

<br/>

## Deploy & CI/CD
 <img src="https://img.shields.io/badge/amazonapigateway-FF4F8B?style=for-the-badge&logo=amazonapigateway&logoColor=white"/> <img src="https://img.shields.io/badge/lambda-FF9900?style=for-the-badge&logo=awslambda&logoColor=white"/>  <img src="https://img.shields.io/badge/docker-2496ED?style=for-the-badge&logo=docker&logoColor=white"> <img src="https://img.shields.io/badge/ecr-FC4C02?style=for-the-badge&logo=ecr&logoColor=white"> <img src="https://img.shields.io/badge/codebuild-68A51C?style=for-the-badge&logo=codebuild&logoColor=white"> <img src="https://img.shields.io/badge/codepipeline-527FFF?style=for-the-badge&logo=codepipeline&logoColor=white"> 

<br/>

## Develop Tool

<img src="https://img.shields.io/badge/postman-FF6C37?style=for-the-badge&logo=postman&logoColor=white"> <img src="https://img.shields.io/badge/github-181717?style=for-the-badge&logo=github&logoColor=white"> <img src="https://img.shields.io/badge/git-F05032?style=for-the-badge&logo=git&logoColor=white"> 

<br/>

## Communication Tool

<img src="https://img.shields.io/badge/slack-4A154B?style=for-the-badge&logo=slack&logoColor=white"> <img src="https://img.shields.io/badge/notion-000000?style=for-the-badge&logo=notion&logoColor=white">





<br/>
<br/>

# 🤖 API 명세

- URL: BASE_URL/api/book/{책의 isbn 값}/detail
- Method: `GET`
- 기능 소개: 책의 상세 페이지 정보

<br/>

# 🗣️ Request

## ☝🏻Request Header

```
Content-Type: application/json
```

## ✌🏻Request Params

| Name | Type | Description | Required |
| --- | --- | --- | --- |
| 책의 { isbn 값 } | String | 책의 13자리 isbn 값 | Required |


<br/>

# 🗣️ Response

## ☝🏻Response Body

```json
{
    "code": 200,
    "message": "책의 상세 정보를 가져오는데 성공했습니다.",
    "data": {
        "isbn": "9791140708116",
        "title": "아는 만큼 보이는 백엔드 개발 (한 권으로 보는 백엔드 로드맵과 커리어 가이드)",
        "author": "정우현^이인^김보인",
        "image": "https://shopping-phinf.pstatic.net/main_4519670/45196700648.20240114070834.jpg",
        "publisher": "길벗",
        "publishingDate": "2024-01-19",
        "introduction": "백엔드 개발, 어떻게 시작해야 할지 막막한가요?\n나무가 아닌 숲을 보면 길이 보입니다!\n\n백엔드 로드맵을 따라가며 서버 개발에서 다루는 언어, 기술, 서비스 등을 소개하고 전체 동작 원리를 설명합니다. 처음 개발을 접하며 막막해 하는 입문자를 위해 서버 기초 지식은 물론 각 주제마다 〈추천 프로젝트〉를 제시합니다. 책을 다 읽고 나면 백엔드 개발 전체를 볼 수 있는 시야를 가지게 됩니다.",
        "tableOfContents": "PART 1 처음 만나는 백엔드1장 백엔드 시작하기1.1 웹 개발의 구조1.2 서버의 동작 원리1.3 백엔드 개발자가 하는 일1.4 백엔드 로드맵 소개[정리하기]PART 2 백엔드 로드맵 따라가기2장 CS 기초 지식2.1 네트워크2.2 운영체제2.3 데이터베이스\t2.4 자료구조[정리하기]3장 백엔드 개발 언어와 프레임워크3.1 들어가기 전에3.2 프로그래밍 패러다임3.3 백엔드 개발 언어3.4 백엔드 프레임워크3.5 백엔드 개발 언어와 프레임워크 선택 방법[추천 프로젝트][정리하기]4장 DBMS4.1 DBMS의 개요4.2 RDBMS에서의 CRUD4.3 NoSQL에서의 CRUD[추천 프로젝트][정리하기]5장 API5.1 API의 개요5.2 API의 유형5.3 API 명세서[추천 프로젝트][정리하기]6장 버전 관리 시스템6.1 버전 관리 시스템의 개요6.2 분산 버전 관리 시스템: 깃6.3 웹 기반 버전 관리 저장소: 깃허브[추천 프로젝트][정리하기]7장 클라우드 컴퓨팅7.1 클라우드 컴퓨팅의 개요7.2 클라우드 서비스: AWS7.3 AWS 서버 구축 방법[추천 프로젝트][정리하기]8장 가상화와 컨테이너8.1 가상화와 컨테이너의 개요8.2 컨테이너 플랫폼: 도커8.3 컨테이너 오케스트레이션[추천 프로젝트][정리하기]9장 웹 애플리케이션 아키텍처9.1 웹 애플리케이션 아키텍처의 개요9.2 웹 애플리케이션 아키텍처의 종류[추천 프로젝트][정리하기]10장 테스트와 CI/CD10.1 테스트의 개요10.2 테스트의 종류10.3 테스트 주도 개발10.4 CI/CD[추천 프로젝트][정리하기]11장 백엔드 개발 총정리11.1 프로젝트 소개11.2 프로젝트 생성 및 업로드하기11.3 도커 파일 생성 및 서버 세팅하기11.4 CI/CD 파이프라인 구축 및 배포하기PART 3 백엔드 전문가로 성장하기\t12장 백엔드 커리어 설계하기12.1 백엔드 개발자12.2 아키텍트12.3 DBA12.4 데브옵스 엔지니어12.5 프로젝트 매니저12.6 풀스택 개발자12.7 CTO[정리하기]",
        "publisherBookReview": "백엔드 개발, 어떻게 시작해야 할지 막막한가요?나무가 아닌 숲을 보면 길이 보입니다!백엔드 로드맵을 따라가며 서버 개발에서 다루는 언어, 기술, 서비스 등을 소개하고 전체 동작 원리를 설명합니다. 처음 개발을 접하며 막막해 하는 입문자를 위해 서버 기초 지식은 물론 각 주제마다 '추천 프로젝트'를 제시합니다. 책을 다 읽고 나면 백엔드 개발 전체를 볼 수 있는 시야를 가지게 됩니다.이 책에서 다루는 내용백엔드 시작하기CS 기초 지식 익히기백엔드 개발 언어+프레임워크 알아보기DBMS, API 이해하기버전 관리 시스템 이해하기(깃, 깃허브)클라우드 컴퓨팅 이해하기(AWS)가상화와 컨테이너 이해하기(도커)웹 애플리케이션 이해하기테스트와 CI/CD 이해하기백엔드 커리어 패스 알아보기백엔드 개발에 막 입문했거나 공부 중인 분들께 강력 추천합니다.",
        "price": 21600,
        "purchaseURL": "https://search.shopping.naver.com/book/catalog/45196700648"
    }
}
```


## ✌🏻실패

1. 필요한 값이 없는 경우
    
    ```json
    {
      "code": 400,
      "message": "isbn값이 없습니다.",
      "data": null
    }
    ```
    
2. isbn 값에 매칭되는 책이 없을 경우
    
    ```json
    {
      "code": 404,
      "message": "없는 책입니다.",
      "data": null
    }
    ```
    
3. 서버에러
    
    ```json
    {
      "code": 500,
      "message": "서버 에러",
      "data": null
    }
    ```

