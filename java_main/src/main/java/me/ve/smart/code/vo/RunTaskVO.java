package me.ve.smart.code.vo;

import jakarta.ws.rs.FormParam;
import lombok.Data;

@Data
public class RunTaskVO {

    @FormParam("language")
    private String language;

    @FormParam("inputPath")
    private String inputPath;

    @FormParam("qlpack")
    private String qlpack;

}
