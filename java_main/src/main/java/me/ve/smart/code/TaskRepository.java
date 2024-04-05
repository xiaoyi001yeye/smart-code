package me.ve.smart.code;


import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;

import java.util.List;

public interface TaskRepository extends JpaRepository<TaskEntity, String> {


    Integer countByCurrentStepIn(List<TaskStatusEnum> currentSteps);

    @Query(value = "select * from tasks where current_step='New' order by created_at asc limit 1", nativeQuery = true)
    TaskEntity findNewTask();

}
