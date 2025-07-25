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
        "/add-file-to-task": {
            "post": {
                "description": "Добавляет файл в архив задачи, ограничение — максимум 3 файла на задачу.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tasks"
                ],
                "summary": "Добавить файл к задаче",
                "parameters": [
                    {
                        "description": "Данные для добавления файла",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.AddFileToTaskRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Файл успешно добавлен к задаче",
                        "schema": {
                            "$ref": "#/definitions/handler.AddFileToTaskResponse"
                        }
                    },
                    "400": {
                        "description": "Неверный формат JSON или превышен лимит файлов",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/create-task": {
            "post": {
                "description": "Создаёт задачу, для которой можно добавлять файлы в ZIP архив.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tasks"
                ],
                "summary": "Создание новой задачи",
                "parameters": [
                    {
                        "description": "Путь и имя архива",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.CreateTaskRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Успешный ответ с ID созданной задачи",
                        "schema": {
                            "$ref": "#/definitions/handler.CreateTaskResponse"
                        }
                    },
                    "400": {
                        "description": "Неверный формат JSON",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "503": {
                        "description": "Ошибка создания задачи",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/get": {
            "get": {
                "description": "Возвращает статус задачи и ссылку на архив (если все файлы добавлены).",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tasks"
                ],
                "summary": "Получить статус задачи",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID задачи",
                        "name": "task-id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.TaskStatusResponse"
                        }
                    },
                    "400": {
                        "description": "некорректный ID задачи или задача не найдена",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "handler.AddFileToTaskRequest": {
            "type": "object",
            "properties": {
                "fileName": {
                    "type": "string",
                    "example": "test3"
                },
                "fileURL": {
                    "type": "string",
                    "example": "https://example.com/file.pdf"
                },
                "taskID": {
                    "type": "integer",
                    "example": 1
                }
            }
        },
        "handler.AddFileToTaskResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "файлы успешно добавлен к вашей задаче"
                },
                "taskID": {
                    "type": "integer",
                    "example": 1
                }
            }
        },
        "handler.CreateTaskRequest": {
            "type": "object",
            "properties": {
                "zipArchiveName": {
                    "type": "string",
                    "example": "test1"
                },
                "zipArchivePath": {
                    "type": "string",
                    "example": "G:/GithubRepo/17.07.2025/internal/util"
                }
            }
        },
        "handler.CreateTaskResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "id вашей задачи: "
                },
                "taskID": {
                    "type": "integer",
                    "example": 1
                }
            }
        },
        "handler.TaskStatusResponse": {
            "type": "object",
            "properties": {
                "archiveLink": {
                    "type": "string",
                    "example": "G:/GithubRepo/17.07.2025/internal/util/task_1.zip"
                },
                "status": {
                    "type": "string",
                    "example": "завершена"
                },
                "taskID": {
                    "type": "integer",
                    "example": 1
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api-tasks",
	Schemes:          []string{},
	Title:            "Junior Go-разработчик",
	Description:      "Тестовое задание на позицию Junior Go-разработчик",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
