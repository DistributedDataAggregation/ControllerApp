// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/query": {
            "post": {
                "description": "Queries data with specified grouping and selection",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "query"
                ],
                "summary": "Query data from table",
                "parameters": [
                    {
                        "description": "Query Request",
                        "name": "query",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.HttpQueryRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Query has been processed",
                        "schema": {
                            "$ref": "#/definitions/main.HttpResult"
                        }
                    },
                    "400": {
                        "description": "Invalid request payload",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "main.HttpQueryRequest": {
            "type": "object",
            "properties": {
                "group_columns": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "select": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/main.HttpSelect"
                    }
                },
                "table_name": {
                    "type": "string"
                }
            }
        },
        "main.HttpResult": {
            "type": "object",
            "properties": {
                "processing_time": {
                    "type": "integer"
                },
                "result": {
                    "$ref": "#/definitions/protomodels.QueryResponse"
                }
            }
        },
        "main.HttpSelect": {
            "type": "object",
            "properties": {
                "column": {
                    "type": "string"
                },
                "function": {
                    "type": "string"
                }
            }
        },
        "protomodels.Error": {
            "type": "object",
            "properties": {
                "inner_message": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "protomodels.PartialResult": {
            "type": "object",
            "properties": {
                "count": {
                    "type": "integer"
                },
                "value": {
                    "type": "integer"
                }
            }
        },
        "protomodels.QueryResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "$ref": "#/definitions/protomodels.Error"
                },
                "values": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/protomodels.Value"
                    }
                }
            }
        },
        "protomodels.Value": {
            "type": "object",
            "properties": {
                "grouping_value": {
                    "type": "string"
                },
                "results": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/protomodels.PartialResult"
                    }
                }
            }
        }
    },
    "externalDocs": {
        "description": "OpenAPI",
        "url": "https://swagger.io/resources/open-api/"
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Swagger Distributed data aggregation system API",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
