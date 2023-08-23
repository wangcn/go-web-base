package model

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"mybase/util"
)

type Person struct {
	Id   int64  `gorm:"column:id"`
	Name string `gorm:"column:name"`
}

// --------------------
// 查询参数
// --------------------
type PersonParams struct {
	Id   int64
	Name string
}

func (p *PersonParams) build(sql *gorm.DB) *gorm.DB {
	if p.Id != 0 {
		sql = sql.Where("id = ?", p.Id)
	}
	if p.Name != "" {
		sql = sql.Where("name = ?", p.Name)
	}
	return sql
}

// PersonModel 游戏模型
// 可选择性兼容事务、trace、slave
type PersonModel struct {
	ctx context.Context
	db  *gorm.DB
}

// 实例化模型
func NewPersonModel() *PersonModel {
	model := &PersonModel{
		db: util.DB(util.DBMaster),
	}
	return model
}

// WithContext 用于 trace
func (m *PersonModel) WithContext(ctx context.Context) *PersonModel {
	m.ctx = ctx
	return m
}

// UseSlave 使用 slave
func (m *PersonModel) UseSlave() *PersonModel {
	clone := *m
	dbSlave := util.DB(util.DBSlave)
	if dbSlave != nil {
		clone.db = dbSlave
	}
	return &clone
}

func (m *PersonModel) TableName() string {
	return "person"
}
func (m *PersonModel) SelectOneById(id int64) (*Person, error) {
	tableName := m.TableName()
	sql := m.db.Table(tableName)
	var item Person
	result := sql.Where("id = ?", id).Take(&item)
	if result.Error != nil {
		return nil, result.Error
	}
	return &item, nil
}

func (m *PersonModel) InsertIgnore(dimName string, person *Person) (int64, error) {
	tableName := m.TableName()
	result := m.db.Table(tableName).Clauses(clause.Insert{Modifier: "IGNORE"}).Create(person)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
