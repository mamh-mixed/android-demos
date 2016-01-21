package com.cardinfolink.cashiersdk.listener;

import com.cardinfolink.cashiersdk.model.ResultData;

public interface CashierListener {
    void onResult(ResultData resultData);

    void onError(int errorCode);
}
