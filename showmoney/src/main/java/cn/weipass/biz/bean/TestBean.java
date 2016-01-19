package cn.weipass.biz.bean;

/**
 * Created by apple on 16/1/19.
 */
public class TestBean {


    /**
     * 下单支付返回结果
     * 
     * txndir : A
     * busicd : REFD
     * respcd : R6
     * inscd : 99911888
     * mchntid : 999118880000338
     * terminalid : phonechan@foxmail.com
     * txamt : 000000000001
     * errorDetail : CANCEL_TIME_ERROR
     * orderNum : 16011917140429378
     * origOrderNum : 16011612403234497
     * sign : a96588036a44cb27f22e31f865322491e83b7c2f
     */

    private String txndir;
    private String busicd;
    private String respcd;
    private String inscd;
    private String mchntid;
    private String terminalid;
    private String txamt;
    private String errorDetail;
    private String orderNum;
    private String origOrderNum;
    private String sign;

    public void setTxndir(String txndir) {
        this.txndir = txndir;
    }

    public void setBusicd(String busicd) {
        this.busicd = busicd;
    }

    public void setRespcd(String respcd) {
        this.respcd = respcd;
    }

    public void setInscd(String inscd) {
        this.inscd = inscd;
    }

    public void setMchntid(String mchntid) {
        this.mchntid = mchntid;
    }

    public void setTerminalid(String terminalid) {
        this.terminalid = terminalid;
    }

    public void setTxamt(String txamt) {
        this.txamt = txamt;
    }

    public void setErrorDetail(String errorDetail) {
        this.errorDetail = errorDetail;
    }

    public void setOrderNum(String orderNum) {
        this.orderNum = orderNum;
    }

    public void setOrigOrderNum(String origOrderNum) {
        this.origOrderNum = origOrderNum;
    }

    public void setSign(String sign) {
        this.sign = sign;
    }

    public String getTxndir() {
        return txndir;
    }

    public String getBusicd() {
        return busicd;
    }

    public String getRespcd() {
        return respcd;
    }

    public String getInscd() {
        return inscd;
    }

    public String getMchntid() {
        return mchntid;
    }

    public String getTerminalid() {
        return terminalid;
    }

    public String getTxamt() {
        return txamt;
    }

    public String getErrorDetail() {
        return errorDetail;
    }

    public String getOrderNum() {
        return orderNum;
    }

    public String getOrigOrderNum() {
        return origOrderNum;
    }

    public String getSign() {
        return sign;
    }
}
