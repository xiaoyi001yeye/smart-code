package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
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
	_ "github.com/lib/pq"
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
var outputPath = "./results"

func main() {
	loc := time.FixedZone("CST", 8*60*60)
	time.Local = loc
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)
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
			log.Println("task scan.")
			rows, err := db.Query("SELECT task_id,input_path,output_path,code_language,qlpack FROM public.tasks WHERE current_step in ('New')  ORDER BY created_at ASC LIMIT 1")
			if err != nil {
				log.Fatal("查询新数据出错:", err)
			}
			defer rows.Close()
			rowCount := 0
			for rows.Next() {
				rowCount++
				var taskID, inputPath, outputPath, codeLanguage, qlpack string
				err := rows.Scan(&taskID, &inputPath, &outputPath, &codeLanguage, &qlpack)
				if err != nil {
					log.Println("扫描数据出错:", err)
					continue
				}
				log.Printf("Process task:taskID=%s, inputPath=%s, outputPath=%s,language=%s,qlpack=%s\n", taskID, inputPath, outputPath, codeLanguage, qlpack)
				updateTaskStatus(db, taskID, TaskStatusDoing)

				cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
				if err != nil {
					updateTaskStatus(db, taskID, TaskStatusFailed)
					log.Fatalf("Error connecting to Docker: %+v\n", err)
				} else {
					log.Println("NewClientWithOpts OK.")
				}

				createDatabaseEnv := []string{
					"CODEQL_CLI_ARGS=database create --language=" + codeLanguage + " /opt/results/source_db -s /opt/src",
				}

				upgradeDatabaseEnv := []string{
					"CODEQL_CLI_ARGS=database upgrade /opt/results/source_db",
				}

				analyzeEnv := []string{
					"CODEQL_CLI_ARGS=database analyze /opt/results/source_db --format=sarifv2 --output=/opt/results/issues.sarif " + codeLanguage + "-" + qlpack + ".qls",
				}

				err1 := createAndStartContainer(cli, createDatabaseEnv, inputPath, outputPath, taskID)
				if err1 != nil {
					updateTaskStatus(db, taskID, TaskStatusFailed)
					log.Fatalf("An error occurred: %+v\n", err1)
				} else {
					log.Println("create Database OK.")
				}
				err2 := createAndStartContainer(cli, upgradeDatabaseEnv, inputPath, outputPath, taskID)
				if err2 != nil {
					updateTaskStatus(db, taskID, TaskStatusFailed)
					log.Fatalf("An error occurred: %+v\n", err2)
				} else {
					log.Println("upgrade Database OK.")
				}
				err3 := createAndStartContainer(cli, analyzeEnv, inputPath, outputPath, taskID)
				if err3 != nil {
					updateTaskStatus(db, taskID, TaskStatusFailed)
					log.Fatalf("An error occurred: %+v\n", err3)
				} else {
					log.Println("analyze Database OK.")
				}
				updateTaskStatus(db, taskID, TaskStatusDone)
				log.Printf("task %s finished.", taskID)

			}

			if err := rows.Err(); err != nil {
				log.Fatalf("An error occurred: %+v\n", err)
			}
			if rowCount == 0 {
				log.Println("No executable task found.")
			}
			time.Sleep(10 * time.Second)
		}
	}()

	log.Printf("Starting server on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Printf("Error starting server: %s", err)
	}

}

func updateTaskStatus(db *sql.DB, taskID string, status TaskStatus) {
	_, err := db.Exec("UPDATE tasks SET current_step = $1 WHERE task_id = $2", status, taskID)
	if err != nil {
		log.Printf("更新任务状态出错: %s", err)
	}
}

func createAndStartContainer(cli *client.Client, commandEnv []string, inputPath, outputPath, taskID string) error {

	curr_container, err := cli.ContainerInspect(context.Background(), containerName)
	if err != nil {
		if errors.Is(err, types.ContainerNotFound{}) {
			fmt.Printf("Container %s not found.\n", containerName)
		} else {
			return fmt.Errorf("error inspecting container: %w", err)
		}
	}
	if curr_container.State.Running {
		log.Println("container is running.")
	} else {

		if err := cli.ContainerRemove(context.Background(), curr_container.ID, container.RemoveOptions{}); err != nil {
			return fmt.Errorf("error removing container:%s", err)
		}

		log.Printf("container %s has been deleted.", curr_container.ID)
	}
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
		return fmt.Errorf("error creating container for task %s: %s", taskID, err)
	}
	containerId := createResponse.ID
	log.Printf("create new Container ID: %s", containerId)
	err = cli.ContainerStart(context.Background(), containerId, container.StartOptions{})
	if err != nil {
		cli.ContainerRemove(context.Background(), containerId, container.RemoveOptions{}) // 尝试清理容器
		return fmt.Errorf("error starting container for task %s: %w", taskID, err)
	}

	waitResponseCh, errCh := cli.ContainerWait(context.Background(), containerId, container.WaitConditionNotRunning)
	select {
	case result := <-waitResponseCh:
		if result.StatusCode != 0 {
			var output string
			jsonOutput, _ := cli.ContainerInspect(context.Background(), containerId)
			if jsonOutput.ContainerJSONBase != nil && jsonOutput.ContainerJSONBase.State != nil {
				state := jsonOutput.ContainerJSONBase.State
				output = fmt.Sprintf("Status: %s, Error: %s, OOMKilled: %v, Dead: %v", state.Status, state.Error, state.OOMKilled, state.Dead)
			}
			return fmt.Errorf("Container error.id:%s, code:%s,output:%s", containerId, int(result.StatusCode), output)
		}
	case err := <-errCh:
		return fmt.Errorf("container for task %s encountered an error: %w", taskID, err)
	}
	log.Println("WaitConditionNotRunning OK.")

	err = cli.ContainerRemove(context.Background(), containerId, container.RemoveOptions{})
	if err != nil {
		return fmt.Errorf("error removing container for task %s: %w", taskID, err)
	}

	return nil

}

func runTask(w http.ResponseWriter, r *http.Request) {
	inputPath := r.FormValue("inputPath")
	// outputPath := r.FormValue("outputPath")
	language := r.FormValue("language")
	qlpack := r.FormValue("qlpack")
	if inputPath == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("参数inputPath未传递"))
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
	err = db.QueryRow("SELECT COUNT(*) FROM tasks WHERE current_step not in('Done','Failed')").Scan(&hasUncompletedTask)
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
	if err != nil {
		log.Fatal(err)
	}
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
