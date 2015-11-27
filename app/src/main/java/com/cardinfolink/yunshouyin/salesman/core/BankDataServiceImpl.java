package com.cardinfolink.yunshouyin.salesman.core;


import android.content.Context;
import android.os.AsyncTask;

import com.cardinfolink.yunshouyin.salesman.api.BankDataApi;
import com.cardinfolink.yunshouyin.salesman.api.BankDataApiImpl;
import com.cardinfolink.yunshouyin.salesman.api.QuickPayConfigStorage;
import com.cardinfolink.yunshouyin.salesman.api.QuickPayException;
import com.cardinfolink.yunshouyin.salesman.db.SalesmanDB;
import com.cardinfolink.yunshouyin.salesman.model.Bank;
import com.cardinfolink.yunshouyin.salesman.model.City;
import com.cardinfolink.yunshouyin.salesman.model.Province;
import com.cardinfolink.yunshouyin.salesman.model.SubBank;
import com.cardinfolink.yunshouyin.salesman.utils.SalesmanApplication;

import java.util.ArrayList;
import java.util.Collection;
import java.util.Hashtable;
import java.util.Iterator;
import java.util.List;
import java.util.Map;


/**
 * BankDataService接口的实现子类
 * Created by mamh on 15-11-24.
 */
public class BankDataServiceImpl implements BankDataService {
    private BankDataApi bankDataApi;
    private QuickPayConfigStorage quickPayConfigStorage;
    private SalesmanDB salesmanDB;//使用数据库存储省份城市银行信息

    public BankDataServiceImpl(QuickPayConfigStorage quickPayConfigStorage) {
        this.salesmanDB = SalesmanDB.getInstance(SalesmanApplication.getInstance());
        this.bankDataApi = new BankDataApiImpl(quickPayConfigStorage);
        this.quickPayConfigStorage = quickPayConfigStorage;
    }

    public BankDataServiceImpl(Context context, QuickPayConfigStorage quickPayConfigStorage) {
        this.salesmanDB = SalesmanDB.getInstance(context);
        this.bankDataApi = new BankDataApiImpl(quickPayConfigStorage);
        this.quickPayConfigStorage = quickPayConfigStorage;
    }

    @Override
    public void getProvince(final QuickPayCallbackListener<List<Province>> listener) {
        new AsyncTask<Void, Integer, AsyncTaskResult<List<Province>>>() {
            @Override
            protected AsyncTaskResult<List<Province>> doInBackground(Void... params) {
                try {
                    List<Province> provinceList = salesmanDB.loadProvince();
                    if (provinceList == null || provinceList.size() <= 0) {
                        provinceList = bankDataApi.getProvince();
                    }
                    saveProvinces(provinceList);//save to database
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

    private synchronized void saveProvinces(List<Province> provinceList) {
        for (Province p : provinceList) {
            salesmanDB.saveProvince(p);
        }
    }

    @Override
    public void getCity(final String province, final QuickPayCallbackListener<List<City>> listener) {

        new AsyncTask<Void, Integer, AsyncTaskResult<List<City>>>() {
            @Override
            protected AsyncTaskResult<List<City>> doInBackground(Void... params) {
                try {
                    List<City> cityList = salesmanDB.loadCity(province);//要查出某个省下面的所有城市
                    if (cityList == null || cityList.size() <= 0) {
                        cityList = bankDataApi.getCity(province);
                    }
                    saveCities(cityList);
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

    private synchronized void saveCities(List<City> cities) {
        for (City c : cities) {
            salesmanDB.saveCity(c);
        }
    }

    @Override
    public void getBank(final QuickPayCallbackListener<List<Bank>> listener) {
        new AsyncTask<Void, Integer, AsyncTaskResult<List<Bank>>>() {
            @Override
            protected AsyncTaskResult<List<Bank>> doInBackground(Void... params) {
                try {
                    List<Bank> bankList = salesmanDB.loadBank();
                    if (bankList == null || bankList.size() <= 0) {
                        Map<String, Bank> bankMap = bankDataApi.getBank();//这个返回的是个map
                        Collection<Bank> bankss = bankMap.values();
                        bankList = new ArrayList<Bank>();
                        for (Bank b : bankss) {
                            bankList.add(b);
                        }
                    }

                    saveBanks(bankList);
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

    private synchronized void saveBanks(Collection<Bank> bankList) {
        for (Bank b : bankList) {
            salesmanDB.saveBank(b);
        }
    }

    @Override
    public void getBranchBank(final String cityCode, final String bankId, final QuickPayCallbackListener<List<SubBank>> listener) {
        new AsyncTask<Void, Integer, AsyncTaskResult<List<SubBank>>>() {
            @Override
            protected AsyncTaskResult<List<SubBank>> doInBackground(Void... params) {
                try {
                    List<SubBank> subBankList = salesmanDB.loadBranchBank(cityCode, bankId);
                    if (subBankList == null || subBankList.size() <= 0) {
                        subBankList = bankDataApi.getBranchBank(cityCode, bankId);
                    }
                    saveBranchBanks(bankId, subBankList);//把相应的大银行行号也存入
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

    private synchronized void saveBranchBanks(String bankId, List<SubBank> subBankList) {
        for (SubBank b : subBankList) {
            b.setBankId(bankId);
            salesmanDB.saveBranchBank(b);
        }
    }
}
