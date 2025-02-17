package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/yesoreyeram/grafana-infinity-datasource/pkg/infinity"
)

func (host *PluginHost) getRouter() *mux.Router {
	router := mux.NewRouter()
	router.Handle("/graphql", host.getGraphQLHandler())                                          // NOT IN USE YET
	router.HandleFunc("/proxy", host.withDatasourceHandlerFunc(GetProxyHandler)).Methods("POST") // DEPRECATED
	router.NotFoundHandler = http.HandlerFunc(host.withDatasourceHandlerFunc(defaultHandler))
	return router
}

func (host *PluginHost) withDatasourceHandlerFunc(getHandler func(d *instanceSettings) http.HandlerFunc) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		client, err := getInstanceFromRequest(host.im, r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		h := getHandler(client)
		h.ServeHTTP(rw, r)
	}
}

func (host *PluginHost) getGraphQLHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		client, err := getInstanceFromRequest(host.im, r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error())) //nolint
			return
		}
		schemaConfig := graphql.SchemaConfig{
			Query: graphql.NewObject(
				graphql.ObjectConfig{
					Name: "RootQuery",
					Fields: graphql.Fields{
						"hello": &graphql.Field{
							Type: graphql.String,
							Args: graphql.FieldConfigArgument{
								"to": &graphql.ArgumentConfig{
									Type:         graphql.String,
									Description:  "Name of the entity to say hello",
									DefaultValue: "World",
								},
							},
							Resolve: func(p graphql.ResolveParams) (interface{}, error) {
								return fmt.Sprintf("Hello %s! 👋", p.Args["to"]), nil
							},
						},
						"query": &graphql.Field{
							Type: graphql.String,
							Args: graphql.FieldConfigArgument{
								"type": &graphql.ArgumentConfig{
									Type:         graphql.String,
									Description:  "query type. Can be one of json, csv, xml",
									DefaultValue: "json",
								},
								"url": &graphql.ArgumentConfig{
									Type:         graphql.String,
									Description:  "URL to be queried",
									DefaultValue: "https://jsonplaceholder.typicode.com/users",
								},
							},
							Resolve: func(p graphql.ResolveParams) (interface{}, error) {
								res, _, _, err := client.client.GetResults(infinity.Query{
									Type: p.Args["type"].(string),
									URL:  p.Args["url"].(string),
								}, map[string]string{})
								return res, err
							},
						},
						"time": &graphql.Field{
							Type: graphql.String,
							Args: graphql.FieldConfigArgument{
								"format": &graphql.ArgumentConfig{
									Type:         graphql.String,
									Description:  "Format for the time in golang time layout. https://yourbasic.org/golang/format-parse-string-time-date-example/",
									DefaultValue: "2006-01-02T15:04:05Z07:00",
								},
							},
							Resolve: func(p graphql.ResolveParams) (interface{}, error) {
								return time.Now().Format(p.Args["format"].(string)), nil
							},
						},
					},
				},
			),
		}
		schema, _ := graphql.NewSchema(schemaConfig)
		h := handler.New(&handler.Config{
			Schema:   &schema,
			Pretty:   true,
			GraphiQL: true,
		})
		h.ServeHTTP(w, r)
	})
}

func GetProxyHandler(client *instanceSettings) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var query infinity.Query
			err := json.NewDecoder(r.Body).Decode(&query)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}
			if query.Source == "url" {
				response, _, _, err := client.client.GetResults(query, map[string]string{})
				if err != nil {
					http.Error(rw, err.Error(), http.StatusInternalServerError)
					return
				}
				fmt.Fprintf(rw, "%s", response)
				return
			}
			http.Error(rw, "unknown query", http.StatusNotImplemented)
			return
		}
		http.Error(rw, "500 - Something bad happened! Invalid query.", http.StatusInternalServerError)
	}
}

func defaultHandler(client *instanceSettings) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		http.Error(rw, "not a known resource call", http.StatusInternalServerError)
	}
}
