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
        "/api/address/geocode": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Return a list of addresses provided geo coordinates",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "address"
                ],
                "summary": "Search by coordinates",
                "parameters": [
                    {
                        "description": "coordinates",
                        "name": "query",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/service.RequestAddressGeocode"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/service.ResponseAddress"
                        }
                    },
                    "400": {
                        "description": "bad request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "internal error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/address/search": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Return a list of addresses provided street name",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "address"
                ],
                "summary": "Search by street name",
                "parameters": [
                    {
                        "description": "street name",
                        "name": "query",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/service.RequestAddressSearch"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/service.ResponseAddress"
                        }
                    },
                    "400": {
                        "description": "bad request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "internal error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/login": {
            "post": {
                "description": "Authenticate user provided their email and password",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "authenticate user",
                "parameters": [
                    {
                        "description": "user credentials",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/auth.JSONResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/auth.JSONResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/auth.JSONResponse"
                        }
                    }
                }
            }
        },
        "/api/register": {
            "post": {
                "description": "Register new user provided email address and passport",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "register user",
                "parameters": [
                    {
                        "description": "user credentials",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.User"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/auth.JSONResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/auth.JSONResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "auth.JSONResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 400
                },
                "data": {},
                "message": {
                    "type": "string",
                    "example": "status bad request"
                }
            }
        },
        "auth.User": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "example": "admin@example.com"
                },
                "password": {
                    "type": "string",
                    "example": "password"
                }
            }
        },
        "service.Address": {
            "type": "object",
            "properties": {
                "city": {
                    "type": "string"
                },
                "house": {
                    "type": "string"
                },
                "lat": {
                    "type": "string"
                },
                "lon": {
                    "type": "string"
                },
                "street": {
                    "type": "string"
                }
            }
        },
        "service.RequestAddressGeocode": {
            "type": "object",
            "properties": {
                "lat": {
                    "type": "string",
                    "example": "55.753214"
                },
                "lng": {
                    "type": "string",
                    "example": "37.642589"
                }
            }
        },
        "service.RequestAddressSearch": {
            "type": "object",
            "properties": {
                "query": {
                    "type": "string",
                    "example": "Подкопаевский переулок"
                }
            }
        },
        "service.ResponseAddress": {
            "type": "object",
            "properties": {
                "addresses": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/service.Address"
                    }
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "2.0.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Geoservice API",
	Description:      "Geoservice with swagger docs and authentication",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
