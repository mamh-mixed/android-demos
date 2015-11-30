package com.cardinfolink.yunshouyin.core;


import android.os.AsyncTask;

import com.cardinfolink.yunshouyin.api.BankDataApi;
import com.cardinfolink.yunshouyin.api.BankDataApiImpl;
import com.cardinfolink.yunshouyin.api.QuickPayConfigStorage;
import com.cardinfolink.yunshouyin.api.QuickPayException;
import com.cardinfolink.yunshouyin.model.Bank;
import com.cardinfolink.yunshouyin.model.City;
import com.cardinfolink.yunshouyin.model.Province;
import com.cardinfolink.yunshouyin.model.SubBank;

import java.util.List;
import java.util.Map;

//TODO: 参照QuickPayServiceImpl,在这里加入缓存效果
public class BankDataServiceImpl implements BankDataService {
    private BankDataApi bankDataApi;
    private QuickPayConfigStorage quickPayConfigStorage;


    public BankDataServiceImpl(QuickPayConfigStorage quickPayConfigStorage) {
        this.bankDataApi = new BankDataApiImpl(quickPayConfigStorage);
        this.quickPayConfigStorage = quickPayConfigStorage;
    }

    @Override
    public void getProvince(final QuickPayCallbackListener<List<Province>> listener) {
        new AsyncTask<Void, Integer, AsyncTaskResult<List<Province>>>() {
            @Override
            protected AsyncTaskResult<List<Province>> doInBackground(Void... params) {
                try {
                    List<Province> provinceList = bankDataApi.getProvince();
                    return new AsyncTaskResult<List<Province>>(provinceList);
                } catch (QuickPayException ex) {
                    return new AsyncTaskResult<List<Province>>(ex);
                }
            }

            @Override
            protected void onPostExecute(AsyncTaskResult<List<Province>> result) {
                if (result.getException() != null) {
                    listener.onFailure(result.getException());
                } else {
                    listener.onSuccess(result.getResult());
                }
            }
        }.execute();
    }

    /**
     * 传人一个省份，会去服务器查询出这个省份下的所有的城市。
     * @param province
     * @param listener
     */
    @Override
    public void getCity(final String province, final QuickPayCallbackListener<List<City>> listener) {
        new AsyncTask<Void, Integer, AsyncTaskResult<List<City>>>() {
            @Override
            protected AsyncTaskResult<List<City>> doInBackground(Void... params) {
                try {
                    List<City> cityList = bankDataApi.getCity(province);
                    return new AsyncTaskResult<List<City>>(cityList);
                } catch (QuickPayException ex) {
                    return new AsyncTaskResult<List<City>>(ex);
                }
            }

            @Override
            protected void onPostExecute(AsyncTaskResult<List<City>> result) {
                if (result.getException() != null) {
                    listener.onFailure(result.getException());
                } else {
                    listener.onSuccess(result.getResult());
                }
            }
        }.execute();
    }

    @Override
    public void getBank(QuickPayCallbackListener<Map<String, Bank>> listener) {

    }

    @Override
    public void getBranchBank(String city_code, String bank_id, QuickPayCallbackListener<List<SubBank>> listener) {

    }
}
