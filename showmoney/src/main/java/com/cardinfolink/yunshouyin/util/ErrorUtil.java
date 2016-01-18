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
        if ("user_lock".equals(error)) {
            return ShowMoneyApp.getResString(R.string.alert_error_user_lock);
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

        if ("DATA_FORMAT_ERROR".equals(error)) {
            return ShowMoneyApp.getResString(R.string.alert_error_data_format_error);
        }
        if ("QRCODE_INVALID".equals(error)) {
            return ShowMoneyApp.getResString(R.string.alert_error_qrcode_invalid);
        }

        if ("NO_CHANNEL".equals(error)) {
            return ShowMoneyApp.getResString(R.string.alert_error_no_channel);
        }

        if ("success".equals(error)) {
            return ShowMoneyApp.getResString(R.string.alert_error_success);
        }

        if ("params_empty".equals(error)) {
            return ShowMoneyApp.getResString(R.string.alert_error_params_empty);
        }

        if ("params_format_error".equals(error)) {
            return ShowMoneyApp.getResString(R.string.alert_error_params_format_error);
        }

        if ("user_already_improved".equals(error)) {
            return ShowMoneyApp.getResString(R.string.alert_error_user_already_improved);
        }

        if ("merId_no_exist".equals(error)) {
            return ShowMoneyApp.getResString(R.string.alert_error_merId_no_exist);
        }

        if ("company_login_name_error".equals(error)) {
            return ShowMoneyApp.getResString(R.string.alert_error_company_login_name_error);
        }
        if ("user_data_error".equals(error)) {
            return ShowMoneyApp.getResString(R.string.alert_error_user_data_error);
        }
        if ("accessToken_error".equals(error)) {
            return ShowMoneyApp.getResString(R.string.alert_error_accessToken_error);
        }
        if ("INVALID_REPORT_TYPE".equals(error)) {
            return ShowMoneyApp.getResString(R.string.alert_error_INVALID_REPORT_TYPE);
        }
        if ("CODE_PAYTYPE_NOT_MATCH".equals(error)) {
            return ShowMoneyApp.getResString(R.string.alert_error_CODE_PAYTYPE_NOT_MATCH);
        }
        return error;

    }


}
