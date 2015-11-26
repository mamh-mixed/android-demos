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
    public void getCity(final String province, final QuickPayCallbackListener<List<City>> listener) {

        new AsyncTask<Void, Integer, AsyncTaskResult<List<City>>>() {
            @Override
            protected AsyncTaskResult<List<City>> doInBackground(Void... params) {
                try {
                    List<City> city = bankDataApi.getCity(province);
                    // TODO: 15-11-24 在这里做缓存会好些？？
                    return new AsyncTaskResult<List<City>>(city);
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
    public void getBank(final QuickPayCallbackListener<Map<String, Bank>> listener) {
        new AsyncTask<Void, Integer, AsyncTaskResult<Map<String, Bank>>>() {
            @Override
            protected AsyncTaskResult<Map<String, Bank>> doInBackground(Void... params) {
                try {
                    Map<String, Bank> bank = bankDataApi.getBank();
                    // TODO: 15-11-24 在这里做缓存会好些？？
                    return new AsyncTaskResult<Map<String, Bank>>(bank);
                } catch (QuickPayException ex) {
                    return new AsyncTaskResult<Map<String, Bank>>(ex);
                }
            }

            @Override
            protected void onPostExecute(AsyncTaskResult<Map<String, Bank>> result) {
                if (result.getException() != null) {
                    listener.onFailure(result.getException());
                } else {
                    listener.onSuccess(result.getResult());
                }
            }
        }.execute();
    }

    @Override
    public void getBranchBank(final String cityCode, final String bankId, final QuickPayCallbackListener<List<SubBank>> listener) {
        new AsyncTask<Void, Integer, AsyncTaskResult<List<SubBank>>>() {
            @Override
            protected AsyncTaskResult<List<SubBank>> doInBackground(Void... params) {
                try {
                    List<SubBank> branchBank = bankDataApi.getBranchBank(cityCode, bankId);
                    // TODO: 15-11-24 在这里做缓存会好些？？
                    return new AsyncTaskResult<List<SubBank>>(branchBank);
                } catch (QuickPayException ex) {
                    return new AsyncTaskResult<>(ex);
                }
            }

            @Override
            protected void onPostExecute(AsyncTaskResult<List<SubBank>> result) {
                if (result.getException() != null) {
                    listener.onFailure(result.getException());
                } else {
                    listener.onSuccess(result.getResult());
                }
            }
        }.execute();
    }
}
