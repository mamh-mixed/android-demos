package com.cardinfolink.cashiersdk.util;

import java.math.BigDecimal;


public class TxamtUtil {

    public static String getTxamtUtil(String txamt) {
        String str = txamt;
        try {
            double i = Double.parseDouble(str);
            BigDecimal bg = new BigDecimal(i);
            double j = bg.setScale(2, BigDecimal.ROUND_HALF_UP).doubleValue();
            j = j * 100;
            long num = (long) j;
            str = "" + num;
            int k = 12 - str.length();
            String sum = "";
            for (int l = 0; l < k; l++) {
                sum = sum + "0";
            }
            sum = sum + str;
            return sum;


        } catch (Exception e) {

            e.printStackTrace();
        }
        return null;
    }


    public static String getNormal(String txamt) {
        String str = txamt;
        if (str != null) {
            try {
                String sum = "";
                int index = 0;
                char c = str.charAt(index);
                while (c == '0') {
                    index++;
                    c = str.charAt(index);
                }
                sum = str.substring(index);
                double i = Double.parseDouble(sum);
                i = i / 100;
                BigDecimal bg = new BigDecimal(i);
                double j = bg.setScale(2, BigDecimal.ROUND_HALF_UP).doubleValue();
                sum = "" + j;
                return sum;


            } catch (Exception e) {

                e.printStackTrace();
            }
        }
        return null;
    }

}
