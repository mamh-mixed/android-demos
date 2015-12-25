package com.cardinfolink.cashiersdk.listener;

public interface CommunicationListener {
    void onResult(String result);

    void onError(int error);
}
