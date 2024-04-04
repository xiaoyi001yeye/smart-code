package me.ve.smart.code;


import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;

public interface TaskRepository extends JpaRepository<TaskEntity, String> {


    Integer countByCurrentStepIn(List<String> currentSteps);

}
