package host 

import (
	"fmt"
	"encoding/json"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

	"ops-backend/pkg/aliyun"
	"ops-backend/pkg/models"
	"ops-backend/view/asset/region"
)

type EcsInstanceResponse struct {
	RequestId  string  							 `json:"RequestId"`
	TotalCount int     							 `json:"TotalCount"`
	PageNumber int     							 `json:"PageNumber"`
	PageSize   int     							 `json:"PageSize"`
	Instances  ecs.InstancesInDescribeInstances  `json:"Instances"`
}

func SyncHost () {
	regions := &[]region.Region{}
	if err := models.FindAll(regions); err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, r := range *regions {
		SyncRegionHost(r.Name)
	}
}

func SyncRegionHost(regionName string) {
	client, err := aliyun.GetEcsClient(regionName)
	request := ecs.CreateDescribeInstancesRequest()
	request.Scheme = "https"
	request.RegionId = regionName 

	response, err := client.DescribeInstances(request)
	if err != nil {
		fmt.Println(err.Error())
		return 
	}
	baseResponse := response.BaseResponse
	if !baseResponse.IsSuccess() {
		fmt.Println("request faild")
		return 
	}

	content := baseResponse.GetHttpContentBytes()
	var ecsInstanceResponse EcsInstanceResponse 
	err = json.Unmarshal(content, &ecsInstanceResponse)
	if err != nil {
		return
	}
	instances := ecsInstanceResponse.Instances.Instance
	for _, instance := range instances {
		var priIp, pubIp string
		pubIpList := instance.PublicIpAddress.IpAddress
		if len(pubIpList) > 0 {
			pubIp = pubIpList[0]	
		}
		priIpList := instance.NetworkInterfaces.NetworkInterface 
		if len(priIpList) > 0 {
			priIp = priIpList[0].PrimaryIpAddress	
		}
		host := Host{
			InstanceId: instance.InstanceId,
			RegionId: instance.RegionId,
			HostName: instance.HostName,
			Memory: instance.Memory,
			Cpu: instance.Cpu,
			OSName: instance.OSName,
			PublicIpAddress: pubIp,
			PrimaryIpAddress: priIp,
			Status: instance.Status,
			ExpiredTime: instance.ExpiredTime,
		}
		host.CreateOrUpdate()
	} 
}