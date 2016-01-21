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
        if ("NO_ENHANCE_LIMIT_AMT".equals(error)) {
            return ShowMoneyApp.getResString(R.string.alert_error_no_enhance_limit_amt);
        }
        if ("C0".equals(error)) {//字段不能为空
            return ShowMoneyApp.getResString(R.string.coupon_c0);
        }
        if ("C1".equals(error)) {//卡券已被核销
            return ShowMoneyApp.getResString(R.string.coupon_c1);
        }

        if ("C2".equals(error)) {//卡券已过期
            return ShowMoneyApp.getResString(R.string.coupon_c2);
        }
        if ("C3".equals(error)) {//无效的卡券
            return ShowMoneyApp.getResString(R.string.coupon_c3);
        }
        if ("C4".equals(error)) {//券状态异常
            return ShowMoneyApp.getResString(R.string.coupon_c4);
        }
        if ("C5".equals(error)) {//未到卡券使用时间
            return ShowMoneyApp.getResString(R.string.coupon_c5);
        }
        if ("C6".equals(error)) {//商户不能使用该卡券
            return ShowMoneyApp.getResString(R.string.coupon_c6);
        }
        if ("C7".equals(error)) {//金额达不到满足优惠条件的最小金额
            return ShowMoneyApp.getResString(R.string.coupon_c7);
        }
        if ("C8".equals(error)) {//不能用此支付类型
            return ShowMoneyApp.getResString(R.string.coupon_c8);
        }
        if ("C9".equals(error)) {//未找到原验证记录
            return ShowMoneyApp.getResString(R.string.coupon_c9);
        }


        return error;

    }


}
