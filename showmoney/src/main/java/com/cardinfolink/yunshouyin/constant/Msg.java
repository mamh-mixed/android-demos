package com.cardinfolink.yunshouyin.constant;

public class Msg {
    /**
     * 发送这个四个不同的消息 会跳转到
     */
    public static final int MSG_GO_SCAN_CODE_VIEW = 91;
    public static final int MSG_GO_TICKET_VIEW = 92;
    public static final int MSG_GO_BILL_VIEW = 93;
    public static final int MSG_GO_MY_SETTING_VIEW = 94;

    public static final int MSG_SCAN_CODE_VIEW_CLEAR_INPUT_OUTPUT = 120;

    public static final int MSG_FROM_SCANCODE_SUCCESS = 131;

    public static final int MSG_FROM_DIGLOG_SECOND = 141;
    public static final int MSG_FROM_DIGLOG_CLOSE = 142;

    public static final int MSG_FROM_SERVER_TIMEOUT = 161;
    public static final int MSG_FROM_SERVER_TRADE_SUCCESS = 163;
    public static final int MSG_FROM_SERVER_TRADE_FAIL = 164;
    public static final int MSG_FROM_SERVER_TRADE_NOPAY = 165;

    //关单用到的
    public static final int MSG_FROM_SERVER_CLOSEBILL_SUCCESS = 500;
    public static final int MSG_FROM_SERVER_CLOSEBILL_DOING = 501;
    public static final int MSG_FROM_SERVER_CLOSEBILL_FAIL = 502;
    public static final int MSG_FROM_SEARCHING_POLLING = 503;

    public static final int MSG_CREATE_QR_SUCCESS = 600;
    public static final int MSG_CREATE_QR_FAIL = 601;

    //卡券用到的
    public static final int MSG_FROM_SERVER_COUPON_SUCCESS = 602;
    public static final int MSG_FROM_SERVER_COUPON_FAIL = 603;
    public static final int MSG_COUPON_CANCEL = 604;

    public static final int MSG_FINISH_BIG_SCANCODEVIEW = 605;

    /**
     * 退款的时候 使用 刷新账单界面
     */
    public static final int MSG_REFRESH_BILL_LIST_VIEW = 702;
}
