package me.ve.smart.code;

import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.Id;
import jakarta.persistence.Table;
import lombok.Data;

import java.util.Date;

@Entity
@Table(name = "tasks")
@Data
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

    @Column(name = "task_type")
    private String taskType;

    @Column(name = "current_step")
    private String currentStep;

    @Column(name = "created_at")
    private Date createdAt;

    @Column(name = "updated_at")
    private Date updatedAt;


}
