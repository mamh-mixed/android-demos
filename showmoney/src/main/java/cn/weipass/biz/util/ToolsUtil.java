package cn.weipass.biz.util;

import android.content.Context;

import com.cardinfolink.cashiersdk.model.ResultData;
import com.cardinfolink.yunshouyin.data.SessonData;
import com.cardinfolink.yunshouyin.data.TradeBill;
import com.cardinfolink.yunshouyin.util.Log;

import cn.weipass.pos.sdk.IPrint;
import cn.weipass.pos.sdk.Printer;

/**
 * Created by apple on 16/1/19.
 */
public class ToolsUtil {

    public static final String TAG = "ToolsUtil";
    public static final int mediumSize = 10 * 2;

    private static final int TYPE_RECEIPT_PAY = 100;//支付凭条
    private static final int TYPE_RECEIPT_REFUND = 101;//退款凭条
    private static final int TYPE_RECEIPT_TICKET = 102;//核券凭条


    /**
     * 获取凭条类型
     *
     * @param data
     * @return
     */
    private static int getReceiptType(ResultData data) {

        if (data == null || data.busicd == null) {
            return 0;
        } else {

            switch (data.busicd) {
                case "VERI":
                    return TYPE_RECEIPT_TICKET;
                case "REFD":
                    return TYPE_RECEIPT_REFUND;
                case "PURC":
                    return TYPE_RECEIPT_PAY;
//                case "PAUT":
//                    return TYPE_RECEIPT_PAY;
                default:
                    return 0;
            }
        }


    }

    /**
     * 打印机的错误信息
     *
     * @param what
     * @param info
     * @return
     */
    public static String getPrintErrorInfo(int what, String info) {
        String message = "";
        switch (what) {
            case IPrint.EVENT_CONNECT_FAILD:
                message = "连接打印机失败";
                break;
            case IPrint.EVENT_CONNECTED:
                // Log.e("subscribe_msg", "连接打印机成功");
                break;
            case IPrint.EVENT_PAPER_JAM:
                message = "打印机卡纸";
                break;
            case IPrint.EVENT_UNKNOW:
                message = "打印机未知错误";
                break;
            case IPrint.EVENT_OK:
                // 回调函数中不能做UI操作，所以可以使用runOnUiThread函数来包装一下代码块
                // Log.e("subscribe_msg", "打印机正常");
                break;
            case IPrint.EVENT_NO_PAPER:
                message = "打印机缺纸";
                break;
            case IPrint.EVENT_HIGH_TEMP:
                message = "打印机高温";
                break;
        }

        return message;
    }


    /**
     * 设置空白长度
     *
     * @param size
     * @return
     */
    public static String getBlankBySize(int size) {
        String resultStr = "";
        for (int i = 0; i < size; i++) {
            resultStr += " ";
        }
        return resultStr;
    }


    public static boolean isLetter(char c) {
        int k = 0x80;
        return c / k == 0 ? true : false;
    }

    /**
     * 得到一个字符串的长度,显示的长度,一个汉字或日韩文长度为2,英文字符长度为1
     *
     * @param String s 需要得到长度的字符串
     * @return int 得到的字符串长度
     */
    public static int length(String s) {
        if (s == null)
            return 0;
        char[] c = s.toCharArray();
        int len = 0;
        for (int i = 0; i < c.length; i++) {
            len++;
            if (!isLetter(c[i])) {
                len++;
            }
        }
        return len;
    }

