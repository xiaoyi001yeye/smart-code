package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/google/uuid"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

type TaskDescription struct {
	InputPath  string `json:"inputPath"`
	OutputPath string `json:"outputPath"`
	TaskId     string `json:"TaskId"`
	Language   string `json:"language"`
	Qlpack     string `json:"qlpack"`
}

var CodeQLImageId = "mcr.microsoft.com/cstsectools/codeql-container:latest"
var containerName = "codeql-container"

func main() {
	http.HandleFunc("/run", runTask)
	http.HandleFunc("/", welcome)
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	dbConfig := DefaultDbConfig()
	db, err := sql.Open("postgres", DbConnectionString(dbConfig))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	go func() {
		for {
			rows, err := db.Query("SELECT task_id,input_path,output_path,languge,qlpack FROM tasks WHERE current_step in ('New')  ORDER BY created_at ASC LIMIT 1")
			if err != nil {
				log.Println("查询新数据出错:", err)
				time.Sleep(5 * time.Second)
				continue
			}
			defer rows.Close()
			for rows.Next() {
				var taskID string
				var inputPath string
				var outputPath string
				var codeLanguage string
				var qlpack string
				// 读取数据并进行处理
				err := rows.Scan(&taskID, &inputPath, &outputPath, &codeLanguage, &qlpack)
				if err != nil {
					log.Println("扫描数据出错:", err)
					continue
				}

				// 处理数据的逻辑
				fmt.Printf("处理新数据:taskID=%s, inputPath=%s, outputPath=%s,language=%s,qlpack=%s\n", taskID, inputPath, outputPath, codeLanguage, qlpack)

				// 标记数据已处理
				_, err = db.Exec("UPDATE tasks SET current_step = 'Doing' WHERE task_id = $1", taskID)
				if err != nil {
					log.Println("更新数据状态出错:", err)
				}

				cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
				if err != nil {
					log.Printf("Error connecting to Docker: %s", err.Error())
					return
				}

				createDatabaseEnv := []string{
					fmt.Sprintf("CODEQL_CLI_ARGS=database create --language=%s /opt/results/source_db -s /opt/src", codeLanguage),
				}

				upgradeDatabaseEnv := []string{
					fmt.Sprintf("CODEQL_CLI_ARGS=database upgrade /opt/results/source_db"),
				}

				analyzeEnv := []string{
					fmt.Sprintf("CODEQL_CLI_ARGS=database analyze /opt/results/source_db --format=sarifv2 --output=/opt/results/issues.sarif %s-%s.qls", codeLanguage, qlpack),
				}

				err1 := createAndStartContainer(cli, createDatabaseEnv, inputPath, outputPath, taskID)
				if err1 != nil {
					updateTaskStatus(db, taskID, TaskStatusFailed)
					log.Printf("An error occurred: %s", err)
					continue
				}
				err2 := createAndStartContainer(cli, upgradeDatabaseEnv, inputPath, outputPath, taskID)
				if err2 != nil {
					updateTaskStatus(db, taskID, TaskStatusFailed)
					log.Printf("An error occurred: %s", err)
					continue
				}
				err3 := createAndStartContainer(cli, analyzeEnv, inputPath, outputPath, taskID)
				if err3 != nil {
					updateTaskStatus(db, taskID, TaskStatusFailed)
					log.Printf("An error occurred: %s", err)
					continue
				}
				updateTaskStatus(db, taskID, TaskStatusDone)
				log.Printf("task %s finished.", taskID)

			}
		}
	}()

	fmt.Printf("Starting server on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("Error starting server: %s", err)
	}

}

func updateTaskStatus(db *sql.DB, taskID string, status TaskStatus) {
	_, err := db.Exec("UPDATE tasks SET status = ? WHERE task_id = ?", status, taskID)
	if err != nil {
		log.Printf("更新任务状态出错: %s", err)
	}
}

func createAndStartContainer(cli *client.Client, commandEnv []string, inputPath, outputPath, taskID string) error {
	containerConfig := &container.Config{
		Image: CodeQLImageId,
		Env:   commandEnv,
	}

	mounts := []mount.Mount{
		{
			Type:   mount.TypeBind,
			Source: inputPath,
			Target: "/opt/src",
		},
		{
			Type:   mount.TypeBind,
			Source: outputPath,
			Target: "/opt/results",
		},
	}
	hostConfig := &container.HostConfig{
		// 这里可以设置需要的主机配置,例如端口映射、卷挂载等
		Mounts: mounts,
	}
	// 定义网络配置
	networkingConfig := &network.NetworkingConfig{
		// 网络配置,可以指定容器应该连接的网络等
	}

	// 定义平台信息,如果不需要特定平台,可以传 nil
	//platform := &v1.Platform{
	//	OS:           "linux",
	//	Architecture: "amd64",
	//}
	platform := &v1.Platform{}

	createResponse, err := cli.ContainerCreate(context.Background(), containerConfig, hostConfig, networkingConfig, platform, containerName)
	if err != nil {
		return fmt.Errorf("error creating container for task %s: %w", taskID, err)
	}
	containerId := createResponse.ID

	err = cli.ContainerStart(context.Background(), containerId, container.StartOptions{})
	if err != nil {
		cli.ContainerRemove(context.Background(), containerId, container.RemoveOptions{}) // 尝试清理容器
		return fmt.Errorf("error starting container for task %s: %w", taskID, err)
	}

	waitResponseCh, errCh := cli.ContainerWait(context.Background(), containerId, container.WaitConditionNotRunning)
	select {
	case status := <-waitResponseCh:
		if status.StatusCode != 0 {
			return fmt.Errorf("container for task %s exited with non-zero status code: %d", taskID, status.StatusCode)
		}
	case err := <-errCh:
		return fmt.Errorf("container for task %s encountered an error: %w", taskID, err)
	}

	err = cli.ContainerRemove(context.Background(), containerId, container.RemoveOptions{})
	if err != nil {
		return fmt.Errorf("error removing container for task %s: %w", taskID, err)
	}

	return nil

}

func runTask(w http.ResponseWriter, r *http.Request) {
	inputPath := r.FormValue("inputPath")
	outputPath := r.FormValue("outputPath")
	language := r.FormValue("language")
	qlpack := r.FormValue("qlpack")
	if inputPath == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("参数inputPath未传递"))
		return
	}
	if outputPath == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("参数outputPath未传递"))
		return
	}
	if language == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("参数language未传递"))
		return
	}
	if qlpack == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("参数qlpack未传递"))
		return
	}

	dbConfig := DefaultDbConfig()
	db, err := sql.Open("postgres", DbConnectionString(dbConfig))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("数据库连接失败"))
		return
	}
	defer db.Close()

	// 查询是否有未完成的任务
	var hasUncompletedTask bool
	err = db.QueryRow("SELECT COUNT(*) FROM codeql_tasks WHERE current_step not in('Done','Failed')").Scan(&hasUncompletedTask)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("查询未完成任务失败"))
		return
	}

	if hasUncompletedTask {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("有未完成的任务，无法开始新任务"))
		return
	}

	taskID := uuid.New().String()

	stmt, err := db.Prepare("INSERT INTO tasks (task_id, input_path, output_path, code_language, qlpack, task_type, current_step) VALUES ($1, $2, $3, $4, $5, $6, $7)")

	_, err = stmt.Exec(taskID,
		inputPath,
		outputPath,
		language,
		qlpack,
		TaskTypeAnalyse,
		TaskStatusNew)

	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(TaskDescription{TaskId: taskID})

}

func welcome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the codeql Task Runner!\n")
}
