package aliyun

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

	"ops-backend/pkg/setting"
) 

func GetEcsClient (regionId string) (*ecs.Client, error) {
	return ecs.NewClientWithAccessKey(regionId, setting.AliyunSetting.AccessKey, setting.AliyunSetting.AccessSecret)
}
