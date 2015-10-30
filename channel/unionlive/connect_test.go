package unionlive

// http://d.umq.me/PosService/CouponsPurchaseService.ashx?t=PurchaseCoupons  电子券验证
// http://d.umq.me/PosService/CouponsPurchaseService.ashx?t=QueryPurchaseCouponsResult 电子券验证结果查询
// http://d.umq.me/PosService/CouponsPurchaseService.ashx?t=QueryPurchaseLog 商户券验证流水查询
/*
http://d.umq.me/PosService/CouponsPurchaseService.ashx?t=PurchaseCoupons  电子券验证
样例报文:{"header":{"version":"1.0","transType":"W412","transDirect":"Q","sessionId":"747f6cf7-dadf-46ef-83e9-d3c0a87b3dbf","merchantId":"182000001000000","submitTime":"20130501201012","clientTraceNo":"497540"},"body":{"couponsNo":"1809706004000705","termId":"00000667","termSn":"9e908a255b3e5989","amount":"1"}}

http://d.umq.me/PosService/CouponsPurchaseService.ashx?t=QueryPurchaseCouponsResult 电子券验证结果查询
样例报文:{"header":{"version":"1.0","transType":"W394","transDirect":"Q","sessionId":"747f6cf7-dadf-46ef-83e9-d3c0a87b3dbf","merchantId":"182000001000000","submitTime":"20130501201012","clientTraceNo":"497540"},"body":{"couponsNo":"1810068010108200","termId":"00000667","termSn":"9e908a255b3e5989","amount":"1","oldClientTraceNo":"497540","oldSubmitTime":"20151021144017"}}
http://d.umq.me/PosService/CouponsPurchaseService.ashx?t=QueryPurchaseLog 商户券验证流水查询
样例报文:{"header":{"version":"1.0","transType":"W395","transDirect":"Q","sessionId":"747f6cf7-dadf-46ef-83e9-d3c0a87b3dbf","merchantId":"182000001000000","submitTime":"20130501201012","clientTraceNo":"497540"},"body":{"termId":"00000667","termSn":"9e908a255b3e5989","pageIndex":"1"}}

http://d.umq.me/PosService/CouponsPurchaseService.ashx?t=PurchaseCoupons&channelId=182000899000001
http://d.umq.me/PosService/CouponsPurchaseService.ashx?t=QueryPurchaseCouponsResult&channelId=182000899000001 电子券验证结果查询
http://d.umq.me/PosService/CouponsPurchaseService.ashx?t=QueryPurchaseLog&channelId=182000899000001

其中商户编号:182000001000000
终端编号:00000667
终端硬件序列号:9e908a255b3e5989

已使用的消费的券号:
1808700004000875
1805702004000605
1806706004000405
1802702004000305
1806704504000709
1802708104000667

//提供高勃测试的券码
1802708004000235
1804700004000125
1802709204000517
1807709204000407
1803704604000700
1803700604000630
1809707604000570
1804704604000430
1809703604000320
1805709604000220
//提供给罗强测试的券码
1802704104000597
1808707104000287
1804703104000187
1803700404000144
1804700404000863

可用于消费的券号:
1808709404000713
1809707504000801
1805703504000701
1801709504000601
1806706504000501
1806703504000431
1806702004000834
1801707004000744
1801705004000614
1805703004000514
1804709004000314
1806707004000929
1801707004000679
1806704004000579
1806700004000329
1802707104000831
1809706104000220
1805702104000120
1804702004000066
1805701004000779
1803707004000948
1804708104000050
1808700004000875
1805702004000605
1806706004000405
1802702004000305
1802708004000235
1804700004000125
1802709204000517
1807709204000407
1803704604000700
1803700604000630
1809707604000570
1804704604000430
1809703604000320
1805709604000220
1806704504000709
1802708104000667
1802704104000597
1808707104000287
1804703104000187
1803700404000144
1804700404000863
1808709404000713
1809707504000801
1805703504000701
1801709504000601
1806706504000501
1806703504000431
1806702004000834
1801707004000744
1801705004000614
1805703004000514
1804709004000314
*/
