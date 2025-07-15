package models

import (
	"encoding/json"
	"time"

	"github.com/uptrace/bun"
)

// Collection represents a Postman collection
type Collection struct {
	bun.BaseModel `bun:"table:collections,alias:c"`

	ID          int64     `bun:"id,pk,autoincrement" json:"id"`
	Name        string    `bun:"name,notnull" json:"name"`
	Description string    `bun:"description" json:"description"`
	Schema      string    `bun:"schema" json:"schema"`
	Variables   JSONMap   `bun:"variables,type:jsonb" json:"variables"`
	Auth        JSONMap   `bun:"auth,type:jsonb" json:"auth,omitempty"`
	Events      JSONMap   `bun:"events,type:jsonb" json:"events,omitempty"`
	Items       JSONMap   `bun:"items,type:jsonb" json:"items,omitempty"`
	PostmanID   string    `bun:"postman_id" json:"_postman_id,omitempty"`
	ExporterID  string    `bun:"exporter_id" json:"_exporter_id,omitempty"`
	CreatedAt   time.Time `bun:"created_at,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt   time.Time `bun:"updated_at,notnull,default:current_timestamp" json:"updated_at"`

	Requests []*Request `bun:"rel:has-many,join:id=collection_id" json:"requests,omitempty"`
}

// Request represents an API request within a collection
type Request struct {
	bun.BaseModel `bun:"table:requests,alias:r"`

	ID           int64             `bun:"id,pk,autoincrement" json:"id"`
	CollectionID int64             `bun:"collection_id,notnull" json:"collection_id"`
	Name         string            `bun:"name,notnull" json:"name"`
	Description  string            `bun:"description" json:"description"`
	FolderPath   string            `bun:"folder_path" json:"folder_path,omitempty"`
	URL          JSONMap           `bun:"url,type:jsonb" json:"url"`
	Method       string            `bun:"method,notnull" json:"method"`
	Headers      map[string]string `bun:"headers,type:jsonb" json:"headers,omitempty"`
	Params       JSONMap           `bun:"params,type:jsonb" json:"params,omitempty"`
	Body         JSONMap           `bun:"body,type:jsonb" json:"body,omitempty"`
	Auth         JSONMap           `bun:"auth,type:jsonb" json:"auth,omitempty"`
	Events       JSONMap           `bun:"events,type:jsonb" json:"events,omitempty"`
	Responses    JSONMap           `bun:"responses,type:jsonb" json:"responses,omitempty"`
	PostmanID    string            `bun:"postman_id" json:"_postman_id,omitempty"`
	CreatedAt    time.Time         `bun:"created_at,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt    time.Time         `bun:"updated_at,notnull,default:current_timestamp" json:"updated_at"`

	Collection *Collection `bun:"rel:belongs-to,join:collection_id=id" json:"collection,omitempty"`
}

