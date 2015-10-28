package com.cardinfolink.cashiersdk.listener;

import com.cardinfolink.cashiersdk.model.ResultData;

public interface CashierListener {
    public void onResult(ResultData resultData);

    public void onError(int errorCode);
}