    /**
     * \n 代码换行
     *
     * @param context
     * @param printer
     */
    public static void printNormal(Context context, Printer printer, ResultData resultData) {

        Log.i(TAG, "print info 打印内容:" + resultData);
        // 标准打印，每个字符打印所占位置可能有一点出入（尤其是英文字符）
        String mediumSpline = "";
        for (int i = 0; i < mediumSize - 5; i++) {
            mediumSpline += "-";
        }

        String achd = "";
        switch (resultData.chcd) {
            case "WXP":
                achd = "微信";
                break;
            case "ALP":
                achd = "支付宝";
                break;
            case "ULIVE":
                achd = "";
                break;
            default:
                achd = "";
                break;
        }

        printer.printText("签购单（顾客存根）\n" + mediumSpline,
                Printer.FontFamily.SONG, Printer.FontSize.LARGE,
                Printer.FontStyle.NORMAL, Printer.Gravity.CENTER);
        printer.printText("商户名称:" + SessonData.loginUser.getMerName(),
                Printer.FontFamily.SONG, Printer.FontSize.MEDIUM,
                Printer.FontStyle.NORMAL, Printer.Gravity.LEFT);
        printer.printText("商户编号:" + SessonData.loginUser.getUsername(),
                Printer.FontFamily.SONG, Printer.FontSize.MEDIUM,
                Printer.FontStyle.NORMAL, Printer.Gravity.LEFT);
        printer.printText("收银员:" + SessonData.loginUser.getUsername(),
                Printer.FontFamily.SONG, Printer.FontSize.MEDIUM,
                Printer.FontStyle.NORMAL, Printer.Gravity.LEFT);

        switch (getReceiptType(resultData)) {
            case TYPE_RECEIPT_PAY:
                printer.printText("交易类型:" + achd + "    扫码付",
                        Printer.FontFamily.SONG, Printer.FontSize.MEDIUM,
                        Printer.FontStyle.NORMAL, Printer.Gravity.LEFT);
                printer.printText("日期时间:" + resultData.expDate,
                        Printer.FontFamily.SONG, Printer.FontSize.MEDIUM,
                        Printer.FontStyle.NORMAL, Printer.Gravity.LEFT);
                printer.printText("交易账号:" + resultData.consumerAccount,
                        Printer.FontFamily.SONG, Printer.FontSize.MEDIUM,
                        Printer.FontStyle.NORMAL, Printer.Gravity.LEFT);
                printer.printText("渠道订单号:" + resultData.channelOrderNum,
                        Printer.FontFamily.SONG, Printer.FontSize.MEDIUM,
                        Printer.FontStyle.NORMAL, Printer.Gravity.LEFT);
                printer.printText("商家订单号:" + resultData.orderNum,
                        Printer.FontFamily.SONG, Printer.FontSize.MEDIUM,
                        Printer.FontStyle.NORMAL, Printer.Gravity.LEFT);
                printer.printText("金额: RMB " + resultData.txamt,
                        Printer.FontFamily.SONG, Printer.FontSize.MEDIUM,
                        Printer.FontStyle.NORMAL, Printer.Gravity.LEFT);
                break;
            case TYPE_RECEIPT_REFUND:
                printer.printText("交易类型:" + achd + "    退款",
                        Printer.FontFamily.SONG, Printer.FontSize.MEDIUM,
                        Printer.FontStyle.NORMAL, Printer.Gravity.LEFT);
                printer.printText("日期时间:" + resultData.expDate,
                        Printer.FontFamily.SONG, Printer.FontSize.MEDIUM,
                        Printer.FontStyle.NORMAL, Printer.Gravity.LEFT);
                printer.printText("交易账号:" + resultData.consumerAccount,
                        Printer.FontFamily.SONG, Printer.FontSize.MEDIUM,
                        Printer.FontStyle.NORMAL, Printer.Gravity.LEFT);
                printer.printText("渠道订单号:" + resultData.channelOrderNum,
                        Printer.FontFamily.SONG, Printer.FontSize.MEDIUM,
                        Printer.FontStyle.NORMAL, Printer.Gravity.LEFT);
                printer.printText("商家订单号:" + resultData.orderNum,
                        Printer.FontFamily.SONG, Printer.FontSize.MEDIUM,
                        Printer.FontStyle.NORMAL, Printer.Gravity.LEFT);
                printer.printText("原交易订单号:" + resultData.origOrderNum,
                        Printer.FontFamily.SONG, Printer.FontSize.MEDIUM,
                        Printer.FontStyle.NORMAL, Printer.Gravity.LEFT);
                printer.printText("金额: RMB -" + resultData.txamt,
                        Printer.FontFamily.SONG, Printer.FontSize.MEDIUM,
                        Printer.FontStyle.NORMAL, Printer.Gravity.LEFT);

                break;
            case TYPE_RECEIPT_TICKET:
                printer.printText("交易类型:" + "   卡券核销",
                        Printer.FontFamily.SONG, Printer.FontSize.MEDIUM,
                        Printer.FontStyle.NORMAL, Printer.Gravity.LEFT);
                printer.printText("日期时间:" + resultData.expDate,
                        Printer.FontFamily.SONG, Printer.FontSize.MEDIUM,
                        Printer.FontStyle.NORMAL, Printer.Gravity.LEFT);
                printer.printText("商家订单号:" + resultData.orderNum,
                        Printer.FontFamily.SONG, Printer.FontSize.MEDIUM,
                        Printer.FontStyle.NORMAL, Printer.Gravity.LEFT);
                printer.printText("卡券号:" + resultData.orderNum,
                        Printer.FontFamily.SONG, Printer.FontSize.MEDIUM,
                        Printer.FontStyle.NORMAL, Printer.Gravity.LEFT);
                printer.printText("详情:" + resultData.cardId,
                        Printer.FontFamily.SONG, Printer.FontSize.MEDIUM,
                        Printer.FontStyle.NORMAL, Printer.Gravity.LEFT);

                break;
            default:
                break;
        }

        printer.printText("\n\n\n\n\n",
                Printer.FontFamily.SONG, Printer.FontSize.LARGE,
                Printer.FontStyle.NORMAL, Printer.Gravity.LEFT);
    }
}
