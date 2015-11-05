package com.cardinfolink.yunshouyin.salesman.utils;

import com.cardinfolink.yunshouyin.salesman.model.SAServerPacket;

public interface CommunicationListenerV2 {
    // 对state==success的packet进行进一步操作
    void onResult(SAServerPacket serverPacket);

    // 只处理错误码
    void onError(String error);
}
