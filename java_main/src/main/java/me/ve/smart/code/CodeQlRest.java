package me.ve.smart.code;

import jakarta.ws.rs.GET;
import jakarta.ws.rs.Path;
import jakarta.ws.rs.Produces;
import jakarta.ws.rs.core.MediaType;

@Path("/codql")
public class CodeQlRest {
    /**
     *
     */
    @GET
    @Produces(MediaType.TEXT_PLAIN)
    @Path("/welcome")
    public String welcome() {
        return "welcome to codeql java.";
    }

}
