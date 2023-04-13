package main

import (
	"bytes"
	"context"
	"io"
	"log"

	"entgo.io/ent/dialect"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/strrl/kubernetes-auditing-dashboard/ent"
	"github.com/strrl/kubernetes-auditing-dashboard/ent/migrate"
	"github.com/strrl/kubernetes-auditing-dashboard/gql"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	auditv1 "k8s.io/apiserver/pkg/apis/audit/v1"
)

func main() {
	entClient, err := ent.Open(dialect.SQLite, "file:data.db?_fk=1")
	if err != nil {
		log.Fatal(err)
	}
	scheme := runtime.NewScheme()
	if err := auditv1.AddToScheme(scheme); err != nil {
		log.Fatal(err)
	}
	defer entClient.Close()
	ctx := context.Background()
	// Run the automatic migration tool to create all schema resources.
	if err := entClient.Schema.Create(ctx, migrate.WithGlobalUniqueID(true)); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	codec := json.NewSerializerWithOptions(json.DefaultMetaFactory, scheme, scheme, json.SerializerOptions{Yaml: false, Pretty: false, Strict: false})
	app := gin.Default()
	apiGroup := app.Group("/api")
	apiGroup.POST("/audit-webhook", func(c *gin.Context) {
		requestBody, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Println(err)
		}
		eventList := auditv1.EventList{}
		_, _, err = codec.Decode(requestBody, nil, &eventList)
		if err != nil {
			log.Println(err)
		}

		var entities []*ent.AuditEventCreate
		for _, event := range eventList.Items {
			buffer := bytes.Buffer{}
			err := codec.Encode(&event, &buffer)
			if err != nil {
				log.Fatal(err)
			}
			raw, err := io.ReadAll(&buffer)
			if err != nil {
				log.Fatal(err)
			}
			item := entClient.AuditEvent.Create().
				SetStage(string(event.Stage)).
				SetAuditID(string(event.AuditID)).
				SetVerb(event.Verb).
				SetUserAgent(event.UserAgent).
				SetLevel(string(event.Level)).
				SetRequestTimestamp(event.RequestReceivedTimestamp.Time).
				SetStageTimestamp(event.StageTimestamp.Time).
				SetRaw(string(raw))

			if event.ObjectRef != nil {
				item.SetNamespace(event.ObjectRef.Namespace).
					SetName(event.ObjectRef.Name).
					SetApiVersion(event.ObjectRef.APIVersion).
					SetApiGroup(event.ObjectRef.APIGroup).
					SetResource(event.ObjectRef.Resource).
					SetSubResource(event.ObjectRef.Subresource)
			}
			entities = append(entities, item)
		}
		_, err = entClient.AuditEvent.CreateBulk(entities...).Save(ctx)
		if err != nil {
			log.Fatal(err)
		}
		println(len(eventList.Items))
		c.Status(200)
	})
	apiGroup.GET("/playground", gin.WrapF(playground.Handler("", "/api/query")))
	graphqlServer := handler.NewDefaultServer(gql.NewExecutableSchema(
		gql.Config{
			Resolvers: gql.NewResolver(entClient),
		}))
	apiGroup.Any("/query", gin.WrapH(graphqlServer))
	app.Run("0.0.0.0:23333")
}
