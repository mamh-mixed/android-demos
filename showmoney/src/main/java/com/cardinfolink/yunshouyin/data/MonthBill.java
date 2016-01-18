package com.cardinfolink.yunshouyin.data;

/**
 * 月账单，用来记录一个月账单的汇总
 * Created by mamh on 15-12-26.
 */
public class MonthBill {

    private String currentYear;//记录一下当前的年份

    private String currentMonth;//记录一下这是哪个月的账单


    /**
     * {
     * "state": "success",
     * <p/>
     * <p/>
     * "total": "0.00",
     * "count": 25,
     * "totalRecord": 1471,
     * "size": 100,
     * "refdcount": 25,
     * "refdtotal": "175.18",
     * "txn": [这里是一个数组]
     */
    //下面这个几个都是 返回的json里面带的字段
    private String total;
    private int count;//返回的条数

    private int totalRecord;
    private int size;

    private int refdcount;//这个是和退款相关的
    private String refdtotal;

    private String nextMonth;//新的接口，返回下个月份，拉取账单的时候下次拉取的时候就会用到这个字段


    public MonthBill(String currentYear, String currentMonth) {
        this.currentYear = currentYear;
        this.currentMonth = currentMonth;
    }

    public String getCurrentYear() {
        return currentYear;
    }

    public void setCurrentYear(String currentYear) {
        this.currentYear = currentYear;
    }

    public String getCurrentMonth() {
        return currentMonth;
    }

    public void setCurrentMonth(String currentMonth) {
        this.currentMonth = currentMonth;
    }

    public MonthBill() {

    }

    public int getCount() {
        return count;
    }

    public void setCount(int count) {
        this.count = count;
    }

    public String getTotal() {
        return total;
    }

    public void setTotal(String total) {
        this.total = total;
    }

    public int getRefdcount() {
        return refdcount;
    }

    public void setRefdcount(int refdcount) {
        this.refdcount = refdcount;
    }

    public String getRefdtotal() {
        return refdtotal;
    }

    public void setRefdtotal(String refdtotal) {
        this.refdtotal = refdtotal;
    }

    public int getSize() {
        return size;
    }

    public void setSize(int size) {
        this.size = size;
    }

    public int getTotalRecord() {
        return totalRecord;
    }

    public void setTotalRecord(int totalRecord) {
        this.totalRecord = totalRecord;
    }

    public String getNextMonth() {
        return nextMonth;
    }

    public void setNextMonth(String nextMonth) {
        this.nextMonth = nextMonth;
    }
}