// OpenAPISpec represents an OpenAPI specification
type OpenAPISpec struct {
	bun.BaseModel `bun:"table:openapi_specs,alias:o"`

	ID          int64     `bun:"id,pk,autoincrement" json:"id"`
	Title       string    `bun:"title,notnull" json:"title"`
	Description string    `bun:"description" json:"description"`
	Version     string    `bun:"version,notnull" json:"version"`
	Content     JSONMap   `bun:"content,type:jsonb" json:"content"`
	CreatedAt   time.Time `bun:"created_at,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt   time.Time `bun:"updated_at,notnull,default:current_timestamp" json:"updated_at"`
}

// JSONMap is a helper type for JSON columns
type JSONMap map[string]any

// Scan implements the sql.Scanner interface.
func (j *JSONMap) Scan(src any) error {
	if src == nil {
		*j = make(JSONMap)
		return nil
	}

	var sourceBytes []byte
	switch v := src.(type) {
	case string:
		sourceBytes = []byte(v)
	case []byte:
		sourceBytes = v
	default:
		*j = make(JSONMap)
		return nil
	}

	err := json.Unmarshal(sourceBytes, j)
	if err != nil {
		*j = make(JSONMap)
		return nil
	}

	return nil
}

// PostmanCollection represents the full structure of a Postman collection
type PostmanCollection struct {
	Info     CollectionInfo  `json:"info"`
	Item     []PostmanItem   `json:"item"`
	Variable []KeyValuePair  `json:"variable,omitempty"`
	Auth     json.RawMessage `json:"auth,omitempty"`
	Event    []PostmanEvent  `json:"event,omitempty"`
	Schema   string          `json:"schema,omitempty"`
}

// CollectionInfo holds collection metadata
type CollectionInfo struct {
	Name           string `json:"name"`
	Description    string `json:"description,omitempty"`
	Schema         string `json:"schema,omitempty"`
	PostmanID      string `json:"_postman_id,omitempty"`
	ExporterID     string `json:"_exporter_id,omitempty"`
	CollectionLink string `json:"_collection_link,omitempty"`
}

// PostmanItem represents a folder or request in a Postman collection
type PostmanItem struct {
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Item        []PostmanItem     `json:"item,omitempty"`
	Request     *PostmanRequest   `json:"request,omitempty"`
	Response    []PostmanResponse `json:"response,omitempty"`
	Event       []PostmanEvent    `json:"event,omitempty"`
	Variable    []KeyValuePair    `json:"variable,omitempty"`
	Auth        json.RawMessage   `json:"auth,omitempty"`
	PostmanID   string            `json:"id,omitempty"`
}

// PostmanRequest represents a request in a Postman collection
type PostmanRequest struct {
	URL         any             `json:"url"`
	Method      string          `json:"method"`
	Header      []KeyValuePair  `json:"header,omitempty"`
	Body        PostmanBody     `json:"body,omitzero"`
	Description string          `json:"description,omitempty"`
	Auth        json.RawMessage `json:"auth,omitempty"`
}

// PostmanBody represents the body of a Postman request
type PostmanBody struct {
	Mode       string          `json:"mode,omitempty"`
	Raw        string          `json:"raw,omitempty"`
	URLEncoded []KeyValuePair  `json:"urlencoded,omitempty"`
	FormData   []KeyValuePair  `json:"formdata,omitempty"`
	GraphQL    json.RawMessage `json:"graphql,omitempty"`
	File       json.RawMessage `json:"file,omitempty"`
	Options    json.RawMessage `json:"options,omitempty"`
}

// KeyValuePair represents key-value pairs like headers, params, etc.
type KeyValuePair struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Description any    `json:"description,omitempty"`
	Type        string `json:"type,omitempty"`
	Disabled    bool   `json:"disabled,omitempty"`
	Name        string `json:"name,omitempty"`
}

// PostmanResponse represents an response in a Postman collection
type PostmanResponse struct {
	Name        string            `json:"name"`
	OriginalReq json.RawMessage   `json:"originalRequest,omitempty"`
	Status      string            `json:"status,omitempty"`
	Code        int               `json:"code,omitempty"`
	Header      []KeyValuePair    `json:"header,omitempty"`
	Body        string            `json:"body,omitempty"`
	Cookie      []json.RawMessage `json:"cookie,omitempty"`
	PreviewType string            `json:"_postman_previewlanguage,omitempty"`
	PostmanID   string            `json:"id,omitempty"`
}

// PostmanEvent represents event scripts in Postman
type PostmanEvent struct {
	Listen   string        `json:"listen"`
	Script   PostmanScript `json:"script"`
	Disabled bool          `json:"disabled,omitempty"`
}

// PostmanScript represents a script in a Postman event
type PostmanScript struct {
	Type string   `json:"type"`
	Exec []string `json:"exec"`
	ID   string   `json:"id,omitempty"`
	Src  string   `json:"src,omitempty"`
}

// URLObject represents a URL in Postman
type URLObject struct {
	Raw      string         `json:"raw,omitempty"`
	Protocol string         `json:"protocol,omitempty"`
	Host     []string       `json:"host,omitempty"`
	Path     []string       `json:"path,omitempty"`
	Port     string         `json:"port,omitempty"`
	Query    []KeyValuePair `json:"query,omitempty"`
	Hash     string         `json:"hash,omitempty"`
	Variable []KeyValuePair `json:"variable,omitempty"`
}
