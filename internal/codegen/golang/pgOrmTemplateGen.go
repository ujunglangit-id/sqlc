package golang

var (
	templateModel = `
{{define "entityFile"}}
package types //generated by sqlc-pgorm

type {{.StructData.Name}} struct {
	tableName struct{} {{$.Q}}pg:{{.StructData.Name}},alias:_{{.StructData.Name}}{{$.Q}}
	{{- range .Fields}}
	{{- if .Comment}}
	{{comment .Comment}}{{else}}
	{{- end}}
	{{.StructData.Name}} {{.Type}}
	{{- end}}
}
{{end}}

{{define "payloadFile"}}
package types //generated by sqlc-pgorm
		
type {{.StructData.Name}}Payload struct {
	{{- range .Fields}}
	{{- if .Comment}}
	{{comment .Comment}}{{else}}
	{{- end}}
	{{.StructData.Name}} {{.Type}} {{$.Q}}{{.Tag}}{{$.Q}} 
	{{- end}}
}
{{end}}`

	templateRepo = `
{{define "repoInterfaceFile"}}
package repository //generated by sqlc-pgorm

import (
	types "{{.StructData.ProjectPath}}/internal/model/types/{{.StructData.Name}}"
	"context"
)

type {{.StructData.Name}}Repository interface {
	Submit(ctx context.Context, data types.{{.StructData.Name}}Entity) error
	SubmitMultiple(ctx context.Context, data []*types.{{.StructData.Name}}Entity) error
	UpdateByPK(ctx context.Context, data types.{{.StructData.Name}}Entity) error
	DeleteByPK(ctx context.Context, data types.{{.StructData.Name}}Entity) error
	GetList(ctx context.Context, start, limit int) (data []types.{{.StructData.Name}}Entity, count int, err error)
	GetByID(ctx context.Context, id int) (data types.{{.StructData.Name}}Entity, err error)
}
{{end}}

{{define "repoImplFile"}}
package postgre //generated by sqlc-pgorm

import (
	types "{{.StructData.ProjectPath}}/internal/model/types/{{.StructData.Name}}"
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
)

type {{.StructData.Name}}Data struct {
	MasterDB *pg.DB
	SlaveDB  *pg.DB
}

func New{{.StructData.Name}}Wrapper(master, slave *pg.DB) *{{.StructData.Name}}Data {
	return &{{.StructData.Name}}Data{
		MasterDB: master,
		SlaveDB:  slave,
	}
}

func (d *{{.StructData.Name}}Data) Submit(ctx context.Context, data types.{{.StructData.Name}}Entity) (err error) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span, ctx = opentracing.StartSpanFromContext(ctx, "repo.Submit")
		defer span.Finish()
	}
	_, err = d.MasterDB.ModelContext(ctx, &data).Insert()
	if err != nil {
		log.Errorf("[{{.StructData.Name}}RepoImpl][Submit] error Submit : %+v", err)
	}
	return
}

func (d *{{.StructData.Name}}Data) SubmitMultiple(ctx context.Context, data []*types.{{.StructData.Name}}Entity) (err error) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span, ctx = opentracing.StartSpanFromContext(ctx, "repo.SubmitMultiple")
		defer span.Finish()
	}
	_, err = d.MasterDB.ModelContext(ctx, &data).Insert()
	if err != nil {
		log.Errorf("[{{.StructData.Name}}RepoImpl][SubmitMultiple] error SubmitMultiple : %+v", err)
	}
	return
}

func (d *{{.StructData.Name}}Data) GetByID(ctx context.Context, id int) (data types.{{.StructData.Name}}Entity, err error) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span, ctx = opentracing.StartSpanFromContext(ctx, "repo.GetByID")
		defer span.Finish()
	}
	err = d.SlaveDB.ModelContext(ctx, &data).Where("id = ?", 1).Select()
	if err != nil {
		log.Errorf("[{{.StructData.Name}}RepoImpl][GetByID] error GetByID : %+v", err)
	}
	return
}

func (d *{{.StructData.Name}}Data) GetList(ctx context.Context, offset, limit int) (data []types.{{.StructData.Name}}Entity, count int, err error) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span, ctx = opentracing.StartSpanFromContext(ctx, "repo.GetList")
		defer span.Finish()
	}
	count, err = d.SlaveDB.ModelContext(ctx, &data).Offset(offset).Limit(limit).SelectAndCount()
	if err != nil {
		log.Errorf("[{{.StructData.Name}}RepoImpl][GetList] error GetList : %+v", err)
	}
	return
}

func (d *{{.StructData.Name}}Data) UpdateByPK(ctx context.Context, data types.{{.StructData.Name}}Entity) (err error) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span, ctx = opentracing.StartSpanFromContext(ctx, "repo.UpdateByPK")
		defer span.Finish()
	}
	_, err = d.MasterDB.ModelContext(ctx, &data).WherePK().Update()
	if err != nil {
		log.Errorf("[{{.StructData.Name}}RepoImpl][UpdateByPK] error UpdateByPK : %+v", err)
	}
	return
}

func (d *{{.StructData.Name}}Data) DeleteByPK(ctx context.Context, data types.{{.StructData.Name}}Entity) (err error) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span, ctx = opentracing.StartSpanFromContext(ctx, "repo.DeleteByPK")
		defer span.Finish()
	}
	_, err = d.MasterDB.ModelContext(ctx, &data).WherePK().Delete()
	if err != nil {
		log.Errorf("[{{.StructData.Name}}RepoImpl][DeleteByPK] error DeleteByPK : %+v", err)
	}
	return
}

{{end}}`
)