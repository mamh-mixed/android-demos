package com.cardinfolink.yunshouyin.salesman.core;


import android.os.AsyncTask;

import com.cardinfolink.yunshouyin.salesman.api.BankDataApi;
import com.cardinfolink.yunshouyin.salesman.api.BankDataApiImpl;
import com.cardinfolink.yunshouyin.salesman.api.QuickPayConfigStorage;
import com.cardinfolink.yunshouyin.salesman.api.QuickPayException;
import com.cardinfolink.yunshouyin.salesman.model.Bank;
import com.cardinfolink.yunshouyin.salesman.model.City;
import com.cardinfolink.yunshouyin.salesman.model.SubBank;

import java.util.List;
import java.util.Map;


/**
 * BankDataService接口的实现子类
 * Created by mamh on 15-11-24.
 */
public class BankDataServiceImpl implements BankDataService {
    private BankDataApi bankDataApi;
    private QuickPayConfigStorage quickPayConfigStorage;

    public BankDataServiceImpl(QuickPayConfigStorage quickPayConfigStorage) {
        this.bankDataApi = new BankDataApiImpl(quickPayConfigStorage);
        this.quickPayConfigStorage = quickPayConfigStorage;
    }

    @Override
    public void getProvince(final QuickPayCallbackListener<List<String>> listener) {
        new AsyncTask<Void, Integer, AsyncTaskResult<List<String>>>() {
            @Override
            protected AsyncTaskResult<List<String>> doInBackground(Void... params) {
                try {
                    List<String> province = bankDataApi.getProvince();
                    // TODO: 15-11-24 在这里做缓存会好些？？
                    return new AsyncTaskResult<List<String>>(province);
                } catch (QuickPayException ex) {
                    return new AsyncTaskResult<List<String>>(ex);
                }
            }

            @Override
            protected void onPostExecute(AsyncTaskResult<List<String>> result) {
                if (result.getException() != null) {
                    listener.onFailure(result.getException());
                } else {
                    listener.onSuccess(result.getResult());
                }
            }
        }.execute();

    }

    @Override
    public void getCity(String province, QuickPayCallbackListener<List<City>> listener) {

    }

    @Override
    public void getBank(QuickPayCallbackListener<Map<String, Bank>> listener) {

    }

    @Override
    public void search(String cityCode, String bankId, QuickPayCallbackListener<List<SubBank>> listener) {

    }
}
