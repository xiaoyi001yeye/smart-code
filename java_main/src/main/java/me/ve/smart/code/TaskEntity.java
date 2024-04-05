package me.ve.smart.code;

import jakarta.persistence.*;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.util.Date;

@Entity
@Table(name = "tasks")
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class TaskEntity {


    @Id
    @Column(name = "task_id")
    private String taskId;

    @Column(name = "input_path")
    private String inputPath;

    @Column(name = "output_path")
    private String outputPath;

    @Column(name = "code_language")
    private String language;

    @Column(name = "qlpack")
    private String qlpack;

    @Enumerated(EnumType.STRING)
    @Column(name = "task_type")
    private TaskTypeEnum taskType;

    @Enumerated(EnumType.STRING)
    @Column(name = "current_step")
    private TaskStatusEnum currentStep;

    @Column(name = "created_at")
    private Date createdAt;

    @Column(name = "updated_at")
    private Date updatedAt;


}
