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

import java.util.ArrayList;
import java.util.Collection;
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
     *
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

    /**
     * 获取所有的银行信息，返回一个list
     *
     * @param listener
     */
    @Override
    public void getBank(final QuickPayCallbackListener<List<Bank>> listener) {
        new AsyncTask<Void, Integer, AsyncTaskResult<List<Bank>>>() {
            @Override
            protected AsyncTaskResult<List<Bank>> doInBackground(Void... params) {
                try {
                    Map<String, Bank> bankMap = bankDataApi.getBank();//这个返回的是个map
                    Collection<Bank> bankCollection = bankMap.values();
                    List<Bank> bankList = new ArrayList<Bank>();
                    for (Bank b : bankCollection) {
                        bankList.add(b);
                    }
                    return new AsyncTaskResult<List<Bank>>(bankList);
                } catch (QuickPayException ex) {
                    return new AsyncTaskResult<List<Bank>>(ex);
                }
            }

            @Override
            protected void onPostExecute(AsyncTaskResult<List<Bank>> result) {
                if (result.getException() != null) {
                    listener.onFailure(result.getException());
                } else {
                    listener.onSuccess(result.getResult());
                }
            }
        }.execute();
    }

    /**
     * 传人cityCode和bankId，去服务器查询分行的信息
     * @param cityCode
     * @param bankId
     * @param listener
     */
    @Override
    public void getBranchBank(final String cityCode, final String bankId, final QuickPayCallbackListener<List<SubBank>> listener) {
        new AsyncTask<Void, Integer, AsyncTaskResult<List<SubBank>>>() {
            @Override
            protected AsyncTaskResult<List<SubBank>> doInBackground(Void... params) {
                try {
                    List<SubBank> subBankList = bankDataApi.getBranchBank(cityCode, bankId);
                    return new AsyncTaskResult<List<SubBank>>(subBankList);
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
