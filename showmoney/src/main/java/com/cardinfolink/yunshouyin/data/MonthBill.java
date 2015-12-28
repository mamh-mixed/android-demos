package com.cardinfolink.yunshouyin.data;

/**
 * 月账单，用来记录一个月账单的汇总
 * Created by mamh on 15-12-26.
 */
public class MonthBill {

    private String currentYear;//记录一下当前的年份

    private String currentMonth;//记录一下这是哪个月的账单

    private int count;//返回的条数
    private String total;

    private int refdcount;//这个是和退款相关的
    private String refdtotal;

    private int size;

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
}
