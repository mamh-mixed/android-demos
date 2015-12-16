package com.cardinfolink.yunshouyin.util;

import com.cardinfolink.yunshouyin.R;

public class ErrorUtil {

    public static String getErrorString(String error) {
        if (error.equals("user_no_activate")) {
            return ShowMoneyApp.getResString(R.string.alert_error_user_no_activate);
        }

        if (error.equals("username_password_error")) {
            return ShowMoneyApp.getResString(R.string.alert_error_username_password_error);
        }

        if (error.equals("sign_fail")) {
            return ShowMoneyApp.getResString(R.string.alert_error_sign_fail);
        }

        if (error.equals("username_no_exist")) {
            return ShowMoneyApp.getResString(R.string.alert_error_username_no_exist);
        }

        if (error.equals("username_exist")) {
            return ShowMoneyApp.getResString(R.string.alert_error_username_exist);
        }

        if (error.equals("system_error")) {
            return ShowMoneyApp.getResString(R.string.alert_error_system_error);
        }
        if (error.equals("old_password_error")) {
            return ShowMoneyApp.getResString(R.string.alert_error_old_password_error);
        }
        if (error.equals("user_has_three_times")) {
            return ShowMoneyApp.getResString(R.string.alert_error_user_has_three_times);
        }
        if (error.equals("user_has_two_times")) {
            return ShowMoneyApp.getResString(R.string.alert_error_user_has_two_times);
        }
        if (error.equals("user_has_one_times")) {
            return ShowMoneyApp.getResString(R.string.alert_error_user_has_one_times);
        }
        return error;

    }


}
