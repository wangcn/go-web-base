package user

import (
	"context"

	"go.uber.org/zap"

	"mybase/model"
	"mybase/util"
)

type ProfileInfo struct {
	Id   int64  `json:"id"`
	Name string `json:"market"`
}

type ProfileInfos []*ProfileInfo

type Profile struct{}

func (s *Profile) GetInfoById(ctx context.Context, id int64) (*ProfileInfo, error) {
	logger := util.LogWithContext(ctx)
	personModel := model.NewPersonModel()

	// 查结果
	person, err := personModel.SelectOneById(id)
	if err != nil {
		logger.With(zap.Error(err)).Error("get person info failed")
		return nil, err
	}

	// 拼装结果
	result := &ProfileInfo{
		Id:   person.Id,
		Name: person.Name,
	}

	return result, nil
}
