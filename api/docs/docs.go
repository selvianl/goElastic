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
        "/configs": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "configs"
                ],
                "summary": "List configs",
                "parameters": [
                    {
                        "type": "string",
                        "description": "pagination limit",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "active page",
                        "name": "page",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/api.PaginatedResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        " count": {
                                            "type": "integer"
                                        },
                                        "Results": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/models.SortConfig"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            },
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "configs"
                ],
                "summary": "Create config",
                "parameters": [
                    {
                        "description": "Create Config",
                        "name": "applicant",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.ConfigInput"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.SortConfig"
                        }
                    }
                }
            }
        },
        "/configs/{id}": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "configs"
                ],
                "summary": "Get config",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Config id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.SortConfig"
                        }
                    }
                }
            },
            "put": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "configs"
                ],
                "summary": "Update config",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Config ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Update Config",
                        "name": "address",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.ConfigInputUpdate"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.SortConfig"
                        }
                    }
                }
            },
            "delete": {
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "configs"
                ],
                "summary": "Delete Config",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Config id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    }
                }
            }
        },
        "/items/": {
            "post": {
                "tags": [
                    "items"
                ],
                "summary": "List Items",
                "parameters": [
                    {
                        "description": "search input",
                        "name": "esSearch",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.FilterParams"
                        }
                    },
                    {
                        "type": "string",
                        "description": "pagination limit",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "active page",
                        "name": "page",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/api.PaginatedResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        " count": {
                                            "type": "integer"
                                        },
                                        "Results": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/models.ItemOutput"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "api.PaginatedResponse": {
            "type": "object",
            "properties": {
                "count": {
                    "type": "integer"
                },
                "results": {}
            }
        },
        "models.ConfigInput": {
            "type": "object",
            "required": [
                "sort_option",
                "sort_order"
            ],
            "properties": {
                "is_active": {
                    "type": "boolean"
                },
                "sort_option": {
                    "type": "string"
                },
                "sort_order": {
                    "type": "string"
                }
            }
        },
        "models.ConfigInputUpdate": {
            "type": "object",
            "properties": {
                "is_active": {
                    "type": "boolean"
                },
                "sort_option": {
                    "type": "string"
                },
                "sort_order": {
                    "type": "string"
                }
            }
        },
        "models.FilterCondition": {
            "type": "object",
            "properties": {
                "field_name": {
                    "type": "string"
                },
                "operation": {
                    "type": "string"
                },
                "value": {
                    "type": "string"
                }
            }
        },
        "models.FilterParams": {
            "type": "object",
            "properties": {
                "conditions": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.FilterCondition"
                    }
                }
            }
        },
        "models.ItemOutput": {
            "type": "object",
            "properties": {
                "click": {
                    "type": "integer"
                },
                "item_id": {
                    "type": "string"
                },
                "locale": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "purchase": {
                    "type": "integer"
                }
            }
        },
        "models.SortConfig": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "is_active": {
                    "type": "boolean"
                },
                "sort_option": {
                    "type": "string"
                },
                "sort_order": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
