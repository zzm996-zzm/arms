// Package http API.
// @title arms
// @version 1.1
// @description arms测试
// @termsOfService https://github.com/swaggo/swag

// @contact.name zhangzhiming
// @contact.email zzmhaoshuai@icloud.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /
// @query.collection.format multi

// @securityDefinitions.basic BasicAuth

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @x-extension-openapi {"example": "value on a json format"}
package http

import (
	_ "github.com/zzm996-zzm/arms/app/http/swagger"
)
