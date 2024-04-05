package me.ve.smart.code;

import lombok.extern.slf4j.Slf4j;
import org.springframework.boot.CommandLineRunner;
import org.springframework.stereotype.Service;

import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

@Service
@Slf4j
public class TaskService implements CommandLineRunner {

    final TaskRepository taskRepository;

    final DockerService dockerService;

    public TaskService(TaskRepository taskRepository, DockerService dockerService) {
        this.taskRepository = taskRepository;
        this.dockerService = dockerService;
    }

    @Override
    public void run(String... args) throws Exception {
        ExecutorService executorService = Executors.newSingleThreadExecutor();
        executorService.execute(() -> {
            TaskEntity taskEntity = taskRepository.findNewTask();
            if (taskEntity != null) {
                taskEntity.setCurrentStep(TaskStatusEnum.Doing);
                taskRepository.save(taskEntity);
                try {
                    dockerService.analyse(taskEntity);
                    taskEntity.setCurrentStep(TaskStatusEnum.Done);
                    taskRepository.save(taskEntity);
                } catch (Exception e) {
                    taskEntity.setCurrentStep(TaskStatusEnum.Failed);
                    taskRepository.save(taskEntity);
                }
            }
            try {
                Thread.sleep(5000);
            } catch (InterruptedException e) {
                e.printStackTrace();
            }
        });
    }
}
