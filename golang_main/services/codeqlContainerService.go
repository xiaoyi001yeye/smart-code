// codeqlContainerService.go

package services

import (
	"context"
	"log"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// CodeQLContainerService 结构体
type CodeQLContainerService struct {
	Cli    *client.Client
}

// NewCodeQLContainerService 创建一个新的 CodeQLContainerService 实例
func NewCodeQLContainerService() (*CodeQLContainerService,err) {
	var service *CodeQLContainerService
	var err error
	
	Cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersion("1.40"))
	if err != nil {
		return nil, err
	}
	service = &CodeQLContainerService{
		Cli:    Cli
	}
	return service,nil
}

// RunScan 启动 CodeQL 扫描器容器
// func (service *CodeQLContainerService) RunScan(ctx context.Context, targetImage string) error {
// 	containerName := "codeql-container"

// 	// 配置容器创建参数
// 	resp, err := service.Cli.ContainerCreate(
// 		ctx,
// 		&types.ContainerCreateConfig{
// 			Name: containerName,
// 			Image: service.Image,
// 			Cmd:   []string{"codeql", "database", "create", "--languages=java", "--upload", "my-codeql-database"},
// 			Volumes: map[string]struct{}{
// 				service.Volume: {},
// 			},
// 		},
// 		types.ContainerHostConfig{
// 			Binds: []types.MountPoint{
// 				{Source: service.Volume, Target: "/codeql-run"},
// 			},
// 		},
// 		nil, // 网络配置
// 		nil, // 平台配置
// 	)
// 	if err != nil {
// 		return err
// 	}

// 	// 启动容器
// 	err = service.Cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
// 	if err != nil {
// 		return err
// 	}
// 	defer service.Cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{Force: true})

// 	// 等待容器完成扫描
// 	watcher, err := service.Cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true, Follow: true})
// 	if err != nil {
// 		return err
// 	}
// 	defer watcher.Close()

// 	// 这里可以添加逻辑来处理日志输出或检查扫描状态
// 	<-ctx.Done() // 假设上下文被取消时，我们希望停止等待

// 	return nil
// }


func (service *CodeQLContainerService) GetContainerStatus() (types.ContainerState, error) {
	containerInfo, err := service.Cli.ContainerInspect(context.Background(),"codeql-container")
	if err != nil {
		return types.ContainerState{}, err
	}
	return containerInfo.State, nil
}