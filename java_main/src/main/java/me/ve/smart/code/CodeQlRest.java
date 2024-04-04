package me.ve.smart.code;


import com.google.common.collect.Lists;
import me.ve.smart.code.vo.RunTaskVO;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

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
    public ResponseEntity runTask(@RequestParam RunTaskVO runTask) {
        taskRepository.countByCurrentStepIn(Lists.newArrayList("d"));
        return ResponseEntity.ok().build();

    }

}
