package com.cardinfolink.yunshouyin.salesman.utils;

/**
 * 将服务器code转换成用户友好的消息
 */
public class ErrorUtil {
    public static String getErrorString(String error) {
        if (error.equals("user_no_activate")) {
            return "账号未激活!";
        }

        if (error.equals("username_password_error")) {
            return "用户名或密码错误!";
        }

        if (error.equals("sign_fail")) {
            return "签名错误,报文被串改!";
        }

        if (error.equals("username_no_exist")) {
            return "用户名不存在!";
        }

        if (error.equals("username_exist")) {
            return "用户名已存在!";
        }

        if (error.equals("system_error")) {
            return "对不起,系统出现错误!";
        }
        if (error.equals("old_password_error")) {
            return "原密码错误!";
        }

        if (error.equals("accessToken_error")) {
            return "登录过期,请重新登录!";
        }

        //new added, 这些错误都是在何时触发的,需要梳理清楚.
        if (error.equals("params_empty")) {
            return "参数不能为空!";
        }

        if (error.equals("user_already_improved")) {
            return "用户数据已经更新了!";
        }
        if (error.equals("merId_no_exist")) {
            return "商户ID不存在!";
        }
        if (error.equals("user_data_error")) {
            return "用户数据错误!";
        }

        return error;
    }
}
