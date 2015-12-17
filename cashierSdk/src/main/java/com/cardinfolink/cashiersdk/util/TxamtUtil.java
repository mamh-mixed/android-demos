package com.cardinfolink.cashiersdk.util;

import java.math.BigDecimal;


public class TxamtUtil {

    public static String getTxamtUtil(String txamt) {
        try {
            // 这里乘以100，只取了小数点后两位。
            //所以这里乘以了100 然后取成long类型的。最后还要转换成12位的字符串。
            BigDecimal bg = new BigDecimal(txamt);
            BigDecimal bg100 = new BigDecimal("100"); //两个BigDecimal相乘，
            long longNum = bg.multiply(bg100).longValue();
            String strNum = String.format("%012d", longNum); //这里前面做了限制不能传很大的数过来，最大9位吧
            return strNum;
        } catch (Exception e) {
            e.printStackTrace();
        }
        return null;
    }


    public static String getNormal(String txamt) {
        if (txamt != null) {
            try {
                BigDecimal bg = new BigDecimal(txamt);
                BigDecimal bg100 = new BigDecimal("100");
                return bg.divide(bg100).toString();
            } catch (Exception e) {
                e.printStackTrace();
            }
        }
        return null;
    }

}
