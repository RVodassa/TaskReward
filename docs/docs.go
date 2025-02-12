// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "API Support",
            "email": "assadov.spb@bk.ru"
        },
        "license": {
            "name": "Free open source"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/login": {
            "post": {
                "description": "Возвращает JWT токен для доступа к защищенным маршрутам.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Аутентификация пользователя",
                "parameters": [
                    {
                        "description": "Логин и пароль",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.AuthRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Успешная аутентификация",
                        "schema": {
                            "$ref": "#/definitions/api.LoginResponse"
                        }
                    },
                    "400": {
                        "description": "Ошибка клиента",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    },
                    "403": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Ошибка на сервере",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/auth/register": {
            "post": {
                "description": "Создает нового пользователя, возвращает информацию о новом пользователе.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Регистрация пользователя",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID реферала, если нет укажите 0",
                        "name": "referID",
                        "in": "query",
                        "required": true
                    },
                    {
                        "description": "Логин и пароль",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.AuthRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Успешная регистрация",
                        "schema": {
                            "$ref": "#/definitions/api.StatusUserResponse"
                        }
                    },
                    "400": {
                        "description": "Ошибка клиента",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    },
                    "403": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Ошибка на сервере",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/users/leaderboard": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Возвращает топ 10 лидеров по балансу",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Получить список лидеров",
                "responses": {
                    "200": {
                        "description": "Успешно",
                        "schema": {
                            "$ref": "#/definitions/api.LeaderBoardResponse"
                        }
                    },
                    "400": {
                        "description": "Ошибка клиента",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    },
                    "403": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Ошибка на сервере",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/users/tasks/activetasks": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "возвращает список активных задач",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Tasks"
                ],
                "summary": "Получить список активных задач",
                "responses": {
                    "200": {
                        "description": "Успешно",
                        "schema": {
                            "$ref": "#/definitions/api.GetAllTasksResponse"
                        }
                    },
                    "400": {
                        "description": "Ошибка клиента",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    },
                    "403": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Ошибка на сервере",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/users/{userID}/status": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Возвращает информацию о пользователе в случае успешной операции",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Получить информацию о пользователе по ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID пользователя",
                        "name": "userID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Успешно",
                        "schema": {
                            "$ref": "#/definitions/api.StatusUserResponse"
                        }
                    },
                    "400": {
                        "description": "Ошибка клиента",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    },
                    "403": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Ошибка на сервере",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/users/{userID}/tasks/{taskID}/complete": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Возвращает информацию о выполненной задаче",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Tasks"
                ],
                "summary": "Выполнить задачу",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID пользователя",
                        "name": "userID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "ID задачи",
                        "name": "taskID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Успешно",
                        "schema": {
                            "$ref": "#/definitions/api.TaskCompletedResponse"
                        }
                    },
                    "400": {
                        "description": "Ошибка клиента",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    },
                    "403": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Ошибка на сервере",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "api.AuthRequest": {
            "type": "object",
            "properties": {
                "login": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "api.ErrorResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "status": {
                    "type": "boolean"
                }
            }
        },
        "api.GetAllTasksResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "status": {
                    "type": "boolean"
                },
                "tasks": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Task"
                    }
                }
            }
        },
        "api.LeaderBoardResponse": {
            "type": "object",
            "properties": {
                "listLeader": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.User"
                    }
                },
                "message": {
                    "type": "string"
                },
                "status": {
                    "type": "boolean"
                }
            }
        },
        "api.LoginResponse": {
            "type": "object",
            "properties": {
                "jwtoken": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                },
                "status": {
                    "type": "boolean"
                }
            }
        },
        "api.StatusUserResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "status": {
                    "type": "boolean"
                },
                "user": {
                    "$ref": "#/definitions/models.User"
                }
            }
        },
        "api.TaskCompletedResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "status": {
                    "type": "boolean"
                },
                "task": {
                    "$ref": "#/definitions/models.Task"
                }
            }
        },
        "models.Task": {
            "type": "object",
            "properties": {
                "bonus": {
                    "type": "integer"
                },
                "completed_at": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "status": {
                    "description": "\"не завершено\", \"завершено\"",
                    "type": "string"
                },
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "models.User": {
            "type": "object",
            "properties": {
                "balance": {
                    "type": "integer"
                },
                "created_at": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "login": {
                    "type": "string"
                },
                "refer_id": {
                    "type": "integer"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "description": "Укажите свой токен 'Bearer JWT_TOKEN'.",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "TaskReward API",
	Description:      "API для работы с пользователями и выполнения задач",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
