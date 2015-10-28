package com.cardinfolink.yunshouyin.salesman.models;

import com.cardinfolink.yunshouyin.salesman.utils.HttpCommunicationUtil;
import com.cardinfolink.yunshouyin.salesman.utils.ParamsUtil;

public class SessonData {
    public static User loginUser = new User();
    public static int position_view = 0;

    public static String getAccessToken() {
        return loginUser.getAccessToken();
    }

    private static String uploadToken;

    //TODO: 过期怎么办
    public static String getUploadToken() {
        if (uploadToken == null || uploadToken.equals("")) {
            try {
                SAServerPacket serverPacket = HttpCommunicationUtil.getServerPacket(ParamsUtil.getUploadToken_SA(getAccessToken()));
                uploadToken = serverPacket.getUploadToken();
            } catch (Exception ex) {
                ex.printStackTrace();
            }

        }
        return uploadToken;
    }

    public static User registerUser;
}
