package com.cardinfolink.cashiersdk.listener;

public interface CommunicationListener {
    public void onResult(String result);

    public void onError(int error);
}
