package com.example.model;

import java.io.Serializable;
import java.util.ArrayList;

public class AppInfo implements Serializable {

    private static final long serialVersionUID = 1L;

    public AppInfo() {
        super();
        this.appName = "";
        this.packageName = "";
        this.versionName = "";
        this.versionCode = 0;
        this.launchTimes = 0;
        this.killTimes = 0;
        this.resultArrayList = new ArrayList<KillLaunchInfo>();
    }

    public AppInfo(String appName, String packageName, String versionName, int versionCode) {
        super();
        this.appName = appName;
        this.packageName = packageName;
        this.versionName = versionName;
        this.versionCode = versionCode;
        this.launchTimes = 0;
        this.killTimes = 0;
        this.resultArrayList = new ArrayList<KillLaunchInfo>();
    }

    public String getAppName() {
        return appName;
    }

    public void setAppName(String appName) {
        this.appName = appName;
    }

    public String getPackageName() {
        return packageName;
    }

    public void setPackageName(String packageName) {
        this.packageName = packageName;
    }

    public String getVersionName() {
        return versionName;
    }

    public void setVersionName(String versionName) {
        this.versionName = versionName;
    }

    public int getVersionCode() {
        return versionCode;
    }

    public void setVersionCode(int versionCode) {
        this.versionCode = versionCode;
    }

    public int getLaunchTimes() {
        return launchTimes;
    }

    public void setLaunchTimes(int launchTimes) {
        this.launchTimes = launchTimes;
    }

    public int getKillTimes() {
        return killTimes;
    }

    public void setKillTimes(int killTimes) {
        this.killTimes = killTimes;
    }

    public ArrayList<KillLaunchInfo> getResultArrayList() {
        return resultArrayList;
    }

    public void setResultArrayList(ArrayList<KillLaunchInfo> resultArrayList) {
        this.resultArrayList = resultArrayList;
    }

    private String appName;
    private String packageName;
    private String versionName;
    private int versionCode;
    private int launchTimes = 0;
    private int killTimes = 0;
    private ArrayList<KillLaunchInfo> resultArrayList;
}
