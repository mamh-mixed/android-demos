package com.example.model;

import java.io.Serializable;

public class KillLaunchInfo implements Serializable {

    private static final long serialVersionUID = 1L;
    private String killorlaunch;
    private String killorlaunchdatetime;
    private long difftime;

    // app kill flag.
    public static final String killFlag = "Kill";;
    // app launch flag.
    public static final String launchFlag = "Launch";

    public KillLaunchInfo() {
        this.killorlaunch = "";
        this.killorlaunchdatetime = "";
        this.difftime = 0;
    }

    public KillLaunchInfo(String kl, String kldatetime) {
        this.killorlaunch = kl;
        this.killorlaunchdatetime = kldatetime;
    }

    public String getKillorlaunch() {
        return killorlaunch;
    }

    public void setKillorlaunch(String killorlaunch) {
        this.killorlaunch = killorlaunch;
    }

    public String getKillorlaunchdatetime() {
        return killorlaunchdatetime;
    }

    public void setKillorlaunchdatetime(String killorlaunchdatetime) {
        this.killorlaunchdatetime = killorlaunchdatetime;
    }

    public long getDifftime() {
        return difftime;
    }

    public void setDifftime(long difftime) {
        this.difftime = difftime;
    }
}
