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
                "summary": "Query data from a table",
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
        },
        "/status": {
            "get": {
                "description": "Checks controller status",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "health check"
                ],
                "summary": "Health check endpoint",
                "responses": {
                    "200": {
                        "description": "Health check passed",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/tables": {
            "get": {
                "description": "Returns the names of non-empty folders in the data path",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tables"
                ],
                "summary": "List available tables",
                "responses": {
                    "200": {
                        "description": "List of table names",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
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
        },
        "/tables/columns": {
            "get": {
                "description": "Returns column names and their types for a given table. Filters out columns of unsupported types.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tables"
                ],
                "summary": "Get table columns",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Table name",
                        "name": "name",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "List of columns with their types",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "array",
                                "items": {
                                    "$ref": "#/definitions/main.ParquetColumnInfo"
                                }
                            }
                        }
                    },
                    "400": {
                        "description": "Bad request with error message",
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
        },
        "/tables/select-columns": {
            "get": {
                "description": "Returns column names and their types for a given table. Filters out columns of types unsupported for aggregations.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tables"
                ],
                "summary": "Get table columns that can be aggregated",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Table name",
                        "name": "name",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "List of columns with their types",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "array",
                                "items": {
                                    "$ref": "#/definitions/main.ParquetColumnInfo"
                                }
                            }
                        }
                    },
                    "400": {
                        "description": "Bad request with error message",
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
        },
        "/tables/upload": {
            "post": {
                "description": "Uploads a Parquet file (max 10 MB) to a table with a given name. If the table does not exist, it is created.",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tables"
                ],
                "summary": "Upload file to a table with given name",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Table name",
                        "name": "name",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "file",
                        "description": "File to upload (must have .parquet extension)",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "File uploaded successfully",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid table name or file",
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
        "main.HttpAggregateFunction": {
            "type": "string",
            "enum": [
                "Minimum",
                "Maximum",
                "Average",
                "Sum",
                "Count"
            ],
            "x-enum-varnames": [
                "Minimum",
                "Maximum",
                "Average",
                "Sum",
                "Count"
            ]
        },
        "main.HttpError": {
            "type": "object",
            "properties": {
                "inner_message": {
                    "description": "Inner error message",
                    "type": "string"
                },
                "message": {
                    "description": "Error message",
                    "type": "string"
                }
            }
        },
        "main.HttpPartialResult": {
            "type": "object",
            "properties": {
                "aggregation": {
                    "description": "The aggregate function.",
                    "allOf": [
                        {
                            "$ref": "#/definitions/main.HttpAggregateFunction"
                        }
                    ]
                },
                "double_value": {
                    "description": "Double result (nullable).",
                    "type": "number"
                },
                "float_value": {
                    "description": "Float result (nullable).",
                    "type": "number"
                },
                "int_value": {
                    "description": "Integer result (nullable).",
                    "type": "integer"
                },
                "is_null": {
                    "description": "Indicates if the result is null.",
                    "type": "boolean"
                },
                "result_type": {
                    "description": "Type of result: \"INT\", \"FLOAT\", \"DOUBLE\".",
                    "allOf": [
                        {
                            "$ref": "#/definitions/main.HttpResultType"
                        }
                    ]
                }
            }
        },
        "main.HttpQueryRequest": {
            "type": "object",
            "properties": {
                "group_columns": {
                    "description": "The names of the columns on which grouping will be performed.",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "select": {
                    "description": "A list of objects describing the columns and the aggregate functions to be executed on them.",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/main.HttpSelect"
                    }
                },
                "table_name": {
                    "description": "The name of the table on which the query will be executed.",
                    "type": "string"
                }
            }
        },
        "main.HttpQueryResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "description": "Information about the query processing error (if any)",
                    "allOf": [
                        {
                            "$ref": "#/definitions/main.HttpError"
                        }
                    ]
                },
                "values": {
                    "description": "List of results of performed aggregations for individual combinations of grouping column values",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/main.HttpValue"
                    }
                }
            }
        },
        "main.HttpResult": {
            "type": "object",
            "properties": {
                "processing_time": {
                    "description": "The total time to process the query in milliseconds.",
                    "type": "integer"
                },
                "result": {
                    "description": "The result of the processed query.",
                    "allOf": [
                        {
                            "$ref": "#/definitions/main.HttpQueryResponse"
                        }
                    ]
                }
            }
        },
        "main.HttpResultType": {
            "type": "string",
            "enum": [
                "INT",
                "FLOAT",
                "DOUBLE",
                "UNKNOWN"
            ],
            "x-enum-varnames": [
                "IntResult",
                "FloatResult",
                "DoubleResult",
                "UnknownResult"
            ]
        },
        "main.HttpSelect": {
            "type": "object",
            "properties": {
                "column": {
                    "description": "The name of the aggregated column.",
                    "type": "string"
                },
                "function": {
                    "description": "The aggregate function.",
                    "allOf": [
                        {
                            "$ref": "#/definitions/main.HttpAggregateFunction"
                        }
                    ]
                }
            }
        },
        "main.HttpValue": {
            "type": "object",
            "properties": {
                "grouping_value": {
                    "description": "Grouping value, subsequent grouping column values ​​separated by the '|' character",
                    "type": "string"
                },
                "results": {
                    "description": "List of results of given aggregations for the grouping value.",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/main.HttpPartialResult"
                    }
                }
            }
        },
        "main.ParquetColumnInfo": {
            "type": "object",
            "properties": {
                "name": {
                    "description": "The name of the column.",
                    "type": "string"
                },
                "type": {
                    "description": "The type of the column",
                    "allOf": [
                        {
                            "$ref": "#/definitions/main.ParquetColumnType"
                        }
                    ]
                }
            }
        },
        "main.ParquetColumnType": {
            "type": "string",
            "enum": [
                "INT",
                "DOUBLE",
                "FLOAT",
                "STRING",
                "BOOL",
                "UNSUPPORTED"
            ],
            "x-enum-varnames": [
                "INT",
                "DOUBLE",
                "FLOAT",
                "STRING",
                "BOOL",
                "UNSUPPORTED"
            ]
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
