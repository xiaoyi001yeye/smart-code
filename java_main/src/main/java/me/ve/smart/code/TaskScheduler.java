package me.ve.smart.code;


import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;


public class TaskScheduler {


    public void start() {
        ExecutorService executorService = Executors.newSingleThreadExecutor();
        executorService.execute(() -> {

            for (; ; ) {
                System.out.println(1);
                try {
                    Thread.sleep(5000);
                } catch (InterruptedException e) {
                    throw new RuntimeException(e);
                }
            }


        });
    }
}
