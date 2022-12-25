// configuration.go
package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"strings"
)

const MysqlConfig = "okc:okcpassword@(localhost:3306)/xen_stats?charset=utf8mb4&parseTime=True&loc=Local"

func main() {
	// connect db
	db, err := gorm.Open(mysql.Open(MysqlConfig))
	if err != nil {
		panic(fmt.Errorf("cannot establish db connection: %w", err))
	}

	// gen instance
	g := gen.NewGenerator(gen.Config{
		OutPath: "./query",

		Mode: gen.WithDefaultQuery | gen.WithQueryInterface,

		FieldNullable: true, // generate pointer when field is nullable

		FieldCoverable: false, // generate pointer when field has default value, to fix problem zero value cannot be assign: https://gorm.io/docs/create.html#Default-Values

		FieldSignable:     false, // detect integer field's unsigned type, adjust generated data type
		FieldWithIndexTag: false, // generate with gorm index tag
		FieldWithTypeTag:  true,  // generate with gorm column type tag
	})

	g.UseDB(db)

	dataMap := map[string]func(detailType string) (dataType string){
		"tinyint":   func(detailType string) (dataType string) { return "int64" },
		"smallint":  func(detailType string) (dataType string) { return "int64" },
		"mediumint": func(detailType string) (dataType string) { return "int64" },
		"bigint":    func(detailType string) (dataType string) { return "int64" },
		"int":       func(detailType string) (dataType string) { return "int64" },
	}
	g.WithDataTypeMap(dataMap)

	jsonField := gen.FieldJSONTagWithNS(func(columnName string) (tagContent string) {
		toStringField := `balance, `
		if strings.Contains(toStringField, columnName) {
			return columnName + ",string"
		}
		return columnName
	})
	autoUpdateTimeField := gen.FieldGORMTag("update_time", "column:update_time;type:int unsigned;autoUpdateTime")
	autoCreateTimeField := gen.FieldGORMTag("create_time", "column:create_time;type:int unsigned;autoCreateTime")
	softDeleteField := gen.FieldType("delete_time", "soft_delete.DeletedAt")
	fieldOpts := []gen.ModelOpt{jsonField, autoCreateTimeField, autoUpdateTimeField, softDeleteField}

	User := g.GenerateModel("user")
	allModel := g.GenerateAllTable(fieldOpts...)

	//Score := g.GenerateModel("score",
	//	append(
	//		fieldOpts,
	//		// user 一对多 address 关联, 外键`uid`在 address 表中
	//		gen.FieldRelate(field.HasMany, "user", User, &field.RelateConfig{GORMTag: "foreignKey:UID"}),
	//	)...,
	//)

	g.ApplyBasic(User)
	g.ApplyBasic(allModel...)

	g.Execute()
}
