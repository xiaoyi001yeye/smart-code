package me.ve.smart.code;

import com.github.dockerjava.api.DockerClient;
import com.github.dockerjava.api.command.CreateContainerResponse;
import com.github.dockerjava.api.model.Bind;
import com.github.dockerjava.api.model.HostConfig;
import com.github.dockerjava.api.model.Volume;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

@Service
@Slf4j
public class DockerService {

    private static final String CODEQL_IMAGE_NAME = "mcr.microsoft.com/cstsectools/codeql-container:latest";

    private static final String CODEQL_CONTAINER_NAME = "codeql-container";

    private final DockerClient dockerClient;

    public DockerService(DockerClient dockerClient) {
        this.dockerClient = dockerClient;
    }


    public void analyse(TaskEntity task) {
        String createDatabaseEnv = "CODEQL_CLI_ARGS=database create --language=" + task.getLanguage() + " /opt/results/source_db -s /opt/src";
        String upgradeDatabaseEnv = "CODEQL_CLI_ARGS=database upgrade /opt/results/source_db";
        String analyzeEnv = "CODEQL_CLI_ARGS=database analyze /opt/results/source_db --format=sarifv2 --output=/opt/results/issues.sarif " + task.getLanguage() + "-" + task.getQlpack() + ".qls";
        createAndStartContainer(task.getInputPath(), task.getOutputPath(), createDatabaseEnv);
        createAndStartContainer(task.getInputPath(), task.getOutputPath(), upgradeDatabaseEnv);
        createAndStartContainer(task.getInputPath(), task.getOutputPath(), analyzeEnv);

    }


    public void createAndStartContainer(String inputPath, String outputPath, String envStr) {
        Volume volume1 = new Volume(inputPath);
        Volume volume2 = new Volume(outputPath);
        CreateContainerResponse containerResponse = dockerClient.createContainerCmd(CODEQL_IMAGE_NAME)
                .withName(CODEQL_CONTAINER_NAME)
                .withEnv(envStr)
                .withHostConfig(HostConfig.newHostConfig().withBinds(new Bind("/opt/src", volume1), new Bind("/opt/results", volume2)))
                .exec();
        String containerId = containerResponse.getId();
        dockerClient.startContainerCmd(containerId).exec();
        Integer statusCode = dockerClient.waitContainerCmd(containerId).start().awaitStatusCode();
    }
}
