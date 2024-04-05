package me.ve.smart.code.vo;


import jakarta.validation.constraints.NotNull;
import lombok.Data;


@Data
public class RunTaskVO {

    @NotNull(message = "language cannot be null")
    private String language;

    @NotNull(message = "inputPath cannot be null")
    private String inputPath;

    @NotNull(message = "outputPath cannot be null")
    private String outputPath;

    @NotNull(message = "qlpack cannot be null")
    private String qlpack;

}
