package me.ve.smart.code;


import com.google.common.collect.Lists;
import jakarta.validation.Valid;
import me.ve.smart.code.vo.RunTaskVO;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.UUID;

import static me.ve.smart.code.TaskStatusEnum.*;

@RestController
@RequestMapping("/codeql")
public class CodeQlRest {

    @Autowired
    TaskRepository taskRepository;

    /**
     * This is the method that handles HTTP GET requests. The return type send back a text
     */


    @GetMapping("/welcome")
    public ResponseEntity welcome() {

        return ResponseEntity.ok()
                .body("welcome to codeql java.")
                ;

    }


    @PostMapping("/runTask")
    public ResponseEntity<String> runTask(@Valid @RequestBody RunTaskVO runTask) {
        Integer count = taskRepository.countByCurrentStepIn(Lists.newArrayList(Done, Failed));
        if (count > 0) {
            ResponseEntity.internalServerError().body("有未完成的任务，无法开始新任务");
        }
        String newTaskId = UUID.randomUUID().toString();
        taskRepository.save(TaskEntity.builder()
                .taskId(newTaskId)
                .taskType(TaskTypeEnum.Analyse)
                .currentStep(New)
                .inputPath(runTask.getInputPath())
                .language(runTask.getLanguage())
                .qlpack(runTask.getQlpack())
                .build());

        return ResponseEntity.ok().body(newTaskId);

    }

}
