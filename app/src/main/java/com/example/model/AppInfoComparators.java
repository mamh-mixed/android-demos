package com.example.model;

import java.text.ParseException;
import java.text.SimpleDateFormat;
import java.util.Comparator;
import java.util.Date;

public class AppInfoComparators implements Comparator<Object> {

    @Override
    public int compare(Object lhs, Object rhs) {
        if (lhs instanceof AppInfo && rhs instanceof AppInfo) {
            return ((AppInfo) rhs).getLaunchTimes() - ((AppInfo) lhs).getLaunchTimes();
        } else if (lhs instanceof KillLaunchInfo && rhs instanceof KillLaunchInfo) {
            SimpleDateFormat df = new SimpleDateFormat("MM-dd HH:mm:ss.SSS");
            try {
                Date lhsDate = df.parse(((KillLaunchInfo) lhs).getKillorlaunchdatetime());
                Date rhsDate = df.parse(((KillLaunchInfo) rhs).getKillorlaunchdatetime());
                return (int) (lhsDate.getTime() - rhsDate.getTime());
            } catch (ParseException e) {
                e.printStackTrace();
            }
        }

        return 0;
    }
}
