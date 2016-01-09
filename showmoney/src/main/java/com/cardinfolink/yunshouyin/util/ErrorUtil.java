package com.cardinfolink.yunshouyin.util;

import com.cardinfolink.yunshouyin.R;

public class ErrorUtil {

    public static String getErrorString(String error) {
        if ("user_no_activate".equals(error)) {
            return ShowMoneyApp.getResString(R.string.alert_error_user_no_activate);
        }

        if ("username_password_error".equals(error)) {
            return ShowMoneyApp.getResString(R.string.alert_error_username_password_error);
        }

        if ("sign_fail".equals(error)) {
            return ShowMoneyApp.getResString(R.string.alert_error_sign_fail);
        }

        if ("username_no_exist".equals(error)) {
            return ShowMoneyApp.getResString(R.string.alert_error_username_no_exist);
        }

        if ("username_exist".equals(error)) {
            return ShowMoneyApp.getResString(R.string.alert_error_username_exist);
        }

        if ("system_error".equals(error)) {
            return ShowMoneyApp.getResString(R.string.alert_error_system_error);
        }
        if ("old_password_error".equals(error)) {
            return ShowMoneyApp.getResString(R.string.alert_error_old_password_error);
        }
        if ("user_has_three_times".equals(error)) {
            return ShowMoneyApp.getResString(R.string.alert_error_user_has_three_times);
        }
        if ("user_has_two_times".equals(error)) {
            return ShowMoneyApp.getResString(R.string.alert_error_user_has_two_times);
        }
        if ("user_has_one_times".equals(error)) {
            return ShowMoneyApp.getResString(R.string.alert_error_user_has_one_times);
        }
        return error;

    }


}
