// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import (
	"github.com/swaggo/swag"
	"os"
)

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "LazyJean",
            "email": "lazyjean@example.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/v1/login": {
            "post": {
                "responses": {}
            }
        },
        "/v1/register": {
            "post": {
                "responses": {}
            }
        },
        "/v1/words": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "分页获取单词列表",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "words"
                ],
                "summary": "获取单词列表",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer 令牌",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "default": 1,
                        "description": "页码",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 10,
                        "description": "每页数量",
                        "name": "perPage",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/handler.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "allOf": [
                                                {
                                                    "$ref": "#/definitions/handler.ListResponse"
                                                },
                                                {
                                                    "type": "object",
                                                    "properties": {
                                                        "items": {
                                                            "type": "array",
                                                            "items": {
                                                                "$ref": "#/definitions/dto.WordResponseDTO"
                                                            }
                                                        }
                                                    }
                                                }
                                            ]
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "401": {
                        "description": "未授权",
                        "schema": {
                            "$ref": "#/definitions/handler.Response"
                        }
                    },
                    "500": {
                        "description": "服务器内部错误",
                        "schema": {
                            "$ref": "#/definitions/handler.Response"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "创建一个新的单词记录",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "words"
                ],
                "summary": "创建新单词",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer 令牌",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "单词信息",
                        "name": "word",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.WordCreateDTO"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/handler.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/dto.WordResponseDTO"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "请求参数错误",
                        "schema": {
                            "$ref": "#/definitions/handler.Response"
                        }
                    },
                    "401": {
                        "description": "未授权",
                        "schema": {
                            "$ref": "#/definitions/handler.Response"
                        }
                    },
                    "500": {
                        "description": "服务器内部错误",
                        "schema": {
                            "$ref": "#/definitions/handler.Response"
                        }
                    }
                }
            }
        },
        "/v1/words/search": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "根据条件搜索单词",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "words"
                ],
                "summary": "搜索单词",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer 令牌",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "单词文本",
                        "name": "text",
                        "in": "query"
                    },
                    {
                        "type": "array",
                        "items": {
                            "type": "string"
                        },
                        "collectionFormat": "csv",
                        "description": "标签列表",
                        "name": "tags",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "最小难度",
                        "name": "minDifficulty",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "最大难度",
                        "name": "maxDifficulty",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 1,
                        "description": "页码",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 10,
                        "description": "每页数量",
                        "name": "pageSize",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/handler.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "allOf": [
                                                {
                                                    "$ref": "#/definitions/handler.ListResponse"
                                                },
                                                {
                                                    "type": "object",
                                                    "properties": {
                                                        "items": {
                                                            "type": "array",
                                                            "items": {
                                                                "$ref": "#/definitions/dto.WordResponseDTO"
                                                            }
                                                        }
                                                    }
                                                }
                                            ]
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "请求参数错误",
                        "schema": {
                            "$ref": "#/definitions/handler.Response"
                        }
                    },
                    "401": {
                        "description": "未授权",
                        "schema": {
                            "$ref": "#/definitions/handler.Response"
                        }
                    },
                    "500": {
                        "description": "服务器内部错误",
                        "schema": {
                            "$ref": "#/definitions/handler.Response"
                        }
                    }
                }
            }
        },
        "/v1/words/{id}": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "通过ID获取单词的详细信息",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "words"
                ],
                "summary": "获取单词详情",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer 令牌",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "单词ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/handler.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/dto.WordResponseDTO"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "请求参数错误",
                        "schema": {
                            "$ref": "#/definitions/handler.Response"
                        }
                    },
                    "401": {
                        "description": "未授权",
                        "schema": {
                            "$ref": "#/definitions/handler.Response"
                        }
                    },
                    "404": {
                        "description": "单词不存在",
                        "schema": {
                            "$ref": "#/definitions/handler.Response"
                        }
                    },
                    "500": {
                        "description": "服务器内部错误",
                        "schema": {
                            "$ref": "#/definitions/handler.Response"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "删除指定ID的单词",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "words"
                ],
                "summary": "删除单词",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer 令牌",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "单词ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.Response"
                        }
                    },
                    "400": {
                        "description": "请求参数错误",
                        "schema": {
                            "$ref": "#/definitions/handler.Response"
                        }
                    },
                    "401": {
                        "description": "未授权",
                        "schema": {
                            "$ref": "#/definitions/handler.Response"
                        }
                    },
                    "404": {
                        "description": "单词不存在",
                        "schema": {
                            "$ref": "#/definitions/handler.Response"
                        }
                    },
                    "500": {
                        "description": "服务器内部错误",
                        "schema": {
                            "$ref": "#/definitions/handler.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dto.WordCreateDTO": {
            "type": "object",
            "required": [
                "text",
                "translation"
            ],
            "properties": {
                "examples": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "Hello",
                        " world!"
                    ]
                },
                "phonetic": {
                    "type": "string",
                    "example": "həˈləʊ"
                },
                "tags": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "common",
                        "greeting"
                    ]
                },
                "text": {
                    "type": "string",
                    "example": "hello"
                },
                "translation": {
                    "type": "string",
                    "example": "你好"
                }
            }
        },
        "dto.WordResponseDTO": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string",
                    "example": "2025-01-26 18:00:00"
                },
                "examples": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "Hello",
                        " world!"
                    ]
                },
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "phonetic": {
                    "type": "string",
                    "example": "həˈləʊ"
                },
                "tags": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "common",
                        "greeting"
                    ]
                },
                "text": {
                    "type": "string",
                    "example": "hello"
                },
                "translation": {
                    "type": "string",
                    "example": "你好"
                },
                "updated_at": {
                    "type": "string",
                    "example": "2025-01-26 18:00:00"
                }
            }
        },
        "handler.ListResponse": {
            "type": "object",
            "properties": {
                "items": {
                    "description": "列表数据"
                },
                "page": {
                    "description": "当前页码",
                    "type": "integer"
                },
                "page_size": {
                    "description": "每页数量",
                    "type": "integer"
                },
                "total": {
                    "description": "总记录数",
                    "type": "integer"
                }
            }
        },
        "handler.Response": {
            "type": "object",
            "properties": {
                "code": {
                    "description": "业务状态码",
                    "type": "integer"
                },
                "data": {
                    "description": "响应数据"
                },
                "message": {
                    "description": "响应消息",
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "Bearer": {
            "description": "Bearer token for authentication",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             getSwaggerHost(),
	BasePath:         "/api",
	Schemes:          getSwaggerSchemes(),
	Title:            "生词本 API",
	Description:      "生词本服务 API 文档",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

// getSwaggerHost 根据环境返回适当的 host
func getSwaggerHost() string {
	if env := os.Getenv("GIN_MODE"); env != "release" {
		return "localhost:9000"
	}
	return "sla2.leeszi.cn"
}

// getSwaggerSchemes 根据环境返回适当的 schemes
func getSwaggerSchemes() []string {
	if env := os.Getenv("GIN_MODE"); env != "release" {
		return []string{"http"}
	}
	return []string{"https"}
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
